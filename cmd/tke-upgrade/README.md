## How To Upgrate TKEStack

### git clone
```shell script
git clone http://git.code.oa.com/eetke/upgrade.git
```

### find installer data in /opt/tke-install/data and copy to tke-upgrade
```shell script
cp /opt/tke-installer/data/tke.json cmd/tke-upgrade/app/data/tke.json
cp /opt/tke-installer/data/oidc_client_secret cmd/tke-upgrade/app/data/oidc_client_secret
```

### modify environment variables(ex: VERSION) in upgrade.sh 
```shell script
vi ./cmd/tke-upgrade/upgrade.sh
```

### execute upgrate
```shell script
./cmd/tke-upgrade/upgrade.sh
```



