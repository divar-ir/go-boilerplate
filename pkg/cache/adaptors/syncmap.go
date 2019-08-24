package adaptors

import (
	"context"
	"github.com/cafebazaar/go-boilerplate/pkg/cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

type syncMap struct {
	data   sync.Map
	logger *logrus.Logger
}

func NewSynMapAdaptor(logger *logrus.Logger) cache.Layer {
	logger.Info("new cache initialized")
	return &syncMap{
		data:   sync.Map{},
		logger: logger,
	}
}

func (cache *syncMap) Get(_ context.Context, key string) ([]byte, error) {
	value, ok := cache.data.Load(key)
	if !ok {
		cache.logger.Infof("key %s not found", key)
		return nil, errors.New("not found")
	}
	cache.logger.Infof("get key %s successfully", key)
	return value.([]byte), nil
}
func (cache *syncMap) Delete(_ context.Context, key string) error {
	cache.logger.Infof("key %s is deleted", key)
	cache.data.Delete(key)
	return nil
}

func (cache *syncMap) Set(_ context.Context, key string, value []byte) error {
	cache.logger.Infof("key %s is updated", key)
	cache.data.Store(key, value)
	return nil
}

func (cache *syncMap) Clear(_ context.Context) error {
	cache.logger.Infof("clearing")
	cache.data = sync.Map{}
	return nil
}
