# OIDC config and rollback scripts

## How to use oidc.sh script
`oidc.sh` script use to config the OIDC login.
When the first tile use the `oidc.sh`.You should fill in the configuration file `oidc.conf`then execute the `oidc.sh`with the `oidc.conf`:
```
chmod +x oidc.sh
./oidc.sh oidc.conf
```
After that,it will store the `oidc.conf` in configmap, you can check:
```
kubectl get cm -n tke oidc-config
```
Except for the first time，you can execute the script without the script:
```
./oidc.sh
```
## How to use rollback.sh script
`rollback.sh` use to rollback OIDC login, after rollback you can use tkeanywhere to login.
```
chmod +x rollback.sh
./rollback.sh
```

