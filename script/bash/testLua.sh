# send
go run main.go -l 172.17.249.65 -s 229.0.0.1:9999 -p 0 -c
# receive
go run main.go -l 172.17.249.65 -r 229.0.0.1:9999 -proto script/lua/test.lua