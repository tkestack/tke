# Installer Optimization

## Current Installer Steps

0. Execute pre install hook
1. Load images (about 400s)
2. Tag images
3. Setup local registry
4. Push images (about 1370s)
5. Generate certificates for TKE components
6. Create global cluster (3 masters about 720s)
7. Write kubeconfig
8. Execute post cluster ready hook
9. Prepare front proxy certificates
10. Create namespace for install TKE
11. Prepare certificates
12. Prepare baremetal provider config
13. Install etcd
14. Patch platform versions in cluster info
15. Install tke-auth-api
16. Install tke-auth-controller
17. Install tke-platform-api
18. Install tke-platform-controller
19. Install tke-registry-api
20. Install tke-registry-controller
21. Install tke-business-api
22. Install tke-business-controller
23. Install InfluxDB
24. Install tke-monitor-api
25. Install tke-monitor-controller
26. Install tke-notify-api
27. Install tke-notify-controller
28. Install tke-logagent-api
29. Install tke-logagent-controller
30. Install tke-application-api
31. Install tke-application-controller
32. Install tke-gateway
33. Register tke api into global cluster
34. Import resource to TKE platform
35. Prepare push images to TKE registry
36. Push images to registry (about 3000s)
37. Set global cluster hosts
38. Import charts
39. Import Expansion Charts
40. Execute post install hook

All steps cost about 6400 seconds. Time cost may be different in different enviroments.

> - Most steps don't have time cost description, these steps' time cost is under 60s.
> - For all in one mode, installer node as global master node, there are 2 extra steps after step 31: `Prepare images before stop local registry` and `Stop local registry to give up 80/443 for tke-gateway`. These steps cost about 20 seconds.

## How to optimize installer steps

### 1. Push fewer images to local registry

We don't need to push all images to local registry. We should push essential images to local registry, like K8s related images, tke-platform-xx. After global cluster is created, we could swith local registry to tke-registry and push all images to tke-registry.

### 2. Stop using local registry

The temporary local registry is created to provide images for creating global cluster and TKE components. We could load these images to global node and remove local registry related steps.

### 3. TKE components should be independeng of installer

We should allow user to deploy TKE components in existing cluster through helm. And installer could convert input config to chart values to make a unified UX. Except essential components: `tke-auth`, `tke-platform`, and `tke-gateway`, user could ignore other components to accelerate installation.


### 4. Run images push as a job

As we have seen before, most time cost is images push related steps. `Load images` costs 400 seconds, `Push images (local)` costs 1370 seconds, and `Push images to registry` costs 3000 seconds. All steps cost about 6400 seconds, and images push related steps' ratio is about 75%. If we move these steps after TKEStack is deployed and run them as a job, it will make UX more friendly. To make it possible, we should seprate images packge to manage some images in installer steps and manage other images through K8s job when TKE components is running.