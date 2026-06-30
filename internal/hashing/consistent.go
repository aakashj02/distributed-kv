package hashing

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// HashRing manages the distribution of keys across nodes.
type HashRing struct {
	mu       sync.RWMutex
	replicas int               
	keys     []uint32          
	hashMap  map[uint32]string 
}

func NewHashRing(replicas int) *HashRing {
	return &HashRing{
		replicas: replicas,
		hashMap:  make(map[uint32]string),
	}
}

func (hr *HashRing) AddNode(node string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	for i := 0; i < hr.replicas; i++ {
		hash := crc32.ChecksumIEEE([]byte(strconv.Itoa(i) + node))
		hr.keys = append(hr.keys, hash)
		hr.hashMap[hash] = node
	}
	
	sort.Slice(hr.keys, func(i, j int) bool {
		return hr.keys[i] < hr.keys[j]
	})
}

func (hr *HashRing) GetNode(key string) string {
	if len(hr.keys) == 0 {
		return ""
	}

	hr.mu.RLock()
	defer hr.mu.RUnlock()

	hash := crc32.ChecksumIEEE([]byte(key))

	idx := sort.Search(len(hr.keys), func(i int) bool {
		return hr.keys[i] >= hash
	})

	if idx == len(hr.keys) {
		idx = 0
	}

	return hr.hashMap[hr.keys[idx]]
}