package main

import (
	"context"
	"time"

	"github.com/samyfodil/taubyte-llama-satellite/sdk"
	"github.com/taubyte/vm-orbit/satellite"
)

func (l *llama) W_stop(
	ctx context.Context, module satellite.Module,
	id uint32,
) uint32 {
	res := l.getResponse(id)
	if res == nil {
		return uint32(sdk.ErrorPredictionNotFound)
	}
	res.ctxC()
	return uint32(sdk.ErrorNone)
}

func (l *llama) W_token(
	ctx context.Context, module satellite.Module,
	id uint32,
	bufferPtr uint32, bufferSize uint32, // buffer for the token
	writtenPtr uint32,
	_ttl int64, // how many nanoseconds to wait
) uint32 {
	res := l.getResponse(id)
	if res == nil {
		return uint32(sdk.ErrorPredictionNotFound)
	}
	ttl := 5 * time.Minute
	if _ttl > 0 {
		ttl = time.Duration(_ttl)
	}
	cleanup := false
	defer func() {
		if cleanup {
			res.ctxC()
			l.cleanupResponse(id)
		}
	}()
	select {
	case <-res.ctx.Done():
		cleanup = true
		return uint32(sdk.ErrorPredictionStopped)
	case <-ctx.Done():
		cleanup = true
		return uint32(sdk.ErrorPredictionKilled)
	case <-time.After(ttl):
		return uint32(sdk.ErrorPredictionTimeout)
	case token, ok := <-res.stream:
		if !ok {
			cleanup = true
			if res.err != nil {
				data := []byte(res.err.Error())
				written := len(data)
				if len(data) > int(bufferSize) {
					written = int(bufferSize)
				}
				module.MemoryWrite(bufferPtr, data[:written])
				module.WriteUint32(writtenPtr, uint32(written))
				return uint32(sdk.ErrorWith)
			} else {
				module.WriteUint32(writtenPtr, 0)
				return uint32(sdk.ErrorEOF)
			}
		} else {
			data := []byte(token)
			written := len(data)
			if len(data) > int(bufferSize) {
				written = int(bufferSize)
			}
			module.MemoryWrite(bufferPtr, data[:written])
			module.WriteUint32(writtenPtr, uint32(written))
			return uint32(sdk.ErrorNone)
		}
	}
}

func (l *llama) W_predict(
	ctx context.Context, module satellite.Module,
	textPtr uint32, textLen uint32, // text to process
	seed int64,
	topK int64,
	topP float64,
	stopWordsSlicePtr,
	stopWordsSliceSize uint32,
	tokens int32,
	temperature float64,
	penalty float64,
	repeat int32,
	batch int32,
	nKeep int32,
	predictIdPtr uint32,
) uint32 {
	var err error

	stopWords, err := module.ReadStringSlice(stopWordsSlicePtr, stopWordsSliceSize)
	if err != nil {
		return uint32(sdk.ErrorParsingStopPrompt)
	}

	req := &request{
		params: sdk.Params{
			Seed:        seed,
			TopK:        topK,
			TopP:        topP,
			StopWords:   stopWords,
			Tokens:      tokens,
			Temperature: temperature,
			Penalty:     penalty,
			Repeat:      repeat,
			Batch:       batch,
			NKeep:       nKeep,
		},
	}

	req.text, err = module.ReadString(textPtr, textLen)
	if err != nil {
		return uint32(sdk.ErrorGettingPredictionText)
	}

	id, _err := l.push(l.ctx, req)
	if _err != sdk.ErrorNone {
		return uint32(_err)
	}

	if _, err := module.WriteUint32(predictIdPtr, id); err != nil {
		l.cleanupResponse(id)
		return uint32(sdk.ErrorReturnId)
	}

	return uint32(sdk.ErrorNone)
}
