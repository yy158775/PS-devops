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

# docker文件系统





# container容器是如何实现的

容器平台都有什么特点，都有哪些特点在里面。

## namespace

命名空间（namespaces）是 Linux 为我们提供的用于分离进程树、网络接口、挂载点以及进程间通信等资源的方法。

Linux 的命名空间机制提供了以下七种不同的命名空间，包括  CLONE_NEWCGROUP、CLONE_NEWIPC、CLONE_NEWNET、CLONE_NEWNS、CLONE_NEWPID、CLONE_NEWUSER 和 CLONE_NEWUTS，通过这七个选项我们能在创建新的进程时设置新进程应该在哪些资源上与宿主机器进行隔离。

​	在新的容器内部执行 ps 命令打印出了非常干净的进程列表，只有包含当前 ps -ef 在内的三个进程，在宿主机器上的几十个进程都已经消失不见了。

​	当前的 Docker 容器成功将容器内的进程与宿主机器中的进程隔离，如果我们在宿主机器上打印当前的全部进程时，会得到下面三条与 Docker 相关的结果：

​	Docker 的容器就是使用上述技术实现与宿主机器的进程隔离，当我们每次运行 docker run 或者 docker start 时，都会在下面的方法中创建一个用于设置进程间隔离的 Spec：

# docker的用处

编写



# 为什么要用docker

​	Docker 的出现一定是因为目前的后端在开发和运维阶段确实需要一种虚拟化技术解决开发环境和生产环境环境一致的问题，通过 Docker 我们可以将程序运行的环境也纳入到版本控制中，排除因为环境造成不同运行结果的可能。

我在本地写代码，可以跑起来，但是到了服务器上可能跑步起来？？



大家需要注意，**Docker本身并不是容器**，它是创建容器的工具，是应用容器引擎。



​	通过 `docker run`  命令指定一个容器创建镜像时，实际上是在该镜像之上创建一个空的可读写的文件系统层级，可以将这个文件系统层级当成一个临时的镜像来对待，而命令中所指的模版镜像则可以称之为父镜像。

​	父镜像的内容都是以只读的方式挂载进来的，容器会读取共享父镜像的内容，用户所做的所有修改都是在文件系统中，不会对父镜像造成任何影响。当然用户可以通过其他一些手段使修改持久化到父镜像中，这个我们后面会详细介绍到。

​	简而言之，镜像就是一个固定的不会变化的模版文件，容器是根据这个模版创建出来的，容器会在模版的基础上做一些修改，这些修改本身并不会影响到模版，我们还可以根据模版（镜像）创建出来更多的容器。

如果有必要，我们是可以修改模版（镜像）的。

# docker build

> ​	The `docker build` command builds Docker images from a Dockerfile and a “context”. 
>
> ​	A build’s context is the set of files located in the specified `PATH` or `URL`. 
>
> ​	The build process can refer to any of the files in the context. For example, your build can use a [*COPY*](https://docs.docker.com/engine/reference/builder/#copy) instruction to reference a file in the context.



# docker 网络

​	docker容器的四种网络模式：bridge 桥接模式、host 模式、container 模式和 none 模式 

