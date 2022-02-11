package cache

import "container/list"

type lru struct {
	maxBytes  int64
	usedBytes int64
	l         *list.List
	cache     map[string]*list.Element
	// optional and executed when an entry is purged
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

func (l *lru) Get(key string) (Value, bool) {
	if e, ok := l.cache[key]; ok {
		l.l.MoveToBack(e)
		return e.Value.(*entry).value, true
	}
	return nil, false
}

func (l *lru) Add(key string, value Value) {
	if e, ok := l.cache[key]; ok {
		exist := e.Value.(*entry)
		exist.value = value
		l.l.MoveToBack(e)
		l.usedBytes += int64(value.Len()) - int64(exist.value.Len())
	} else {
		e := l.l.PushBack(&entry{
			key:   key,
			value: value,
		})
		l.cache[key] = e
		l.usedBytes += int64(len(key)) + int64(value.Len())
	}
	for l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.removeOldest()
	}
}

func (l *lru) Len() int {
	return l.l.Len()
}

func (l *lru) removeOldest() {
	e := l.l.Front()
	if e != nil {
		l.l.Remove(e)
		v := e.Value.(*entry)
		delete(l.cache, v.key)
		l.usedBytes -= int64(len(v.key)) + int64(v.value.Len())
		if l.OnEvicted != nil {
			l.OnEvicted(v.key, v.value)
		}
	}
}

func NewLru(maxBytes int64, onEvicted func(string, Value)) Cache {
	return &lru{
		maxBytes:  maxBytes,
		l:         list.New(),
		cache:     map[string]*list.Element{},
		OnEvicted: onEvicted,
	}
}
