package main

import (
	"context"
	"time"

	gollama "github.com/go-skynet/go-llama.cpp"
	"github.com/samyfodil/taubyte-llama-satellite/sdk"
)

func New(ctx context.Context, model, filename string) (*llama, error) {
	l := &llama{
		model:     model,
		modelFile: filename,
		requests:  make(chan *request, RequestsQueueSize),
		responses: make(map[uint32]*response),
		ready:     make(chan error, WorkersCount),
	}
	l.ctx, l.ctxC = context.WithCancel(ctx)
	return l, l.start(WorkersCount)
}

func (l *llama) push(ctx context.Context, req *request) (uint32, sdk.Error) {
	select {
	case l.requests <- req:
		response := l.allocateResponse(ctx)
		req.response = response
		return response.id, sdk.ErrorNone
	default:
		return 0, sdk.ErrorMaximumCapacity
	}
}

func (l *llama) allocateResponse(ctx context.Context) *response {
	l.responsesLock.Lock()
	defer l.responsesLock.Unlock()
	l.responsesLastest++
	r := &response{
		id:     l.responsesLastest,
		stream: make(chan string, TokenBufferSize),
	}
	r.ctx, r.ctxC = context.WithTimeout(ctx, 5*time.Minute)
	l.responses[r.id] = r
	return r
}

func (l *llama) cleanupResponse(id uint32) {
	l.responsesLock.Lock()
	defer l.responsesLock.Unlock()
	delete(l.responses, id)
}

func (l *llama) getResponse(id uint32) *response {
	l.responsesLock.RLock()
	defer l.responsesLock.RUnlock()
	return l.responses[id]
}

func toGoLlamaCppOptions(p *sdk.Params) []gollama.PredictOption {
	opts := make([]gollama.PredictOption, 0, 16)

	if p.Seed != 0 {
		opts = append(opts, gollama.SetSeed(int(p.Seed)))
	}

	if p.TopK != 0 {
		opts = append(opts, gollama.SetTopK(int(p.TopK)))
	}

	if p.TopP != 0 {
		opts = append(opts, gollama.SetTopP(p.TopP))
	}

	if len(p.StopWords) != 0 {
		opts = append(opts, gollama.SetStopWords(p.StopWords...))
	}

	if p.Tokens < 1 {
		p.Tokens = int32(DefaultTokens)
	}
	opts = append(opts, gollama.SetTokens(int(p.Tokens)))

	if p.Temperature != 0 {
		opts = append(opts, gollama.SetTemperature(p.Temperature))
	}

	if p.Penalty != 0 {
		opts = append(opts, gollama.SetPenalty(p.Penalty))
	}

	if p.Repeat != 0 {
		opts = append(opts, gollama.SetRepeat(int(p.Repeat)))
	}

	if p.Batch != 0 {
		opts = append(opts, gollama.SetBatch(int(p.Batch)))
	}

	if p.NKeep != 0 {
		opts = append(opts, gollama.SetNKeep(int(p.NKeep)))
	}

	return opts
}

func (l *llama) startModelWorker() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()

		ai, err := gollama.New(
			l.modelFile, gollama.SetContext(ContextSize),
			gollama.SetGPULayers(GPULayers),
			gollama.EnableF16Memory,
		)
		if err != nil {
			l.ready <- err
			return
		}
		defer ai.Free()

		l.ready <- nil

		for req := range l.requests {
			select {
			case <-l.ctx.Done():
				return
			default:
				select {
				case <-req.response.ctx.Done():
					l.cleanupResponse(req.response.id)
				default:
					stream := req.response.stream
					opts := append(
						toGoLlamaCppOptions(&req.params),
						gollama.EnableF16KV,
						gollama.SetThreads(PredictionThreads),
						gollama.SetTokenCallback(
							func(token string) bool {
								select {
								case <-req.response.ctx.Done():
									return false
								default:
									stream <- token
									return true
								}
							}),
					)
					if Debug {
						opts = append(opts, gollama.Debug)
					}
					_, err := ai.Predict(
						req.text,
						opts...,
					)
					if err != nil {
						req.response.err = err
					}
					close(stream)
				}
			}
		}
	}()
}

func (l *llama) start(workers int) error {
	l.wg.Add(workers)
	for i := 0; i < workers; i++ {
		l.startModelWorker()
	}

	for i := 0; i < workers; i++ {
		if err := <-l.ready; err != nil {
			return err
		}
	}

	return nil
}

func (l *llama) Wait() {
	l.wg.Wait()
}

func (l *llama) Shutdown() {
	l.ctxC()
	l.Wait()
}

func (l *llama) Kill() {
	l.ctxC()
}
