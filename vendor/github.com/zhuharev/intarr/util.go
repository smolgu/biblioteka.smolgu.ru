package intarr

import "encoding/binary"

func Uint64ToBytes(i uint64) []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, i)
	return data
}

func BytesToUint64(bts []byte) uint64 {
	return binary.BigEndian.Uint64(bts)
}
func InArray(value int64, array []int64) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
