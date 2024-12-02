package linkedmap

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/quintans/faults"
)

var errEOA = errors.New("End of Array")

type MapJSON Map[string, any]

func NewJSON() *MapJSON {
	return (*MapJSON)(New[string, any]())
}

func (m *MapJSON) Unwrap() *Map[string, any] {
	return (*Map[string, any])(m)
}

// MarshalJSON implements the json.Marshaller interface, so it could be serialized.
// When serializing, the keys of the map will keep the order they are added.
func (m MapJSON) MarshalJSON() ([]byte, error) {
	var out bytes.Buffer

	out.WriteString("{")
	idx := 0
	for it := m.list.Head(); it != nil; it = it.Next() {
		if idx > 0 {
			out.WriteString(",")
		}

		esc := strings.Replace(it.Value.key, `"`, `\"`, -1)
		out.WriteString(`"` + esc + `"`)

		out.WriteString(":")

		// marshal the value
		b, err := json.Marshal(it.Value.value)
		if err != nil {
			return []byte{}, faults.Wrap(err)
		}
		out.WriteString(string(b))

		idx++
	}

	out.WriteString("}")
	return out.Bytes(), nil
}

func (m *MapJSON) UnmarshalJSON(b []byte) error {
	// Using Decoder to parse the bytes.
	in := bytes.TrimSpace(b)
	dec := json.NewDecoder(bytes.NewReader(in))

	t, err := dec.Token()
	if err != nil {
		return faults.Wrap(err)
	}

	// must open with a delim token '{'
	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return faults.Errorf("expect JSON object open with '{'")
	}

	err = m.parseObject(dec)
	if err != nil {
		return faults.Wrap(err)
	}

	t, err = dec.Token() //'}'
	if err != nil {
		return faults.Wrap(err)
	}
	if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return faults.Errorf("expect JSON object close with '}'")
	}

	return nil
}

func (m *MapJSON) parseObject(dec *json.Decoder) error {
	om := (*Map[string, any])(m)
	for dec.More() { // Loop until it has no more tokens
		t, err := dec.Token()
		if err != nil {
			return faults.Wrap(err)
		}

		key, ok := t.(string)
		if !ok {
			return faults.Errorf("key must be a string, got %T\n", t)
		}

		val, err := parseValue(dec)
		if err != nil {
			return faults.Wrap(err)
		}
		om.Set(key, val)
	}
	return nil
}

func parseValue(dec *json.Decoder) (any, error) {
	t, err := dec.Token()
	if err != nil {
		return nil, faults.Wrap(err)
	}

	switch tok := t.(type) {
	case json.Delim:
		switch tok {
		case '[': // If it's an array
			return parseArray(dec)
		case '{': // If it's a map
			om := (*MapJSON)(New[string, any]())
			err := om.parseObject(dec)
			if err != nil {
				return nil, faults.Wrap(err)
			}
			_, err = dec.Token() // }
			if err != nil {
				return nil, faults.Wrap(err)
			}
			return om, nil
		case ']':
			return nil, faults.Wrap(errEOA)
		case '}':
			return nil, faults.New("unexpected '}'")
		default:
			return nil, faults.Errorf("Unexpected delimiter: %q", tok)
		}
	default:
		return tok, nil
	}
}

func parseArray(dec *json.Decoder) ([]any, error) {
	ret := []any{}
	for {
		v, err := parseValue(dec)
		if errors.Is(err, errEOA) {
			return ret, nil
		}
		if err != nil {
			return nil, faults.Wrap(err)
		}
		ret = append(ret, v)
	}
}
