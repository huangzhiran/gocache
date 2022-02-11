package cache_test

import (
	"reflect"
	"testing"

	"github.com/huangzhiran/gocache/cache"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestLruGet(t *testing.T) {
	c := cache.NewLru(10, nil)
	c.Add("key1", String("1111"))
	v, ok := c.Get("key1")
	if !ok || string(v.(String)) != "1111" {
		t.Fatalf("cache hit key1=1111 failed")
	}
	if _, ok := c.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestAutoEvict(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "v1", "v2", "v3"
	c := cache.NewLru(int64(len(k1+k2+v1+v2)), nil)
	c.Add(k1, String(v1))
	c.Add(k2, String(v2))
	c.Add(k3, String(v3))
	if _, ok := c.Get(k1); ok || c.Len() != 2 {
		t.Fatalf("auto evict failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, v cache.Value) {
		keys = append(keys, key)
	}
	c := cache.NewLru(10, callback)
	c.Add("key1", String("123456"))
	c.Add("k2", String("k2"))
	c.Add("k3", String("k3"))
	c.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}
	if !reflect.DeepEqual(keys, expect) {
		t.Fatalf("callback failed, get %v, expect %v", keys, expect)
	}
}
