# mping
[executable] multicast test for multicast testing ASM & SSM sending and receiving

[可执行] 组播测试工具（指定源组播和任意源组播shou'fa）

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
下方给出了go编译的示例

```bash
# generally
// windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/mping.exe main.go
// linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/mping main.go
// linux arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o build/mping-arm main.go
 
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

## golang库未完善

☞ On Windows, the ReadBatch and WriteBatch methods of PacketConn are not implemented.

☞ On Windows, the ReadBatch and WriteBatch methods of RawConn are not implemented.

☞ This package is not implemented on JS, NaCl and Plan 9.

☞ On Windows, the JoinSourceSpecificGroup, LeaveSourceSpecificGroup, ExcludeSourceSpecificGroup and IncludeSourceSpecificGroup methods of PacketConn and RawConn are not implemented.

☞ On Windows, the ReadFrom and WriteTo methods of RawConn are not implemented.

☞ On Windows, the ControlMessage for ReadFrom and WriteTo methods of PacketConn is not implemented.

详情见https://godoc.org/golang.org/x/net/ipv4#NewPacketConn
