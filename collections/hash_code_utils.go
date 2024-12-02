package collections

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
)

type Hasher interface {
	HashCode() int
}

// Collected methods which allow easy implementation of <code>hashCode</code>.
//
// Example use case:
//    h := NewHash(field1).
//    	Add(field2)
//    return h.Result();

// An initial value for a hashCode, to which is added contributions
// from fields. Using a non-zero value decreases collisions of hashCode
// values.
const (
	HASH_SEED    = Seed(23)
	prime_number = 31
)

func numberHashCode(aObject any) Seed {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, aObject)
	if err != nil {
		return 0
	}
	return bytesHashCode(buf.Bytes())
}

func bytesHashCode(what []byte) Seed {
	h := 0
	for _, v := range what {
		// h = 31*h + int(v)
		h = (h << 5) - h + int(v)
	}
	return Seed(h)
}

func firstTerm(aSeed Seed) Seed {
	return prime_number * aSeed
}

type Seed int

func (s Seed) Int() int {
	return int(s)
}

func (s Seed) HashAny(a any) Seed {
	result := s.hash(a)
	if result == 0 {
		result = s
	}
	return result
}

func (s Seed) hash(a any) Seed {
	if a == nil {
		return firstTerm(s)
	}

	switch t := a.(type) {
	case Hasher:
		return firstTerm(s) + Seed(t.HashCode())
	case *bool, bool, []bool,
		*int8, int8, []int8, *uint8, uint8, []uint8,
		*int16, int16, []int16, *uint16, uint16, []uint16,
		int, *int32, int32, []int32, *uint32, uint32, []uint32,
		*int64, int64, []int64, *uint64, uint64, []uint64,
		*float32, float32, []float32,
		*float64, float64, []float64:
		return firstTerm(s) + numberHashCode(a)
	case *string:
		return s.HashString(*t)
	case string:
		return s.HashString(t)
	case []string:
		for _, a := range t {
			s = s.HashString(a)
		}
		return s
	default:
		v := reflect.ValueOf(a)
		k := v.Kind()
		if k == reflect.Ptr {
			// tries pointer element
			o := v.Elem().Interface()
			result := s.hash(o)
			if result == 0 {
				// no luck with the pointer. lets use pointer address value
				result = firstTerm(s) + Seed(v.Pointer())
			}
			return result
		}
	}

	return 0
}

func (s Seed) HashBool(b bool) Seed {
	if b {
		return firstTerm(s) + 1
	}
	return firstTerm(s)
}

func (s Seed) HashInt(i int) Seed {
	return firstTerm(s) + Seed(i)
}

func (s Seed) HashUint(i uint) Seed {
	return s.HashInt(int(i))
}

func (s Seed) HashInt64(i int64) Seed {
	result := firstTerm(s) + Seed(i)
	a := i >> 32
	return firstTerm(result) + Seed(a)
}

func (s Seed) HashUint64(i uint64) Seed {
	return s.HashInt64(int64(i))
}

func (s Seed) HashFloat32(aFloat float32) Seed {
	bits := math.Float32bits(aFloat)
	return firstTerm(s) + Seed(bits)
}

func (s Seed) HashFloat64(f float64) Seed {
	bits := math.Float64bits(f)
	return s.HashUint64(bits)
}

func (s Seed) HashString(aString string) Seed {
	return s.HashBytes([]byte(aString))
}

func (s Seed) HashBytes(aBytes []byte) Seed {
	return firstTerm(s) + bytesHashCode(aBytes)
}
