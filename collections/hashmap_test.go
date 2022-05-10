package collections_test

import (
	"strconv"
	"testing"

	"github.com/quintans/dstruct/collections"
	"github.com/stretchr/testify/require"
)

const value1 = "World"

func TestPutAndGet(t *testing.T) {
	r := require.New(t)

	hm := collections.NewHashMap[string, string](
		collections.Equals[string],
		collections.HashCode[string],
	)
	k := "Hello"
	hm.Put(k, value1)
	v, ok := hm.Get(k)
	r.Equal(value1, v)
	r.True(ok)
	_, ok = hm.Get("nothing")
	r.False(ok)
}

func TestResize(t *testing.T) {
	r := require.New(t)

	hm := collections.NewHashMap[string, int](collections.Equals[string], collections.HashCode[string])
	loop := 20
	// Insert
	for i := 0; i < loop; i++ {
		hm.Put("Hello"+strconv.Itoa(i), i*10)
	}
	r.Equal(loop, hm.Size())
	// Check
	for i := 0; i < loop; i++ {
		v, _ := hm.Get("Hello" + strconv.Itoa(i))
		k := i * 10
		require.Equal(t, k, v)
	}
	// Delete
	for i := 0; i < loop; i++ {
		v, _ := hm.Delete("Hello" + strconv.Itoa(i))
		k := i * 10
		r.Equal(k, v)
	}
	r.Equal(0, hm.Size())
}

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

func TestHashMapIterator(t *testing.T) {
	r := require.New(t)

	hm := collections.NewHashMap[Foo, string](
		collections.Equals[Foo],
		func(a Foo) int { return collections.HASH_SEED.HashString(a.Name).Int() },
	)
	hm.Put(foos[0], "chilling")
	hm.Put(foos[1], "working")
	hm.Put(foos[2], "stressing")
	hm.Put(foos[3], "rebelling")

	tmp := []Foo{}
	for it := hm.Iterator(); it.HasNext(); {
		tmp = append(tmp, it.Next().Key)
	}
	r.ElementsMatch(foos, tmp)

	tmp = []Foo{}
	hm.ForEach(func(f Foo, s string) {
		tmp = append(tmp, f)
	})
	r.ElementsMatch(foos, tmp)

	hm.Delete(foos[3])

	foos2 := foos[:3]
	tmp = []Foo{}
	for it := hm.Iterator(); it.HasNext(); {
		tmp = append(tmp, it.Next().Key)
	}
	r.ElementsMatch(foos2, tmp)

	tmp = []Foo{}
	hm.ForEach(func(f Foo, s string) {
		tmp = append(tmp, f)
	})
	r.ElementsMatch(foos2, tmp)
}
