# docker 网络

​	Docker的本地网络实现其实就是利用了Linux上的网络命名空间和虚拟网络设备（特别是veth pair）。

​	Docker中的网络接口默认都是虚拟的接口。

​	虚拟接口的最大优势就是**转发效率极高**。

​	这是因为Linux通过在内核中进行数据复制来实现虚拟接口之间的数据转发，即发送接口的发送缓存中的数据包将被直接复制到接收接口的接收缓存中，而无需通过外部物理网络设备进行交换。对于本地系统和容器内系统来看，虚拟接口跟一个正常的以太网卡相比并无区别，只是它速度要快得多。 

​	Docker容器网络就很好地利用了Linux虚拟网络技术，在本地主机和容器内分别创建一个虚拟接口，并让它们彼此连通（这样的一对接口叫做veth pair）。

# 网卡

https://docs.docker.com/network/bridge/

bridge networking 

​	With bridged networking, Oracle VM VirtualBox uses a device driver on  your *host* system that **filters data from your  physical network adapter**. 

​	This driver is therefore called a   *net filter* driver. This enables      Oracle VM VirtualBox to intercept data from the physical network and      inject data into it, effectively creating a new network interface      in software. When a guest is using such a new software interface,      it looks to the host system as though the guest were physically      connected to the interface using a network cable. The host can      send data to the guest through that interface and receive data      from it. This means that you can set up routing or bridging      between the guest and the rest of your network.    