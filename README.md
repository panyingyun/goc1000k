# A C1000k demo

## windows 7下启动Server
server.exe

打印输出：

Connect [1000] [4000] [4000] [127.0.0.1:54686]

注：括号中分别是 当前客户端连接数、收包数量、发包数量、当前活跃的Client IP

## windows 7下启动Client
client.exe

Connect [127.0.0.1:54656]

注：括号中是 当前活跃的Client IP

## 统计

netstat -an | grep -i "9999" | wc -l 

1000

