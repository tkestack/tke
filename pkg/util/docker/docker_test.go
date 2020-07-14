package docker

import (
	"fmt"
	"strconv"
	"testing"
)

func TestDockerGetNameArchTag(t *testing.T) {
	imageTestCases := []struct {
		// image is the image name component of testcase
		image string
		// name is the string representation for the image name (without arch)
		name string
		// arch is the arch info for the image name
		arch string
		// tag is the tag for the image name
		tag string
		// err is the error expected from Parse, or nil
		err error
	}{
		{
			image:	"localhost:5000/library/test-amd64:v1.2.3",
			name:	"localhost:5000/library/test",
			arch:	"amd64",
			tag:	"v1.2.3",
			err:	nil,
		},
		{
			image:	"127.0.0.1/tke/test-aaa:1.2.3",
			name:	"127.0.0.1/tke/test-aaa",
			arch:	"",
			tag:	"1.2.3",
			err:	nil,
		},
		{
			image:	"tke/test-arm64:1-2-3",
			name:	"tke/test",
			arch:	"arm64",
			tag:	"1-2-3",
			err:	nil,
		},
		{
			image:	"test-aaa:1.2.3-amd64",
			name:	"test-aaa",
			arch:	"",
			tag:	"1.2.3-amd64",
			err:	nil,
		},
		{
			image:	"test:5000/repo@sha256:ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			name:	"",
			arch:	"",
			tag:	"",
			err:	fmt.Errorf("image: %v does not have a name and tag", "test:5000/repo@sha256:ffff"),
		},
		{
			image:	"tketest:<none>",
			name:	"",
			arch:	"",
			tag:	"",
			err:	fmt.Errorf("image %s is invalid", "tketest:<none>"),
		},
		{
			image:	"tketest",
			name:	"",
			arch:	"",
			tag:	"",
			err:	fmt.Errorf("image: %v does not have a name and tag", "tketest"),
		},
		{
			image:	"tketest/test/test",
			name:	"",
			arch:	"",
			tag:	"",
			err:	fmt.Errorf("image: %v does not have a name and tag", "tketest"),
		},
		{
			image:	"tketest/test1/test2/test3:v1.23",
			name:	"tketest/test1/test2/test3",
			arch:	"",
			tag:	"v1.23",
			err:	nil,
		},
	}

	docker := New()
	for _, testcase := range imageTestCases {
		failf := func(format string, v ...interface{}) {
			t.Logf(strconv.Quote(testcase.image)+": "+format, v...)
			t.Fail()
		}

		name, arch, tag, err := docker.GetNameArchTag(testcase.image)
		if testcase.err != nil {
			if err == nil {
				failf("missing expected error: %v", testcase.err)
			}
		}

		if testcase.err == nil {
			if err != nil {
				failf("unexpected error: %v", err)
			}
		}

		if name != testcase.name {
			failf("mismatched name: got %q, expected %q", name, testcase.name)
		}

		if arch != testcase.arch {
			failf("mismatched arch: got %q, expected %q", arch, testcase.arch)
		}

		if tag != testcase.tag {
			failf("mismatched tag: got %q, expected %q", tag, testcase.tag)
		}
	}
}
