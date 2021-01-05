/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2021 Tencent. All Rights Reserved.
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

package util

import (
	"bytes"
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	platformv1client "tkestack.io/tke/api/client/clientset/versioned/typed/platform/v1"
	utilcrypto "tkestack.io/tke/pkg/util/crypto"
)

const (
	SSHSecretName = "tke-ssh-key"
	SSHKeyName    = "aes.key"
)

type sshInfo struct {
	password   []byte
	privateKey []byte
	passPharse []byte
}

func CreateAesKeySecret(ctx context.Context, k8sClient kubernetes.Interface) (secret *corev1.Secret, err error) {
	secret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      SSHSecretName,
			Namespace: "tke",
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			SSHKeyName: []byte(utilcrypto.NewAesKey()),
		},
	}
	secret, err = k8sClient.CoreV1().Secrets(metav1.NamespaceSystem).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func EncryptAllSSH(ctx context.Context, platformClient platformv1client.PlatformV1Interface) (err error) {
	k8sClient, err := BuildExternalClientSetWithName(context.Background(), platformClient, "global")
	if err != nil {
		return err
	}
	secret, err := k8sClient.CoreV1().Secrets(metav1.NamespaceSystem).Get(ctx, SSHKeyName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	aesKey := string(secret.Data[SSHKeyName])

	clusters, err := platformClient.Clusters().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, cls := range clusters.Items {
		for i, mc := range cls.Spec.Machines {
			ssh := sshInfo{
				password:   mc.Password,
				privateKey: mc.PrivateKey,
				passPharse: mc.PassPhrase,
			}
			newSSH, err := encryptSSH(ssh, aesKey)
			if err != nil {
				return err
			}
			if bytes.Equal(ssh.password, newSSH.password) && bytes.Equal(ssh.password, newSSH.password) {
				continue
			}
			cls.Spec.Machines[i].Password = newSSH.password
			cls.Spec.Machines[i].PrivateKey = newSSH.privateKey
			cls.Spec.Machines[i].PassPhrase = newSSH.passPharse
		}
		_, err = platformClient.Clusters().Update(ctx, &cls, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}

	mcs, err := platformClient.Machines().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, mc := range mcs.Items {
		ssh := sshInfo{
			password:   mc.Spec.Password,
			privateKey: mc.Spec.PrivateKey,
			passPharse: mc.Spec.PassPhrase,
		}
		newSSH, err := encryptSSH(ssh, aesKey)
		if err != nil {
			return err
		}
		if bytes.Equal(ssh.password, newSSH.password) && bytes.Equal(ssh.password, newSSH.password) {
			continue
		}
		mc.Spec.Password = newSSH.password
		mc.Spec.PrivateKey = newSSH.privateKey
		mc.Spec.PassPhrase = newSSH.passPharse
		_, err = platformClient.Machines().Update(ctx, &mc, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func encryptSSH(ssh sshInfo, aesKey string) (sshInfo, error) {
	// ssh password case
	if len(ssh.password) != 0 {
		_, err := utilcrypto.AesDecrypt(string(ssh.password), aesKey)
		if err == nil {
			return ssh, nil
		}
		encryptPwd, err := utilcrypto.AesEncrypt(string(ssh.password), aesKey)
		if err != nil {
			return ssh, err
		}
		ssh.password = []byte(encryptPwd)
		return ssh, nil
	}
	// ssh public/private key case
	_, err := utilcrypto.AesDecrypt(string(ssh.privateKey), aesKey)
	if err == nil {
		return ssh, nil
	}
	encryptPrivateKey, err := utilcrypto.AesEncrypt(string(ssh.privateKey), aesKey)
	if err != nil {
		return ssh, err
	}
	ssh.privateKey = []byte(encryptPrivateKey)
	// ssh use pass phrase case
	if len(ssh.passPharse) != 0 {
		encryptPassPhrase, err := utilcrypto.AesEncrypt(string(ssh.passPharse), aesKey)
		if err != nil {
			return ssh, err
		}
		ssh.passPharse = []byte(encryptPassPhrase)
	}
	return ssh, nil
}
