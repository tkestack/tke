module tkestack.io/tke

go 1.12

replace (
	github.com/deislabs/oras => github.com/deislabs/oras v0.8.0
	k8s.io/client-go => k8s.io/client-go v0.17.0
)

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/Azure/azure-storage-blob-go v0.0.0-20181022225951-5152f14ace1c // indirect
	github.com/Azure/go-autorest v13.3.0+incompatible // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/OneOfOne/xxhash v1.2.5 // indirect
	github.com/StackExchange/wmi v0.0.0-20180725035823-b12b22c5341f // indirect
	github.com/VividCortex/ewma v1.1.1 // indirect
	github.com/aws/aws-sdk-go v1.20.18
	github.com/beevik/etree v1.1.0 // indirect
	github.com/biogo/store v0.0.0-20160505134755-913427a1d5e8 // indirect
	github.com/bitly/go-simplejson v0.5.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/casbin/casbin v1.8.1
	github.com/cenk/backoff v2.0.0+incompatible // indirect
	github.com/certifi/gocertifi v0.0.0-20180905225744-ee1a9a0726d2 // indirect
	github.com/chartmuseum/storage v0.5.0
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/cockroachdb/cockroach v0.0.0-20170608034007-84bc9597164f // indirect
	github.com/cockroachdb/cockroach-go v0.0.0-20181001143604-e0a95dfd547c // indirect
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/coreos/go-oidc v2.1.0+incompatible
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/coreos/prometheus-operator v0.34.0
	github.com/dexidp/dex v0.0.0-20190620162747-157c359f3e86
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7
	github.com/elastic/gosigar v0.9.0 // indirect
	github.com/elazarl/go-bindata-assetfs v1.0.0 // indirect
	github.com/emicklei/go-restful v2.11.1+incompatible
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/fatih/color v1.7.0
	github.com/fatih/structtag v1.0.0 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/fsnotify/fsnotify v1.4.7
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/getsentry/raven-go v0.1.2 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-openapi/inflect v0.19.0
	github.com/go-openapi/spec v0.19.4
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/snappy v0.0.1
	github.com/google/gofuzz v1.0.0
	github.com/googleapis/gax-go v2.0.2+incompatible // indirect
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/gorilla/handlers v1.4.2 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/gosuri/uitable v0.0.1
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645 // indirect
	github.com/hashicorp/consul v1.4.4 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/go-plugin v0.0.0-20181212150838-f444068e8f5a // indirect
	github.com/hashicorp/go-rootcerts v0.0.0-20160503143440-6bb64b370b90 // indirect
	github.com/hashicorp/go-uuid v1.0.1
	github.com/hashicorp/serf v0.8.2 // indirect
	github.com/howeyc/fsnotify v0.9.0
	github.com/influxdata/influxdb v0.0.0-20170331210902-15e594fc09f1 // indirect
	github.com/influxdata/influxdb1-client v0.0.0-20190402204710-8ff2fc3824fc
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgx v3.2.0+incompatible // indirect
	github.com/jinzhu/configor v1.1.1
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/joho/godotenv v1.3.0
	github.com/json-iterator/go v1.1.8
	github.com/knz/strtime v0.0.0-20181018220328-af2256ee352c // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leanovate/gopter v0.2.4 // indirect
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/lib/pq v1.0.0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de
	github.com/lightstep/lightstep-tracer-go v0.15.6 // indirect
	github.com/lovoo/gcloud-opentracing v0.3.0 // indirect
	github.com/miekg/dns v1.1.8 // indirect
	github.com/minio/minio-go/v6 v6.0.27-0.20190529152532-de69c0e465ed // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/montanaflynn/stats v0.0.0-20180911141734-db72e6cae808 // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/mozillazg/go-cos v0.12.0 // indirect
	github.com/oklog/oklog v0.0.0-20170918173356-f857583a70c3 // indirect
	github.com/olekukonko/tablewriter v0.0.1 // indirect
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/opentracing-contrib/go-stdlib v0.0.0-20170113013457-1de4cc2120e7 // indirect
	github.com/opentracing/basictracer-go v1.0.0 // indirect
	github.com/opentracing/opentracing-go v1.0.2 // indirect
	github.com/openzipkin/zipkin-go v0.2.2 // indirect
	github.com/parnurzeal/gorequest v0.2.15
	github.com/pborman/uuid v1.2.0
	github.com/peterbourgon/g2s v0.0.0-20170223122336-d4e7ad98afea // indirect
	github.com/petermattis/goid v0.0.0-20170504144140-0ded85884ba5 // indirect
	github.com/pkg/errors v0.8.1
	github.com/pkg/sftp v1.10.0
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/prometheus/alertmanager v0.17.0
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/prometheus v2.5.0+incompatible // indirect
	github.com/rlmcpherson/s3gof3r v0.5.0 // indirect
	github.com/rs/cors v1.6.0
	github.com/rubyist/circuitbreaker v2.2.1+incompatible // indirect
	github.com/russellhaering/goxmldsig v0.0.0-20180430223755-7acd5e4a6ef7 // indirect
	github.com/samuel/go-zookeeper v0.0.0-20161028232340-1d7be4effb13 // indirect
	github.com/sasha-s/go-deadlock v0.0.0-20161201235124-341000892f3d // indirect
	github.com/segmentio/ksuid v1.0.2
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.107+incompatible
	github.com/thoas/go-funk v0.4.0
	github.com/ugorji/go v1.1.7 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20191028145041-f83a4685e152
	golang.org/x/net v0.0.0-20191028085509-fe3aa8a45271
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.24.0
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/fsnotify/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/go-playground/validator.v9 v9.28.0
	gopkg.in/square/go-jose.v2 v2.3.1
	gopkg.in/yaml.v2 v2.2.4
	gotest.tools v2.2.0+incompatible
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
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
)
