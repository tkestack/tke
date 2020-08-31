# CronHPA

Cron Horizontal Pod Autoscaler(CronHPA)使我们能够使用[crontab](https://en.wikipedia.org/wiki/Cron)模式定期自动扩容工作负载（那些支持扩展子资源的负载，例如deployment、statefulset）。

CronHPA使用[Cron](https://en.wikipedia.org/wiki/Cron)格式进行编写，周期性地在给定的调度时间对工作负载进行扩缩容。

## CronHPA 资源结构

CronHPA定义了一个新的CRD，cron-hpa-controller是该CRD对应的controller/operator，它解析CRD中的配置，根据系统时间信息对相应的工作负载进行扩缩容操作。

```go
// CronHPA represents a set of crontabs to set target's replicas.
type CronHPA struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired identities of pods in this cronhpa.
	Spec CronHPASpec `json:"spec,omitempty"`

	// Status is the current status of pods in this CronHPA. This data
	// may be out of date by some window of time.
	Status CronHPAStatus `json:"status,omitempty"`
}

// A CronHPASpec is the specification of a CronHPA.
type CronHPASpec struct {
	// scaleTargetRef points to the target resource to scale
	ScaleTargetRef autoscalingv2.CrossVersionObjectReference `json:"scaleTargetRef" protobuf:"bytes,1,opt,name=scaleTargetRef"`

	Crons []Cron `json:"crons" protobuf:"bytes,2,opt,name=crons"`
}

type Cron struct {
	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	Schedule string `json:"schedule" protobuf:"bytes,1,opt,name=schedule"`

	TargetReplicas int32 `json:"targetReplicas" protobuf:"varint,2,opt,name=targetReplicas"`
}

// CronHPAStatus represents the current state of a CronHPA.
type CronHPAStatus struct {
	// Information when was the last time the schedule was successfully scheduled.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty" protobuf:"bytes,2,opt,name=lastScheduleTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CronHPAList is a collection of CronHPA.
type CronHPAList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CronHPA `json:"items"`
}
```

## 使用CronHPA

### 安装CronHPA

CronHPA为TKEStack扩展组件，需要在【平台管理】-> 【扩展组件】里安装该组件。

### 示例1：指定deployment每周五20点扩容到60个实例，周日23点缩容到30个实例

```yaml
apiVersion: extensions.tkestack.io/v1
kind: CronHPA
metadata:
  name: example-cron-hpa	# CronHPA 名
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment	# CronHPA操作的负载类型
    name: demo-deployment	# CronHPA操作的负载类型名
  crons:
    - schedule: "0 20 * * 5"	# Crontab语法格式
      targetReplicas: 60			# 负载目标pod数量
    - schedule: "0 23 * * 7"
      targetReplicas: 30
```

### 示例2：指定deployment每天8点到9点，19点到21点扩容到60，其他时间点恢复到10

```yaml
apiVersion: extensions.tkestack.io/v1
kind: CronHPA
metadata:
  name: web-servers-cronhpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: web-servers
  crons:
    - schedule: "0 8 * * *"
      targetReplicas: 60
    - schedule: "0 9 * * *"
      targetReplicas: 10
    - schedule: "0 19 * * *"
      targetReplicas: 60
    - schedule: "0 21 * * *"
      targetReplicas: 10
```

### 查看cronhpa

```shell
# kubectl get cronhpa
NAME               AGE
example-cron-hpa   104s

# kubectl get cronhpa example-cron-hpa -o yaml
apiVersion: extensions.tkestack.io/v1
kind: CronHPA
...
spec:
  crons:
  - schedule: 0 20 * * 5
    targetReplicas: 60
  - schedule: 0 23 * * 7
    targetReplicas: 30
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: demo-deployment

```

### 删除cronhpa

```shell
kubectl delete cronhpa example-cron-hpa
```
