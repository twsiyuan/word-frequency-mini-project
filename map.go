package main

import "errors"

// MapKeyer is used in Mapper
type MapKeyer interface {
	Equal(MapKeyer) bool
	Hash() int
	Clone() MapKeyer
}

// Mapper in map interface
type Mapper interface {
	Count() int
	Get(key MapKeyer) interface{}
	Set(key MapKeyer, value interface{})
	Keys([]interface{}) []interface{}
	Values([]interface{}) []interface{}
}

// NewMap is create map structure
func NewMap(capacity int) (Mapper, error) {
	if capacity <= 0 {
		return nil, errors.New("Bad capacity")
	}
	m := mymap{
		cap: make([][]mapValue, capacity),
	}
	return &m, nil
}

type mapValue struct {
	Key   MapKeyer
	Value interface{}
}

type mymap struct {
	cap   [][]mapValue
	count int
}

func (m mymap) Count() int {
	return m.count
}

func (m mymap) index(key MapKeyer) int {
	hash := key.Hash()
	return hash % len(m.cap)
}

func (m mymap) Get(key MapKeyer) interface{} {
	data := m.cap[m.index(key)]
	if data != nil {
		for _, d := range data {
			if key.Equal(d.Key) {
				return d.Value
			}
		}
	}
	return nil
}

func (m *mymap) Set(key MapKeyer, value interface{}) {
	capIdx := m.index(key)
	data := m.cap[capIdx]
	mvalue := mapValue{
		Key:   key,
		Value: value,
	}
	if data == nil {
		data = make([]mapValue, 1)
		data[0] = mvalue
		m.cap[capIdx] = data
		m.count++
	} else {
		sliceIdx := -1
		for idx, d := range data {
			if key.Equal(d.Key) {
				sliceIdx = idx
				break
			}
		}
		if sliceIdx < 0 {
			data = append(data, mapValue{
				Key:   key,
				Value: value,
			})
			m.count++
		} else {
			data[sliceIdx] = mvalue
		}
	}
}

func (m mymap) Keys(output []interface{}) []interface{} {
	for i := 0; i < len(m.cap); i++ {
		for k := 0; k < len(m.cap[i]); k++ {
			output = append(output, m.cap[i][k].Key)
		}
	}
	return output
}

func (m mymap) Values(output []interface{}) []interface{} {
	for i := 0; i < len(m.cap); i++ {
		for k := 0; k < len(m.cap[i]); k++ {
			output = append(output, m.cap[i][k].Value)
		}
	}
	return output
}
