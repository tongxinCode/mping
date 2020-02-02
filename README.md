# mping
[executable] multicast_test_go for multicast testing ASM &amp; SSM sending and receiving

# 使用说明
二进制文件和EXE文件选择相应的文件进行运行
运行命令
    ./programe -h
查看帮助

# 核心代码
    main.go
    ./multicast
    --broadcaster.go
    --listener.go
## 编译说明
给出了go编译的示例，在build文件中
建议使用gox

# 注意事项

需要
    import "golang.org/x/net/ipv4"

可以自行GitHub golang/net查找
