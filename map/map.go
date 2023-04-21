package internal

import "sync"

// Map 是对 sync.Map 的一个泛型封装
// 要注意，K 必须是 comparable 的，并且谨慎使用指针作为 K。
// 使用指针的情况下，两个 key 是否相等，仅仅取决于它们的地址
// 注意，key 不存在和 key 存在但是值恰好为零值（如 nil），是两码事
type Map[K comparable, V any] struct {
	m sync.Map
}

// Load 加载键值对
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	var anyVal any
	anyVal, ok = m.m.Load(key)
	if anyVal != nil {
		value = anyVal.(V)
	}
	return
}

// Store 存储键值对
func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

// LoadOrStore 加载或者存储一个键值对
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	var anyVal any
	anyVal, loaded = m.m.LoadOrStore(key, value)
	if anyVal != nil {
		actual = anyVal.(V)
	}
	return
}

// LoadAndDelete 加载并且删除一个键值对
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	var anyVal any
	anyVal, loaded = m.m.LoadAndDelete(key)
	if anyVal != nil {
		value = anyVal.(V)
	}
	return
}

// Delete 删除键值对
func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// Range 遍历, f 不能为 nil
// 传入 f 的时候，K 和 V 直接使用对应的类型，如果 f 返回 false，那么就会中断遍历
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool {
		var (
			k K
			v V
		)
		if value != nil {
			v = value.(V)
		}
		if key != nil {
			k = key.(K)
		}
		return f(k, v)
	})
}

// SyncMap 是对 map 的一个泛型封装
// 要注意，K 必须是 comparable 的，并且谨慎使用指针作为 K。
// 使用指针的情况下，两个 key 是否相等，仅仅取决于它们的地址
// 注意，key 不存在和 key 存在但是值恰好为零值（如 nil），是两码事
type SyncMap[K comparable, V any] struct {
	lock sync.RWMutex
	Map  map[K]V
}

// Load 加载键值对
func (s *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	v, ok := s.Map[key]
	return v, ok
}

// Store 存储键值对
func (s *SyncMap[K, V]) Store(key K, value V) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Map[key] = value
}

// LoadOrStore 加载或者存储一个键值对
func (s *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	s.lock.RLock()
	actual, ok := s.Map[key]
	s.lock.RUnlock()
	if ok {
		return actual, true
	}

	s.lock.Lock()
	s.Map[key] = value
	s.lock.Unlock()
	return value, false
}

// LoadAndDelete 加载并且删除一个键值对
func (s *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	var val V
	s.lock.RLock()
	val, ok := s.Map[key]
	if !ok {
		return val, false
	}
	s.lock.RUnlock()

	s.lock.Lock()
	defer s.lock.Unlock()
	value, _ = s.Map[key]
	delete(s.Map, key)
	return value, true
}

// Delete 删除键值对
func (s *SyncMap[K, V]) Delete(key K) {
	s.lock.RLock()
	if _, ok := s.Map[key]; !ok {
		return
	}
	s.lock.RUnlock()

	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.Map, key)
}
