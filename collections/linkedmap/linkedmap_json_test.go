package linkedmap_test

import (
	"encoding/json"
	"testing"

	"github.com/quintans/ds/collections/linkedmap"
	"github.com/stretchr/testify/require"
)

func TestSerialisation(t *testing.T) {
	s := `{"one":"1","two":"2","three":"3","sub":{"zero":"0","minus-one":"-1","minus-two":"-2"},"array":[1,2,3]}`
	om := linkedmap.NewJSON()
	err := json.Unmarshal([]byte(s), om)
	require.NoError(t, err)

	order := []string{}
	om.Unwrap().Range(func(key string, _ any, _ int) bool {
		order = append(order, key)
		return true
	})

	require.Equal(t, []string{"one", "two", "three", "sub", "array"}, order)

	j, err := json.Marshal(om)
	require.NoError(t, err)

	require.Equal(t, s, string(j))
}
