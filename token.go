package main

type token []byte

func (t token) Equal(d MapKeyer) bool {
	v, ok := d.(token)
	if !ok {
		return false
	}
	if t == nil && v == nil {
		return true
	} else if t == nil && v != nil {
		return false
	} else if t != nil && v == nil {
		return false
	} else if len(t) != len(v) {
		return false
	}

	for i := 0; i < len(t); i++ {
		if t[i] != v[i] {
			return false
		}
	}
	return true
}

func (t token) String() string {
	return string(t)
}

func (t token) Hash() int {
	if t == nil {
		return 0
	}
	hash := 0
	for i, b := range t {
		hash += (i + 1) * int(b)
	}
	return hash
}

func (t token) Clone() MapKeyer {
	b := make([]byte, len(t))
	copy(b, t)
	return token(b)
}

// Compare, true means the order in back of v, otherwise order in front of v. (Alphabet order)
func (t token) Compare(v token) bool {
	for i := 0; i < len(t); i++ {
		if i >= len(v) {
			return true
		}
		if t[i] != v[i] {
			return t[i] < v[i]
		}
	}
	return true
}
