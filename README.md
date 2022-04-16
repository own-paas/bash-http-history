bash-http-history 
============
## 关于
将bash输入的历史命令写入http服务器的2种方法:

第一种: 直接给bash源码打补丁

第二种: 通过ebpf uprobe获取bash readline输入

## bash path用法:
~~~
下载bash rpm源码包
https://mirrors.aliyun.com/centos-vault/7.9.2009/os/Source/SPackages/bash-4.2.46-34.el7.src.rpm
https://mirrors.aliyun.com/centos-vault/centos/8-stream/BaseOS/Source/SPackages/bash-4.4.20-3.el8.src.rpm
~~~

~~~
安装源码rpm包
# rpm -ivh bash-4.4.20-3.el8.src.rpm

拷贝补丁到~/rpmbuild/SOURCES/bash-http-history.patch
# cp patch/bash-4.4.20-3.el8_patch/bash-http-history.patch  ~/rpmbuild/SOURCES/bash-http-history.patch

修改~/rpmbuild/SPECS/bash.spec加入补丁并重新编译
# grep bash-http-history  bash.spec 
Patch158: bash-http-history.patch
# rpmbuild -ba bash.spec

安装编译好的rpm包
# rpm -Uvh bash-4.4.20-4.el8.x86_64.rpm
~~~

~~~
运行测试http服务器,重新打开一个新的bash窗口输入命令观察输出
# cd example
# go run server.go
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] PUT    /logger/:type             --> main.setupRouter.func1 (3 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :6666
<nil> {  bash_history 192.168.1.100 32350 21772 10620 0 root /dev/pts/2 /root/rpmbuild/BUILD/bash-4.2 ifconfog 1650124616}
[GIN] 2022/04/16 - 23:56:56 | 200 |     123.885µs |       127.0.0.1 | PUT      "/logger/bash_history"
<nil> {  bash_history 192.168.1.100 32350 21772 10620 0 root /dev/pts/2 /root/rpmbuild/BUILD/bash-4.2 ip a 1650124626}
[GIN] 2022/04/16 - 23:57:06 | 200 |      86.157µs |       127.0.0.1 | PUT      "/logger/bash_history"
<nil> {  bash_history 192.168.1.100 32350 21772 10620 0 root /dev/pts/2 /root/rpmbuild/BUILD/bash-4.2 ls 1650124627}
[GIN] 2022/04/16 - 23:57:07 | 200 |      49.104µs |       127.0.0.1 | PUT      "/logger/bash_history"
<nil> {  bash_history 192.168.1.100 32350 21772 10620 0 root /dev/pts/2 /root/rpmbuild/BUILD/bash-4.2 cd .. 1650124636}
[GIN] 2022/04/16 - 23:57:16 | 200 |      68.693µs |       127.0.0.1 | PUT      "/logger/bash_history"
<nil> {  bash_history 192.168.1.100 32350 21772 10620 0 root /dev/pts/2 /root/rpmbuild/BUILD ls 1650124636}
[GIN] 2022/04/16 - 23:57:16 | 200 |      65.726µs |       127.0.0.1 | PUT      "/logger/bash_history"
<nil> {  bash_history 192.168.1.100 32350 21772 10620 0 root /dev/pts/2 /root/rpmbuild/BUILD cd .. 1650124638}
[GIN] 2022/04/16 - 23:57:18 | 200 |      65.836µs |       127.0.0.1 | PUT      "/logger/bash_history"
<nil> {  bash_history 192.168.1.100 32350 21772 10620 0 root /dev/pts/2 /root/rpmbuild ls 1650124638}
[GIN] 2022/04/16 - 23:57:18 | 200 |     114.602µs |       127.0.0.1 | PUT      "/logger/bash_history"
<nil> {  bash_history 192.168.1.100 32350 21772 10620 0 root /dev/pts/2 /root/rpmbuild exit 1650124720}
[GIN] 2022/04/16 - 23:58:40 | 200 |      69.229µs |       127.0.0.1 | PUT      "/logger/bash_history"
~~~


## ebpf uprobe用法:
~~~
运行ebpf uprobe程序,重新打开一个新的bash窗口输入命令观察输出
# cd ebpf_uprobe
# go run bash_readline.go
            HOATNAME	       TTY	    CLIENT	       PID	      PPID	       UID	  USERNAME	       PWD	   COMMAND	        TS
    ebpf.example.com	/dev/pts/1	192.168.1.100	    134138	    134137	         0	      root	     /root	      ip a	1650127258
    ebpf.example.com	/dev/pts/1	192.168.1.100	    134138	    134137	         0	      root	     /root	   whoami 	1650127262
~~~
## License

* Apache License Version 2.0

