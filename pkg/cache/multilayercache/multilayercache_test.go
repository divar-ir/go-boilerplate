package multilayercache

import (
	"context"
	"github.com/cafebazaar/go-boilerplate/pkg/cache/adaptors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetGet(t *testing.T) {

	var err error

	multiLayerCache := New(adaptors.NewSynMapAdaptor(logrus.New()))

	key := "test-key"
	value := "test-value"
	err = multiLayerCache.Set(context.Background(), key, []byte(value))
	assert.NoError(t, err, "fail to set data")

	gottenValue, err := multiLayerCache.Get(context.Background(), key)
	assert.NoError(t, err, "fail to get data")

	assert.Equal(t, value, string(gottenValue), "gotten value is not equal to set value")
}
