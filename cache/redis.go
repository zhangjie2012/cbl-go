package cache

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	redis "github.com/go-redis/redis/v7"
)

var (
	NotExist                = fmt.Errorf("key not exist")
	CounterZero             = fmt.Errorf("counter zero")
	ErrUnLockTicketNotMatch = fmt.Errorf("unlock ticket not match")
)

var (
	// for compose key
	appName       string = "not_set"
	disLockModule string = "_dislock_"
	mqModule      string = "_mq_"
	counterModule string = "_counter_"
	setModule     string = "_set_"

	once        sync.Once
	redisClient *redis.Client = nil
)

// InitCache init cache, only init once
func InitCache(app string, addr string, password string, db int) error {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})
		appName = app
		redisClient = client
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		return err
	}

	return nil
}

func CloseCache() error {
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			return err
		}
		redisClient = nil
	}

	return nil
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

// ----------------------------------------------------------------------------
// common built-in type wrapper
// ----------------------------------------------------------------------------

// SetObject set object, object must be json marshaled
func SetObject(key string, value interface{}, expire time.Duration) error {
	realKey := composeKey(key)

	bs, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = redisClient.Set(realKey, bs, expire).Result()
	if err != nil {
		return err
	}

	return nil
}

// GetObject get object, object must be json unmarshaled
func GetObject(key string, value interface{}) error {
	realKey := composeKey(key)

	bs, err := redisClient.Get(realKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return NotExist
		}
		return err
	}

	if err := json.Unmarshal(bs, value); err != nil {
		return err
	}

	return nil
}

// TTL seconds resolution
// - The command returns -1 if the key exists but has no associated expire.
// - The command returns -2 if the key does not exist.
func TTL(key string) time.Duration {
	realKey := composeKey(key)
	d, err := redisClient.TTL(realKey).Result()
	if err != nil {
		return 0
	}
	return d
}

// PTTL milliseconds resolution
// - The command returns -1 if the key exists but has no associated expire.
// - The command returns -2 if the key does not exist.
func PTTL(key string) time.Duration {
	realKey := composeKey(key)
	d, err := redisClient.PTTL(realKey).Result()
	if err != nil {
		return 0
	}
	return d
}

func Del(key string) error {
	realKey := composeKey(key)
	_, err := redisClient.Del(realKey).Result()
	return err
}

func SetString(key string, value string, expire time.Duration) error {
	realKey := composeKey(key)

	_, err := redisClient.Set(realKey, []byte(value), expire).Result()
	if err != nil {
		return err
	}

	return nil
}

func GetString(key string) (string, error) {
	realKey := composeKey(key)

	bs, err := redisClient.Get(realKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return "", NotExist
		}
		return "", err
	}
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
			return 0, NotExist
		}
		return 0, err
	}

	return value, nil
}

func SetInt64(key string, value int64, expire time.Duration) error {
	return SetString(key, strconv.FormatInt(value, 10), expire)
}

func GetInt64(key string) (int64, error) {
	realKey := composeKey(key)

	value, err := redisClient.Get(realKey).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, NotExist
		}
		return 0, err
	}
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
			return 0, NotExist
		}
		return 0, err
	}

	return value, nil
}

func SetBool(key string, b bool, expire time.Duration) error {
	if b {
		return SetInt(key, 1, expire)
	} else {
		return SetInt(key, 0, expire)
	}
}

func GetBool(key string) (bool, error) {
	value, err := GetInt(key)
	if err != nil {
		return false, err
	}
	if value == 1 {
		return true, err
	} else {
		return false, err
	}
}

// -----------------------------------------------------------------------------
// distributed lock
// -----------------------------------------------------------------------------

// TryLock if lock failure, max wait "timeout" duration (retry lock)
func TryLock(name string, ticket string, expire time.Duration, timeout time.Duration) bool {
	t := time.NewTimer(timeout)
	for {
		select {
		case <-t.C:
			return false
		default:
			result := Lock(name, ticket, expire)
			if result {
				return true
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func Lock(name string, ticket string, expire time.Duration) bool {
	lockKey := composeKey2(disLockModule, name)
	result := redisClient.SetNX(lockKey, ticket, expire).Val()
	return result
}

func UnLock(name string, ticket string) error {
	lockKey := composeKey2(disLockModule, name)
	v, err := redisClient.Get(lockKey).Result()
	if err != nil {
		return err
	}

	// just can unlock itself
	if v == ticket {
		_, err := redisClient.Del(lockKey).Result()
		return err
	} else {
		return ErrUnLockTicketNotMatch
	}
}

// -----------------------------------------------------------------------------
// message queue
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

	// result[0] is key name
	return []byte(result[1]), nil
}

func MQLen(key string) int64 {
	mqKey := composeKey2(mqModule, key)
	count, err := redisClient.LLen(mqKey).Result()
	if err != nil {
		return 0
	}
	return count
}

// MQDel delete mq return count, mq key can not use `Del` delete, they have different compose method
func MQDel(key string) int64 {
	mqKey := composeKey2(mqModule, key)
	count, err := redisClient.Del(mqKey).Result()
	if err != nil {
		return 0
	}
	return count
}

// -----------------------------------------------------------------------------
// atomic counter
// -----------------------------------------------------------------------------

// CounterIncr atomic increment 1, return inc result value
func CounterIncr(key string, expire time.Duration) (int64, error) {
	aKey := composeKey2(counterModule, key)

	pipe := redisClient.TxPipeline()
	incr := pipe.Incr(aKey)
	pipe.Expire(aKey, expire)
	_, err := pipe.Exec()

	return incr.Val(), err
}

// CounterIncrBy atomic increment n, return incrby result value
func CounterIncrBy(key string, n int64, expire time.Duration) (int64, error) {
	aKey := composeKey2(counterModule, key)

	pipe := redisClient.TxPipeline()
	incr := pipe.IncrBy(aKey, n)
	pipe.Expire(aKey, expire)
	_, err := pipe.Exec()

	return incr.Val(), err
}

// CounterDecr atomic decrement 1, return decr result value
func CounterDecr(key string) (int64, error) {
	aKey := composeKey2(counterModule, key)
	return redisClient.Decr(aKey).Result()
}

// CounterDecrMinZero atomic decrement, min value is 0
func CounterDecrMinZero(key string) (int64, error) {
	aKey := composeKey2(counterModule, key)

	script := `
local v = redis.call("GET", KEYS[1])
if v == false then
   return -2
end

if tonumber(v) > 0 then
   return redis.call("DECR", KEYS[1])
else
   return -1
end
	`
	result, err := redisClient.Eval(script, []string{aKey}).Int64()
	if err != nil {
		return 0, err
	}
	if result == -2 {
		return 0, NotExist
	}
	if result == -1 {
		return 0, CounterZero
	}
	return result, nil
}

// CounterDecrBy atomic decrement n, return decr result value
func CounterDecrBy(key string, n int64) (int64, error) {
	aKey := composeKey2(counterModule, key)
	return redisClient.DecrBy(aKey, n).Result()
}

// CounterReset reset counter to 0
func CounterReset(key string, expire time.Duration) error {
	aKey := composeKey2(counterModule, key)
	_, err := redisClient.Set(aKey, "0", expire).Result()
	return err
}

// CounterDel delete counter
func CounterDel(key string) {
	aKey := composeKey2(counterModule, key)
	redisClient.Del(aKey)
}

// CounterGet get counter value
func CounterGet(key string) (int64, error) {
	aKey := composeKey2(counterModule, key)
	value, err := redisClient.Get(aKey).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, NotExist
		}
		return 0, err
	}
	return value, nil
}

// -----------------------------------------------------------------------------
// Set wrapper
// SS for Set String
// -----------------------------------------------------------------------------

// SSMembers get all members slice
func SSMembers(key string) ([]string, error) {
	aKey := composeKey2(setModule, key)
	values, err := redisClient.SMembers(aKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, NotExist
		}
		return nil, err
	}
	return values, nil
}

// SSAdd add members to Set
func SSAdd(key string, members ...string) error {
	aKey := composeKey2(setModule, key)
	t := []interface{}{}
	for _, v := range members {
		t = append(t, v)
	}
	_, err := redisClient.SAdd(aKey, t...).Result()
	return err
}

// SSRem remove members from Set
func SSRem(key string, members ...string) error {
	aKey := composeKey2(setModule, key)
	t := []interface{}{}
	for _, v := range members {
		t = append(t, v)
	}
	_, err := redisClient.SRem(aKey, t...).Result()
	return err
}

// SSCount get member count
func SSCount(key string) int64 {
	aKey := composeKey2(setModule, key)
	count, err := redisClient.SCard(aKey).Result()
	if err != nil {
		return 0
	}
	return count
}

// SSIsMember check set if include member
func SSIsMember(key string, member string) bool {
	aKey := composeKey2(setModule, key)
	ok, err := redisClient.SIsMember(aKey, member).Result()
	if err != nil {
		return false
	}
	return ok
}

// SSRandomN random get N members
func SSRandomN(key string, count int64) []string {
	aKey := composeKey2(setModule, key)
	values, err := redisClient.SRandMemberN(aKey, count).Result()
	if err != nil {
		return []string{}
	}
	return values
}

func SSDelete(key string) {
	aKey := composeKey2(setModule, key)
	redisClient.Del(aKey)
}

// SS_TTL seconds resolution
// - The command returns -1 if the key exists but has no associated expire.
// - The command returns -2 if the key does not exist.
func SS_TTL(key string) time.Duration {
	aKey := composeKey2(setModule, key)
	d, err := redisClient.TTL(aKey).Result()
	if err != nil {
		return 0
	}
	return d
}

// SS_TTL milliseconds resolution
// - The command returns -1 if the key exists but has no associated expire.
// - The command returns -2 if the key does not exist.
func SS_PTTL(key string) time.Duration {
	aKey := composeKey2(setModule, key)
	d, err := redisClient.PTTL(aKey).Result()
	if err != nil {
		return 0
	}
	return d
}

func SSExpire(key string, d time.Duration) error {
	aKey := composeKey2(setModule, key)
	return redisClient.Expire(aKey, d).Err()
}
