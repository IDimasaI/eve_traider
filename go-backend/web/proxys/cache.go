package proxy

import (
	"sync"
	"time"
)

type InMemoryCache struct {
	Ttl   time.Time
	Value any
}
type StorageType int

const (
	InMemory StorageType = iota

//	FileSystem
//
// Database
)

type Proxy struct {
	mu          sync.RWMutex
	Cache       map[string]*InMemoryCache
	TypeStorage StorageType
}

func NewProxy(typeStorage StorageType) *Proxy {
	return &Proxy{
		Cache:       make(map[string]*InMemoryCache),
		TypeStorage: typeStorage,
	}
}

func (proxy *Proxy) Get(key string) (any, bool) {
	if proxy == nil {
		return nil, false
	}

	proxy.mu.RLock()
	defer proxy.mu.RUnlock()

	value, exists := proxy.Cache[key]
	if !exists {
		return nil, false
	}

	// Проверка TTL (если требуется)
	if value.Ttl.Before(time.Now()) {
		delete(proxy.Cache, key)
		return nil, false
	}

	return value.Value, true
}

func (proxy *Proxy) Set(key string, value any, ttl time.Duration) {
	if proxy == nil {
		return
	}

	proxy.mu.Lock()
	defer proxy.mu.Unlock()

	proxy.Cache[key] = &InMemoryCache{
		Ttl:   time.Now().Add(ttl),
		Value: value,
	}

}

func (proxy *Proxy) Delete(key string) {
	if proxy == nil {
		return
	}

	proxy.mu.Lock()
	defer proxy.mu.Unlock()

	delete(proxy.Cache, key)
}

func (proxy *Proxy) GetAll() map[string]*InMemoryCache {
	if proxy == nil {
		return nil
	}

	proxy.mu.RLock()
	defer proxy.mu.RUnlock()

	return proxy.Cache
}

// Дополнительная типобезопасная функция для удобства
func GetAs[T any](proxy *Proxy, key string) (T, bool) {
	var zero T
	value, ok := proxy.Get(key)
	if !ok {
		return zero, false
	}

	v, ok := value.(T)
	if !ok {
		return zero, false
	}

	return v, true
}

func (proxy *Proxy) Vacuum() {
	if proxy == nil {
		return
	}

	proxy.mu.Lock()
	defer proxy.mu.Unlock()

	for key, value := range proxy.Cache {
		if value.Ttl.Before(time.Now()) {
			delete(proxy.Cache, key)
		}
	}
}
