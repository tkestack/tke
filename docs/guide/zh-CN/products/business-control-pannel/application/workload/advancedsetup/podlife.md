# Pod 生命周期和重启策略

Pod 在整个生命周期中被系统定义为各种状态，熟悉 Pod 的各种状态对于理解如何设置Pod的调度策略、重启策略是很有必要的。

| Pod 状态  | 描述                                                         |
| --------- | ------------------------------------------------------------ |
| Pending   | APIServer 已经创建 Pod，但 Pod 内还有一个或多个容器的镜像没有创建，包括正在下载镜像的过程 |
| Running   | Pod 内所有容器均已创建，且至少有一个容器处于运行状态、正在启动状态或正在重启状态 |
| Succeeded | Pod 内所有容器均成功执行后退出，且不会再重启                 |
| Failed    | Pod 内所有容器均已退出，但至少有一个容器退出为失败状态       |
| Unknown   | 由于某种原因无法获取该 Pod 状态，可能由于网络不通畅导致      |

Pod 的重启策略（RestartPolicy）应用于 Pod 内的所有容器，并且仅在 Pod 所处的 Node 上由 kubelet 进行判断和重启操作。当某个容器异常退出或者健康检查失败时，kubelet 将根据 RestartPolicy 的设置来进行相应的操作。

Pod 的重启策略包括 Always、OnFailure 和 Never，默认值为 Always。

* Always：当容器失效时，由kubelet自动重启该容器。

* OnFailure：当容器终止运行且退出码不为 0 时，由 kubelet 自动重启该容器。

* Never：不论容器运行状态如何，kubelet 都不会重启该容器。

kubelet 重启失效容器的时间间隔以 sync-frequency 乘以 2n 来计 算，例如1、2、4、8倍等，最长延时 5min，并且在成功重启后的 10min 后重置该时间。

Pod 的重启策略与控制方式息息相关，当前可用于管理 Pod 的控制器包括ReplicationController、Job、DaemonSet 及直接通过 kubelet 管理（静态Pod）。每种控制器对 Pod 的重启策略要求如下。

* RC 和 DaemonSet：必须设置为Always，需要保证该容器持续运行。

* Job：OnFailure或Never，确保容器执行完成后不再重启。

* kubelet：在 Pod 失效时自动重启它，不论将 RestartPolicy 设置为什么值，也不会对Pod进行健康检查。

