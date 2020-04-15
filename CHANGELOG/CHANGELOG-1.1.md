# [1.1.0](https://github.com/tkestack/tke/compare/v1.0.1...v1.1.0) (2020-02-14)


### Bug Fixes

* accessing the cluster resource with the wrong address ([b18be8d](https://github.com/tkestack/tke/commit/b18be8dac71194c475c08853207614877319126e))
* add ShouldDeleteDuringUpdate check and add list related policy/role ([7692bf5](https://github.com/tkestack/tke/commit/7692bf5d37013715cdb534ac9a2b0eec6108dc02))
* add tke-gateway dns name to installer ([893c14f](https://github.com/tkestack/tke/commit/893c14fb5d15882703cf6f1891d79411ecde428d))
* auth processUpdate when phase is not terminating ([#105](https://github.com/tkestack/tke/issues/105)) ([e36a99f](https://github.com/tkestack/tke/commit/e36a99f7ccbd79ce37bc0832f5f03d665290e321))
* auth-api and auth-controller json config ([#124](https://github.com/tkestack/tke/issues/124)) ([af0666a](https://github.com/tkestack/tke/commit/af0666a24d1b5f56987544f7ce889143f039c39d))
* authz handler type error ([930cb63](https://github.com/tkestack/tke/commit/930cb638bb89d15876f8dbe995d31a9f08d9e394))
* binding check subjects and auth resource delete ([b85e454](https://github.com/tkestack/tke/commit/b85e4544bf9672fbb7de7e90ef79ad3764de527e))
* change etcd client package to github.com/etcd-io ([182e831](https://github.com/tkestack/tke/commit/182e8312a3f3787aa7c692fd8ced63c01e9a947a))
* change imported package name ([a736311](https://github.com/tkestack/tke/commit/a73631124cbdd3e708fa00b59aa6d42367c9fd5c))
* create cluster conflict ([2743663](https://github.com/tkestack/tke/commit/27436637631ea77dc87bdefb96b0041def9c6f5e))
* default registry domain ([1eb2bf7](https://github.com/tkestack/tke/commit/1eb2bf7f5ca3756aba70dfef0b5fb1f175f30626))
* deleted vGPU and resources are not cleaned up ([392e9e6](https://github.com/tkestack/tke/commit/392e9e6cb35a8ec6d774c525571faaa263ec36c5))
* high version of k8s requires configuring cni version ([272c55d](https://github.com/tkestack/tke/commit/272c55d69fb00266d5d03322d293a303a4ccc4ca))
* intall failed when disable monitoring, close [#109](https://github.com/tkestack/tke/issues/109) ([7a9a214](https://github.com/tkestack/tke/commit/7a9a2141daff5f4e2df2a457b7721d61b7c965c1))
* k8s version ([14e8e72](https://github.com/tkestack/tke/commit/14e8e72bc0fcf3ab19fc698ad8871eeecf7b517a))
* lint ([3ebf999](https://github.com/tkestack/tke/commit/3ebf99901417f42ca364e5b6ce6529a0a1ae816b))
* lint ([02d0b2f](https://github.com/tkestack/tke/commit/02d0b2fb418969009713271d91ca25a25e324cfc))
* provider config when use thirdparty registry ([c446a1e](https://github.com/tkestack/tke/commit/c446a1e1ebcfe907236e607ddfe3f09fa6e583e8))
* regenerate codes for auth api ([fe8cf01](https://github.com/tkestack/tke/commit/fe8cf0116da3ef89377c72a7925bc73f55f28edb))
* release provider-res version ([d9d7194](https://github.com/tkestack/tke/commit/d9d719465a444a96895267efa09faf7345879db8))
* remove exec log of ssh ([41f29d8](https://github.com/tkestack/tke/commit/41f29d895833413d30f40136fa4084ec97bf4d1f))
* remove get user groups from local for generic idp ([82a6864](https://github.com/tkestack/tke/commit/82a686491b7a47c60352e3f4f295b30ab490d4d5))
* remove unused error handling ([b9e4a63](https://github.com/tkestack/tke/commit/b9e4a63d84953b235f39160d2f8c6673a45b81e7))
* remove useless images ([11a305d](https://github.com/tkestack/tke/commit/11a305d11de01228767e133335fd8b29f703b84f))
* resolve conflict with chartmuseum and kubernetes ([1e11413](https://github.com/tkestack/tke/commit/1e1141398bfeea6806d5dbf650f262e6eb66c7ad))
* role sync and some unnecessary logs ([e90e637](https://github.com/tkestack/tke/commit/e90e6373c9bbc9ca70a725041950ba5198af9263))
* set /etc/hosts when using third party registry, close [#121](https://github.com/tkestack/tke/issues/121) ([e811b17](https://github.com/tkestack/tke/commit/e811b172fa8d8a069f7aad2d4d64586b428a4aa9))
* set /etc/hosts when using third party registry, close [#121](https://github.com/tkestack/tke/issues/121) ([704478d](https://github.com/tkestack/tke/commit/704478d2005eb0ecb9e0c8a453d2a23348fadfde))
* solve the problem of pulling business cluster resources error ([3ed15d7](https://github.com/tkestack/tke/commit/3ed15d71b9696da5ceb46c51338e9d44e04b08a0))
* typo in prometheus controller ([f6080ba](https://github.com/tkestack/tke/commit/f6080ba39f451e22f9a9a0453191429aa6c696f9))
* update dex to forked dex ([12dfcd0](https://github.com/tkestack/tke/commit/12dfcd090e12bac5ae720073db77992f847cc714))
* Update running locally guide ([62061e3](https://github.com/tkestack/tke/commit/62061e352afa7825ef331138c6f203c1be13ce43))


### Features

* add aws s3 object storage configuration for chartmeseum ([9c80272](https://github.com/tkestack/tke/commit/9c802726f7e369c6bc49b8254632f1b1b17207c4))
* add category and policy config for auth controller ([f2485f0](https://github.com/tkestack/tke/commit/f2485f0505309150e8cf074e224db2cd3db82e0f))
* add chartmeseum support ([525b0ea](https://github.com/tkestack/tke/commit/525b0eadb5b5734a2ae0e7f66f875130cff7836d))
* aggregate tke-auth ([1549970](https://github.com/tkestack/tke/commit/1549970be8ed7ebc50f2e6a04d5ec537f6b51946))
* allow only one installation process, close [#101](https://github.com/tkestack/tke/issues/101) ([4ed5465](https://github.com/tkestack/tke/commit/4ed5465351dc9adcae407a93149d9fee0201c160))
* build assets for gateway ([8f92a29](https://github.com/tkestack/tke/commit/8f92a2937dae953854edb721f62e47372bb38af1))
* change the apiVersion for clusterVersion that over 1.12, and fix the casing letter of filname ([b4f5e74](https://github.com/tkestack/tke/commit/b4f5e74dc8099a26627db5b63d9b6ecfb04f07a6))
* cluster and machine support TableConvertor ([36a5b87](https://github.com/tkestack/tke/commit/36a5b87f7fec5f79616b7a1182744b70335304a7))
* compatible with 1.16 version API downgrade processing ([d532919](https://github.com/tkestack/tke/commit/d532919a81a0bf52abca66b98c7c233a4caf03f0))
* installer excludes using itself during installation, close [#116](https://github.com/tkestack/tke/issues/116) ([2122af6](https://github.com/tkestack/tke/commit/2122af660111d03f2ed9fbf57004d08d605cdb93))
* modify the request path of users ([#119](https://github.com/tkestack/tke/issues/119)) ([7b88a64](https://github.com/tkestack/tke/commit/7b88a6445ba85ffc3bf98994af73526fb96afc5d)), closes [#3](https://github.com/tkestack/tke/issues/3)
* regenerate api codes for auth ([1ef3bad](https://github.com/tkestack/tke/commit/1ef3bad83c6030e4ebcd86b21a7a7e1027dd88ee))
* remove chartmuseum temporarily ([ce6b533](https://github.com/tkestack/tke/commit/ce6b533a89de4614d861c80ada21da6175c77997))
* rename tke-generate-images to generate-images ([0168599](https://github.com/tkestack/tke/commit/01685995b6b0ee637e36e02af45b2b27467a0b0f))
* set max upload size to 20g ([b4b9953](https://github.com/tkestack/tke/commit/b4b9953e94083e501be4c9261d47c128a5859603))
* set registry domain default value ([42f8854](https://github.com/tkestack/tke/commit/42f8854354f0f5b63c6d9ec576accd2b48ecec86))
* support installer to listen on all network interface, close [#95](https://github.com/tkestack/tke/issues/95) ([15eda37](https://github.com/tkestack/tke/commit/15eda37498d6ffd4735d0f2b61dacc188f447ee6))
* support k8s version 1.14.10 1.16.6 ([ee2192c](https://github.com/tkestack/tke/commit/ee2192ce9f84b3a44670e79c51973aae5e8e900c))
* support logcollector and modify the resourceversion for that ([ba42fab](https://github.com/tkestack/tke/commit/ba42fabeff2ea872685bb557e9cd380b390b4a57))
* support self-defined notify webhook ([f021e6d](https://github.com/tkestack/tke/commit/f021e6db25ffc914754918ce24cd695a62eb70a4))
* ui support 1.16.6 ([a9ca35b](https://github.com/tkestack/tke/commit/a9ca35b800af147efbcac87881000edf3d2ba100))
* update code-generator and regenerate codes ([4f160ac](https://github.com/tkestack/tke/commit/4f160ac8d7bf08e8fc6f4d8d67e45308537784ac))
* write kubeconfig when post cluster ready, close [#120](https://github.com/tkestack/tke/issues/120) ([b5cc722](https://github.com/tkestack/tke/commit/b5cc72287a710c8c261e12c54a7645402a5b91bd))


### Performance Improvements

* update apiserver to 0.17.0 ([5c1d91a](https://github.com/tkestack/tke/commit/5c1d91aad008c7a7cadd7b82baefc3d265008c55))
* update code generator to 1.17.0 ([3e35d73](https://github.com/tkestack/tke/commit/3e35d7397a9dd936935b13de040d53e07813d921))


### Reverts

* Revert "Integrate rook in global cluster" ([9b40ad9](https://github.com/tkestack/tke/commit/9b40ad9edd5b95eaf23093d500b8a10e68200dee))
