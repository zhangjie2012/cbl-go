package cache

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
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
	defer CloseCache()

	os.Exit(m.Run())
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
			Username: "Bob",
			Phonenum: "136****1234",
			Age:      45,
			Deposit:  10000000.89,
		}
	)

	err := SetObject(key, &value, 10*time.Millisecond)
	assert.Nil(t, err)

	gValue := ValueT{}
	err = GetObject(key, &gValue)
	assert.Nil(t, err)

	assert.EqualValues(t, value, gValue)
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
		err error
		key = "TestDel"
	)
	err = SetString(key, "whatever", 0)
	assert.Nil(t, err)

	err = Del(key)
	assert.Nil(t, err)

	v, err := GetString(key)
	assert.Equal(t, "", v)
	assert.Equal(t, NotExist, err)
}

func TestSetGetInt64(t *testing.T) {
	var (
		key         = "TestSetGetInt64"
		value int64 = 123456789
		err   error
	)

	err = SetInt64(key, value, 10*time.Millisecond)
	require.Nil(t, err)

	gValue, err := GetInt64(key)
	require.Nil(t, err)

	assert.Equal(t, gValue, value)
}

func TestSetGetString(t *testing.T) {
	var (
		err   error
		key   = "TestSetGetString"
		value = "hello"
	)

	err = SetString(key, value, 10*time.Millisecond)
	require.Nil(t, err)

	gValue, err := GetString(key)
	require.Nil(t, err)

	assert.Equal(t, gValue, value)
}

func TestSetGetInt(t *testing.T) {
	var (
		err   error
		key   = "TestSetGetString"
		value = 1234
	)

	err = SetInt(key, value, 10*time.Millisecond)
	require.Nil(t, err)

	gValue, err := GetInt(key)
	require.Nil(t, err)

	assert.Equal(t, value, gValue)
}

func TestSetGetFloat64(t *testing.T) {
	var (
		err   error
		key           = "TestSetGetString"
		value float64 = 3.1415926
	)

	err = SetFloat64(key, value, 10*time.Millisecond)
	require.Nil(t, err)

	gValue, err := GetFloat64(key)
	require.Nil(t, err)

	assert.Equal(t, value, gValue)
}

func TestDisLock0(t *testing.T) {
	var (
		key    = "awesomelock0"
		ticket = "ticket_0"
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
		ticket = "ticket_1"
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
	var (
		err   error
		key   = "mqtest"
		value = "awesome_values"
	)

	err = MQPush(key, []byte(value))
	require.Nil(t, err)

	bs, err := MQPop(key)
	require.Nil(t, err)
	assert.EqualValues(t, value, string(bs))

	bs, err = MQPop(key)
	require.Equal(t, NotExist, err)
	assert.EqualValues(t, "", string(bs))
}

// producer and consumer
func TestMQ1(t *testing.T) {
	var (
		err   error
		key   = "mqtest1"
		count = 1000
		wg    = sync.WaitGroup{}
		in    = []string{}
		out   = []string{}
	)

	MQDel(key)

	// producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			d := fmt.Sprintf("message_%d", i)
			err = MQPush(key, []byte(d))
			require.Nil(t, err)
			in = append(in, d)
		}
	}()

	// consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < count; {
			o, err := MQPop(key)
			if err != nil && err != NotExist {
				t.Log(err.Error())
				break
			}
			if err == NotExist {
				// t.Log("wait producer ...")
				time.Sleep(1 * time.Microsecond)
				continue
			}
			out = append(out, string(o))
			i++
		}
	}()

	wg.Wait()

	assert.EqualValues(t, in, out)
}

// producer and consumer
func TestMQ2(t *testing.T) {
	key := "mqtest2"
	count := 1000
	wg := sync.WaitGroup{}

	MQDel(key)

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

	MQDel(key)

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

func TestCounter(t *testing.T) {
	var (
		err    error
		key    = "global.counter"
		expire = 1 * time.Second
	)

	v, err := CounterIncr(key, expire)
	require.Nil(t, err)
	assert.EqualValues(t, 1, v)

	v, err = CounterIncrBy(key, 999, expire)
	require.Nil(t, err)
	assert.EqualValues(t, 1000, v)

	v, err = CounterDecr(key)
	require.Nil(t, err)
	assert.EqualValues(t, 999, v)

	v, err = CounterGet(key)
	require.Nil(t, err)
	assert.EqualValues(t, 999, v)

	v, err = CounterDecrBy(key, 99)
	require.Nil(t, err)
	assert.EqualValues(t, 900, v)

	err = CounterReset(key, expire)
	require.Nil(t, err)

	v, err = CounterIncr(key, expire)
	require.Nil(t, err)
	assert.EqualValues(t, 1, v)

	v, err = CounterDecrMinZero(key)
	require.Nil(t, err)

	v, err = CounterDecrMinZero(key)
	require.Equal(t, CounterZero, err)
	assert.EqualValues(t, 0, v)

	v, err = CounterDecrMinZero(key)
	require.Equal(t, CounterZero, err)
	assert.EqualValues(t, 0, v)

	key = "global.notexist.key"
	v, err = CounterDecrMinZero(key)
	require.Equal(t, NotExist, err)
	assert.EqualValues(t, 0, v)

	CounterDel(key)
}

func TestLua(t *testing.T) {
	key := "global.counter"
	CounterDecrMinZero(key)
}

// TestSetS test string set
func TestSetS(t *testing.T) {
	var (
		err    error
		key    = "TestSetS"
		values = []string{}
		count  int64
		ok     bool
	)

	values, err = SSMembers(key)
	require.Nil(t, err)
	assert.Equal(t, 0, len(values))

	err = SSAdd(key, "hello1", "hello2", "hello3")
	require.Nil(t, err)

	values, err = SSMembers(key)
	require.Nil(t, err)
	assert.EqualValues(t, 3, len(values))

	count = SSCount(key)
	assert.EqualValues(t, 3, count)

	ok = SSIsMember(key, "hello2")
	assert.Equal(t, true, ok)

	ok = SSIsMember(key, "hello4")
	assert.Equal(t, false, ok)

	err = SSRem(key, "hello2")
	assert.Nil(t, err)

	count = SSCount(key)
	assert.EqualValues(t, 2, count)

	values = SSRandomN(key, 2)
	assert.Equal(t, 2, len(values))

	d := SS_TTL(key)
	assert.EqualValues(t, -1, d)

	// err = SSExpire(key, 10*time.Second)
	// assert.Nil(t, err)

	// d = SS_TTL(key)
	// assert.EqualValues(t, -1, d)

	SSDelete(key)
}
