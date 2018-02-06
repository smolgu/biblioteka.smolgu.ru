package intarr

import (
	"github.com/cznic/sortutil"
)

func Sort(s Slice) {
	sortutil.Uint64Slice(s).Sort()
}

func (s Slice) Sort() {
	sortutil.Uint64Slice(s).Sort()
}
