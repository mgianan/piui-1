package main

import (
	"net/http/httputil"
	"sync"
)

type bytePool struct {
	pool sync.Pool
}

var _ httputil.BufferPool = &bytePool{}

func newBytePool(size int) *bytePool {
	return &bytePool{
		sync.Pool{
			New: func() interface{} {
				return make([]byte, size)
			},
		},
	}
}

func (p *bytePool) Get() []byte {
	return p.pool.Get().([]byte)
}

func (p *bytePool) Put(bs []byte) {
	p.pool.Put(bs)
}
