module tkestack.io/tke

go 1.12

replace (
	// wait https://github.com/chartmuseum/storage/pull/34 to be merged
	github.com/chartmuseum/storage => github.com/choujimmy/storage v0.5.1-0.20191225102245-210f7683d0a6
	github.com/deislabs/oras => github.com/deislabs/oras v0.8.0
	// wait https://github.com/dexidp/dex/pull/1607 to be merged
	github.com/dexidp/dex => github.com/choujimmy/dex v0.0.0-20191225100859-b1cb4b898bb7
	k8s.io/client-go => k8s.io/client-go v0.17.0
)

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/Azure/go-autorest v13.3.1+incompatible // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/aws/aws-sdk-go v1.25.7
	github.com/bitly/go-simplejson v0.5.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/casbin/casbin/v2 v2.1.2
	github.com/chartmuseum/storage v0.5.0
	github.com/coreos/go-oidc v2.1.0+incompatible
	github.com/coreos/prometheus-operator v0.34.0
	github.com/deislabs/oras v0.8.0 // indirect
	github.com/dexidp/dex v0.0.0-20191223120519-789272a0c18f
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7
	github.com/emicklei/go-restful v2.11.1+incompatible
	github.com/fatih/color v1.7.0
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.4.7
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/go-openapi/inflect v0.19.0
	github.com/go-openapi/spec v0.19.4
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/snappy v0.0.1
	github.com/google/gofuzz v1.0.0
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/gorilla/mux v1.7.3
	github.com/gosuri/uitable v0.0.1
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/howeyc/fsnotify v0.9.0
	github.com/influxdata/influxdb1-client v0.0.0-20190402204710-8ff2fc3824fc
	github.com/jinzhu/configor v1.1.1
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.8
	github.com/kr/fs v0.1.0 // indirect
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/parnurzeal/gorequest v0.2.15
	github.com/pkg/errors v0.8.1
	github.com/pkg/sftp v1.10.0
	github.com/prometheus/alertmanager v0.17.0
	github.com/prometheus/client_golang v1.2.1
	github.com/rs/cors v1.6.0
	github.com/segmentio/ksuid v1.0.2
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.107+incompatible
	github.com/thoas/go-funk v0.4.0
	github.com/ugorji/go v1.1.7 // indirect
	go.etcd.io/etcd v0.0.0-20191023171146-3cf2f69b5738
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20191028145041-f83a4685e152
	golang.org/x/net v0.0.0-20191028085509-fe3aa8a45271
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.24.0
	gopkg.in/go-playground/validator.v9 v9.28.0
	gopkg.in/ldap.v2 v2.5.1
	gopkg.in/square/go-jose.v2 v2.3.1
	gopkg.in/yaml.v2 v2.2.4
	helm.sh/chartmuseum v0.11.0
	k8s.io/api v0.17.0
	k8s.io/apiextensions-apiserver v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/apiserver v0.17.0
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/cluster-bootstrap v0.17.0
	k8s.io/component-base v0.17.0
	k8s.io/helm v2.16.1+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.17.0
	k8s.io/kube-openapi v0.0.0-20191107075043-30be4d16710a
	k8s.io/kubectl v0.17.0 // indirect
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
)
