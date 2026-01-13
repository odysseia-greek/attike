package gateway

import (
	"encoding/json"
	"strconv"
)

func parseOptInt32(raw json.RawMessage) (*int32, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}

	// number?
	var n int64
	if err := json.Unmarshal(raw, &n); err == nil {
		v := int32(n)
		return &v, nil
	}

	// string?
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	if s == "" {
		return nil, nil
	}
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil, err
	}
	v := int32(i)
	return &v, nil
}

func parseOptInt64(raw json.RawMessage) (*int64, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}

	// number?
	var n int64
	if err := json.Unmarshal(raw, &n); err == nil {
		return &n, nil
	}

	// string?
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	if s == "" {
		return nil, nil
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return &i, nil
}
