# This is the offline-pot ansible 'hosts' template file.
# Create an all group that contains the masters,workers and installer groups
[all:children]
masters
workers
installer
ceph
db
redis
logs
monitor
ingress
nfs
minio
salt
sgikes

# define global variables
[all:vars]
#SSH user, this user should allow ssh based auth without requiring a password
ansible_ssh_user=root
ansible_ssh_pass=*******
ansible_port=22

app_env_flag='shtu' # Named after current customer

# offline-pot path, exec install-tke-installer.sh will be auto set. don't change!!!
dpl_dir='/data/offline-pot'

# remote image registry url, exec builder.sh will be auto set. don't change!!!
remote_img_registry_url='reg.xx.yy.com'

# download enable or disable proxy script domain for devnetcloud
proxy_domain='download.devenv.xxx.com'

# switch
check_data_disk_size_switch='true' # check data disk size switch,the value must be true or false
data_disk_init_swith='true' # when the data disk is raw device,need init the data disk and mount to /data dir,the value must be true or false
check_ex_net='true' # check whether can access out server,the value must be true or false.
recover_disk_to_raw_switch='false' # recover disk to raw switch,the value must be true or false

deploy_ceph='false' # when deploy ceph need to check ceph disk whethere is raw device,the value must be true or false
use_calico='false' # pod network component whether use calico,default unuse,the value must be true or false.current not support calico
stress_perfor='true' # whether stress performance switch,default true, the value must be true or false.
deploy_redis='true' # whether deploy redis server,default true, the value must be true or false.
deploy_mysql='true' # whether deploy mysql server,default true, the value must be true or false.
deploy_prometheus='true' # whether deploy prometheus,default true, the value must be true or false.
deploy_helmtiller='true' # whether deploy helm tiller,default true, the value must be true or false.
deploy_nginx_ingress='true' # whether deploy nginx ingress controller,default true, the value must be true or false.
deploy_kafka='true' # whether deploy kafka, default true, the value must be true or false.
deploy_elk='true'  # whether deploy elk, default true, the value must be true or false; resources are not enough set false
deploy_nfs='true' # whether deploy nfs, default true, the value must be true or false
deploy_minio='false' # whether deploy minio, default true, the value must be true or false
deploy_business='true' # whether deploy business, defaut true, the value must be true or false
deploy_salt_minion='true' # whether deploy salt-minion for business cicd, default true, the value must be true or false
deploy_postgres='true' # whether dpeloy postgres, default true, the value must be true or false
deploy_sgikes='false' # whether dpeloy sgikes, default false, the value must be true or false; sgikes current for t*pd

# check items variables
check_path='/data' # check data disk dir, or for mount disk
check_data_disk_size='200' # check data disk size whether meets requirements
data_disk_name='vdb' # check data disk whethere is raw for data disk init,lsblk to get
ceph_disk_name='vdc' # check ceph disk whethere is raw device,lsblk to get
check_ceph_disk_size='200' # check ceph disk whether meets requirements.
disk_io_latency='25000' #  disk randread io latency reference
disk_iops='300' # disk io qps reference
fio_size='10G' # fio test io performance file size

# check out server domains, just support two; domain1 is https, domain2 is http
out_ser_domain1='api.xx.yy.com'
out_ser_domain2='reg.xx.yy.com'
 

# create lvm and mount disk variables
vg_name='vg_data' # required, volume group name
lv_name='lv_data' # required, logical volume name
disk_device_name='/dev/vdb' # data disk device name,fdisk -l to get
filesystem='xfs' # optional, default is 'xfs'
partition_cmd='n\np\n1\n\n\nt\n8e\nw' # when create partition failed need check system create partition proccess with manual

# offline registry config
# if tkestack not deploy, registry_domain's value must be "registry.tke.com" for harbor
registry_domain='registry.tke.com' # offline registry domain

# tke config
tke_replicas="1" # tke components's replicas number, default 1, Adjust according to the actual situation 
tke_ha_type="third" # tke ha type , third: will be use lb, tke: will be use keepalived, none: is not ha deploy
k8s_version='1.16.6' # kubernetes version, need forllow tkestack install pkg
net_interface='eth1'  # network insterface name, need all machine is the same name
cluster_cidr='172.16.0.0/19' # tke cluster's pod and service network cidr,default is 172.16.0.0/19.
tke_vip='172.17.0.6' # kubernetes master's ha vip address(lb or float ip), ha will be use
tke_vport='6443' # k8s api port, will be config on lb
max_cluster_service_num=256 # cluster service max number
max_node_pod_num=256 # per node max pod number
docker_data_root='/data/docker' # docker data root dir
kubelet_root_dir='/data/kubelet'
tke_admin_user='admin' # tke controller platform admin user name
tke_pwd='admin' # tke controller platform admin user password
tke_console_domain='console.tke.com' # tke console domain
ipvs='true' # whether enable ipvc,default true,must be true or false

# salt_minion configs
salt_master_domain='s.master.tke.com' # salt master's domain
salt_master_port='44506' # salt master's port

# redis configs
redis_mode='master-slave' # redis deploy mode, support master-slave or cluster.
REDIS_PORT='10101' # set redis listen port
REDIS_PASS='redis_P@s5' # set redis password

# redis-cluster config
redis_image_tag='6.0.5-debian-10-r2' # redis image tag
redis_nodes='6' # redis cluster node number, include master and slave
redis_replicas='1' # redis cluster replicas number
use_aof='no' # redis data whethere use AOF persistence, yes or no
redis_persistence='false' # redis data whethere use persistence, true or false
redis_data_dir="/data/redis" # redis persistence data dir
redis_exporter_img_tag='1.6.1-debian-10-r28' # redis exporter image tag
redis_sysctl_img_tag='buster' # redis sysctl image tag 
redis_taints='true' # redis node whetere add taints for not allow other server pod  schedule, true not allow; true or false
redis_cluster_client_img_tag='6.0.5-debian-10-r0' # redis cluster client image tag

# mysql configs
MYSQL_BUFFER='521M' # mysql buffer  value
MYSQL_PORT='3306' # mysql listen port
MYSQL_DATADIR='/data/mysql' # mysql data dir
MYSQL_PASS='mysql_P@s5'
# mysql mode, wx please set:
# sql_mode=STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
# git please set:
# sql_mode=STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
# t*pd , please set:
#sql_mode=
# common's please set ' ', default is wx's value
MYSQL_MODE='sql_mode=ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'

# postgres configs, postgres will be deploy db hosts group's second node 
POSTGRES_DATA_DIR='/data/postgres/data'
POSTGRES_PASSWORD='pgdb_P@s5'
POSTGRES_IMAGE='library/postgres:12-alpine'

# ingress config
ingress_replica='2' # ingress will be deploy number,default 2
ingress_host_network='true' # when has LoadBalancer ip for nginx-ingress and just has one master node set false, default true;
ingress_svc_type='ClusterIP' # when has LoadBalancer ip for nginx-ingress set LoadBalancer, otherwise set ClusterIP
ingress_lb_ip='' # set LoadBalancer ip for nginx-ingress

# kafka and zookeeper configs
kafka_data='/data/kafka' # kafka and zookeeper data save dir 
kafka_limit_cpu='1' # Adjust according to the actual situation,1 eq 1c
kafka_limit_mem='2Gi' # Adjust according to the actual situation
kafka_request_cpu='500m' # Adjust according to the actual situation, 1000m eq 1c
kafka_request_mem='1Gi' # Adjust according to the actual situation
kafka_heap_options='-Xmx2G -Xms2G' # Adjust according to the actual situation
zk_heap_size='2G' # Adjust according to the actual situation
kafka_image_name='library/cp-kafka' # Adjust according to the actual situation
kafka_image_tag='5.0.1' # Adjust according to the actual situation
zk_image_name='library/zookeeper' # Adjust according to the actual situation
zk_image_tag='3.5.5' # Adjust according to the actual situation
kafka_manager_image_name='library/kafka-manager' # Adjust according to the actual situation
kafka_manager_image_tag='1.3.3.22' # Adjust according to the actual situation
kafka_manager_username='admin' # Adjust according to the actual situation
kafka_manager_pwd='kafka_Mgr' # Adjust according to the actual situation 

# elk configs
es_data='/data/es' # save es data dir, Adjust according to the actual situation
logstash_replicas='1' # logstash replicas number, Adjust according to the actual situation
logstash_mem_limit='2Gi' # logstash memory limit, Adjust according to the actual situation
logstash_mem_req='2Gi' # logstash memory request, Adjust according to the actual situation
es_java_opts='-Xmx1g -Xms1g' # es jave options, Adjust according to the actual situation
es_request_cpu='100m' # Adjust according to the actual situation,1000m eq 1c
es_request_mem='2Gi' # Adjust according to the actual situation
es_limit_cpu='1' # Adjust according to the actual situation,1 eq 1c
es_limit_mem='2Gi' # Adjust according to the actual situation
kibana_request_cpu='100m' # Adjust according to the actual situation
kibana_request_mem='500Mi' # Adjust according to the actual situation
kibana_limit_cpu='1' # Adjust according to the actual situation
kibana_limit_mem='1Gi' # Adjust according to the actual situation
es_pwd='es_Mgr2020' # Adjust according to the actual situation
es_uname='elastic' # Adjust according to the actual situation

# nfs config
nfs_data='/data/nfsdata' # nfs data dir
nfs_app_list='("wx-nfs" "wx-uni" "wx-web")' # need nfs storage app's name,must be shell array
nfs_pv_storage_size="2Gi" # nfs pv storage size, app's pvc must be match this size 
is_create_pv="false" # default will be not create pv, must be true or false

# minio config
minio_img_name='library/minio' 
minio_img_tag='RELEASE.2019-12-17T23-16-33Z'
minio_mcimg_name='library/mc'
minio_mcimg_tag='edge'
minio_mount_path='/data/minio' # minio data  dir
minio_cpu_request='250m' # minio cpu request
minio_mem_request='256Mi' # minio memory request
minio_domain='minio.pot.tke.com'

# sg-ik-es 
sg_ik_repository="library" # sg ik es registry uri
sg_ik_busyboxversion="1.29.3" # sg ik busybox image tag
sg_ik_elkversion="6.8.0" # sg ik elasticsearch image tag
sg_ik_sgversion="25.1.ik" # sg ik searchguard image tag
sg_ik_sgkibanaversion="18.3" # sg ik kibana image tag
sg_ik_heapSize="2g" # sg ik es heap size, please adjust according to the actual situation
sg_ik_cpu_limit="1" # sg ik es cpu limit, please adjust according to the actual situation
sg_ik_mem_limit="4Gi" # sg ik es memory limit, please adjust according to the actual situation
sg_ik_cpu_req="500m" # sg ik es cpu request, please adjust according to the actual situation
sg_ik_mem_req="2Gi" # sg ik es memory request, please adjust according to the actual situation
data_size="500Gi" # es data node pvc size, please adjust according to the actual situation
master_data_size="200Gi" # es master node pvc size, please adjust according to the actual situation
sg_ik_kibana_cpu_limit="500m" # sg ik kibana cpu limit, please adjust according to the actual situation
sg_ik_kibana_mem_limit="1Gi" # sg ik kibana memory limit, please adjust according to the actual situation
sg_ik_kibana_cpu_req="100m" # sg ik kibana cpu request, please adjust according to the actual situation
sg_ik_kibana_mem_req="500Mi" # sg ik kibana memory request, please adjust according to the actual situation
sg_ik_ingress_class="nginx" # sg ik ingress class 
sg_ik_kibana_domain="sgik-kibana.t*pd.tke.com" # sg ik kibana domain
sg_ik_es_data="/data/sg-ik-es" # sg ik es data dir

# harbor's config
harbor_http="80" # harbor proxy http port
harbor_https="443" # harbor proxy https port
harbor_admin_password="Harbor12345" # harbor admin password
harbor_db_password="root123" # db password
harbor_data_volume="/data/registry" # harbor registry dir

# Create installer group
[installer]
127.0.0.1

# Create masters group
[masters]
172.17.0.44

# Create workers group
[workers]

# Create db group, must be two nodes
[db]

# create redis group, when master-slave mode must be two nodes; cluster mode less three nodes and 
# persistence current just support three nodes
[redis]

# Create ceph group
[ceph]

# Create logs group, must be three nodes
[logs]

# create monitor group
[monitor]

# create ingress controller group, when ha deploy, ingress group config masters's group second and third nodes
# and when just has one master node, the ingress group could not include master's node
[ingress]

# create nfs server group, must be one node
[nfs]

# create minio group, must be four nodes
[minio]

# create salt-minion group, there are several centers to define several nodes,defualt is the first of master's group
# salt-minion node must be can exec kubelet cmd
[salt]
172.17.0.44

# sg ik elasticsearch, current for t*pd,three nodes or six nodes, when six node es master is 0~2, es data is 3~5
[sgikes]

