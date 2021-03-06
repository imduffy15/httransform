package ca

import (
	"hash/fnv"
	"sync"

	"github.com/karlseguin/ccache"
)

var (
	signRequestPool = sync.Pool{
		New: func() interface{} {
			return &signRequest{
				response: make(chan *signResponse),
			}
		},
	}

	signResponsePool = sync.Pool{
		New: func() interface{} {
			return &signResponse{}
		},
	}

	hashPool = sync.Pool{
		New: func() interface{} {
			return fnv.New32a()
		},
	}
)

type signRequest struct {
	host     string
	response chan *signResponse
}

type signResponse struct {
	item ccache.TrackedItem
	err  error
}
