package main

type frequency struct {
	Token token
	Count int
}

type frequencySorter []*frequency

func (fs frequencySorter) Len() int {
	return len(fs)
}

func (fs frequencySorter) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func (fs frequencySorter) Less(i, j int) bool {
	// Sort by Frequency count DESC, Token (Alphabet) ASC
	vi := fs[i]
	vj := fs[j]
	if vi.Count == vj.Count {
		return !vi.Token.Compare(vj.Token)
	}
	return fs[i].Count > fs[j].Count
}
