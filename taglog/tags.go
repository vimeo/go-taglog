package taglog

// A map type specific to tags. The value type must be either string or []string.
// Users should avoid modifying the map directly and instead use the provided
// functions.
type Tags map[string]interface{}

// Add one or more values to a key.
func (t Tags) Add(key string, value ...string) {
	for _, v := range value {
		switch vs := t[key].(type) {
		case nil:
			t[key] = v
		case string:
			t[key] = []string{vs, v}
		case []string:
			t[key] = append(vs, v)
		}
	}
}

// Add one or more values to a key, merging any duplicate values.
func (t Tags) Merge(key string, value ...string) {
	for _, v := range value {
		current := t.GetAll(key)
		found := false
		for _, cv := range current {
			if v == cv {
				found = true
				break
			}
		}
		if !found {
			t.Add(key, v)
		}
	}
}

// Append one or more values to a key. This the same as Add() and is only
// provided to couple with Pop() for code clarity.
func (t Tags) Push(key string, value ...string) {
	t.Add(key, value...)
}

// Remove the last value for a key
func (t Tags) Pop(key string) {
	switch vs := t[key].(type) {
	case nil:
		return
	case string:
		delete(t, key)
	case []string:
		if len(vs) <= 1 {
			delete(t, key)
		} else if len(vs) == 2 {
			t[key] = vs[0]
		} else {
			t[key] = vs[:len(vs)-1]
		}
	}
}

// Set one or more values for a key. Any existing values are discarded.
func (t Tags) Set(key string, value ...string) {
	delete(t, key)
	t.Add(key, value...)
}

// Get the first value for a key. If the key does not exist, an empty string is
// returned.
func (t Tags) Get(key string) string {
	switch vs := t[key].(type) {
	case string:
		return vs
	case []string:
		return vs[0]
	}
	return ""
}

// Get all the values for a key. If the key does not exist, a nil slice is
// returned.
func (t Tags) GetAll(key string) []string {
	switch vs := t[key].(type) {
	case string:
		return []string{vs}
	case []string:
		return vs
	}
	return nil
}

// Delete a key.
func (t Tags) Del(key string) {
	delete(t, key)
}

// Delete all keys.
func (t Tags) DelAll() {
	for k, _ := range t {
		delete(t, k)
	}
}

// Export all tags as a map of string slices.
func (t Tags) Export() map[string][]string {
	tags := make(map[string][]string)
	for k, v := range t {
		switch vs := v.(type) {
		case string:
			tags[k] = []string{vs}
		case []string:
			ts := make([]string, len(vs))
			copy(ts, vs)
			tags[k] = ts
		}
	}
	return tags
}

// Import tags from a map of string slices.
func (t Tags) Import(tags map[string][]string) {
	for k, v := range tags {
		t.Merge(k, v...)
	}
}

// Copy tags. Performs a deep copy of all tag values.
func (t Tags) Copy() Tags {
	out := make(Tags)
	out.Import(t.Export())
	return out
}
