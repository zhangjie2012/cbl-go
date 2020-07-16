package cbl

import (
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/zhangjie2012/cbl-go/cache"
)

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.TraceLevel)

	var (
		name      string = "cblcache"
		redisAddr string = "localhost:6379"
		password  string = ""
		db        int    = 0
	)

	if err := cache.InitCache(name, redisAddr, password, db); err != nil {
		fmt.Printf("redis init failure, err=%s", err)
		return
	}

	ec := m.Run()

	cache.CloseCache()

	os.Exit(ec)
}

func TestToString(t *testing.T) {
	y, _ := ConvYangYin(2020, 6, 25)
	t.Log(y.ToString1())
	t.Log(y.ToString2())
	t.Log(y.ToString3())
}

func TestConvYinYang(t *testing.T) {
	y, err := ConvYinYang(2020, 5, 0, 5)
	t.Log(y, err) // 2020-06-25 <nil>
}

func TestConvYangYin(t *testing.T) {
	y, err := ConvYangYin(2020, 6, 26)
	t.Log(y.ToString1(), err)
}
