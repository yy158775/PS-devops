# docker原理略

​	Docker是一个开源的软件项目,让用户程序部署在一个相对隔离的环境运行，借此在Linux操作系统上提供一层额外的抽象，以及操作系统层虚拟化的自动管理机制。

​	需要额外指出的是，Docker并不等于容器（containers），Docker只是容器的一种，其他的种类的容器还有Kata container，Rocket container等等。

​	一句话概括起来Docker就是利用Namespace做资源隔离，用Cgroup做资源限制，利用Union  FS做容器文件系统的轻量级虚拟化技术。

​	Docker容器的本质还是一个直接运行在宿主机上面的特殊进程，看到的文件系统是隔离后的，但是操作系统内核是共享宿主机OS，所以说Docker是轻量级的虚拟化技术。

![image-20210929203113009](/home/yy/.config/Typora/typora-user-images/image-20210929203113009.png)

```
int pid = clone(call_function, stack_size, CLONE_NEWPID | SIGCHLD, NULL);
```

![image-20210929203838075](/home/yy/.config/Typora/typora-user-images/image-20210929203838075.png)

​	说白了，我们用kubernetes去管理Docker集群，即可以将Docker看成Kubernetes内部使用的低级别组件。另外，kubernetes不仅仅支持Docker，还支持Rocket，这是另一种容器技术。希望我这篇文章中简单的描述能让你对两者有所理解和认识。



# k8s

一个K8S系统，通常称为一个**K8S集群（Cluster）**。



这个集群主要包括两个部分：

- **一个Master节点（主节点）**
- **一群Node节点（计算节点）**

![image-20210929204640042](/home/yy/.config/Typora/typora-user-images/image-20210929204640042.png)

虚拟化，容器化