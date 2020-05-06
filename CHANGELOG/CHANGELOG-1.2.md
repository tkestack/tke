## [1.2.4](https://github.com/tkestack/tke/compare/v1.2.3...v1.2.4) (2020-04-26)


### Bug Fixes

* Add master toleration for alertmanager ([f6cd780](https://github.com/tkestack/tke/commit/f6cd7802592f72168354340588faec0d728d60b5))
* add two variants (unknown & v8) for arm64 images ([3058b0b](https://github.com/tkestack/tke/commit/3058b0b630d130c55334e45d2fe3b24ef4cbba78))
* alarmplicy fetch workloadInstance bug ([#234](https://github.com/tkestack/tke/issues/234)) ([3605b3c](https://github.com/tkestack/tke/commit/3605b3cf1753742e6b9b3fee79f856c83d0ef9f7))
* change label namespaceName to namespace_name for project metrics ([19f80cd](https://github.com/tkestack/tke/commit/19f80cd7bf1349696051119c7a57bdfb70825740))
* chmod with hook's path ([0025cd4](https://github.com/tkestack/tke/commit/0025cd4d181231cff7b2fba6ca851fafd84d6aca))
* cluster dns domain default value ([ecb0277](https://github.com/tkestack/tke/commit/ecb0277770d814891df139519579298776add423))
* extend imporResource timeout ([2d06e70](https://github.com/tkestack/tke/commit/2d06e706f6ebb3ac9068b337f0c81f9b09d77251))
* field selector in list machine when using spec.ip ([7b30b37](https://github.com/tkestack/tke/commit/7b30b37cab3891f0344cd27257b4b737d7db1c8a))
* hook validate spell error ([3d98d50](https://github.com/tkestack/tke/commit/3d98d50265dac23af0143a4d3d860cb7763471f0))
* illegal request to query service event list ([6285f88](https://github.com/tkestack/tke/commit/6285f88763692665abea17f102c09dbbc79d6118))
* label project_name of ns should begin with tkestack.io, and add ([4627865](https://github.com/tkestack/tke/commit/4627865e54f3da720dc0e8738c981834376f6149))
* machine index fields missing spec.ip ([b174f44](https://github.com/tkestack/tke/commit/b174f44e151157900116aff3d048819e7fef1719))
* make Makefile compatible with docker version under 18.03 ([#266](https://github.com/tkestack/tke/issues/266)) ([14b3f65](https://github.com/tkestack/tke/commit/14b3f653edbc2a2a835f8044c874321908db38ff))
* mark master ([2fa91f1](https://github.com/tkestack/tke/commit/2fa91f147585bf9c5051b16acf368a6e28363b92))
* remove request and limit of node-exporter ([e65cc0e](https://github.com/tkestack/tke/commit/e65cc0e298babca4e3b615e632a806c7cef65900))
* **cluster:** AdvertiseAddress in kubeadm ([527636f](https://github.com/tkestack/tke/commit/527636f11fefa515523b5c1f7ea371d184bb5998))
* **cluster:** check url in keepalived vrrp_script ([798c06c](https://github.com/tkestack/tke/commit/798c06c749979705a984b18bdd3e9a2a81dc1db4))
* **cluster:** hang in failed status when health check is not ok ([8d1e4bc](https://github.com/tkestack/tke/commit/8d1e4bc015d69d7d4b623a97f0bd3548378aeba9))
* **cluster:** validation machine clusterName ([019d7af](https://github.com/tkestack/tke/commit/019d7af9ff87b3930b1eb3b93dad67f78bc84e10))
* **tke-installer:** tke-platform-api cert auth ([f929c65](https://github.com/tkestack/tke/commit/f929c654ebc70411f7284c4eabb2a455fa79a6ce))
* **tke-installer:** validate cpu and memory requirements ([225df85](https://github.com/tkestack/tke/commit/225df85deed64d476cb92b12e3ea2469cc1aa61d))
* **tke-installer:** validate error when multiple request body is not same ([0a606e4](https://github.com/tkestack/tke/commit/0a606e44786774c79a20b05a5df2f4a99a36984b))
* remove unnecessary labels for tke project metrics ([4bd16e4](https://github.com/tkestack/tke/commit/4bd16e40662e8d6b813976d67ffd56639f77448b))
* rollback does not exist in deployment of apps/v1 ([23ccebd](https://github.com/tkestack/tke/commit/23ccebdb0fb059f86584cbd65546a312ea44c114))
* set retentionSize to limit prometheus disk usage; set retention ([8810ef2](https://github.com/tkestack/tke/commit/8810ef24ddceba93c8ee8e45e0056b53b9ec6127))
* ValidateEnum ([87300fc](https://github.com/tkestack/tke/commit/87300fcda9c1630132da7d46230abac440687863))
* ValidateHooks ([ed638ac](https://github.com/tkestack/tke/commit/ed638ac5776b11c94b9a1a281bdc0968e1785f65))
* webhook fmt ([#291](https://github.com/tkestack/tke/issues/291)) ([5d04ec7](https://github.com/tkestack/tke/commit/5d04ec7e08436c948a2810a31c57ce1d855f14e0))


### Features

* add alertRepeatInterval option to prometheus addon ([a5a3059](https://github.com/tkestack/tke/commit/a5a3059d8249a2e732a6b8bece5e231e603d3e9d))
* Add metrics for namespace capacity and allocated ([0674196](https://github.com/tkestack/tke/commit/06741964397698ddd2b22c573c303dbc61b2bbe1))
* add more etcd metrics to prometheus config ([bd451ed](https://github.com/tkestack/tke/commit/bd451ed2967948a08318c8e34167ccd6add18a2f))
* Add option message_ttl to tke-notify-api to specify ttl of ([2692bac](https://github.com/tkestack/tke/commit/2692bacc81b1aafb746dedfdc4d8878ed3102434))
* Add Resources and RunOnMaster to spec of prometheus addon ([d6dcdd2](https://github.com/tkestack/tke/commit/d6dcdd261463f9e9bc89da87b4f43b489bf12f98))
* adjust clientset timeout to 30s ([73f35e9](https://github.com/tkestack/tke/commit/73f35e9852c75e4cc492006d32433d3465dee791))
* generate certificate for business namespace ([2fd8b08](https://github.com/tkestack/tke/commit/2fd8b087b77d26ced42226b0e5fc82412ef0edf0))
* generate certificate for business namespace ([7795ce1](https://github.com/tkestack/tke/commit/7795ce110cfa60f62665bd2b8fe270222e3fea33))
* **tke-installer:** the validity period of the certificate is revised to 100 years, close [#260](https://github.com/tkestack/tke/issues/260) ([151d23d](https://github.com/tkestack/tke/commit/151d23d5866834503379a294c7643f712fdf9e7b))
* hook support args ([40c0201](https://github.com/tkestack/tke/commit/40c020141fc28cd5d098160042c08825569b42c3))
* **cluster:** ClusterSpec add ServiceCIDR ([ef33d2d](https://github.com/tkestack/tke/commit/ef33d2d613acd3102e297dcaed4b01145351f8c3))
* reuse the global cluster front proxy cert for tke components ([4085dc6](https://github.com/tkestack/tke/commit/4085dc6f442eb55ee9498c1133125f05d96c2bba))
* **cluster:** support multi network card for control plane ([9eacb15](https://github.com/tkestack/tke/commit/9eacb15395715b43b291d23758fe4940ac05a323))
* **galaxy:** update galaxy and galaxy-ipam version, galaxy-ipam service use clusterIP to avoid ipvs nodePort issue ([b20a314](https://github.com/tkestack/tke/commit/b20a314f22684803a341eb4e552c306e4d12dbd7))
* Improve node list performance ([#240](https://github.com/tkestack/tke/issues/240)) ([9879edf](https://github.com/tkestack/tke/commit/9879edfd1f3eefae80881705eda7daaf1e7dffde))
* platform and auth support webhook authentication ([#290](https://github.com/tkestack/tke/issues/290)) ([084c739](https://github.com/tkestack/tke/commit/084c7391872eb44809c486a38fbdafb3681386f5))
* Send metrics kube_node_labels and kube_namespace_labels to remote ([55b81dc](https://github.com/tkestack/tke/commit/55b81dc99e582363d108a47580b020e07729820f))
* support password grant ([#258](https://github.com/tkestack/tke/issues/258)) ([f604469](https://github.com/tkestack/tke/commit/f604469e3fe072a55c93f186be344130cd11c1b5))
* support taints for master and worker ([af5e105](https://github.com/tkestack/tke/commit/af5e10569ee5a6e258433d8accaca213b4fe3b11))
* **tke-installer:** add cpu and memory validation ([5b24650](https://github.com/tkestack/tke/commit/5b24650a0609540f6deb1ce267f65830de13cd11))
* **tke-installer:** add retry for register apiservices ([dc54a05](https://github.com/tkestack/tke/commit/dc54a050b51ab3dee58f2b74942391746c81df6f))



## [1.2.3](https://github.com/tkestack/tke/compare/v1.2.2...v1.2.3) (2020-04-01)


### Bug Fixes

* add PublicAlternativeNames to tke-auth redirect hosts ([709ed72](https://github.com/tkestack/tke/commit/709ed7222355e4cb53e8055be9d0dc3416c1f56e))
* apiclient clientset ([e54a749](https://github.com/tkestack/tke/commit/e54a74965a8f4e11e2f16b2c46bbab0a3d46ef54))
* apikey description ([2f0ecde](https://github.com/tkestack/tke/commit/2f0ecde6caabdd65bf3311f9d18dbfbccf86facd))
* delete machine in terminating when credential was deleted already. ([50d272c](https://github.com/tkestack/tke/commit/50d272c307460f4a434d4a373096e76021995ac4))
* gpu driver package ([bec85c6](https://github.com/tkestack/tke/commit/bec85c6706c86d643d00109b320824651e0f7985))
* install gpu driver for machine ([3a7c24e](https://github.com/tkestack/tke/commit/3a7c24ee78bf8a70e36b3bba56cdd2e7ddd47b0d))
* install gpu driver for machine ([0b98a43](https://github.com/tkestack/tke/commit/0b98a434c9addd4ac56fe11304501ed43ca606ed))
* ipam check health ([7848567](https://github.com/tkestack/tke/commit/78485674954d49727f86c7d7e5e9d04342acf168))
* ipam check health ([cc7e11a](https://github.com/tkestack/tke/commit/cc7e11acfb31d4fd7110d231e384164e25ff3ebe))
* lose users and groups when update policy spec from console ([#216](https://github.com/tkestack/tke/issues/216)) ([06df0a7](https://github.com/tkestack/tke/commit/06df0a7f441ea87996e3d10f81eb0efe521ddaf6))
* Make sure all files in provider-res are end with tar.gz ([02ade0d](https://github.com/tkestack/tke/commit/02ade0dda8bd840ce996ba6f345e2fd645a1988d))
* update provider-res version in tke ([ab4b7b6](https://github.com/tkestack/tke/commit/ab4b7b6a10eadd6241fa159e1d6cfb5c29b479f7))
* update provider-res version in tke ([00989b6](https://github.com/tkestack/tke/commit/00989b6a1201a5cecbb30ad4fd78d1707dd7c0e4))


### Features

* add assets of png and svg for style ([#219](https://github.com/tkestack/tke/issues/219)) ([8fd98f4](https://github.com/tkestack/tke/commit/8fd98f4cf1f5b818f8c71f3f423d8ef69ba05c05))
* add ca key and etcd key to cluster credential ([94ac217](https://github.com/tkestack/tke/commit/94ac21748f3b9f10e896d44a21d9189a366c7d23))
* allow to install global on the machine running tke-installer, need add EnsureDocker into cluster.spec.features.skipConditions ([91a14de](https://github.com/tkestack/tke/commit/91a14de40f6e463192ab37b7c2e47e98291ae33b))
* make path for /auth/ and apikey password unprotected in gateway ([#217](https://github.com/tkestack/tke/issues/217)) ([dbefdde](https://github.com/tkestack/tke/commit/dbefdde1f1ac013302d17a51943843a3777bd8b7))
* support ha for cluster creation ([6c36511](https://github.com/tkestack/tke/commit/6c36511ef068c557a6f157ac96ea7eb44ee90d95))
* support third party lb can't connect self when bootstrap like tgw and clb etc. ([2243eb4](https://github.com/tkestack/tke/commit/2243eb4774bdc61d3c35b00ac0536e62bbe6240f))
* Use docker buildx to build keepalived & provider-res. Update provider-res version. ([dbde12a](https://github.com/tkestack/tke/commit/dbde12a0d6862878db90b95f602624f27beda23b))



## [1.2.2](https://github.com/tkestack/tke/compare/v1.2.1...v1.2.2) (2020-03-26)


### Bug Fixes

* [#204](https://github.com/tkestack/tke/issues/204) [#211](https://github.com/tkestack/tke/issues/211) ([b905927](https://github.com/tkestack/tke/commit/b9059272a0fb2b367f6f95bbf5f128872c84c3fb))
* get metrics failed if values contain NaN ([2adf668](https://github.com/tkestack/tke/commit/2adf66829213c119e8657e86bfaa5a0f402e4918))
* gpu quota admission ip ([1bde434](https://github.com/tkestack/tke/commit/1bde43422b3f2843d90c3030cab68dfe0f23ee32))
* gpu quota admission service not working when using dns name ([88979a9](https://github.com/tkestack/tke/commit/88979a9c3a2aaf5cc30209fa05f506b12f50ae7a))
* gpu quota admission when annotations is nil ([8609e2f](https://github.com/tkestack/tke/commit/8609e2fd68b5c143092a4275f37193446e2b60e6))



## [1.2.1](https://github.com/tkestack/tke/compare/v1.2.0...v1.2.1) (2020-03-25)


### Bug Fixes

* add list/get/delete option forward to kube-apiserver ([7e6a32e](https://github.com/tkestack/tke/commit/7e6a32eb4681a7f93d874538088ce782d371eac7))
* get rid of empty samples of metric query for thanos backend ([9cd6cc8](https://github.com/tkestack/tke/commit/9cd6cc81ffaadb0ef5462c9d963ae669f26cc9b7))
* GetSourceIP missing port ([6d06206](https://github.com/tkestack/tke/commit/6d06206020cf007b5445bc48d9964956fe410197))
* platform controller registry namespace ([c4a7bbb](https://github.com/tkestack/tke/commit/c4a7bbb1463390a8baec7829660df4befadadaf3))
* update go.mod to track dexidp/dex ([#197](https://github.com/tkestack/tke/issues/197)) ([e3bcb16](https://github.com/tkestack/tke/commit/e3bcb16cde501fb03e4904b4345a6e26cbda9c9d))
* validate users and groups bugs ([#200](https://github.com/tkestack/tke/issues/200)) ([f8a0693](https://github.com/tkestack/tke/commit/f8a0693e79b53de693f685386c9fabf5af81811a))


### Features

* add cluster/project name as group to userinfo ([02f03b7](https://github.com/tkestack/tke/commit/02f03b79d8b1f879d08638d05a5bc32e59d19477))
* add webtty link for container ([ef027c2](https://github.com/tkestack/tke/commit/ef027c2c7b5e01b15ae8106ea88e02f34675b575))
* enhance supportting for changing parent prj ([14040b7](https://github.com/tkestack/tke/commit/14040b726056add95c685aa418fd2246dd276eea))
* support changing parent project ([91b0a25](https://github.com/tkestack/tke/commit/91b0a2569f320c5047cfd6597b08427c6e033244))
* support loadmore for workload ([#207](https://github.com/tkestack/tke/issues/207)) ([055ffc4](https://github.com/tkestack/tke/commit/055ffc4eeac10c7260402025e4297878bb996c25))
* support multi net interface in installer node ([76a0b83](https://github.com/tkestack/tke/commit/76a0b83b37de52bdbe9194241ed5c56242bb09bb))
* support other version cluster when create in running phase ([fc916f4](https://github.com/tkestack/tke/commit/fc916f485abb87192fe723cc80d3bbdb577db71e))
* support removing clusters from business prj ([ce8e962](https://github.com/tkestack/tke/commit/ce8e9622ff9a2e3434f4ecd6522e7b685304c4da))
* use UpdateStatus instead of Update ([fe52855](https://github.com/tkestack/tke/commit/fe52855ef2bc392b032f2def51ef3aa83a737460))
* 完善lbcf ([#205](https://github.com/tkestack/tke/issues/205)) ([3570db1](https://github.com/tkestack/tke/commit/3570db1d6b40ba90922f74d36608ecdf5a909d89))



# [1.2.0](https://github.com/tkestack/tke/compare/v1.1.0...v1.2.0) (2020-03-19)


### Bug Fixes

* abnormal status.used data cause panic ([0cb4914](https://github.com/tkestack/tke/commit/0cb4914e524117b905a5017b1c5ffb83b1244cc1))
* add --iface=$(HOST_IP) arg in flannel ([15036fc](https://github.com/tkestack/tke/commit/15036fc1b7d0dce8cd2d4acdc0159878c72ac464))
* add node name env for galaxy daemonsets ([149a7fd](https://github.com/tkestack/tke/commit/149a7fdefe1ca8cf2f2976df0fa9c3f457b61996))
* arch match regexp ([fa0ac58](https://github.com/tkestack/tke/commit/fa0ac581a9157cf019be03701b2974d862d51dce))
* business can NOT work when there is invalid namespace ([a2d3465](https://github.com/tkestack/tke/commit/a2d3465219457e42959e6f759acaa82b86533446))
* construct tke attr for registry and business api ([#190](https://github.com/tkestack/tke/issues/190)) ([e51bf92](https://github.com/tkestack/tke/commit/e51bf9251183a89a55653bc50a308b7b3f422e48))
* docker manifest arch for arm64 ([cddf568](https://github.com/tkestack/tke/commit/cddf5683633d035d65741f7ac938472b0ba5ebb4))
* docker manifest name ([7873e6a](https://github.com/tkestack/tke/commit/7873e6adeeb81489fb996081f6fcc1470f0a6866))
* docker push manifest error ([0c1e61d](https://github.com/tkestack/tke/commit/0c1e61de29c58d153774785c01583d97237ebf80))
* exit -1 ([7cc530d](https://github.com/tkestack/tke/commit/7cc530d5a9b4810694f52117a387cd55632a5392))
* failed to create notify messages sent by webhook channel ([6ad8fc3](https://github.com/tkestack/tke/commit/6ad8fc3259ed1c61d2504fa52c95206642094e70))
* failed to get metrics from cadvisor of k8s v1.12 ([098ccf2](https://github.com/tkestack/tke/commit/098ccf2770488ad6d5c3204b78d4ae15cc5503ef))
* fix the request method for monitor ([#194](https://github.com/tkestack/tke/issues/194)) ([ab59f5f](https://github.com/tkestack/tke/commit/ab59f5f030cc64eefb1f81904fde3842b0853782))
* gateway's http handler not configured properly and needs to ignore protected path ([04d6ca3](https://github.com/tkestack/tke/commit/04d6ca3c58d75de8b43764f6d68dcb25e60a5778))
* increase resource limit of prometheus and update version to v2.16.0 ([c2eefa8](https://github.com/tkestack/tke/commit/c2eefa8dc4f84cd0f0b38ee7efdba93e7a3b62b7))
* installer e2e test ([91844da](https://github.com/tkestack/tke/commit/91844da146d97c3dd8601a2150ce0d20eb8b54e7))
* installer provider init ([84f0c77](https://github.com/tkestack/tke/commit/84f0c77f4345b3dd8901ae1a6d3bac2400460f90))
* installer registry image without arch ([b6d171a](https://github.com/tkestack/tke/commit/b6d171a899b0bf4c31cd3dbe5304c187f7dd2012))
* IsGPU always false ([91a9d21](https://github.com/tkestack/tke/commit/91a9d21f2a75ca27191a06b78fc62f370284e0be))
* keepalived.conf ([c74b129](https://github.com/tkestack/tke/commit/c74b1290bbe7989aeb7ea68171209dc011a65838))
* kubeadm join ([47747bc](https://github.com/tkestack/tke/commit/47747bc15c2d6193f74e7109f46a0d8bd2337477))
* kubelet find docker credential failed because missing HOME by systemd ([ffbb7c0](https://github.com/tkestack/tke/commit/ffbb7c0eb9982ce9fc577d1feccdc920cca8cfbb))
* manifest overwrite when same name image ([c721057](https://github.com/tkestack/tke/commit/c721057471d103fb86764717b10ad41bbdb1b8b6))
* manifest replace ([367c563](https://github.com/tkestack/tke/commit/367c5633daf3d13f3d5cf8b8cacf914d05fa556b))
* mark master label ([f0dbaf9](https://github.com/tkestack/tke/commit/f0dbaf993ab26d04440e083064863d1f31887460))
* may get duplicate CalculatedNamespaces ([5c9f10a](https://github.com/tkestack/tke/commit/5c9f10a99f9c82e0927ef42c9626189a3e987920))
* messagerequests and templates remains after channel deleted ([cb10732](https://github.com/tkestack/tke/commit/cb10732cc47aed66c3c7b700affad9d7df1cad88))
* missing group_headers in request header ([32e8667](https://github.com/tkestack/tke/commit/32e8667eb4695728d3e6af90299d4f63dc657455))
* missing issue with finalize method in business namespace ([5684206](https://github.com/tkestack/tke/commit/5684206b0f61415b15c603ac31c435ec5a05e76d))
* missing network name ([fed5d6e](https://github.com/tkestack/tke/commit/fed5d6ec6be951dc8f404b9ca2d600afa19ef51e))
* not update cls quota when update business ns ([f057671](https://github.com/tkestack/tke/commit/f057671124c3da8c9a70c78a51670ed9915c0ed5))
* nvidia version ([4f7b3fe](https://github.com/tkestack/tke/commit/4f7b3fe92a957ba87f5af13b2770d62d4cbbd566))
* Only ask users to enable experimental features when necessary ([cab16b2](https://github.com/tkestack/tke/commit/cab16b2d7ac73f1ca7c5375580c813efaa1796dc))
* platform controller flags ([db46fd7](https://github.com/tkestack/tke/commit/db46fd7ea8d1cb42df3b89281bf65216cb3310bd))
* prometheus config adjust ([9488c6a](https://github.com/tkestack/tke/commit/9488c6a1421dc2b5b9998aca2d982cfb5f2678d6))
* provider config in installer ([02d6580](https://github.com/tkestack/tke/commit/02d6580dc02a9a4dd4c2238129da1779422ac1ed))
* provider-res path reference ([121bbd3](https://github.com/tkestack/tke/commit/121bbd38b28fef85b586d63efd797af03343cac1))
* registry ip for baremetal provider config ([3837ca3](https://github.com/tkestack/tke/commit/3837ca3a288ef33c1b22176f4fa46a4bab1c1dc9))
* release upload ([a372875](https://github.com/tkestack/tke/commit/a3728757a966177f2de90ed8f4dc24e9a6c94bad))
* release.sh xargs can't stop when error ([029d8b3](https://github.com/tkestack/tke/commit/029d8b3474731ccce430a14af24247d12f198400))
* remove policy from role spec when delete policy ([#186](https://github.com/tkestack/tke/issues/186)) ([9861e08](https://github.com/tkestack/tke/commit/9861e08443e2ef6fbd87c4642a673891384eb5ab))
* revert setting local timezone for alarms ([cd0b34f](https://github.com/tkestack/tke/commit/cd0b34f2da625a0ccf0fb1f8e883c4c1acdd2418))
* role policyunbinding empty ([#164](https://github.com/tkestack/tke/issues/164)) ([307d793](https://github.com/tkestack/tke/commit/307d793d64d817c32384d218c41718a32b38e9ab))
* set log level of prometheus-beat to info instead of debug ([52721ad](https://github.com/tkestack/tke/commit/52721add09c2716d122b5a166178197148f2320f))
* solve the problem of image push timeout ([05c39c8](https://github.com/tkestack/tke/commit/05c39c8f1acc9e3e13f8f1c20531a309c34bc651))
* spell in nvidia ([a888a57](https://github.com/tkestack/tke/commit/a888a574fd744e37c54d8128a74d83347ec990e9))
* test error ([9fff99d](https://github.com/tkestack/tke/commit/9fff99db9281ee5ca959b3edd5af0ddbeeb67709))
* unable to delete failed business namespace ([f5aa7aa](https://github.com/tkestack/tke/commit/f5aa7aa79c25f42c42e72815059c89275fa0d293))
* unable to delete project when registry disabled ([8c8d96f](https://github.com/tkestack/tke/commit/8c8d96f3e042e37229a007957cd527a5ee412377))
* unable to delete project when registry disabled ([07e43ba](https://github.com/tkestack/tke/commit/07e43ba7b3fd482c354fe7c50f8e953a0148700c))
* upgrade kube-state-metrics to v1.9.5 to support k8s 1.16 ([78ecd11](https://github.com/tkestack/tke/commit/78ecd11f4f45fe1ffb8471e18a0f3066025a1128))
* wrong label value for metrics getting from kube-state-metrics ([b7cf67f](https://github.com/tkestack/tke/commit/b7cf67f99e817d531752e859ad2e64023e24e80e))


### Features

* add clusterVersion and clusterDisplayName for business Namespace and NamespaceList ([72c9d22](https://github.com/tkestack/tke/commit/72c9d22fc08206c54ec1439067dfcccb8bea41ee))
* add experimental for dockerd ([5f42924](https://github.com/tkestack/tke/commit/5f4292418cd3589b924554744b11895287ddf8ed))
* add kubeadm reset -f for join node if error ([d3a1711](https://github.com/tkestack/tke/commit/d3a171173852d55d1f0bca34d43f052fa5354ae1))
* add list roles for policy ([#169](https://github.com/tkestack/tke/issues/169)) ([adeddc4](https://github.com/tkestack/tke/commit/adeddc4621e93647df23cd05ea1ac0e5a7a2dae5))
* add log for enhancing quota managing ([e84bfaa](https://github.com/tkestack/tke/commit/e84bfaadb959a709f3f9b2e75d498546f06e1421))
* add rollback subresource to deployment ([1f05e4e](https://github.com/tkestack/tke/commit/1f05e4e49b379df5c251f69f29ccc61b6cfd8bc2))
* add summary variable for alarms ([4b4aef5](https://github.com/tkestack/tke/commit/4b4aef560357f996927e7170714ee12af58ae062))
* add sync platform admins to auth idp admins ([#188](https://github.com/tkestack/tke/issues/188)) ([f525331](https://github.com/tkestack/tke/commit/f525331eb8624eb389a6a5d873c5769d41b9ad1c))
* add webhook channel type ([1fccdb1](https://github.com/tkestack/tke/commit/1fccdb1072b5519187a4b8a49549bec5c5c1e45e))
* allow CRD in apply interface ([9912950](https://github.com/tkestack/tke/commit/9912950756653c96f0e31e2f1e89c7e378401b1e))
* build arm64 images via make image PLATFORM=linux_arm64 ([7165264](https://github.com/tkestack/tke/commit/7165264bba82e8a910d2454276b1e9f1342148b4))
* build assets for gateway ([#143](https://github.com/tkestack/tke/issues/143)) ([712aa22](https://github.com/tkestack/tke/commit/712aa229db5b70f9082bc8b11ed153796407e9b3))
* clear docker volume when install ([e1afbc5](https://github.com/tkestack/tke/commit/e1afbc5b8d75542709b5d569aa1aa1b127960cb9))
* compatible with 1.16 version stateful API downgrade processing ([9201991](https://github.com/tkestack/tke/commit/92019919cf929687888153b0bfc4790048cf45d2))
* copy certs for installer registry ([4f6d56f](https://github.com/tkestack/tke/commit/4f6d56f0f11dbd3705f343058175823df13d8aa3))
* enable anonymous for docker distribution and chartmeseum ([8c2ea08](https://github.com/tkestack/tke/commit/8c2ea0864d2679076f3c728d95d36fc55ee92a90))
* enable docker cli experimental and prepare certs and registry store ([abbde7a](https://github.com/tkestack/tke/commit/abbde7ae5a548b17f0e79eed585933dd7e16cafa))
* enhance quota managing ([5cb0d1b](https://github.com/tkestack/tke/commit/5cb0d1b53e7f10d813512d6e935d8e789fea6647))
* genetate-images support multi arch ([3c0cb87](https://github.com/tkestack/tke/commit/3c0cb87aea356d778df3ee62c62c10fc0540aec2))
* genetate-images support multi arch ([daf4790](https://github.com/tkestack/tke/commit/daf4790edc65bb8fa89c339924355d822dc94f6b))
* health check for keepalived change nc to curl ([ffe63de](https://github.com/tkestack/tke/commit/ffe63de552f9eae2e74133402bd5de0079fb7094))
* image support multi arch ([2e3d8f9](https://github.com/tkestack/tke/commit/2e3d8f97dc9af67188d8d080b3ac707310c4f7ac))
* installer and galaxy support multi arch ([6ef9026](https://github.com/tkestack/tke/commit/6ef9026de69be8b61d53494a0aa937b9591c8760))
* modify flannel yaml in galaxy to support multi arch ([5405688](https://github.com/tkestack/tke/commit/540568825e8f13729a400bf9d02401966a2bca0b))
* modify the modules import and minors the size of node_modules in third_party ([#168](https://github.com/tkestack/tke/issues/168)) ([d900a79](https://github.com/tkestack/tke/commit/d900a791c24c550cab9b84bd911c6e466cde1c07))
* provider-res supports amd64 & arm64 platforms. Update provider-res version. ([d02d3c6](https://github.com/tkestack/tke/commit/d02d3c637e0ced5e563cd86c47d7efeeb0e5c756))
* refactor reference to provider res for multi arch ([27b6f55](https://github.com/tkestack/tke/commit/27b6f5565cde846672988e90a238aa0e7180b302))
* remove registry-mirrors in daemon.json ([b2e573e](https://github.com/tkestack/tke/commit/b2e573e737be9112672632a6af969ce26e8dfbd4))
* remove third party registry domain hostname validate ([96fe3a4](https://github.com/tkestack/tke/commit/96fe3a48d54efaacd7545c98512829028dc18cae))
* support cluster and machine controller sync period and concurrent in flag and config ([4df559f](https://github.com/tkestack/tke/commit/4df559f6b11e33e255adc649f85877030f01ad4e))
* support ha for cluster ([7a5c3f5](https://github.com/tkestack/tke/commit/7a5c3f5df3c64cb1b78b713e4c2bc0603e69c553))
* support hooks in creating cluster ([67cdcf7](https://github.com/tkestack/tke/commit/67cdcf72d395ea9f4d7f81c7f9bb746a8b5ece28))
* support label for master node ([da51428](https://github.com/tkestack/tke/commit/da51428563bf1b3b77637194551c06fb4b584dca))
* support role manage and group manage ([#176](https://github.com/tkestack/tke/issues/176)) ([e587da2](https://github.com/tkestack/tke/commit/e587da2c74365cf21bdb288edceacd2ddce381d4))
* support skip conditions in creating cluster ([b0a02b9](https://github.com/tkestack/tke/commit/b0a02b95c400461b0d91c13e3f7a7cefabddf30f))
* support the cronhpa & tapp ([#142](https://github.com/tkestack/tke/issues/142)) ([f053fb3](https://github.com/tkestack/tke/commit/f053fb3969d94883af06ee387c625701d2a69e13)), closes [#5](https://github.com/tkestack/tke/issues/5)
* support update auth related resource status with kubectl edit ([#163](https://github.com/tkestack/tke/issues/163)) ([f1de422](https://github.com/tkestack/tke/commit/f1de4224cfecfda0e86b17e6503eaa41be5750c9))
* update coredns etcd version for support multi arch ([72adeef](https://github.com/tkestack/tke/commit/72adeefb4af1d41b65c687a983bd7107eb468622))
* update parent when updating namespace quota ([edc3438](https://github.com/tkestack/tke/commit/edc3438d8ad74c9ca083c95955e52904aa1e72a7))
* upgrade galaxy version to v1.0.2 ([d8da656](https://github.com/tkestack/tke/commit/d8da6569e0200591b5d176a9fe9119bd36d78eef))
* upgrade image version ([ec2a6f8](https://github.com/tkestack/tke/commit/ec2a6f8de50375c6d206e8df27bae3fa4e46c364))



