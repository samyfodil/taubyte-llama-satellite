package main

import (
	"context"
	"sync"

	"github.com/samyfodil/taubyte-llama-satellite/sdk"
)

type response struct {
	ctx  context.Context
	ctxC context.CancelFunc

	id     uint32
	err    error
	stream chan string
}

type request struct {
	text string

	params sdk.Params

	response *response
}

type llama struct {
	ctx              context.Context
	ctxC             context.CancelFunc
	wg               sync.WaitGroup
	ready            chan error
	model            string
	modelFile        string
	requests         chan *request
	responses        map[uint32]*response
	responsesLastest uint32
	responsesLock    sync.RWMutex
}
