# mping
[executable] multicast test for multicast testing ASM & SSM sending and receiving

[可执行] 组播测试工具（指定源组播和任意源组播首发）

# 使用说明
跨平台、系统编译请参照编译说明

二进制文件和EXE文件选择相应的文件进行运行

运行以下命令查看帮助
```bash
go run main.go -h    
# or    
./programe -h
```

# 核心代码
    main.go
    ./multicast
    --broadcaster.go
    --listener.go

## 编译说明
下方给出了go编译的示例，你也可以使用build.sh编译构建

```bash
# generally
// windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o mping.exe main.go
// linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mping main.go
// linux arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o mping main.go
 
# or
# windows
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o mping.exe main.go
# linux
SET CGO_ENABLED=0
SET GOOS=linux 
SET GOARCH=amd64 
go build -o mping main.go
# linux arm
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm
go build -o mping main.go

```

# 注意事项

## 版本介绍
v1.4.0

将程序适配go 1.18，增加go mod适配

v1.5.0

1. complete the realtime send process
2. add the Packet loss rate count function
3. fix some bugs

v1.6.0

**important** fix bugs

v1.7.0

1. complete the limit of sending packets
2. fix the conflict of count and content(message)
3. add the encoding choices of content(message)

v1.8.0
**important** add the lua interface hot-plugin to parse the udp data protocol

## golang库未完善

☞ On Windows, the ReadBatch and WriteBatch methods of PacketConn are not implemented.

☞ On Windows, the ReadBatch and WriteBatch methods of RawConn are not implemented.

☞ This package is not implemented on JS, NaCl and Plan 9.

☞ On Windows, the JoinSourceSpecificGroup, LeaveSourceSpecificGroup, ExcludeSourceSpecificGroup and IncludeSourceSpecificGroup methods of PacketConn and RawConn are not implemented.

☞ On Windows, the ReadFrom and WriteTo methods of RawConn are not implemented.

☞ On Windows, the ControlMessage for ReadFrom and WriteTo methods of PacketConn is not implemented.

详情见https://godoc.org/golang.org/x/net/ipv4#NewPacketConn
