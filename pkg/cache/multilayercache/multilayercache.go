package multilayercache

import (
	"context"
	"sync"

	"github.com/cafebazaar/go-boilerplate/pkg/cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type multilayerCache struct {
	layers []cache.Layer
}

type layerOperator func(layer cache.Layer) (interface{}, error)

func New(layers ...cache.Layer) cache.Layer {
	return &multilayerCache{
		layers: layers,
	}

}

func (mc *multilayerCache) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := mc.performOperation(func(layer cache.Layer) (interface{}, error) {
		return layer.Get(ctx, key)
	})

	if err != nil {
		return nil, err
	}
	return result.([]byte), nil
}

func (mc *multilayerCache) Delete(ctx context.Context, key string) error {
	_, err := mc.performOperation(func(layer cache.Layer) (interface{}, error) {
		err := layer.Delete(ctx, key)
		return nil, err
	})

	return err
}

func (mc *multilayerCache) Set(ctx context.Context, key string, value []byte) error {
	_, err := mc.performOperation(func(layer cache.Layer) (interface{}, error) {
		err := layer.Set(ctx, key, value)
		return nil, err
	})

	return err
}

func (mc *multilayerCache) Clear(ctx context.Context) error {
	_, err := mc.performOperation(func(layer cache.Layer) (interface{}, error) {
		err := layer.Clear(ctx)
		return nil, err
	})

	return err
}

func (mc *multilayerCache) wrapAllErrors(errChannel <-chan error) error {
	var allErrors error
	for err := range errChannel {
		if allErrors == nil {
			allErrors = err
		} else {
			allErrors = errors.Wrap(err, allErrors.Error())
		}
	}

	return allErrors
}

func (mc *multilayerCache) performOperation(operator layerOperator) (interface{}, error) {
	errChannel := make(chan error, len(mc.layers))
	resultChannel := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(len(mc.layers))

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
		close(errChannel)
		close(resultChannel)
	}()

	for _, cacheLayer := range mc.layers {
		go func(layer cache.Layer) {
			defer wg.Done()

			value, err := operator(layer)
			if err != nil {
				errChannel <- err
				return
			}

			select {
			case resultChannel <- value:
			default:
			}
		}(cacheLayer)
	}

	select {
	case value := <-resultChannel:
		go func() {
			for err := range errChannel {
				logrus.WithError(err).Error("failed to get value from multilayer cache")
			}
		}()

		return value, nil

	case <-done:
		return nil, mc.wrapAllErrors(errChannel)
	}
}
