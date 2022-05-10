package collections_test

import (
	"testing"

	"github.com/quintans/dstruct/collections"
	"github.com/stretchr/testify/require"
)

func TestLinkedHashMapIterator(t *testing.T) {
	r := require.New(t)

	lhm := collections.NewLinkedHashMap[Foo, string](collections.Equals[Foo], collections.HashCode[Foo])
	lhm.Put(foos[0], "chilling")
	lhm.Put(foos[1], "working")
	lhm.Put(foos[2], "stressing")
	lhm.Put(foos[3], "rebelling")

	tmp := []Foo{}
	for it := lhm.Iterator(); it.HasNext(); {
		tmp = append(tmp, it.Next().Key)
	}
	r.EqualValues(foos, tmp)

	tmp = []Foo{}
	lhm.ForEach(func(f Foo, s string) {
		tmp = append(tmp, f)
	})
	r.EqualValues(foos, tmp)

	lhm.Delete(foos[3])

	foos2 := foos[:3]
	tmp = []Foo{}
	for it := lhm.Iterator(); it.HasNext(); {
		tmp = append(tmp, it.Next().Key)
	}
	r.EqualValues(foos2, tmp)

	tmp = []Foo{}
	lhm.ForEach(func(f Foo, s string) {
		tmp = append(tmp, f)
	})
	r.EqualValues(foos2, tmp)
}
