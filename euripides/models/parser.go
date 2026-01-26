package models

import (
	"encoding/json"
	"strconv"
	"time"
)

func ParseAttikeTime(s string) (time.Time, error) {
	// common cases you already emit
	layouts := []string{
		time.RFC3339Nano,          // 2026-01-12T19:54:17.687459278Z
		time.RFC3339,              // 2026-01-12T19:54:17Z
		"2006-01-02T15:04:05.000", // 2026-01-12T19:51:32.052
		"2006-01-02T15:04:05",     // if you ever send seconds only
	}

	var lastErr error
	for _, l := range layouts {
		t, err := time.Parse(l, s)
		if err == nil {
			return t.UTC(), nil
		}
		lastErr = err
	}
	return time.Time{}, lastErr
}

func IsEmptyJSON(b json.RawMessage) bool {
	if len(b) == 0 {
		return true
	}
	// trim spaces
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\n' || b[i] == '\t' || b[i] == '\r') {
		i++
	}
	if i == len(b) {
		return true
	}
	return string(b[i:]) == "null"
}

func ParseOptInt32(raw json.RawMessage) (*int32, error) {
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

func ParseOptInt64(raw json.RawMessage) (*int64, error) {
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
