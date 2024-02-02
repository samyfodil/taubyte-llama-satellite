package main

var (
	Debug             = false
	WorkersCount      = 3
	TokenBufferSize   = 2 * 1024
	PredictionThreads = 4
	GPULayers         = 64
	ContextSize       = 512
	RequestsQueueSize = 512
	DefaultTokens     = 512
)
