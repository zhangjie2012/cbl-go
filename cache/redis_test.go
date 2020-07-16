package cache

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhangjie2012/cbl-go"
)

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.TraceLevel)

	var (
		name      string = "cblcache"
		redisAddr string = "localhost:6379"
		password  string = ""
		db        int    = 0
	)

	if err := InitCache(name, redisAddr, password, db); err != nil {
		fmt.Printf("redis init failure, err=%s", err)
		return
	}

	ec := m.Run()

	CloseCache()

	os.Exit(ec)
}

func TestSetGetObject(t *testing.T) {
	type ValueT struct {
		Username string
		Phonenum string
		Age      int
		Deposit  float64
	}

	var (
		key   = "TestSetGetObject"
		value = ValueT{
			Username: "张三",
			Phonenum: "136****1234",
			Age:      45,
			Deposit:  10000000.89,
		}
	)

	if err := SetObject(key, &value, 10*time.Millisecond); err != nil {
		t.Errorf("set object failure, err=%s", err)
		return
	}

	gValue := ValueT{}
	if err := GetObject(key, &gValue); err != nil {
		t.Errorf("get object failure, err=%s", err)
		return
	}

	if value.Username != gValue.Username ||
		value.Phonenum != gValue.Phonenum ||
		value.Age != gValue.Age ||
		value.Deposit != gValue.Deposit {
		t.Errorf("get set not equal, expect=%v, actual=%v", value, gValue)
	}
}

func TestTTL(t *testing.T) {
	var (
		key = "TestTLL"
		ttl = 1234 * time.Millisecond
	)
	SetString(key, "whatever", ttl)
	{
		real := TTL(key)
		t.Logf("set ttl = %d, get ttl = %d", ttl.Milliseconds(), real.Milliseconds())
	}
	{
		real := PTTL(key)
		t.Logf("set ttl = %d, get ttl = %d", ttl.Milliseconds(), real.Milliseconds())
	}
}

func TestDel(t *testing.T) {
	var (
		key = "TestDel"
	)
	SetString(key, "whatever", 0)
	Del(key)
	t.Logf(GetString(key))
}

func TestSetGetString(t *testing.T) {
	var (
		key   = "TestSetGetString"
		value = "hello"
	)

	if err := SetString(key, value, 10*time.Millisecond); err != nil {
		t.Errorf("set string failure, err=%s", err)
		return
	}

	gValue, err := GetString(key)
	if err != nil {
		t.Errorf("get string failure, err=%s", err)
		return
	}
	if value != gValue {
		t.Errorf("get set not equal, expect=%s, real=%s", value, gValue)
		return
	}
}

func TestSetGetInt(t *testing.T) {
	var (
		key   = "TestSetGetString"
		value = 1234
	)

	if err := SetInt(key, value, 10*time.Millisecond); err != nil {
		t.Errorf("set int failure, err=%s", err)
		return
	}

	gValue, err := GetInt(key)
	if err != nil {
		t.Errorf("get string failure, err=%s", err)
		return
	}
	if value != gValue {
		t.Errorf("get set not equal, expect=%d, real=%d", value, gValue)
		return
	}
}

func TestSetGetFloat64(t *testing.T) {
	var (
		key           = "TestSetGetString"
		value float64 = 3.1415926
	)

	if err := SetFloat64(key, value, 10*time.Millisecond); err != nil {
		t.Errorf("set int failure, err=%s", err)
		return
	}

	gValue, err := GetFloat64(key)
	if err != nil {
		t.Errorf("get string failure, err=%s", err)
		return
	}
	if value != gValue {
		t.Errorf("get set not equal, expect=%f, real=%f", value, gValue)
		return
	}
}

func TestDisLock0(t *testing.T) {
	var (
		key    = "awesomelock0"
		ticket = cbl.GenRSessionID()
	)
	r := Lock(key, ticket, 10*time.Second)
	if !r {
		t.Logf("lock failure")
	}
	UnLock(key, ticket)
}

func TestDisLock1(t *testing.T) {
	var (
		key    = "awesomelock1"
		ticket = cbl.GenRSessionID()
	)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		r := Lock(key, ticket, 10*time.Second)
		if !r {
			t.Logf("lock failure")
		}
		t.Logf("thread1 lock|%s|%s", key, ticket)
		time.Sleep(100 * time.Millisecond)
		UnLock(key, ticket)
		t.Logf("thread1 unlock|%s|%s", key, ticket)
	}()

	time.Sleep(5 * time.Millisecond)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			r := Lock(key, ticket, 10*time.Second)
			if !r {
				t.Logf("not get lock, waiting")
				time.Sleep(10 * time.Millisecond)
				continue
			}
			t.Logf("thread2 get lock")
			UnLock(key, ticket)
			t.Logf("thread2 get unlock")
			break
		}
	}()

	wg.Wait()
}

func TestMQ(t *testing.T) {
	key := "mqtest"
	value := "awesome_values"
	if err := MQPush(key, []byte(value)); err != nil {
		t.Log("push message failure")
		return
	}

	bs, err := MQPop(key)
	if err != nil {
		t.Logf("pop message failure|%s", err)
	}
	t.Log(string(bs))

	bs, err = MQPop(key)
	if err != nil {
		t.Logf("pop message failure|%s", err)
	}
	t.Log(string(bs))
}

// producer and consumer
func TestMQ1(t *testing.T) {
	key := "mqtest1"
	count := 1000
	wg := sync.WaitGroup{}

	t.Logf("%s mq len = %d", key, MQDel(key))

	// producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			MQPush(key, []byte(fmt.Sprintf("message_%d", i)))
			time.Sleep(1 * time.Microsecond)
		}
	}()

	// consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < count; {
			_, err := MQPop(key)
			if err != nil && err != NotExist {
				t.Log(err.Error())
				break
			}
			if err == NotExist {
				t.Log("wait producer ...")
				time.Sleep(1 * time.Microsecond)
				continue
			}
			i++
		}
	}()

	wg.Wait()
}

// producer and consumer
func TestMQ2(t *testing.T) {
	key := "mqtest2"
	count := 1000
	wg := sync.WaitGroup{}

	t.Logf("%s mq len = %d", key, MQDel(key))

	// producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			MQPush(key, []byte(fmt.Sprintf("message_%d", i)))
			time.Sleep(1 * time.Microsecond)
		}
	}()

	// consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < count; {
			_, err := MQBlockPop(key, 1*time.Microsecond)
			if err != nil && err != NotExist {
				t.Log(err.Error())
				break
			}
			if err == nil {
				i++
			}
		}
	}()

	wg.Wait()
}

func TestMQ3(t *testing.T) {
	key := "mqtest3"
	count := 1000
	wg := sync.WaitGroup{}

	t.Logf("%s mq len = %d", key, MQDel(key))

	// producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			MQPush(key, []byte(fmt.Sprintf("message_%d", i)))
			time.Sleep(1 * time.Microsecond)
		}
	}()

	// consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < count; {
			if MQLen(key) != 0 {
				MQPop(key)
				i++
			} else {
				t.Log("empty queue")
				time.Sleep(1 * time.Microsecond)
			}
		}
	}()

	wg.Wait()
}
