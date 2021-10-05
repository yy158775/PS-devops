# volumn

​	volumn：存储卷技术绕过存储驱动程序，将数据存储在宿主机可达的地方，跨容器生命周期的持久化。

# kubernetes

共享网络将多个主机

主节点

工作节点

容器调度

负载均衡

消息队列

服务发现

# pod

包含一个或者多个容器，必须运行在同一主机上

# 资源清单

元数据

所需状态



# 命令

```
kubernetes describe pod xxx
kubernetes logs pod xxx
kubectl scale deployment xxx --replicas=4
kubectl autoscale deployment chat-redis --cpu-percent=50 --min=1 --max=3
kubectl get pods
kubectl get hpa
kubectl get svc
kubectl top pod
```

# 结构

API server：

API 对象：有特定的类型，Node,Namespace,,Pod,Service,Deployment

还有负责Pod扩缩容功能，但是使用命令，kubectl毕竟是手动的功能，Kubernetes也为我们提供了一个这样的资源对象Horizontal Pod  Autoscaling Pod水平扩缩容，

## kube-proxy

ClusterIP转发后端POD节点

POD的代理和负载均衡器，确保Node,Service,Pod对象之间有效通信

## Pod

本质上是共享NetWrok,IPC 和UTS名称空间和存储资源的容器集合。

同一Pod中的各个容器，共享网络协议栈，网络设备，路由，IP地址和端口

​	各个Pod对象的IP地址位于同一网络平面内（同一IP网段），各Pod间可通过真实的IP地址直接通信，无须NAT功能介入。

## Service

​	作为一组提供了相同服务的Pod对象的访问入口，由客户端Pod向目标Pod所属的Service对象的IP地址发起访问请求。

```
apiVersion: v1
kind: Service
metadata:
  name: chat-room
  namespace: exp
  labels:
    app: chat-room
spec:
  selector:
    app: chat-room
  ports:
    - name: rpc
      port: 8080
      targetPort: 8080
    - name: response
      port: 8081
      targetPort: 8081
```

## Deployment

期望的副本数量，Pod模板，标签选择器



还有ReplicatSet

![image-20211004144741292](/home/yy/.config/Typora/typora-user-images/image-20211004144741292.png)

![image-20211004144759389](/home/yy/.config/Typora/typora-user-images/image-20211004144759389.png)

## 水平伸缩和滚动更新

# 网络基础

主要有四种网络通信：

- 同一Pod内容器的通信

- 不同Pod之间通信
- Pod与Service通信
- 集群外部的流量与Service间的通信

Pod的IP由插件负责配置，拥有一个独立的网络命名空间。

Service是一个虚拟IP地址，不会添加在任何网络接口设备上。而是由kube-proxy配置在每个工作节点的iptables或者ipvs规则中用于转发，而且仅限于当前节点自身。

## Kubernetes命名空间

使用命名空间来限制资源对象名称的作用域。



# HorizontalAutoScaler

![image-20211005155120565](/home/yy/.config/Typora/typora-user-images/image-20211005155120565.png)

```bash
kubectl autoscale deployment chat-redis --cpu-percent=50 --min=1 --max=3
```



![image-20211005200027901](/home/yy/.config/Typora/typora-user-images/image-20211005200027901.png)

![image-20211005200036898](/home/yy/.config/Typora/typora-user-images/image-20211005200036898.png)

我们可以简单的通过 `kubectl autoscale` 命令来创建一个 HPA 资源对象，`HPA Controller`默认`30s`轮询一次（可通过 `kube-controller-manager` 的`--horizontal-pod-autoscaler-sync-period` 参数进行设置），查询指定的资源中的 Pod 资源使用率，并且与创建时设定的值和指标做对比，从而实现自动伸缩的功能。
