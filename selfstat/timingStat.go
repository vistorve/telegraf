package selfstat

import (
	"hash/fnv"
	"sync"
)

type timingStat struct {
	measurement string
	field       string
	metadata    map[string]string
	key         uint64
	v           int64
	prev        int64
	count       int64
	mu          sync.Mutex
}

func (s *timingStat) Incr(v int64) {
	s.mu.Lock()
	s.v += v
	s.count++
	s.mu.Unlock()
}

func (s *timingStat) Set(v int64) {
}

func (s *timingStat) Get() int64 {
	var avg int64
	s.mu.Lock()
	if s.count > 0 {
		s.prev, avg = s.v/s.count, s.v/s.count
		s.v = 0
		s.count = 0
	} else {
		avg = s.prev
	}
	s.mu.Unlock()
	return avg
}

func (s *timingStat) Name() string {
	return s.measurement
}

func (s *timingStat) FieldName() string {
	return s.field
}

// Metadata returns a copy of the timingStat's metadata.
// NOTE this allocates a new map every time it is called.
func (s *timingStat) Tags() map[string]string {
	m := make(map[string]string, len(s.metadata))
	for k, v := range s.metadata {
		m[k] = v
	}
	return m
}

func (s *timingStat) Key() uint64 {
	if s.key == 0 {
		h := fnv.New64a()
		h.Write([]byte(s.measurement))
		for k, v := range s.metadata {
			h.Write([]byte(k + v))
		}
		s.key = h.Sum64()
	}
	return s.key
}
