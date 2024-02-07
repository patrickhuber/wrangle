package envdiff

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
)

// Change is one of Add, Remove or Update
type Change interface {
	change()
}

type Update struct {
	Change   `json:"-"`
	Key      string `json:"key"`
	Previous string `json:"previous"`
	Value    string `json:"value"`
}

type Add struct {
	Change `json:"-"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type Remove struct {
	Change   `json:"-"`
	Key      string `json:"key"`
	Previous string `json:"previous"`
}

func Diff(previous, next map[string]string) []Change {
	var diff []Change
	var keys []string

	// sort keys for deterministic testing
	for key := range previous {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	// add different or missing items to previous
	for _, key := range keys {
		previousValue := previous[key]
		nextValue, inNext := next[key]

		if inNext {
			if previousValue == nextValue {
				continue
			}
			// in previous, in next with different value : it was updated
			diff = append(diff, Update{
				Key:      key,
				Previous: previousValue,
				Value:    nextValue,
			})
		} else {
			// in previous, not in next : it was removed
			diff = append(diff, Remove{
				Key:      key,
				Previous: previousValue,
			})
		}
	}

	keys = keys[:0]

	// sort keys for deterministic testing
	for key := range next {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	// add different or missing items to next
	for _, key := range keys {

		nextValue := next[key]

		// not in previous, in next : it was added
		if _, inPrevious := previous[key]; inPrevious {
			continue
		}
		diff = append(diff, Add{
			Key:   key,
			Value: nextValue,
		})
	}

	return diff
}

func Apply(env map[string]string, changes []Change) {
	for _, change := range changes {
		switch c := change.(type) {
		case Add:
			env[c.Key] = c.Value
		case Remove:
			delete(env, c.Key)
		case Update:
			env[c.Key] = c.Value
		}
	}
}

func Revert(env map[string]string, changes []Change) {
	for _, change := range changes {
		switch c := change.(type) {
		case Add:
			delete(env, c.Key)
		case Remove:
			env[c.Key] = c.Previous
		case Update:
			env[c.Key] = c.Previous
		}
	}
}

func Encode(changes []Change) (string, error) {
	jsonData, err := json.Marshal(changes)
	if err != nil {
		return "", err
	}

	// gzip
	gzBuf := &bytes.Buffer{}
	writer := gzip.NewWriter(gzBuf)
	_, err = writer.Write(jsonData)
	if err != nil {
		return "", err
	}
	if err = writer.Close(); err != nil {
		return "", err
	}

	// base64 encode (reuse json buf)
	return base64.URLEncoding.EncodeToString(gzBuf.Bytes()), nil
}

func Decode(encoded string) ([]Change, error) {

	// base64 decode
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	// gunzip
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	gzbuf := &bytes.Buffer{}
	_, err = gzbuf.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	if err = reader.Close(); err != nil {
		return nil, err
	}

	// json decode
	return unmarshallJSON(gzbuf.Bytes())
}

func unmarshallJSON(buf []byte) ([]Change, error) {
	arr := []any{}

	err := json.Unmarshal(buf, &arr)
	if err != nil {
		return nil, err
	}

	var changes []Change
	for _, item := range arr {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("unknown type %T found while unmarshalling json", item)
		}

		v, vok := m["value"]
		p, pok := m["previous"]
		if !pok && !vok {
			return nil, fmt.Errorf("one of 'value' or 'previous' required")
		}

		k, kok := m["key"]
		if !kok {
			return nil, fmt.Errorf("missing required attribute 'key'")
		}

		var change Change
		if vok && pok {
			change = Update{Key: k.(string), Previous: p.(string), Value: v.(string)}
		} else if vok {
			change = Add{Key: k.(string), Value: v.(string)}
		} else {
			change = Remove{Key: k.(string), Previous: p.(string)}
		}
		changes = append(changes, change)
	}
	return changes, err
}
