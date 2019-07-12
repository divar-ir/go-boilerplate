package adaptors

import (
	"context"
	"github.com/allegro/bigcache"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBigCacheSetGet(t *testing.T) {
	bigCacheInstance, err := bigcache.NewBigCache(bigcache.Config{
		Shards:             1,
		LifeWindow:         1 * time.Minute,
		MaxEntriesInWindow: 1100,
		MaxEntrySize:       1000,
		Verbose:            true,
		HardMaxCacheSize:   100,
	})
	if !assert.NoError(t, err) {
		t.Fatal("fail to initialize big cache")
	}

	bigCacheAdaptor := NewBigCacheAdaptor(bigCacheInstance)

	if !assert.NoError(t, err) {
		t.Fatal("fail to initialize adaptor")
	}

	key := "test-key"
	value := "test-value"

	err = bigCacheAdaptor.Set(context.Background(), key, []byte(value))
	assert.NoError(t, err, "fail to set data")

	gottenValue, err := bigCacheAdaptor.Get(context.Background(), key)
	assert.NoError(t, err, "fail to get data")
	assert.Equal(t, value, string(gottenValue), "gotten value is not equal to set value")
}
