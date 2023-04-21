package internal

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

/*
goos: darwin
goarch: arm64
pkg: orm/internal
BenchmarkTestMap_Load
BenchmarkTestMap_Load-8   	 6870882	       154.9 ns/op	      48 B/op	       3 allocs/op
PASS
*/
func BenchmarkTestSyncMap_LoadOrStore(b *testing.B) {
	b.ReportAllocs()
	var syncMap = SyncMap[string, string]{
		Map:  map[string]string{},
		lock: sync.RWMutex{},
	}

	b.RunParallel(func(pb *testing.PB) {
		var value string = "value"
		var key string = strconv.Itoa(rand.Intn(100))
		for pb.Next() {
			_, _ = syncMap.LoadOrStore(key, value)
		}
	})

}
func BenchmarkTestSyncMap_LoadAndDelete(b *testing.B) {
	b.ReportAllocs()
	var syncMap = SyncMap[string, string]{
		Map:  map[string]string{},
		lock: sync.RWMutex{},
	}
	var value string = "value"
	var key string = strconv.Itoa(rand.Intn(100))
	for i := 0; i < 100; i++ {
		syncMap.Store(key, value) //

	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = syncMap.LoadAndDelete(key)
		}
	})

}

/*
goos: darwin
goarch: arm64
pkg: orm/internal
BenchmarkTestMap_Load
BenchmarkTestMap_Load-8   	 7281406	       153.0 ns/op	      48 B/op	       3 allocs/op
PASS
*/
func BenchmarkTestMap_Load(b *testing.B) {
	b.ReportAllocs()
	var syncMap = Map[string, string]{}
	b.RunParallel(func(pb *testing.PB) {
		rand.Seed(time.Now().Unix())
		var key string = strconv.Itoa(rand.Intn(1000))
		var value string = "value"
		for pb.Next() {
			syncMap.LoadOrStore(key, value)
		}
	})
}
