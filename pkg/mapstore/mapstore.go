package mapstore

import (
	"encoding/json"
	"io"
	"log"
	"sync"
)

type Collection struct {
	d map[string]interface{}
	sync.RWMutex
}

func (c *Collection) Has(k interface{}) bool {
	c.RLock()
	_, ok := c.d[k.(string)]
	c.RUnlock()
	return ok
}

func (c *Collection) Set(k, v interface{}) {
	c.Lock()
	c.d[k.(string)] = v
	c.Unlock()
}

func (c *Collection) Get(k interface{}) interface{} {
	c.RLock()
	v, ok := c.d[k.(string)]
	if !ok {
		c.RUnlock()
		return nil
	}
	c.RUnlock()
	return v
}

func (c *Collection) Range(iter func(k interface{}) bool) {
	panic("implement me")
}

func (c *Collection) Save(w io.Writer) {
	c.Lock()
	defer c.Unlock()
	b, err := json.Marshal(c.d)
	if err != nil {
		log.Fatalf("error marshaling collection: %v\n", err)
	}
	_, err = w.Write(b)
	if err != nil {
		log.Fatalf("error saving collection: %v\n", err)
	}
}

func (c *Collection) Load(r io.Reader) {
	c.Lock()
	defer c.Unlock()
	var b []byte
	_, err := r.Read(b)
	if err != nil {
		log.Fatalf("error loading collection: %v\n", err)
	}
	err = json.Unmarshal(b, &c.d)
	if err != nil {
		log.Fatalf("error unmarshaling into collection: %v\n", err)
	}
}

func NewCollection() *Collection {
	return &Collection{
		d: make(map[string]interface{}, 0),
	}
}

type MapStore struct {
	d map[string]*Collection
}

func NewMapStore() *MapStore {
	return &MapStore{
		d: make(map[string]*Collection, 0),
	}
}

func (m *MapStore) AddCollection(c string) {
	if _, ok := m.d[c]; ok {
		return
	}
	m.d[c] = NewCollection()
}

func (m *MapStore) GetCollection(c string) *Collection {
	if c, ok := m.d[c]; ok {
		return c
	}
	return nil
}

func (m *MapStore) DelCollection(c string) {
	if _, ok := m.d[c]; ok {
		delete(m.d, c)
	}
	return
}
