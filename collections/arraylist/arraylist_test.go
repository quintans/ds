package arraylist_test

import (
	"testing"

	"github.com/quintans/dstruct/collections/arraylist"
	"github.com/stretchr/testify/require"
)

var unsortedArray = []int{10, 2, 6, 71, 3}

func TestAdd(t *testing.T) {
	list := arraylist.New[int]()
	list.Add(unsortedArray...)
	require.EqualValues(t, unsortedArray, list.ToSlice())
	require.Equal(t, len(unsortedArray), list.Size())
}

func TestGet(t *testing.T) {
	list := arraylist.New[int]()
	list.Add(unsortedArray...)

	testcases := []struct {
		name    string
		idx     int
		want    int
		wantErr bool
	}{
		{
			name: "upper_bound",
			idx:  len(unsortedArray) - 1,
			want: 3,
		},
		{
			name: "lower_bound",
			idx:  0,
			want: 10,
		},
		{
			name:    "upper_out_of_bounds",
			idx:     5,
			wantErr: true,
		},
		{
			name:    "lower_out_of_bounds",
			idx:     -1,
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := list.Get(tc.idx)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.want, v)
		})
	}
}

func TestSet(t *testing.T) {
	list := arraylist.New[int]()
	list.Add(unsortedArray...)

	testcases := []struct {
		name    string
		idx     int
		value   int
		wantErr bool
	}{
		{
			name:  "upper_bound",
			idx:   len(unsortedArray) - 1,
			value: 5,
		},
		{
			name:  "lower_bound",
			idx:   0,
			value: 9,
		},
		{
			name:    "upper_out_of_bounds",
			idx:     5,
			wantErr: true,
		},
		{
			name:    "lower_out_of_bounds",
			idx:     -1,
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := list.Set(tc.idx, tc.value)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				v, _ := list.Get(tc.idx)
				require.Equal(t, tc.value, v)
			}
		})
	}
}

func TestAddAll(t *testing.T) {
	list := arraylist.New[int]()
	list.Add(unsortedArray...)
	list2 := arraylist.New[int]()
	list2.AddAll(list)
	require.EqualValues(t, unsortedArray, list2.ToSlice())
}

func TestIndexOf(t *testing.T) {
	list := arraylist.New[int]()
	list.Add(unsortedArray...)

	i := list.IndexOf(2)
	require.Equal(t, 1, i)
	i = list.IndexOf(5)
	require.Equal(t, -1, i)
}

func TestContains(t *testing.T) {
	list := arraylist.New[int]()
	list.Add(unsortedArray...)

	ok := list.Contains(25)
	require.False(t, ok)
	ok = list.Contains(2)
	require.True(t, ok)
}

func TestDelete(t *testing.T) {
	list := arraylist.New[int]()
	list.Add(unsortedArray...)

	ok := list.Delete(25)
	require.False(t, ok)
	ok = list.Delete(2)
	require.True(t, ok)
}

func TestDeleteAt(t *testing.T) {
	list := arraylist.New[int]()
	list.Add(unsortedArray...)

	testcases := []struct {
		name    string
		idx     int
		want    int
		wantErr bool
	}{
		{
			name: "upper_bound",
			idx:  len(unsortedArray) - 1,
			want: 3,
		},
		{
			name: "lower_bound",
			idx:  0,
			want: 10,
		},
		{
			name:    "upper_out_of_bounds",
			idx:     5,
			wantErr: true,
		},
		{
			name:    "lower_out_of_bounds",
			idx:     -1,
			wantErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			size := list.Size()
			v, err := list.DeleteAt(tc.idx)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.Equal(t, size-1, list.Size())
				require.NoError(t, err)
			}
			require.Equal(t, tc.want, v)
		})
	}
}

func TestForEach(t *testing.T) {
	list := arraylist.New[int]()
	list.Add(unsortedArray...)

	list.ForEach(func(k, v int) {
		require.Equal(t, unsortedArray[k], v)
	})
}

func TestReplaceAll(t *testing.T) {
	list := arraylist.New[int]()
	list.Add(unsortedArray...)

	list.ReplaceAll(func(_, _ int) int {
		return -1
	})

	arr := []int{-1, -1, -1, -1, -1}
	require.Equal(t, arr, list.ToSlice())
}

func TestEnumerator(t *testing.T) {
	list := arraylist.New[int]()
	list.Add(unsortedArray...)

	pos := 0
	for e := list.Iterator(); e.HasNext(); {
		require.Equal(t, unsortedArray[pos], e.Next())
		pos++
	}
}
