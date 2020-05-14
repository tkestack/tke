package util

import (
	"testing"

	"tkestack.io/tke/api/platform"
)

func TestSplitHostAndPath(t *testing.T) {
	type testCase struct {
		name         string
		raw          string
		expectedHost string
		expectedPath string
	}
	cases := []testCase{
		{
			name:         "ip",
			raw:          "1.1.1.1",
			expectedHost: "1.1.1.1",
		},
		{
			name:         "host",
			raw:          "myhost.com",
			expectedHost: "myhost.com",
		},
		{
			name:         "host/path",
			raw:          "myhost.com/path/to/service",
			expectedHost: "myhost.com",
			expectedPath: "path/to/service",
		},
	}
	for _, c := range cases {
		gh, gp := SplitHostAndPath(c.raw)
		if gh != c.expectedHost {
			t.Logf("name: %s, expect host %q, get %q", c.name, c.expectedHost, gh)
		}
		if gp != c.expectedPath {
			t.Logf("name: %s, expect path %q, get %q", c.name, c.expectedPath, gp)
		}
	}
}

func TestClusterHost(t *testing.T) {
	type testCase struct {
		name         string
		cluster      *platform.Cluster
		expectedHost string
		expectedErr  error
	}
	cases := []testCase{
		{
			name: "internal-ip",
			cluster: &platform.Cluster{
				Status: platform.ClusterStatus{
					Addresses: []platform.ClusterAddress{
						{
							Type: platform.AddressInternal,
							Host: "1.1.1.1",
							Port: 6443,
						},
					},
				},
			},
			expectedHost: "1.1.1.1:6443",
		},
		{
			name: "internal-host",
			cluster: &platform.Cluster{
				Status: platform.ClusterStatus{
					Addresses: []platform.ClusterAddress{
						{
							Type: platform.AddressInternal,
							Host: "mycluster.com",
							Port: 6443,
						},
					},
				},
			},
			expectedHost: "mycluster.com:6443",
		},
		{
			name: "internal-host-path",
			cluster: &platform.Cluster{
				Status: platform.ClusterStatus{
					Addresses: []platform.ClusterAddress{
						{
							Type: platform.AddressInternal,
							Host: "proxy.com/to/my/cluster",
							Port: 6443,
						},
					},
				},
			},
			expectedHost: "proxy.com:6443/to/my/cluster",
		},
		{
			name: "advertise-ip",
			cluster: &platform.Cluster{
				Status: platform.ClusterStatus{
					Addresses: []platform.ClusterAddress{
						{
							Type: platform.AddressAdvertise,
							Host: "1.1.1.1",
							Port: 6443,
						},
					},
				},
			},
			expectedHost: "1.1.1.1:6443",
		},
		{
			name: "advertise-host",
			cluster: &platform.Cluster{
				Status: platform.ClusterStatus{
					Addresses: []platform.ClusterAddress{
						{
							Type: platform.AddressAdvertise,
							Host: "mycluster.com",
							Port: 6443,
						},
					},
				},
			},
			expectedHost: "mycluster.com:6443",
		},
		{
			name: "advertise-host-path",
			cluster: &platform.Cluster{
				Status: platform.ClusterStatus{
					Addresses: []platform.ClusterAddress{
						{
							Type: platform.AddressAdvertise,
							Host: "proxy.com/to/my/cluster",
							Port: 6443,
						},
					},
				},
			},
			expectedHost: "proxy.com:6443/to/my/cluster",
		},
		{
			name: "real-ip",
			cluster: &platform.Cluster{
				Status: platform.ClusterStatus{
					Addresses: []platform.ClusterAddress{
						{
							Type: platform.AddressReal,
							Host: "1.1.1.1",
							Port: 6443,
						},
					},
				},
			},
			expectedHost: "1.1.1.1:6443",
		},
		{
			name: "real-host",
			cluster: &platform.Cluster{
				Status: platform.ClusterStatus{
					Addresses: []platform.ClusterAddress{
						{
							Type: platform.AddressReal,
							Host: "mycluster.com",
							Port: 6443,
						},
					},
				},
			},
			expectedHost: "mycluster.com:6443",
		},
		{
			name: "real-host-path",
			cluster: &platform.Cluster{
				Status: platform.ClusterStatus{
					Addresses: []platform.ClusterAddress{
						{
							Type: platform.AddressReal,
							Host: "proxy.com/to/my/cluster",
							Port: 6443,
						},
					},
				},
			},
			expectedHost: "proxy.com:6443/to/my/cluster",
		},
	}
	for _, c := range cases {
		gh, ge := ClusterHost(c.cluster)
		if gh != c.expectedHost {
			t.Logf("name: %s, expect host %q, get %q", c.name, c.expectedHost, gh)
		}
		if (c.expectedErr != nil && ge == nil) || (c.expectedErr == nil && ge != nil) {
			t.Logf("name: %s, expect err %v, get %v", c.name, c.expectedErr, ge)
		}
	}
}
