package config

func New[T any](options ...func(*T)) *T {
	config := new(T)
	for _, o := range options {
		o(config)
	}
	return config
}
