![Go](https://github.com/zhangjie2012/cbl-go/workflows/Go/badge.svg)

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/zhangjie2012/cbl-go)


# cbl-go: Common Basic Library for go

__目前只自己使用，不保证稳定性。__

模块：

- `cache`
  + 基于 redis 的 KV 缓存
  + 基于 redis 的分布式锁
  + 基于 redis 的 message queue
  + 基于 redis 的全局计数器（counter）
- `conv` 类型转换
- `errcode` web 常用错误码
- `gen` token/id 生成器
- `ginext` gin 扩展（标准化）
- `net` 网络扩展
- `prom` prometheus middware
- `wechat_mini` 微信小程序
- `yinyang` 公历和农历转换

安装：

``` shell
go get -u github.com/zhangjie2012/cbl-go
```

运行测试用例：

``` shell
go test -count=1 -v ./...
go test -count=1 -v cache/* -test.run TestMQ
```
