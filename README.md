# mping
[executable] multicast_test_go for multicast testing ASM &amp; SSM sending and receiving

# 使用说明
源代码运行
    go run main.go
编译
    go build main.go
跨平台、系统编译请参照build文件
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

## 编译
需要
    import "golang.org/x/net/ipv4"

可以自行GitHub golang/net查找

## golang库未完善

☞ On Windows, the ReadBatch and WriteBatch methods of PacketConn are not implemented.

☞ On Windows, the ReadBatch and WriteBatch methods of RawConn are not implemented.

☞ This package is not implemented on JS, NaCl and Plan 9.

☞ On Windows, the JoinSourceSpecificGroup, LeaveSourceSpecificGroup, ExcludeSourceSpecificGroup and IncludeSourceSpecificGroup methods of PacketConn and RawConn are not implemented.

☞ On Windows, the ReadFrom and WriteTo methods of RawConn are not implemented.

☞ On Windows, the ControlMessage for ReadFrom and WriteTo methods of PacketConn is not implemented.

详情见https://godoc.org/golang.org/x/net/ipv4#NewPacketConn
