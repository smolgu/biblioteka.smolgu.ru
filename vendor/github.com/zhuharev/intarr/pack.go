package intarr

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/golang/snappy"
)

func Pack(in []uint64) (r []byte, e error) {
	Sort(in)
	w := bytes.NewBuffer(nil)
	if len(in) <= 48 {
		_, e = w.Write([]byte{'n'})
		if e != nil {
			return
		}
		for _, v := range in {
			e = binary.Write(w, binary.BigEndian, v)
			if e != nil {
				return
			}
		}
		return w.Bytes(), nil
		// } else {
		// 	_, e = w.Write([]byte{'p'})
		// 	if e != nil {
		// 		return
		// 	}
		// 	b, e := encode(in)
		// 	if e != nil {
		// 		return nil, e
		// 	}
		// 	_, e = w.Write(b)
		// 	return w.Bytes(), nil
		// }
	}
	return nil, nil
}

func Unpack(data []byte) (arr []uint64, e error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("eror len is %d", len(data))
	}
	r := bytes.NewReader(data)
	var t = make([]byte, 1)
	_, e = r.Read(t)
	if e != nil {
		return
	}
	switch t[0] {
	case 'n':
		l := len(data[1:]) / 4
		for i := 0; i < l; i++ {
			var num uint64
			e = binary.Read(r, binary.BigEndian, &num)
			if e != nil {
				return
			}
			arr = append(arr, num)
		}
		return
	default:
		e = fmt.Errorf("Unknow encoding type")
		return
	}

}

// func compress(in []uint64) (arr []uint64, e error) {
// 	l := len(in)
// 	if len(in)%128 != 0 {
// 		l = l + (128 - len(in)%128)
// 	}
//
// 	arr = make([]uint64, l)
// 	copy(arr, in)
//
// 	compdata := make([]uint64, 2*l)
// 	inpos := cursor.New()
// 	outpos := cursor.New()
// 	codec := bp32.New()
// 	e = codec.Compress(arr, inpos, l, compdata, outpos)
// 	if e != nil {
// 		return nil, e
// 	}
//
// 	return compdata[:outpos.Get()+1], e
// }

// func encode(in []uint64) ([]byte, error) {
//
// 	arr, e := compress(in)
// 	if e != nil {
// 		return nil, e
// 	}
//
// 	w := bytes.NewBuffer(nil)
// 	for _, v := range arr {
// 		e = binary.Write(w, binary.BigEndian, v)
// 		if e != nil {
// 			return nil, e
// 		}
// 	}
// 	return w.Bytes(), nil
// }

// func uncompress(arr []uint64) ([]uint64, error) {
// 	newinpos := cursor.New()
// 	newoutpos := cursor.New()
//
// 	recov := make([]uint64, 10000)
//
// 	codec := bp32.New()
// 	e := codec.Uncompress(arr, newinpos, len(arr)-1, recov, newoutpos)
// 	l := 0
// 	for i, v := range recov[:newoutpos.Get()] {
// 		if v != 0 {
// 			l = i
// 		}
// 	}
// 	return recov[:l+1], e
// }

// func decode(data []byte) ([]uint64, error) {
// 	r := bytes.NewReader(data)
//
// 	var arr []uint64
// 	var num uint64
// 	for e := binary.Read(r, binary.BigEndian, &num); e == nil; e = binary.Read(r, binary.BigEndian, &num) {
// 		arr = append(arr, num)
// 	}
// 	return uncompress(arr)
// }

func is2b(s Slice) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	for _, v := range s {
		e := binary.Write(w, binary.BigEndian, v)
		if e != nil {
			return nil, e
		}
	}
	return w.Bytes(), nil
}

func b2is(b []byte) (s Slice, e error) {
	r := bytes.NewReader(b)
	var i uint64
	for e = binary.Read(r, binary.BigEndian, &i); e == nil; e = binary.Read(r, binary.BigEndian, &i) {
		s = append(s, i)
	}
	if e == io.EOF {
		e = nil
	}
	return
}

func gz(data []byte) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	gw, e := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if e != nil {
		return nil, e
	}
	_, e = gw.Write(data)
	if e != nil {
		return nil, e
	}
	e = gw.Close()
	return w.Bytes(), e
}

func ungz(data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	gr, e := gzip.NewReader(r)
	if e != nil {
		return nil, e
	}
	buf := bytes.NewBuffer(nil)
	_, e = io.Copy(buf, gr)
	return buf.Bytes(), e
}

func PackGzip(s Slice) ([]byte, error) {
	b, e := is2b(s)
	if e != nil {
		return nil, e
	}
	return gz(b)
}

func UnpackGzip(data []byte) (Slice, error) {
	bts, e := ungz(data)
	if e != nil {
		return Slice{}, e
	}
	return b2is(bts)
}

func snap(arr []byte) ([]byte, error) {
	w := bytes.NewBuffer(nil)
	sw := snappy.NewWriter(w)
	_, e := sw.Write(arr)
	if e != nil {
		return nil, e
	}
	return w.Bytes(), nil
}

func PackSnappy(s Slice) ([]byte, error) {
	b, e := is2b(s)
	if e != nil {
		return nil, e
	}
	return snap(b)
}
