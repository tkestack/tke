# [1.7.0](https://github.com/tkestack/tke/compare/v1.6.0...v1.7.0) (2021-06-02)


### Bug Fixes

* **ci:** add env from esecret in release test ([#1229](https://github.com/leoryu/tke/issues/1229)) ([aa15e15](https://github.com/leoryu/tke/commit/aa15e1580b459d5df8c8c057a9cd0ca43ca75c2a))
* **ci:** build images without buildkit ([#1224](https://github.com/leoryu/tke/issues/1224)) ([21f3cce](https://github.com/leoryu/tke/commit/21f3cce13c5d9b5c7a88d8c01dde8ee4a0f009ef))
* **console:** add package-lock.json ([#1171](https://github.com/leoryu/tke/issues/1171)) ([55924b3](https://github.com/leoryu/tke/commit/55924b3f3b5353b031d93c9fa84307aebc19877e))
* **console:** fix  some attribute  not have , let ui error ([#1222](https://github.com/leoryu/tke/issues/1222)) ([8ba52cd](https://github.com/leoryu/tke/commit/8ba52cd0ea7c0ed0f92a09cccd6c1756dcebccab))
* **console:** import package wrong ([#1191](https://github.com/leoryu/tke/issues/1191)) ([e900c2b](https://github.com/leoryu/tke/commit/e900c2b24b36b27a6da4658c50840f23a16c722a))
* **console:** k8s version updgrade and kubeconfig wrong ([#1183](https://github.com/leoryu/tke/issues/1183)) ([d336b92](https://github.com/leoryu/tke/commit/d336b924dbcf15cb62ec3cce847fb9f6cac9dd60))
* **console:** virtual node cpu&memory display wrong ([#1211](https://github.com/leoryu/tke/issues/1211)) ([7ab9db8](https://github.com/leoryu/tke/commit/7ab9db8b8512771359973f54e34cefefcaeb1de9))
* **monitor-controller:** kubelet metrics name changed ([#1187](https://github.com/leoryu/tke/issues/1187)) ([6433051](https://github.com/leoryu/tke/commit/6433051f63129e44ee086e9479ab2953565c8ae5))
* **platform:** clean cilium net before install ([#1280](https://github.com/leoryu/tke/issues/1280)) ([04e8eef](https://github.com/leoryu/tke/commit/04e8eef857d3c05bb5b6d72b692dbf82af8e6b22))
* **platform:** identify version inc info ([#1170](https://github.com/leoryu/tke/issues/1170)) ([34b4fb4](https://github.com/leoryu/tke/commit/34b4fb4b1ead501cc396533d9a621d2456f4e35c))
* **platform:** remove cilium enabled by default ([#1272](https://github.com/leoryu/tke/issues/1272)) ([46a8135](https://github.com/leoryu/tke/commit/46a81354aa796286978e0ba9da32660d540a1f9f))
* **platform:** remove cilium enabled by default ([#1272](https://github.com/leoryu/tke/issues/1272)) ([dc32702](https://github.com/leoryu/tke/commit/dc32702ba279f688f42c3cd865b1b0f67dc523a5))
* **platform:** remove v in same version comparing ([#1172](https://github.com/leoryu/tke/issues/1172)) ([a024c06](https://github.com/leoryu/tke/commit/a024c064880d9180dc8b6d615ffc58b64bb7f903))
* **platform:** set ownerreference for imported cluster credentials ([#1194](https://github.com/leoryu/tke/issues/1194)) ([6966a86](https://github.com/leoryu/tke/commit/6966a86f78d6d834a4d538270509bcdd078f8190))
* **platform:** update cilium according to the latest doc change ([#1324](https://github.com/leoryu/tke/issues/1324)) ([08cdc03](https://github.com/leoryu/tke/commit/08cdc0310499e96c7137adba68d036a50e30cdf2))
* **release, test:** create cluster failed ([#1189](https://github.com/leoryu/tke/issues/1189)) ([9b7f134](https://github.com/leoryu/tke/commit/9b7f134744c57998004b535f7cd9b4195190605f))
* **test:** create docker auth in build vm ([#1198](https://github.com/leoryu/tke/issues/1198)) ([dabd522](https://github.com/leoryu/tke/commit/dabd522ebed41d13199deeeb0a7af5307d8d2f80))
* **test:** retry create cluster in e2e test ([#1188](https://github.com/leoryu/tke/issues/1188)) ([e100e77](https://github.com/leoryu/tke/commit/e100e779e0201c1e600feb74c76b8600358e5b3c))
* **test:** wait cluster condition with more time ([#1199](https://github.com/leoryu/tke/issues/1199)) ([37b4b4d](https://github.com/leoryu/tke/commit/37b4b4dd60ad73b4fb208852924871202fabf037))
* **webconsole:** install helm v3 ([#1203](https://github.com/leoryu/tke/issues/1203)) ([c9bfe5a](https://github.com/leoryu/tke/commit/c9bfe5ad44d139feee39f6f1a0e64ef432ab6665))
* **webconsole:** rewritten clustercredential into cc ([#1177](https://github.com/leoryu/tke/issues/1177)) ([b2a6a71](https://github.com/leoryu/tke/commit/b2a6a712119ec3b00a67b270d5986f217fdbd0fd))


### Features

* **console:** support cilium ([#1273](https://github.com/leoryu/tke/issues/1273)) ([c148420](https://github.com/leoryu/tke/commit/c14842027a746c999f56e286e15a2216c7c7d638))
* **console:** support select networkMode ([#1306](https://github.com/leoryu/tke/issues/1306)) ([153c503](https://github.com/leoryu/tke/commit/153c503eba30897a87bd37f8826e8a4d62e05a41))
* **doc:** add replace service web static file doc ([#1223](https://github.com/leoryu/tke/issues/1223)) ([e783e10](https://github.com/leoryu/tke/commit/e783e10bdbd85966442c88c19267feb6ae0d0ab7))
* **installer:** push chart with force flag ([#1205](https://github.com/leoryu/tke/issues/1205)) ([6a9710b](https://github.com/leoryu/tke/commit/6a9710b9e5183f41f42413142592311fa46409c1))
* **installer:** push charts in installing ([#1182](https://github.com/leoryu/tke/issues/1182)) ([975df20](https://github.com/leoryu/tke/commit/975df203d89528ea407c5bf6ff161d88f1d3893b))
* **installer:** set install application store as default ([#1235](https://github.com/leoryu/tke/issues/1235)) ([c073abf](https://github.com/leoryu/tke/commit/c073abf495a34a081cbbc1195b6bb5be9113c6e3))
* **platform:**  change parameter from string to int ([#1316](https://github.com/leoryu/tke/issues/1316)) ([e641327](https://github.com/leoryu/tke/commit/e64132792686890f67623a68704d16ac933afed3))
* **platform:** add 1.20.4 tke ([#1174](https://github.com/leoryu/tke/issues/1174)) ([c9983b2](https://github.com/leoryu/tke/commit/c9983b28065da47e277666011950164f144c09bc))
* **platform:** add kube vendor in cluster ([#1186](https://github.com/leoryu/tke/issues/1186)) ([e015eae](https://github.com/leoryu/tke/commit/e015eaef8978aa6f54176fe1f6b2d1a38de8c4b6))
* **platform:** add lable on master[0] for cilium underlay ([#1313](https://github.com/leoryu/tke/issues/1313)) ([77875e9](https://github.com/leoryu/tke/commit/77875e9f971aa49bb50261a4b3414efd67d69fd0))
* **platform:** declare default table convertor ([#1162](https://github.com/leoryu/tke/issues/1162)) ([526b9c4](https://github.com/leoryu/tke/commit/526b9c4286cd7007544c3720b34f2a4454fd0f1a))
* **platform:** fix eni ipamd yaml format error ([#1314](https://github.com/leoryu/tke/issues/1314)) ([3743a3a](https://github.com/leoryu/tke/commit/3743a3a17f8966712241ef8520377d343af03bf7))
* **platform:** remove TargetRAMM ([#1167](https://github.com/leoryu/tke/issues/1167)) ([c08f642](https://github.com/leoryu/tke/commit/c08f6424db6605dd8c707a15c5097d29ef01a72b))
* **platform:** stop using internalversion ([#1169](https://github.com/leoryu/tke/issues/1169)) ([e2797cc](https://github.com/leoryu/tke/commit/e2797cce17fc45b008f89a452afcd43b9f5f5813))
* **platform:** stop using ListWatchUtil ([#1166](https://github.com/leoryu/tke/issues/1166)) ([82ba9fe](https://github.com/leoryu/tke/commit/82ba9fedeaac3bbb275bc8dc833db55a72ee0e9d))
* **platform:** support cilium overlay underlay network modes ([#1309](https://github.com/leoryu/tke/issues/1309)) ([3f01691](https://github.com/leoryu/tke/commit/3f01691762dd087492e30588ee1bba6c06f7a4c6))
* **platform:** support cilium to cluster feature ([#1252](https://github.com/leoryu/tke/issues/1252)) ([b5bee04](https://github.com/leoryu/tke/commit/b5bee049a5249178842c9ee8cc32b97bc24d191c))
* **registry:** modify password to request auth ([#1115](https://github.com/leoryu/tke/issues/1115)) ([d4231cd](https://github.com/leoryu/tke/commit/d4231cd40f16053391c50ab3701e55934df3292e))
* **registry:** using etcd storage for chartmuseum ([#1333](https://github.com/leoryu/tke/issues/1333)) ([148121b](https://github.com/leoryu/tke/commit/148121b11434c82278a26dea681d01f77350e9d4))
* **test:** update test for k8s 1.20.4 ([#1178](https://github.com/leoryu/tke/issues/1178)) ([c62404e](https://github.com/leoryu/tke/commit/c62404e47ec990638f62a02a3935d3457239b72d))
* **test:** use githubactions dokcer auth ([#1196](https://github.com/leoryu/tke/issues/1196)) ([0692014](https://github.com/leoryu/tke/commit/06920149256aa14c8b68174a25be8ee8abd439e8))


### Reverts

* Revert "Remove EnsureDisableOffloading" (#1161) ([0e8cf9e](https://github.com/leoryu/tke/commit/0e8cf9e3b1cdd9fba7ac472597d04a73f87c5842)), closes [#1161](https://github.com/leoryu/tke/issues/1161)
