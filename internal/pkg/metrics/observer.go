package metrics

type Observer interface {
	Observe(float64)
	With(labels map[string]string) Observer
}
