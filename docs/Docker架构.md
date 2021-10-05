# Docker架构

![image-20210929165958410](/home/yy/.config/Typora/typora-user-images/image-20210929165958410.png)

Registry：docker hub

Repository：一个Registry可以包含多个Repository

Container：与Image就是对象和类的关系

Image

Tag：每个Repository可以包含多个Tag，每个Tag对应一个image

可以通过 <Repository>:<Tag>来指定是哪个镜像

# Docker程序构成

- dockerd: 就是damon，根据镜像运行容器。

> A self-sufficient runtime for containers.
>
> ​	`dockerd` is the persistent process that manages containers. Docker uses different binaries for the daemon and client. To run the daemon you type `dockerd`.

- docker-compose CLI:

> ```
> webapp:
>   image: examples/web
>   ports:
>     - "8000:8000"
>   volumes:
>     - "/data"
> ```



- docker engine API: 客户端与服务端进行交互的API

> ​	Docker provides an API for interacting with the Docker daemon (called the Docker Engine API), as well as SDKs for Go and Python. 
>
> ​	The Docker daemon can listen for [Docker Engine API](https://docs.docker.com/engine/api/) requests via three different types of Socket: `unix`, `tcp`, and `fd`.

![image-20210929180644066](/home/yy/.config/Typora/typora-user-images/image-20210929180644066.png)

![image-20210929180653465](/home/yy/.config/Typora/typora-user-images/image-20210929180653465.png)

![image-20210929180712811](/home/yy/.config/Typora/typora-user-images/image-20210929180712811.png)





- Docker Engine SDKs: 就是提供给编程语言的，编程语言可以进行使用调用API

    -H, --host list          Daemon socket(s) to connect to

![image-20210929174217381](/home/yy/.config/Typora/typora-user-images/image-20210929174217381.png)



主要感觉你对这个Dokcer 根本不熟悉。客户端，Damon。

​	客户端，就是一个应用程序，我们都可以使用。因此自己写一个程序，里面也可以操纵API。



- containerd: 就是在操作系统和dameon中又加了一层。适用于不同的操作系统平台。相当于设计模式里面的把不变的代码抽出来。针对接口编程。decouples Docker from the OS 解耦

> ​	Containerd was designed to be used by Docker and Kubernetes as well as  any other container platform that wants to abstract away syscalls or OS  specific functionality to run containers on linux, windows, solaris, or  other OSes.  



# container容器是如何实现的

容器平台都有什么特点，都有哪些特点在里面。

namespace 实现资源隔离

cgroup 实现资源限制

写时复制 高效的文件操作

## namespace

命名空间（namespaces）是 Linux 为我们提供的用于分离进程树、网络接口、挂载点以及进程间通信等资源的方法。

Linux 的命名空间机制提供了以下七种不同的命名空间，包括  CLONE_NEWCGROUP、CLONE_NEWIPC、CLONE_NEWNET、CLONE_NEWNS、CLONE_NEWPID、CLONE_NEWUSER 和 CLONE_NEWUTS，通过这七个选项我们能在创建新的进程时设置新进程应该在哪些资源上与宿主机器进行隔离。

​	在新的容器内部执行 ps 命令打印出了非常干净的进程列表，只有包含当前 ps -ef 在内的三个进程，在宿主机器上的几十个进程都已经消失不见了。

​	当前的 Docker 容器成功将容器内的进程与宿主机器中的进程隔离，如果我们在宿主机器上打印当前的全部进程时，会得到下面三条与 Docker 相关的结果：

​	Docker 的容器就是使用上述技术实现与宿主机器的进程隔离，当我们每次运行 docker run 或者 docker start 时，都会在下面的方法中创建一个用于设置进程间隔离的 Spec



## 实验

CLONE_NEWUTS

```C
#define _GNU_SOURCE
#include <sys/types.h>
#include <sys/wait.h>
#include <stdio.h>
#include <sched.h>
#include <signal.h>
#include <unistd.h>
#define STACK_SIZE (1024 * 1024)

static char container_stack[STACK_SIZE];
char* const container_args[] = {
   "/bin/bash",
   NULL
};

// 容器进程运行的程序主函数
int container_main(void *args)
{
   printf("在容器进程中！\n");
   execv(container_args[0], container_args); // 执行/bin/bash   return 1;
}

int main(int args, char *argv[])
{
   printf("程序开始\n");
   // clone 容器进程
   int container_pid = clone(container_main, container_stack + STACK_SIZE, SIGCHLD | CLONE_NEWUTS, NULL);
   // 等待容器进程结束
   waitpid(container_pid, NULL, 0);
   return 0;
}
```

## CGroups

​	但是 namespaces 并不能够为我们提供物理资源上的隔离，比如 CPU、内存、IO  或者网络带宽等，所以如果我们运行多个容器的话，则容器之间就会抢占资源互相影响了，所以对容器资源的使用进行限制就非常重要了，而 Control  Groups（CGroups）技术就能够隔离宿主机上的物理资源。CGroups 由 7 个主要的子系统组成：分别是  cpuset、cpu、cpuacct、blkio、devices、freezer、memory，不同类型资源的分配和管理是由各个 CGroup  子系统负责完成的。

## 文件系统UnionFS

容器和镜像的区别：

容器会在镜像顶层创建一个可写层。所以多个容器可以共享一个镜像。



# docker的用处

​	由于之前我们的后台在开发和运维阶段的环境是不一致的，这就导致了 Docker 的出现，因为我们通过 Docker 可以将程序运行的环境也一起打包到版本控制去了，这样就排除了因为环境不同造成的各种麻烦事情了，也不会出现在本地可以在线上却不行这样的窘境了

# 命令

docker build

> ​	The `docker build` command builds Docker images from a Dockerfile and a “context”. 
>
> ​	A build’s context is the set of files located in the specified `PATH` or `URL`. 
>
> ​	The build process can refer to any of the files in the context. For example, your build can use a [*COPY*](https://docs.docker.com/engine/reference/builder/#copy) instruction to reference a file in the context.

docker run

-d: 守护进程

-it : 交互式的进程

