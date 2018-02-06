package intarr

type Slice []uint64

func New(arr interface{}) (s Slice) {
	switch arr.(type) {
	case []int:
		for _, v := range arr.([]int) {
			s = append(s, uint64(v))
		}
	case []uint64:
		return Slice(arr.([]uint64))
	case []int64:
		for _, v := range arr.([]int64) {
			s = append(s, uint64(v))
		}
	/*case []byte:
	s, e = Decode(arr.([]byte))
	return*/
	default:
		return
	}
	return
}

func (sl Slice) Encode() ([]byte, error) {
	return is2b(sl)
}

func Decode(b []byte) (Slice, error) {
	return b2is(b)
}

func (sl Slice) In(i uint64) bool {
	for _, v := range sl {
		if i == v {
			return true
		}
	}
	return false
}

func (sl Slice) AsInt() (res []int) {
	for _, i := range sl {
		res = append(res, int(i))
	}
	return
}

func (sl Slice) Int64() (res []int64) {
	for _, i := range sl {
		res = append(res, int64(i))
	}
	return
}

func (sl Slice) Remove(num uint64) (res Slice) {
	for i, v := range sl {
		if v == num {
			res = append(sl[:i], sl[i+1:]...)
			return
		}
	}
	return sl
}
