module tkestack.io/tke

go 1.16

replace (
	github.com/NetEase-Object-Storage/nos-golang-sdk => github.com/karuppiah7890/nos-golang-sdk v0.0.0-20191116042345-0792ba35abcc
	github.com/chartmuseum/storage => github.com/leoryu/chartmuseum-storage v0.11.1-0.20211104032734-9da39e8f5170
	github.com/deislabs/oras => github.com/deislabs/oras v0.8.0
	github.com/superedge/superedge => github.com/attlee-wang/superedge v0.8.2
	google.golang.org/grpc => google.golang.org/grpc v1.38.0
	k8s.io/api => k8s.io/api v0.22.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.22.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.3
	// this replace will be removed if https://github.com/kubernetes/kubernetes/pull/104920 is merged in 1.22
	k8s.io/apiserver => github.com/leoryu/k8s-apiserver v0.22.4-0.20211110063743-0341ac1e5801
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.22.3
	k8s.io/client-go => k8s.io/client-go v0.22.3
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.22.3
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.22.3
	k8s.io/code-generator => k8s.io/code-generator v0.22.3
	k8s.io/component-base => k8s.io/component-base v0.22.3
	k8s.io/component-helpers => k8s.io/component-helpers v0.22.3
	k8s.io/controller-manager => k8s.io/controller-manager v0.22.3
	k8s.io/cri-api => k8s.io/cri-api v0.22.3
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.22.3
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.22.3
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.22.3
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.22.3
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.22.3
	k8s.io/kubectl => k8s.io/kubectl v0.22.3
	k8s.io/kubelet => k8s.io/kubelet v0.22.3
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.22.3
	k8s.io/metrics => k8s.io/metrics v0.22.3
	k8s.io/mount-utils => k8s.io/mount-utils v0.22.3
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.22.3
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.22.3
	k8s.io/sample-controller => k8s.io/sample-controller v0.22.3
)

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/Masterminds/semver v1.5.0
	github.com/antihax/optional v1.0.0
	github.com/aws/aws-sdk-go v1.40.37
	github.com/bitly/go-simplejson v0.5.0
	github.com/caddyserver/caddy v1.0.5
	github.com/casbin/casbin/v2 v2.2.1
	github.com/chartmuseum/helm-push v0.9.0
	github.com/chartmuseum/storage v0.11.0
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/coreos/prometheus-operator v0.38.1-0.20200506070354-4231c1d4b313
	github.com/cyphar/filepath-securejoin v0.2.2
	github.com/deckarep/golang-set v1.7.1
	github.com/dexidp/dex v0.0.0-20210802203454-3fac2ab6bc3b
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.2-0.20200708230840-70e0022e42fd+incompatible
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7
	github.com/emicklei/go-restful v2.11.1+incompatible
	github.com/evanphx/json-patch v4.11.0+incompatible
	github.com/fatih/color v1.7.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/go-openapi/inflect v0.19.0
	github.com/gogo/protobuf v1.3.2
	github.com/google/gofuzz v1.1.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/gosuri/uitable v0.0.4
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/go-uuid v1.0.2
	github.com/howeyc/fsnotify v0.9.0
	github.com/imdario/mergo v0.3.12
	github.com/influxdata/influxdb1-client v0.0.0-20191209144304-8bf82d3c094d
	github.com/jinzhu/configor v1.1.1
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.11
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/parnurzeal/gorequest v0.2.15
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.10.1
	github.com/prometheus/alertmanager v0.20.0
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.26.0
	github.com/rs/cors v1.6.0
	github.com/segmentio/ksuid v1.0.3
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/superedge/superedge v0.8.2
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb v1.0.194
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.194
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm v1.0.194
	github.com/thoas/go-funk v0.4.0
	go.etcd.io/etcd/client/pkg/v3 v3.5.0
	go.etcd.io/etcd/client/v3 v3.5.0
	go.uber.org/zap v1.19.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/net v0.0.0-20210903162142-ad29c8ab022f
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac
	google.golang.org/grpc v1.40.0
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/ldap.v2 v2.5.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v2 v2.4.0
	gotest.tools v2.2.0+incompatible
	helm.sh/chartmuseum v0.13.1
	helm.sh/helm/v3 v3.7.1
	istio.io/api v0.0.0-20200715212100-dbf5277541ef
	istio.io/client-go v0.0.0-20200715214203-1ab538406cd1
	k8s.io/api v0.22.3
	k8s.io/apiextensions-apiserver v0.22.3
	k8s.io/apimachinery v0.22.3
	k8s.io/apiserver v0.22.3
	k8s.io/cli-runtime v0.22.3
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/cluster-bootstrap v0.22.3
	k8s.io/component-base v0.22.3
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.22.3
	k8s.io/kube-openapi v0.0.0-20210421082810-95288971da7e // indirect
	k8s.io/kubectl v0.22.3
	k8s.io/kubernetes v1.19.14
	k8s.io/metrics v0.22.3
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a
	sigs.k8s.io/controller-runtime v0.10.3
	sigs.k8s.io/yaml v1.2.0
	yunion.io/x/pkg v0.0.0-20200603123312-ad58e621aec0
)
