module tkestack.io/tke

go 1.12

replace (
	github.com/chartmuseum/storage => github.com/choujimmy/storage v0.0.0-20200507092433-6aea2df34764
	github.com/deislabs/oras => github.com/deislabs/oras v0.8.0
	helm.sh/helm/v3 => helm.sh/helm/v3 v3.2.0
	k8s.io/client-go => k8s.io/client-go v0.18.2
)

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/Masterminds/semver v1.5.0
	github.com/aws/aws-sdk-go v1.29.32
	github.com/bitly/go-simplejson v0.5.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/casbin/casbin/v2 v2.2.1
	github.com/chartmuseum/storage v0.8.0
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/coreos/prometheus-operator v0.38.1-0.20200506070354-4231c1d4b313
	github.com/dexidp/dex v0.0.0-20200408064242-83d8853fd969
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7
	github.com/dovics/domain-role-manager v0.0.0-20200325101749-a44f9c315081
	github.com/emicklei/go-restful v2.11.1+incompatible
	github.com/fatih/color v1.7.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/go-openapi/inflect v0.19.0
	github.com/go-openapi/spec v0.19.4
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/google/gofuzz v1.1.0
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.0
	github.com/gosuri/uitable v0.0.4
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/howeyc/fsnotify v0.9.0
	github.com/imdario/mergo v0.3.8
	github.com/influxdata/influxdb1-client v0.0.0-20190402204710-8ff2fc3824fc
	github.com/jinzhu/configor v1.1.1
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.9
	github.com/kr/fs v0.1.0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.7.1
	github.com/parnurzeal/gorequest v0.2.15
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.10.0
	github.com/prometheus/alertmanager v0.20.0
	github.com/prometheus/client_golang v1.4.0
	github.com/prometheus/common v0.9.1
	github.com/rs/cors v1.6.0
	github.com/segmentio/ksuid v1.0.2
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.107+incompatible
	github.com/thoas/go-funk v0.4.0
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200401174654-e694b7bb0875
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20200414173820-0848c9571904
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/grpc v1.26.0
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/ldap.v2 v2.5.1
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/square/go-jose.v2 v2.4.1
	gopkg.in/yaml.v2 v2.2.8
	helm.sh/chartmuseum v0.12.0
	helm.sh/helm/v3 v3.2.1
	k8s.io/api v0.18.2
	k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/apiserver v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/cluster-bootstrap v0.18.2
	k8s.io/component-base v0.18.2
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.18.2
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	k8s.io/kubectl v0.18.2 // indirect
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89
	sigs.k8s.io/yaml v1.2.0
)
