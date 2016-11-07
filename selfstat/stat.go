package selfstat

import (
	"hash/fnv"
	"sync/atomic"
)

type stat struct {
	measurement string
	field       string
	metadata    map[string]string
	key         uint64
	v           int64
}

func (s *stat) Incr(v int64) {
	atomic.AddInt64(&s.v, v)
}

func (s *stat) Set(v int64) {
	atomic.StoreInt64(&s.v, v)
}

func (s *stat) Get() int64 {
	return atomic.LoadInt64(&s.v)
}

func (s *stat) Name() string {
	return s.measurement
}

func (s *stat) FieldName() string {
	return s.field
}

// Metadata returns a copy of the stat's metadata.
// NOTE this allocates a new map every time it is called.
func (s *stat) Tags() map[string]string {
	m := make(map[string]string, len(s.metadata))
	for k, v := range s.metadata {
		m[k] = v
	}
	return m
}

func (s *stat) Key() uint64 {
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
