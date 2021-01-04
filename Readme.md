
## QuickScan

买了白帽子安全开发这本书来康康,学习golang

这个端口扫描器是[sec-dev-in-action-src](https://github.com/netxfly/sec-dev-in-action-src/blob/main/scanner/)的一个简单修改版,我自己学习用的,就是加了一个扫描结果输出的功能.

## usage

```
go run main.go Scan --host 172.20.10.2,172.20.10.3/16 --port 22 --rate 100  -o 1.txt
```
syn scan

```
go run main.go Scan --host 172.20.10.2,172.20.10.3/16 --port 22 --rate 100  -o 1.txt -m syn -t 10
```
对比了一下另外一个工具,应该都是tcp发包判断,这个虽然有syn扫描功能,会快一点,但是准确率不太行,常规的web服务端口没问题.但是redis,socks这些识别不准.


