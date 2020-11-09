![Go](https://github.com/zhangjie2012/cbl-go/workflows/Go/badge.svg)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/zhangjie2012/cbl-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/zhangjie2012/cbl-go)](https://goreportcard.com/report/github.com/zhangjie2012/cbl-go)

# cbl-go: Common Basic Library for Go

Install:

``` shell
go get -u github.com/zhangjie2012/cbl-go
```

## Libraries include

- Common
  + **generator**: uuid, session id, verfiy code
  + **date expand**
  + **set operations**: uniq, union, intersection, difference
  + **yinyang**: 公历、公里转换 _deprecated_
- Data Structure
  + **MultiError**: muliple error expression
  + **Pair**: key value struct
- **cache**: a redis wrapper, easy to use, include common operator, mq, dislock, set
- **datasize** byte to KB/MB/GB/TB/PB/EB, KiB/MiB/GiB/TiB/PiB/EiB and more elegent to string
