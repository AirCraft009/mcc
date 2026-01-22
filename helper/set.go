package helper

// Set
// https://gist.github.com/Veivel/ab9d590b1f6fe03c2efbdd1f9859af5d#file-set-go
type Set[K comparable] struct {
	data map[K]int
}

func NewSet[K comparable]() Set[K] {
	return Set[K]{data: make(map[K]int)}
}

func (s *Set[K]) IsExist(item K) bool {
	_, ok := s.data[item]
	return ok
}

func (s *Set[K]) Get() []K {
	keys := make([]K, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *Set[K]) Add(item K) Set[K] {
	(s).data[item] = 1
	return *s
}

func (s *Set[K]) Delete(item K) Set[K] {
	delete((s).data, item)
	return *s
}

//Mxsxll

func (s *Set[K]) Size() int {
	return len(s.data)
}
