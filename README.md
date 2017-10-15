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

========================================================================

## 启动Server 

./mserver -p 10000 -n 100

从端口10000~10099 创建100个server 然后查看服务的HTOP

## 启动Client

mclient.exe  -ip 127.0.0.1 -p 10000 -n 100  -cn 10000

端口10000开始 连接100个端口  每个端口创建10000个连接（一共创建100万连接）




