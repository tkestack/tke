module tkestack.io/tke

go 1.12

replace (
	github.com/chartmuseum/storage => github.com/choujimmy/storage v0.5.1-0.20210412121305-660c0e91489b
	github.com/containerd/containerd => github.com/containerd/containerd v1.4.3
	github.com/deislabs/oras => github.com/deislabs/oras v0.8.0
	go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200819165624-17cef6e3e9d5
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	k8s.io/api => k8s.io/api v0.19.7
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.7
	k8s.io/apiserver => k8s.io/apiserver v0.19.7
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.7
	k8s.io/client-go => k8s.io/client-go v0.19.7
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.7
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.7
	k8s.io/code-generator => k8s.io/code-generator v0.19.7
	k8s.io/component-base => k8s.io/component-base v0.19.7
	k8s.io/cri-api => k8s.io/cri-api v0.19.7
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.7
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.7
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.7
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.7
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.7
	k8s.io/kubectl => k8s.io/kubectl v0.19.7
	k8s.io/kubelet => k8s.io/kubelet v0.19.7
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.7
	k8s.io/metrics => k8s.io/metrics v0.19.7
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.7
)

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/Masterminds/semver v1.5.0
	github.com/antihax/optional v0.0.0-20180407024304-ca021399b1a6
	github.com/aws/aws-sdk-go v1.29.32
	github.com/bitly/go-simplejson v0.5.0
	github.com/caddyserver/caddy v1.0.5
	github.com/casbin/casbin/v2 v2.2.1
	github.com/chartmuseum/helm-push v0.9.0
	github.com/chartmuseum/storage v0.8.0
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/coreos/prometheus-operator v0.38.1-0.20200506070354-4231c1d4b313
	github.com/cyphar/filepath-securejoin v0.2.2
	github.com/deckarep/golang-set v1.7.1
	github.com/dexidp/dex v0.0.0-20200408064242-83d8853fd969
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.2-0.20200708230840-70e0022e42fd+incompatible
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7
	github.com/emicklei/go-restful v2.11.1+incompatible
	github.com/evanphx/json-patch v4.9.0+incompatible
	github.com/fatih/color v1.7.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/go-openapi/inflect v0.19.0
	github.com/go-openapi/spec v0.19.4
	github.com/gogo/protobuf v1.3.1
	github.com/google/gofuzz v1.1.0
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.2
	github.com/gosuri/uitable v0.0.4
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/howeyc/fsnotify v0.9.0
	github.com/imdario/mergo v0.3.8
	github.com/influxdata/influxdb1-client v0.0.0-20191209144304-8bf82d3c094d
	github.com/jinzhu/configor v1.1.1
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.10
	github.com/kr/fs v0.1.0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/parnurzeal/gorequest v0.2.15
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.10.0
	github.com/prometheus/alertmanager v0.20.0
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/common v0.10.0
	github.com/rs/cors v1.6.0
	github.com/segmentio/ksuid v1.0.3
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v1.1.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.107+incompatible
	github.com/thoas/go-funk v0.4.0
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200819165624-17cef6e3e9d5
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/grpc v1.28.1
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/ldap.v2 v2.5.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/square/go-jose.v2 v2.4.1
	gopkg.in/yaml.v2 v2.3.0
	gotest.tools v2.2.0+incompatible
	helm.sh/chartmuseum v0.12.0
	helm.sh/helm/v3 v3.4.2
	istio.io/api v0.0.0-20200715212100-dbf5277541ef
	istio.io/client-go v0.0.0-20200715214203-1ab538406cd1
	k8s.io/api v0.19.7
	k8s.io/apiextensions-apiserver v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/apiserver v0.19.7
	k8s.io/cli-runtime v0.19.7
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/cluster-bootstrap v0.19.7
	k8s.io/component-base v0.19.7
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.19.7
	k8s.io/kube-openapi v0.0.0-20200805222855-6aeccd4b50c6
	k8s.io/kubectl v0.19.7
	k8s.io/kubernetes v1.19.7
	k8s.io/metrics v0.19.7
	k8s.io/utils v0.0.0-20200729134348-d5654de09c73
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/yaml v1.2.0
	yunion.io/x/pkg v0.0.0-20200603123312-ad58e621aec0
)
