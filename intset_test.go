package fhirterm

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

const N = 1500

func Test_AddToSet(t *testing.T) {
	s := NewIntset()
	ints := rand.Perm(N)

	for i := 0; i < len(ints); i++ {
		s.Add(int64(ints[i]))
	}

	if s.Len() != N {
		t.Errorf("Incorrect resulting length of set: %d", s.Len())
	}
}

func Test_AddDuplicates(t *testing.T) {
	s := NewIntset()

	if r := s.Add(1); r != true {
		t.Error("Adding new value, received false from s.Add()")
	}

	if r := s.Add(1); r == true {
		t.Error("Adding duplicate, received true from s.Add()")
	}
}

func Test_ToIntSlice(t *testing.T) {
	s := NewIntset()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	r := s.ToIntSlice()
	sort.Ints(r)

	if !reflect.DeepEqual(r, []int{1, 2, 3}) {
		t.Errorf("s.ToIntSlice() returned incorrect value: %d", r)
	}
}

func Test_Equal(t *testing.T) {
	s := NewIntset()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	if !s.Equal(NewIntsetFromSlice([]int64{1, 2, 3})) {
		t.Errorf("s.Equal() returned incorrect value")
	}

	if s.Equal(NewIntsetFromSlice([]int64{1, 2, 8})) {
		t.Errorf("s.Equal() returned incorrect value")
	}

	if s.Equal(NewIntsetFromSlice([]int64{1})) {
		t.Errorf("s.Equal() returned incorrect value")
	}

	if s.Equal(NewIntsetFromSlice([]int64{})) {
		t.Errorf("s.Equal() returned incorrect value")
	}
}

func Test_Remove(t *testing.T) {
	s := NewIntset()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	s.Remove(2)

	r := s.ToIntSlice()
	sort.Ints(r)

	if !reflect.DeepEqual(r, []int{1, 3}) {
		t.Errorf("s.Remove() didn't removed value '2'. Set contains: %d", r)
	}
}

func Test_Union(t *testing.T) {
	s := NewIntset()
	s.Add(1)
	s.Add(2)
	s.Add(3)

	s2 := NewIntset()
	s2.Add(3)
	s2.Add(4)
	s2.Add(5)
	s2.Add(42)

	r := s.Union(s2).ToIntSlice()
	sort.Ints(r)

	if !reflect.DeepEqual(r, []int{1, 2, 3, 4, 5, 42}) {
		t.Errorf("s.Union() returned incorrect value: %d", r)
	}
}

func Test_Intersect(t *testing.T) {
	s := NewIntset()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)
	s.Add(5)
	s.Add(6)

	s2 := NewIntset()
	s2.Add(3)
	s2.Add(4)
	s2.Add(5)
	s2.Add(42)

	r := s.Intersect(s2).ToIntSlice()
	sort.Ints(r)

	if !reflect.DeepEqual(r, []int{3, 4, 5}) {
		t.Errorf("s.Intersect() returned incorrect value: %d", r)
	}
}

func Test_Difference(t *testing.T) {
	s := NewIntset()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)
	s.Add(5)
	s.Add(6)

	s2 := NewIntset()
	s2.Add(3)
	s2.Add(4)
	s2.Add(5)
	s2.Add(42)

	r := s.Difference(s2).ToIntSlice()
	sort.Ints(r)

	if !reflect.DeepEqual(r, []int{1, 2, 6}) {
		t.Errorf("s.Difference() returned incorrect value: %d", r)
	}
}

func Test_Iter(t *testing.T) {
	s := NewIntset()
	s.Add(1)
	s.Add(2)
	s.Add(3)
	s.Add(4)
	s.Add(5)

	r := make([]int, 0)

	for v := range s.Iter() {
		r = append(r, int(v))
	}

	sort.Ints(r)

	if !reflect.DeepEqual(r, []int{1, 2, 3, 4, 5}) {
		t.Errorf("Misbehaving implementation of s.Iter(): %d", r)
	}
}
