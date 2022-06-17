package app

import "strconv"

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) MustIntNotZero() int {
	v, _ := s.Int()
	if v <= 0 {
		return 1
	}
	return v
}

func (s StrTo) MustIntMax100() int {
	v, _ := s.Int()
	if v == 0 {
		return 10
	}
	if v > 100 {
		return 100
	}
	return v
}

func (s StrTo) UInt32() (uint32, error) {
	v, err := strconv.Atoi(s.String())
	return uint32(v), err
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}
