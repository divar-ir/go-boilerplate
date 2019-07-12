package multilayercache

import (
	"context"
	"git.cafebazaar.ir/arcana261/golang-boilerplate/pkg/cache"
	"github.com/pkg/errors"
)

type multilayerCache struct {
	layers []cache.Layer
}

func New(layers ...cache.Layer) cache.Layer {
	return &multilayerCache{
		layers: layers,
	}

}

func (mc *multilayerCache) Get(ctx context.Context, key string) ([]byte, error) {
	if mc == nil {
		return nil, errors.New("cache client is nil")
	}

	var allErrors error

	for _, cacheLayer := range mc.layers {
		value, err := cacheLayer.Get(ctx, key)
		if err == nil {
			return value, nil
		}
		if allErrors == nil {
			allErrors = err
		} else {
			allErrors = errors.Wrap(err, allErrors.Error())
		}
	}

	return nil, errors.Wrap(allErrors, "fail to get cached valued from any cache layers")
}

func (mc *multilayerCache) Delete(ctx context.Context, key string) error {
	if mc == nil {
		return errors.New("multilayer cache is nill")
	}

	for _, cacheLayer := range mc.layers {
		err := cacheLayer.Delete(ctx, key)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (mc *multilayerCache) Set(ctx context.Context, key string, value []byte) error {
	if mc == nil {
		return errors.New("cache client is nil")
	}

	var allErrors error

	for _, cacheLayer := range mc.layers {
		err := cacheLayer.Set(ctx, key, value)
		if err != nil {
			if allErrors == nil {
				allErrors = err
			} else {
				allErrors = errors.Wrap(err, allErrors.Error())
			}
		}
	}

	return allErrors
}

func (mc *multilayerCache) Clear(ctx context.Context) error {
	if mc == nil {
		return errors.New("multilayer cache is nil")
	}
	var allErrors error

	for _, cacheLayer := range mc.layers {
		err := cacheLayer.Clear(ctx)
		if err != nil {
			if allErrors == nil {
				allErrors = err
			} else {
				allErrors = errors.Wrap(err, allErrors.Error())
			}
		}
	}

	return allErrors
}
