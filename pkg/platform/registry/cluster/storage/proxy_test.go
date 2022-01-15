package storage

import (
	"testing"
)

type testCase struct {
	host             string
	proxyPath        string
	expectedURL      string
	expectedRawQuery string
}

func TestMakeURL(t *testing.T) {
	var testCases = []testCase{
		{
			host:             "http://192.168.1.10:8888",
			proxyPath:        "apis/xxx.cloud.tencent.com/v1/namespaces/manifests?labelSelector=a.b.c/d.e=values1&limit=20",
			expectedURL:      "http://192.168.1.10:8888/apis/xxx.cloud.tencent.com/v1/namespaces/manifests?labelSelector=a.b.c/d.e=values1&limit=20",
			expectedRawQuery: "labelSelector=a.b.c/d.e=values1&limit=20",
		},
		{
			host:             "https://192.168.1.10:8888/",
			proxyPath:        "/apis/xxx.cloud.tencent.com/v1/namespaces/manifests?labelSelector=a.b.c/d.e=values1&limit=20",
			expectedURL:      "https://192.168.1.10:8888/apis/xxx.cloud.tencent.com/v1/namespaces/manifests?labelSelector=a.b.c/d.e=values1&limit=20",
			expectedRawQuery: "labelSelector=a.b.c/d.e=values1&limit=20",
		},
		{
			host:             "https://192.168.1.10:8888",
			proxyPath:        "/apis/xxx.cloud.tencent.com/v1/namespaces/manifests",
			expectedURL:      "https://192.168.1.10:8888/apis/xxx.cloud.tencent.com/v1/namespaces/manifests",
			expectedRawQuery: "",
		},

		{
			host:             "http://192.168.1.10:8888/aaa/bbb",
			proxyPath:        "/apis/xxx.cloud.tencent.com/v1/namespaces/manifests?labelSelector=a.b.c/d.e=values1&limit=20",
			expectedURL:      "http://192.168.1.10:8888/aaa/bbb/apis/xxx.cloud.tencent.com/v1/namespaces/manifests?labelSelector=a.b.c/d.e=values1&limit=20",
			expectedRawQuery: "labelSelector=a.b.c/d.e=values1&limit=20",
		},
		{
			host:             "http://192.168.1.10:8888/aaa/bbb",
			proxyPath:        "/apis/xxx.cloud.tencent.com/v1/namespaces/manifests",
			expectedURL:      "http://192.168.1.10:8888/aaa/bbb/apis/xxx.cloud.tencent.com/v1/namespaces/manifests",
			expectedRawQuery: "",
		},
		{
			host:             "http://192.168.1.10",
			proxyPath:        "apis/xxx.cloud.tencent.com/v1/namespaces/manifests?labelSelector=a.b.c/d.e=values1&limit=20",
			expectedURL:      "http://192.168.1.10/apis/xxx.cloud.tencent.com/v1/namespaces/manifests?labelSelector=a.b.c/d.e=values1&limit=20",
			expectedRawQuery: "labelSelector=a.b.c/d.e=values1&limit=20",
		},
	}

	for index, tcase := range testCases {

		uri, err := makeURL(tcase.host, tcase.proxyPath)
		if err != nil || uri.String() != tcase.expectedURL || uri.RawQuery != tcase.expectedRawQuery {
			t.Errorf("not pass test case %d\n", index)
		}
	}
}

func TestNonMakeURL(t *testing.T) {
	var testCases = []testCase{
		{
			host:      "192.168.1.10:8888",
			proxyPath: "/apis/xxx.cloud.tencent.com/v1/namespaces/manifests",
		},
		{
			host:      "",
			proxyPath: "/apis/xxx.cloud.tencent.com/v1/namespaces/manifests",
		},
	}

	for index, tcase := range testCases {
		_, err := makeURL(tcase.host, tcase.proxyPath) //should return an err
		if err == nil {
			t.Errorf("not pass test case %d\n", index)
		}
	}
}
