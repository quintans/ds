package linkedhashmap_test

import (
	"fmt"
	"testing"

	"github.com/quintans/dstruct/collections/linkedhashmap"
	"github.com/stretchr/testify/require"
)

type Foo struct {
	Name string
	Age  int
}

var foos = []Foo{
	{"Martim", 9},
	{"Paulo", 41},
	{"Monica", 33},
	{"Francisca", 15},
}

func TestLinkedHashMapIterator(t *testing.T) {
	r := require.New(t)

	lhm := linkedhashmap.New[Foo, string]()
	lhm.Put(foos[0], "chilling")
	lhm.Put(foos[1], "working")
	lhm.Put(foos[2], "stressing")
	lhm.Put(foos[3], "rebelling")

	entries := lhm.Entries()
	for _, e := range entries {
		fmt.Printf("===> %v: %v\n", e.Key, e.Value)
	}

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
