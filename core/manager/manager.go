package manager

import "fmt"

type Manager struct {
	lastId int
	format string

	data map[string]interface{}
}

func NewManager(format string) *Manager {
	m := &Manager{0, format, map[string]interface{}{}}

	return m
}

func (m *Manager) New(data interface{}) string {
	m.lastId++
	newId := fmt.Sprintf(m.format, m.lastId)

	m.data[newId] = data

	return newId
}

func (m *Manager) NewCustom(key string, data interface{}) {
	m.data[key] = data
}

func (m *Manager) Del(k string) {
	delete(m.data, k)
}

func (m *Manager) Get(k string) interface{} {
	v, _ := m.data[k]

	return v
}

func (m *Manager) Iter(cb func(k string, v interface{})) {
	for k, v := range m.data {
		cb(k, v)
	}
}

func (m *Manager) Keys() []string {
	ret := make([]string, 0, len(m.data))

	for k := range m.data {
		ret = append(ret, k)
	}

	return ret
}
