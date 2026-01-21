package helper

// set
// https://gist.github.com/Veivel/ab9d590b1f6fe03c2efbdd1f9859af5d#file-set-go
type set[K comparable] struct {
	data map[K]int
}

func NewSet[K comparable]() set[K] {
	return set[K]{data: make(map[K]int)}
}

func (s *set[K]) IsExist(item K) bool {
	_, ok := s.data[item]
	return ok
}

func (s *set[K]) Get() []K {
	keys := make([]K, 0, len(s.data))
	for k, _ := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *set[K]) Add(item K) set[K] {
	(s).data[item] = 1
	return *s
}

func (s *set[K]) Delete(item K) set[K] {
	delete((s).data, item)
	return *s
}

//Mxsxll

func (s *set[K]) Size() int {
	return len(s.data)
}
