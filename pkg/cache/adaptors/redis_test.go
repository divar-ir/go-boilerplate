package adaptors

//import (
//	"context"
//	"github.com/go-redis/redis"
//	"github.com/stretchr/testify/assert"
//	"os"
//	"testing"
//	"time"
//)
//
//func TestRedisSetGet(t *testing.T) {
//	redisAddr, ok := os.LookupEnv("MULTILAYERCACHE_REDIS_ADDR")
//	if !ok {
//		redisAddr = "127.0.0.1:6379"
//	}
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: redisAddr,
//		DB:   3,
//	})
//	// Ping Redis
//	err := redisClient.Ping().Err()
//	if !assert.NoError(t, err) {
//		t.Fatal("fail to connect to redis")
//	}
//	redisAdaptor := NewRedisAdaptor(1*time.Minute, redisClient)
//
//	key := "test-key"
//	value := "test-value"
//	err = redisAdaptor.Set(context.Background(), key, []byte(value))
//	assert.NoError(t, err, "fail to set data")
//
//	gottenValue, err := redisAdaptor.Get(context.Background(), key)
//	assert.NoError(t, err, "fail to get data")
//
//	assert.Equal(t, value, string(gottenValue), "gotten value is not equal to set value")
//}
