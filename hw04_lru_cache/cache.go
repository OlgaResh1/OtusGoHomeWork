package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.Mutex
}

type cacheElement struct {
	key   Key
	value interface{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if item, ok := cache.items[key]; ok {
		cache.queue.MoveToFront(item)
		item.Value = cacheElement{key: key, value: value}
		return true
	}

	if cache.queue.Len() == cache.capacity {

		itemBack := cache.queue.Back()

		delete(cache.items, itemBack.Value.(cacheElement).key)
		cache.queue.Remove(itemBack)
	}

	item := cache.queue.PushFront(cacheElement{key: key, value: value})
	cache.items[key] = item

	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if item, ok := cache.items[key]; ok {
		cache.queue.MoveToFront(item)
		return item.Value.(cacheElement).value, true
	}

	return nil, false
}

func (cache *lruCache) Clear() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.queue.Clear()
	cache.items = make(map[Key]*ListItem, 0)
	cache.capacity = 0
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
