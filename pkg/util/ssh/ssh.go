/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package ssh

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
	"gopkg.in/go-playground/validator.v9"
	"k8s.io/apimachinery/pkg/util/wait"
	"tkestack.io/tke/pkg/util/log"
)

const (
	tmpDir = "/tmp"
)

type SSH struct {
	*Config
	authMethods []ssh.AuthMethod
	dialer      sshDialer
}

var _ Interface = &SSH{}

type Config struct {
	User       string `validate:"required"`
	Host       string `validate:"required"`
	Port       int    `validate:"required"`
	Sudo       bool
	Password   string
	PrivateKey []byte
	PassPhrase []byte
	// 150 seconds is longer than the underlying default TCP backoff delay (127
	// seconds). This timeout is only intended to catch otherwise uncaught hangs.
	DialTimeOut time.Duration
	Retry       int
	Proxy       Proxy
}

func (c *Config) addr() string {
	return net.JoinHostPort(c.Host, fmt.Sprintf("%d", c.Port))
}

func New(c *Config) (*SSH, error) {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return nil, err
	}
	if c.Password == "" && c.PrivateKey == nil {
		return nil, errors.New("password or privateKey at least one")
	}

	authMethods := make([]ssh.AuthMethod, 0)
	if c.Password != "" {
		authMethods = append(authMethods,
			ssh.Password(c.Password),
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
				if len(questions) == 1 {
					return []string{c.Password}, nil
				}
				return nil, nil
			}))
	}
	if len(c.PrivateKey) != 0 {
		signer, err := MakePrivateKeySigner(c.PrivateKey, c.PassPhrase)
		if err != nil {
			return nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if c.DialTimeOut == 0 {
		c.DialTimeOut = 5 * time.Second
	}

	if c.User != "root" {
		c.Sudo = true
	}

	return &SSH{
		Config:      c,
		authMethods: authMethods,
		dialer:      &timeoutDialer{&realSSHDialer{}, c.DialTimeOut},
	}, nil
}

func (s *SSH) CheckProxyTunnel() error {
	if s.Proxy != nil {
		return s.Proxy.CheckTunnel()
	}
	return fmt.Errorf("no proxy is set")
}

func (s *SSH) Ping() error {
	_, _, _, err := s.Exec("pwd")

	return err
}

func (s *SSH) CombinedOutput(cmd string) ([]byte, error) {
	stdout, stderr, exit, err := s.Exec(cmd)
	if err != nil {
		return nil, fmt.Errorf("exec cmd %q eror: %w", cmd, err)
	}
	if exit != 0 {
		return nil, fmt.Errorf("exec cmd %q eror: exit code %d: stderr %s", cmd, exit, stderr)
	}
	return []byte(stdout), nil
}

func (s *SSH) Execf(format string, a ...interface{}) (stdout string, stderr string, exit int, err error) {
	return s.Exec(fmt.Sprintf(format, a...))
}

func (s *SSH) Exec(cmd string) (stdout string, stderr string, exit int, err error) {
	if s.Sudo {
		cmd = fmt.Sprintf(`sudo bash << 'EOF'
%s
EOF
`, cmd)
	}
	log.Debugf("[%s] Exec %q", s.addr(), cmd)

	session, closer, err := s.newSession()
	if err != nil {
		return "", "", 0, err
	}
	defer closer()

	// Run the command.
	code := 0
	var bout, berr bytes.Buffer
	session.Stdout, session.Stderr = &bout, &berr
	if err = session.Run(cmd); err != nil {
		// Check whether the command failed to run or didn't complete.
		if exiterr, ok := err.(*ssh.ExitError); ok {
			// If we got an ExitError and the exit code is nonzero, we'll
			// consider the SSH itself successful (just that the command run
			// errored on the Host).
			if code = exiterr.ExitStatus(); code != 0 {
				err = nil
			}
		} else {
			// Some other kind of error happened (e.g. an IOError); consider the
			// SSH unsuccessful.
			err = fmt.Errorf("failed running `%s` on %s@%s: '%v'", cmd, s.User, s.addr(), err)
		}
	}
	return bout.String(), berr.String(), code, err
}

func (s *SSH) CopyFile(src, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	return s.WriteFile(file, dst)
}

func (s *SSH) CopyDir(src, dst string) error {
	files, er := ioutil.ReadDir(src)
	if er != nil {
		return er
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			s.CopyFile(filepath.Join(src, file.Name()), filepath.Join(dst, file.Name()))
		}
	}
	return nil
}

func (s *SSH) WriteFile(src io.Reader, dst string) error {
	tmpfile, err := ioutil.TempFile("", "*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := io.Copy(tmpfile, src); err != nil {
		return err
	}

	if _, err := tmpfile.Seek(0, os.SEEK_SET); err != nil {
		return err
	}

	needWriteFile, err := s.needWriteFile(tmpfile, dst)
	if err != nil {
		return err
	}
	if !needWriteFile {
		log.Debugf("[%s] Skip write %q because already existed", s.addr(), dst)
		return nil
	}

	if _, err := tmpfile.Seek(0, os.SEEK_SET); err != nil {
		return err
	}

	return s.writeFile(tmpfile, dst)
}

func (s *SSH) ReadFile(filename string) ([]byte, error) {
	return s.CombinedOutput(fmt.Sprintf("cat %s", filename))
}

func (s *SSH) Exist(filename string) (bool, error) {
	_, _, exit, err := s.Execf("ls %s", filename)
	if err != nil {
		return false, fmt.Errorf("ssh exec error: %w", err)
	}

	return exit == 0, nil
}

func (s *SSH) LookPath(file string) (string, error) {
	data, err := s.CombinedOutput(fmt.Sprintf("which %s", file))
	return string(data), err
}

func (s *SSH) ReadDir(dir string) (string, error) {
	data, err := s.CombinedOutput(fmt.Sprintf("ls %s", dir))
	return string(data), err
}

func (s *SSH) writeFile(src io.Reader, dst string) error {
	log.Debugf("[%s] Write data to %q", s.addr(), dst)

	sftpClient, closer, err := s.newSFTPClient()
	if err != nil {
		return err
	}
	defer closer()

	realDst := dst
	dst = path.Join(tmpDir, ksuid.New().String(), dst)
	err = sftpClient.MkdirAll(path.Dir(dst))
	if err != nil {
		return err
	}
	dstFile, err := sftpClient.Create(dst)
	if err != nil {
		return fmt.Errorf("create file error:%s:%s", dst, err)
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(src)
	if err != nil {
		return err
	}

	if err := dstFile.Close(); err != nil {
		return err
	}

	_, err = s.CombinedOutput(fmt.Sprintf("mkdir -p $(dirname %s) && mv %s %s && rm -rf $(dirname %s)", realDst, dst, realDst, dst))
	if err != nil {
		return err
	}

	return err
}

func (s *SSH) needWriteFile(src io.Reader, dst string) (bool, error) {
	srcHash := md5.New()
	if _, err := io.Copy(srcHash, src); err != nil {
		return false, err
	}

	hashFile := tmpDir + dst + ".md5"
	buffer := new(bytes.Buffer)
	buffer.WriteString(fmt.Sprintf("%x %s\n", srcHash.Sum(nil), dst))
	err := s.writeFile(buffer, hashFile)
	if err != nil {
		return false, err
	}

	_, err = s.CombinedOutput(fmt.Sprintf("md5sum --check --status %s", hashFile))
	if err == nil { // means dst exist and same as src
		return false, nil
	}

	return true, nil
}

func (s *SSH) newSFTPClient() (*sftp.Client, func(), error) {
	client, closer, err := s.newClient()
	if err != nil {
		return nil, nil, err
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return nil, nil, err
	}

	return sftpClient,
		func() {
			closer()
			sftpClient.Close()
		},
		nil
}

// newClient returns ssh session and closer which need defer run!
func (s *SSH) newSession() (*ssh.Session, func(), error) {
	client, closer, err := s.newClient()
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}

	return session,
		func() {
			closer()
			session.Close()
		},
		nil
}

// newClient returns ssh client and closer which need defer run!
func (s *SSH) newClient() (*ssh.Client, func(), error) {
	config := &ssh.ClientConfig{
		User:            s.User,
		Auth:            s.authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var client *ssh.Client
	var closer func()
	var err error

	if s.Proxy != nil {
		client, closer, err = s.proxyClientConn(config)
		// if retry is 0, loop will not stop
		if err != nil && s.Retry != 0 {
			wait.Poll(5*time.Second, time.Duration(s.Retry)*5*time.Second, func() (bool, error) {
				if client, closer, err = s.proxyClientConn(config); err != nil {
					return false, nil
				}
				err = nil
				return true, nil
			})
		}
	} else {
		client, err = s.dialer.Dial("tcp", s.addr(), config)
		// if retry is 0, loop will not stop
		if err != nil && s.Retry != 0 {
			wait.Poll(5*time.Second, time.Duration(s.Retry)*5*time.Second, func() (bool, error) {
				if client, err = s.dialer.Dial("tcp", s.addr(), config); err != nil {
					return false, nil
				}
				err = nil
				return true, nil
			})
		}
	}

	if err != nil {
		return nil, nil, err
	}

	return client,
		func() {
			if closer != nil {
				closer()
			}
			client.Close()
		},
		nil
}

func (s *SSH) proxyClientConn(config *ssh.ClientConfig) (*ssh.Client, func(), error) {
	conn, closer, err := s.Proxy.ProxyConn(s.addr())
	if err != nil {
		return nil, nil, err
	}
	nconn, chans, reqs, err := ssh.NewClientConn(conn, s.addr(), config)
	if err != nil {
		return nil, nil, err
	}
	return ssh.NewClient(nconn, chans, reqs), closer, nil

}

// Interface to allow mocking of ssh.Dial, for testing SSH
type sshDialer interface {
	Dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error)
}

// Real implementation of sshDialer
type realSSHDialer struct{}

var _ sshDialer = &realSSHDialer{}

func (d *realSSHDialer) Dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	dialer := proxy.FromEnvironmentUsing(&net.Dialer{Timeout: config.Timeout})
	conn, err := dialer.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	c, chans, reqs, err := ssh.NewClientConn(conn, addr, config)
	if err != nil {
		return nil, err
	}
	conn.SetReadDeadline(time.Time{})
	return ssh.NewClient(c, chans, reqs), nil
}

// timeoutDialer wraps an sshDialer with a timeout around Dial(). The golang
// ssh library can hang indefinitely inside the Dial() call (see issue #23835).
// Wrapping all Dial() calls with a conservative timeout provides safety against
// getting stuck on that.
type timeoutDialer struct {
	dialer  sshDialer
	timeout time.Duration
}

func (d *timeoutDialer) Dial(network, addr string, config *ssh.ClientConfig) (*ssh.Client, error) {
	config.Timeout = d.timeout
	return d.dialer.Dial(network, addr, config)
}
