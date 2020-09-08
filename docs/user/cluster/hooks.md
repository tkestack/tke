# Using hooks in tkestack

HuXiaoLiang QianChenglong

2020-09-6

The admin user is able to deploy the customized hooks to tkestack so that inject customized logic during tkestack life cycle. The hook is a set of customized shell based scripts, which invoked by framework during tkestack life cycle

For example, admin want to add custom logic during cluster creation, such as installing software packages or modifying system configurations, the hook mechanism will be address this requirement

## Hook point

The hook script MUST be placed to `installer` container from installer node or `tke-platform-api` and `tke-platform-controller` container from `global` node so that it is invoked by framework. Admin user should make the decision which hook point get used according to below guideline:

- **Bootstrap cluster hook**: The hook placed in `installer` container：Inject customized logic during installer node and global node setup,  there are 3 types hook for this case,  Admin can put hook scripts in the `/app/hooks/`directory of the installer container. It should be noted that the names of the hooks must be corresponding.

    - pre-install which can be used to modify installer related configurations for customization
    - post-cluster-ready which can be used to customize tke components
    - post-install which can be used to automate your own deployment tasks or adjust global cluster or tke configuration

- **Business cluster hook**: The hook placed in `tke-platform` component: Inject customized logic during business cluster and node life cycle,  there are 10 types hook include `cluster` scope and `node` scope for this case, the hook type is self-explanatory,  you should know what it did from hook type.  Admin can put these hook script to `configmap` or`pv`, then mount it `tke-platform` deployment later. 
    - HookPreInstall
    - HookPostInstall
    - HookPreUpgrade
    - HookPostUpgrade
    - HookPreClusterInstall
    - HookPostClusterInstall
    - HookPreClusterUpgrade
    - HookPostClusterUpgrade
    - HookPreClusterDelete
    - HookPostClusterDelete

## Hook scope (only works for business cluster hook)

The scope is hook script wording scope, currently, tkestack support `cluster` and `node` scope：
- `cluster` scope: The hook will be executed only on first master node from business cluster during cluster life cycle:
   - HookPreClusterInstall
    - HookPostClusterInstall
    - HookPreClusterUpgrade
    - HookPostClusterUpgrade
    - HookPreClusterDelete
    - HookPostClusterDelete
- `node` scope: The hook will be executed on all business cluster nodes include master node during node life cycle
    - HookPreInstall
    - HookPostInstall
    - HookPreUpgrade
    - HookPostUpgrade

## Pass parameter to hook (only works for business cluster hook)

In some scenarios, admin want to pass parameters to hook script so that control script execution flow. To achieve this,  admin should pre-define the key-value pair to `cluster` object `labels` or `annotations`section, then framework will combine them together and pass them to hook script,  you can use `$@` to get all of them.  The multiple key-value pairs separated by `blank`, please refer to `Example 2` for the details. Note that it is only supported for `cluster`scope hook now

## Example 1: 

The hook is single file, it needs to be accessible in tke-platform-api and tke-platform-controller, and can be provided through hostPath, volumeMounts, etc.

- tke-platform-api will validate the files which specified by cluster
- tke-platform-controller will invoke cluster provider for copy the files to every cluster node

The following is demonstrated through ConfigMap

1. Create ConfigMap

    `kubectl -ntke create configmap hooks --from-file=.`
    ![](../../images/2020-03-31-10-34-13.png)

2. Mount

    ```bash
    kubectl -ntke edit deploy tke-platform-api
    kubectl -ntke edit deploy tke-platform-controller
    ```
    ![](../../images/2020-03-31-10-43-23.png)
    ![](../../images/2020-03-31-10-43-01.png)
    ![](../../images/2020-03-31-10-41-50.png)
    ![](../../images/2020-03-31-10-45-44.png)

3. Declaring required files for hooks

The console does not support this feature currently, you can create a cluster through kubectl.

In your global cluster or setup kubeconfig for tke:

`kubectl create -f cluster.yaml`

cluster.yaml:
```yaml
{
    "apiVersion": "platform.tkestack.io/v1",
    "kind": "Cluster",
    "spec": {
        "clusterCIDR": "10.244.0.0/16",
        "type": "Baremetal",
        "version": "1.14.10",
        "machines": [
            {
                "ip": "1.2.3.4",
                "port": 22,
                "username": "root",
                "password": "MTIzNDU2" // echo -n "123456" | base64
            }
        ],
        "features": {
            "files": [
                {
                    "src": "hooks/pre-install", // source file in tke-platform-controller
                    "dst": "/tmp/hooks/pre-install" // destinate file in node
                },
                                {
                    "src": "hooks/post-install",
                    "dst": "/tmp/hooks/post-install"
                }
            ],
            "hooks": {
                "PreInstall": "/tmp/hooks/pre-install", // reference to destinate file which defined in files
                "PostInstall": "/tmp/hooks/post-install"
            }
        }
    }
}
```

## Example 2

The hook scripts from a directory, it needs to be accessible in tke-platform-api and tke-platform-controller, and can be provided through hostPath, volumeMounts, etc.

- tke-platform-api will validate the files which specified by cluster
- tke-platform-controller will invoke cluster provider for copy the files from directory to every cluster node

The following is demonstrated through ConfigMap

1.  Prepare a set of hooks from directory
```
root@VM-0-127-ubuntu:~# tree hooks/         
hooks/
├── post
│   ├── cluster-post-install
│   └── post-install
└── pre
    ├── cluster-pre-install
    └── pre-install
```

```
root@VM-0-127-ubuntu:~# find hooks/ -type f -print | xargs -I file sh -c 'ls file; cat file'  
hooks/pre/pre-install
#!/bin/bash
echo "pre node install"
hooks/pre/cluster-pre-install
#!/bin/bash
echo "the cluster labels and annotations key=value pair:$@"
echo "pre cluster install" 
hooks/post/cluster-post-install
#!/bin/bash
echo "the cluster labels and annotations key=value pair:$@"
echo "post cluster install"
hooks/post/post-install
#!/bin/bash
echo "post node install"
```

2. Copy hook directory to `tke-platform-api` and `tke-platform-controller`. In production env, you should use `pv` to persist the hook directory and `pvc` to bind `deployment`
  
 ```
 docker cp hooks $tke-platform-api-container-id:/app
 docker cp hooks $tke-platform-controller-container-id:/app
 ```

3. Declaring required files for hooks

The console does not support this feature currently, you can create a cluster through kubectl.

In your global cluster or setup kubeconfig for tke:

`kubectl create -f cluster.json`

cluster.json:
```json
root@VM-0-127-ubuntu:~# cat cluster3 
{
    "apiVersion": "platform.tkestack.io/v1",
    "kind": "Cluster",
    "metadata": {
       "labels": {
          "location": "xa"
       },
       "annotations": {
          "node.alpha.kubernetes.io/test": "hello"
       }
    },
    "spec": {
        "clusterCIDR": "10.244.0.0/16",
        "tenantID": "default",
        "type": "Baremetal",
        "version": "1.18.3",
        "machines": [
            {
                "ip": "10.0.0.9",
                "port": 22,
                "username": "root",
                "password": "TGV0bWVpbjEyMywuMQ=="
            }
        ],
        "features": {
            "files": [
                {
                    "src": "hooks/pre", 
                    "dst": "/tmp/hooks/pre"
                },
                {
                    "src": "hooks/post",
                    "dst": "/tmp/hooks/post"
                }
            ],
            "hooks": {
                "PreInstall": "/tmp/hooks/pre/pre-install",
                "PostInstall": "/tmp/hooks/post/post-install",
                "PreClusterInstall": "/tmp/hooks/pre/cluster-pre-install",
                "PostClusterInstall": "/tmp/hooks/post/cluster-post-install"
            }
        }
    }
}
```

## Notes

1. The source file must be a **regular** file!
2. The destinate file parent directories will be created on demand.
3. Every node including master and worker node will copy file and run the hook file!
4. If you wish run only once, your hook file need to guarantee idempotence!
5. If hook file occur error which exit code is's 0, the hook file will retry and cluster conditions will be hang in EnsureXXXHook.