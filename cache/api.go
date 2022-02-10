package cache

type Cache interface {
	Get(key string) (v Value, ok bool)
	Add(key string, v Value)
	Len() int
}

type Value interface {
	Len() int
}
