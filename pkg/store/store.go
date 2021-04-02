package store

import "io"

type Store interface {
	Has(k interface{}) bool
	Set(k, v interface{})
	Get(k interface{}) interface{}
	Range(iter func(k interface{}) bool)
	Save(w io.Writer)
	Load(r io.Reader)
}
