package hashmap

import (
	"errors"
	"strconv"
	"sync"
	"unicode/utf8"
)

type HashMap struct {
	m   map[string]string
	mux *sync.Mutex
}

func New() *HashMap {
	return &HashMap{m: make(map[string]string), mux: &sync.Mutex{}}
}

func (m *HashMap) Set(key, value string) int {
	m.mux.Lock()
	defer m.mux.Unlock()

	if value, found := m.m[key]; found {
		m.m[key] = value
		return 0
	}

	m.m[key] = value
	return 1
}

func (m *HashMap) Setx(key, value string) int {
	m.mux.Lock()
	defer m.mux.Unlock()

	if _, found := m.m[key]; found {
		return 0
	}

	m.m[key] = value
	return 1
}

func (m *HashMap) SetBulk(fields []string) (string, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if len(fields)%2 != 0 || len(fields) == 0 {
		return "", errors.New("invalid number of arguments")
	}

	for i := 0; i < len(fields); i += 2 {
		m.m[fields[i]] = fields[i+1]
	}

	return "OK", nil
}

func (m *HashMap) Get(key string) string {
	m.mux.Lock()
	defer m.mux.Unlock()

	return m.m[key]
}

func (m *HashMap) GetBulk(keys ...string) []string {
	m.mux.Lock()
	defer m.mux.Unlock()

	results := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		results[i] = m.m[keys[i]]
	}

	return results
}

func (m *HashMap) Keys() []string {
	m.mux.Lock()
	defer m.mux.Unlock()

	keys := make([]string, len(m.m))
	i := 0
	for k := range m.m {
		keys[i] = k
		i++
	}

	return keys
}

func (m *HashMap) Vals() []string {
	m.mux.Lock()
	defer m.mux.Unlock()

	vals := make([]string, len(m.m))
	i := 0
	for _, v := range m.m {
		vals[i] = v
		i++
	}

	return vals
}

func (m *HashMap) Fields() []string {
	m.mux.Lock()
	defer m.mux.Unlock()

	paris := make([]string, len(m.m)*2)
	i := 0

	for key, val := range m.m {
		paris[i] = key
		i++
		paris[i] = val
		i++
	}
	return paris
}

func (m *HashMap) Delete(keys ...string) int {
	m.mux.Lock()
	defer m.mux.Unlock()

	deleted := 0
	for _, key := range keys {
		if _, found := m.m[key]; found {
			delete(m.m, key)
			deleted++
		}
	}

	return deleted
}

func (m *HashMap) Exists(key string) int {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, found := m.m[key]

	if found {
		return 1
	}

	return 0
}

func (m *HashMap) IncrementBy(key, amount string) (int, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	// Convert value to int
	val, err := strconv.Atoi(amount)

	if err != nil {
		return 0, errors.New("invalid integer or integer out of range")
	}

	target, found := m.m[key]

	if !found {
		m.m[key] = amount
		return val, nil
	}

	// we have a hit
	targetVal, err := strconv.Atoi(target)

	if err != nil {
		return 0, errors.New("invalid type, excepted integer")
	}

	newVal := targetVal + val
	m.m[key] = strconv.Itoa(newVal)

	return newVal, nil
}

func (m *HashMap) IncrementByFloat(key, amount string) (float64, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	// Convert value to int
	val, err := strconv.ParseFloat(amount, 64)

	if err != nil {
		return 0, errors.New("invalid float or float out of range")
	}

	target, found := m.m[key]

	if !found {
		m.m[key] = strconv.FormatFloat(val, 'f', 6, 64)
		return val, nil
	}

	// we have a hit
	targetVal, err := strconv.ParseFloat(target, 64)

	if err != nil {
		return 0, errors.New("invalid type, excepted float")
	}

	newVal := targetVal + val
	m.m[key] = strconv.FormatFloat(newVal, 'f', 6, 64)

	return newVal, nil
}

func (m *HashMap) Len() int {
	m.mux.Lock()
	defer m.mux.Unlock()

	return len(m.m)
}

func (m *HashMap) FLen(key string) int {
	m.mux.Lock()
	defer m.mux.Unlock()

	return utf8.RuneCountInString(m.m[key])
}
