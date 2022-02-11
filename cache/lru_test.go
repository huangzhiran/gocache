package cache_test

import (
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
