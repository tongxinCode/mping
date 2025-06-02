# mping

组播测试工具，用于测试任意源组播（ASM）和指定源组播（SSM）的发送和接收功能。

## 程序架构

```
mping/
├── main.go              # 主程序入口
└── multicast/           # 组播核心实现
    ├── broadcaster.go   # 组播发送实现
    └── listener.go      # 组播接收实现
```

## 使用说明

1. 跨平台编译请参考下方编译说明
2. 根据您的操作系统选择对应的二进制文件运行
3. 运行以下命令查看帮助信息：

```bash
go run main.go -h    
# 或    
./mping -h
```

## 编译说明

支持多种平台的编译方式，您也可以使用 `build.sh` 脚本进行编译：

```bash
# 通用编译方式
# Windows 64位
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o mping.exe main.go

# Linux 64位
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mping main.go

# Linux ARM
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o mping main.go

# Windows 环境变量方式
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o mping.exe main.go

# Linux 环境变量方式
SET CGO_ENABLED=0
SET GOOS=linux 
SET GOARCH=amd64 
go build -o mping main.go

# Linux ARM 环境变量方式
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm
go build -o mping main.go
```

## 版本历史

### v1.8.3
- 优化 Lua 接口使用方式，提高数据包生成效率
- 新增无日志模式选项

### v1.8.2
- 使用 Lua 接口热插拔功能生成随机 UDP 数据

### v1.8.1
- 适配 Go 1.20
- 更新依赖，消除安全告警

### v1.8.0
- 新增 Lua 接口热插拔功能，用于解析 UDP 数据协议

### v1.7.0
- 完善数据包发送限制功能
- 修复计数和内容（消息）冲突问题
- 增加内容（消息）编码选项

### v1.6.0
- 重要：修复多个程序缺陷

### v1.5.0
- 完善实时发送流程
- 新增丢包率统计功能
- 修复部分程序缺陷

### v1.4.0
- 适配 Go 1.18
- 增加 Go modules 支持

## 平台限制说明

由于 golang 标准库的限制，以下功能在某些平台上可能不可用：

- Windows 平台：
  - PacketConn 的 ReadBatch 和 WriteBatch 方法未实现
  - RawConn 的 ReadBatch 和 WriteBatch 方法未实现
  - PacketConn 和 RawConn 的以下方法未实现：
    - JoinSourceSpecificGroup
    - LeaveSourceSpecificGroup
    - ExcludeSourceSpecificGroup
    - IncludeSourceSpecificGroup
  - RawConn 的 ReadFrom 和 WriteTo 方法未实现
  - PacketConn 的 ReadFrom 和 WriteTo 方法的 ControlMessage 未实现

- 其他限制：
  - JS、NaCl 和 Plan 9 平台不支持本程序
  - 更多详情请参考 [golang.org/x/net/ipv4#NewPacketConn](https://godoc.org/golang.org/x/net/ipv4#NewPacketConn)
