package fhirterm

type IntsetMap map[int64]bool

type Intset struct {
	M IntsetMap
}

func NewIntset() *Intset {
	m := make(IntsetMap)
	return &Intset{m}
}

func NewIntsetFromSlice(slice []int64) *Intset {
	set := NewIntset()
	set.AddSlice(slice)
	return set
}

func (s *Intset) Add(i int64) bool {
	_, found := s.M[i]
	s.M[i] = true
	return !found
}

func (s *Intset) Iter() <-chan int64 {
	ch := make(chan int64)
	go func() {
		for val, _ := range s.M {
			ch <- val
		}
		close(ch)
	}()

	return ch
}

func (s *Intset) AddSet(other *Intset) {
	for i, _ := range other.M {
		s.M[i] = true
	}
}

func (s *Intset) AddSlice(other []int64) {
	for _, val := range other {
		s.M[val] = true
	}
}

func (s *Intset) Len() int {
	return len(s.M)
}

func (s *Intset) ToInt64Slice() []int64 {
	result := make([]int64, 0, s.Len())

	for val, _ := range s.M {
		result = append(result, val)
	}

	return result
}

func (s *Intset) ToIntSlice() []int {
	result := make([]int, 0, s.Len())

	for val, _ := range s.M {
		result = append(result, int(val))
	}

	return result
}

func (s *Intset) Contains(v int64) bool {
	_, found := s.M[v]

	return found
}

func (s *Intset) Remove(v int64) bool {
	_, found := s.M[v]

	if found {
		delete(s.M, v)
		return true
	} else {
		return false
	}
}

func (s *Intset) Union(other *Intset) *Intset {
	result := NewIntset()
	result.AddSet(s)
	result.AddSet(other)

	return result
}

func (s *Intset) Intersect(other *Intset) *Intset {
	var smaller *Intset
	var larger *Intset

	if other.Len() > s.Len() {
		smaller = s
		larger = other
	} else {
		smaller = other
		larger = s
	}

	result := NewIntset()

	for val, _ := range smaller.M {
		if larger.Contains(val) {
			result.Add(val)
		}
	}

	return result
}

func (s *Intset) Difference(other *Intset) *Intset {
	result := NewIntset()

	for val, _ := range s.M {
		if !other.Contains(val) {
			result.Add(val)
		}
	}

	return result
}

func (s *Intset) Equal(other *Intset) bool {
	if s.Len() != other.Len() {
		return false
	} else {
		for val, _ := range s.M {
			if !other.Contains(val) {
				return false
			}
		}

		return true
	}
}
