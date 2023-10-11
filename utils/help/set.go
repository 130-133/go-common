package help

import "go.uber.org/atomic"

type void struct{}
type IT interface {
	~string | ~int | ~int32 | ~int64 | ~float64 | ~float32 | ~bool
}

type Set[T IT] struct {
	pool map[T]void
	sync atomic.Bool
}

func NewSet[T IT](args ...T) *Set[T] {
	s := Set[T]{
		pool: make(map[T]void),
	}
	for _, arg := range args {
		s.Add(arg)
	}
	return &s
}

func (s *Set[T]) Has(args T) bool {
	_, ok := s.pool[args]
	return ok
}

func (s *Set[T]) Add(args T) *Set[T] {
	if s.sync.Load() {
		return s
	}
	s.sync.CAS(false, true)
	defer s.sync.CAS(true, false)
	s.pool[args] = void{}
	return s
}

func (s *Set[T]) AddAll(args ...T) *Set[T] {
	for _, arg := range args {
		s.Add(arg)
	}
	return s
}

func (s *Set[T]) Remove(args T) *Set[T] {
	if s.sync.Load() {
		return s
	}
	s.sync.CAS(false, true)
	defer s.sync.CAS(true, false)
	delete(s.pool, args)
	return s
}

func (s *Set[T]) Values() []T {
	var res []T
	for k := range s.pool {
		res = append(res, k)
	}
	return res
}
