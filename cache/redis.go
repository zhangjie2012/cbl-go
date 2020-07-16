package cache

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

var (
	NotExist error = fmt.Errorf("key not exist")
)

var (
	// for compose key
	appName       string = "not_set"
	disLockModule string = "_dislock_"
	mqModule      string = "_mq_"

	redisClient *redis.Client = nil
)

// InitCache init cache, only init once
func InitCache(app string, addr string, password string, db int) error {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if _, err := client.Ping().Result(); err != nil {
		return err
	}

	// init
	appName = app
	redisClient = client

	logrus.Infof("init cache success|%s", addr)

	return nil
}

func CloseCache() {
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			logrus.Errorf("close client failure, err=%s", err)
		}
		redisClient = nil
	}

	logrus.Infof("close cache")
}

// C expose redis client for native redis library visit
func C() *redis.Client {
	return redisClient
}

func composeKey(source string) string {
	return fmt.Sprintf("%s:%s", appName, source)
}

func composeKey2(module string, key string) string {
	return fmt.Sprintf("%s:%s.%s", appName, module, key)
}

// SetObject set object, object must be json marshaled
func SetObject(key string, value interface{}, expire time.Duration) error {
	realKey := composeKey(key)

	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}

	val, err := redisClient.Set(realKey, bs, expire).Result()
	if err != nil {
		return err
	}

	logrus.Tracef("cache set object|%s|%s", realKey, val)

	return nil
}

// GetObject get object, object must be json unmarshaled
func GetObject(key string, value interface{}) error {
	realKey := composeKey(key)

	bs, err := redisClient.Get(realKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			logrus.Tracef("cache missing|%s", realKey)
			return NotExist
		}
		return err
	}

	if err := json.Unmarshal(bs, value); err != nil {
		return err
	}

	logrus.Tracef("cache get object|%s", realKey)

	return nil
}

// TTL seconds resolution
//   The command returns -1 if the key exists but has no associated expire.
//   The command returns -2 if the key does not exist.
func TTL(key string) time.Duration {
	realKey := composeKey(key)
	d, err := redisClient.TTL(realKey).Result()
	if err != nil {
		return 0
	}
	return d
}

func SetString(key string, value string, expire time.Duration) error {
	realKey := composeKey(key)

	val, err := redisClient.Set(realKey, []byte(value), expire).Result()
	if err != nil {
		return err
	}

	logrus.Tracef("cache set string|%s|%s", realKey, val)

	return nil
}

func GetString(key string) (string, error) {
	realKey := composeKey(key)

	bs, err := redisClient.Get(realKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			logrus.Tracef("cache missing|%s", realKey)
			return "", NotExist
		}
		return "", err
	}

	logrus.Tracef("cache get string|%s", realKey)

	return string(bs), nil
}

func SetInt(key string, value int, expire time.Duration) error {
	return SetString(key, strconv.Itoa(value), expire)
}

func GetInt(key string) (int, error) {
	realKey := composeKey(key)

	value, err := redisClient.Get(realKey).Int()
	if err != nil {
		if err == redis.Nil {
			logrus.Tracef("cache missing, key=%s", realKey)
			return 0, NotExist
		}
		return 0, err
	}

	logrus.Tracef("cache get int, key=%s", realKey)

	return value, nil
}

func SetFloat64(key string, value float64, expire time.Duration) error {
	return SetString(key, strconv.FormatFloat(value, 'f', -1, 64), expire)
}

func GetFloat64(key string) (float64, error) {
	realKey := composeKey(key)

	value, err := redisClient.Get(realKey).Float64()
	if err != nil {
		if err == redis.Nil {
			logrus.Tracef("cache missing|%s", realKey)
			return 0, NotExist
		}
		return 0, err
	}

	logrus.Tracef("cache get float64|%s", realKey)

	return value, nil
}

// -----------------------------------------------------------------------------
// distribute lock
//   - name: lock key
//   - ticket: lock unique flag, avoid anther process unlock, make sure only one
//           process lock, then unlock it
//   - expire: lock timeout, avoid process dead forget unlock it
// Note: not consider redis server down caused deadlock
// -----------------------------------------------------------------------------
func Lock(name string, ticket string, expire time.Duration) bool {
	lockKey := composeKey2(disLockModule, name)
	result := redisClient.SetNX(lockKey, ticket, expire).Val()
	// logrus.Tracef("distributed lock|%s|%s|%v|%t", name, ticket, expire, result)
	return result
}

func UnLock(name string, ticket string) {
	lockKey := composeKey2(disLockModule, name)
	v, err := redisClient.Get(lockKey).Result()
	if err != nil {
		return
	}

	// just can unlock itself
	if v == ticket {
		redisClient.Del(lockKey)
		// logrus.Tracef("distribute unlock success|%s|%s", name, ticket)
	} else {
		logrus.Tracef("distribute unlock failue|%s|%s|%s", name, v, ticket)
	}
}

// -----------------------------------------------------------------------------
// message queue
//   - redis structure list map to a message queue
//   - right push, left pop
// -----------------------------------------------------------------------------

func MQPush(key string, bs []byte) error {
	mqKey := composeKey2(mqModule, key)
	_, err := redisClient.RPush(mqKey, bs).Result()
	return err
}

func MQPop(key string) ([]byte, error) {
	mqKey := composeKey2(mqModule, key)
	bs, err := redisClient.LPop(mqKey).Bytes()
	if err == redis.Nil {
		return nil, NotExist
	}
	return bs, err
}

// MQBlockPop block pop, in comparison, block pop fast than polling pop
func MQBlockPop(key string, timeout time.Duration) ([]byte, error) {
	// timeout min value is 1s
	if timeout.Seconds() < 1 {
		timeout = time.Second
	}
	mqKey := composeKey2(mqModule, key)
	result, err := redisClient.BLPop(timeout, mqKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil {
		return nil, NotExist
	}

	// result[] is key name
	return []byte(result[1]), nil
}

func MQLen(key string) int64 {
	mqKey := composeKey2(mqModule, key)
	count, err := redisClient.LLen(mqKey).Result()
	if err != nil {
		logrus.Tracef("mqlen error|%s", err)
		return 0
	}
	return count
}

// MQDel delete mq return count
func MQDel(key string) int64 {
	mqKey := composeKey2(mqModule, key)
	count, err := redisClient.Del(mqKey).Result()
	if err != nil {
		logrus.Tracef("mqdel error|%s", err)
		return 0
	}
	return count
}
