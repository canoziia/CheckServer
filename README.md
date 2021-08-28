# CheckServer

检测服务器是否在线

## 使用方法

```
$ cp config.example.ini config.ini
$ go run main.go
```

## 配置文件示例
```ini
[common]
cycle=5
report=1800
log=logs.txt

[mail]
host=smtp.office365.com
port=587
encrypt=SSL
name=Canoziia
username=me@example.com
password=pswd
target=target@example.com

[server1]
name=example
host=example.com
reexectime=300
port=80
mode=tcp
timeout=10
```