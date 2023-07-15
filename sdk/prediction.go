package sdk

import (
	"errors"
	"time"
	"unsafe"

	"github.com/taubyte/go-sdk/utils/codec"
)

type Prediction uint32

type Params struct {
	Seed        int64
	TopK        int64
	TopP        float64
	StopWords   []string
	Tokens      int32
	Temperature float64
	Penalty     float64
	Repeat      int32
	Batch       int32
	NKeep       int32
}

type Option func(*Params)

func With(p *Params) Option {
	return func(p0 *Params) {
		*p0 = *p
	}
}

// SetSeed sets the random seed for sampling text generation.
func WithSeed(s int64) Option {
	return func(p *Params) {
		p.Seed = s
	}
}

// SetTopK sets the value for top-K sampling.
func WithTopK(topK int64) Option {
	return func(p *Params) {
		p.TopK = topK
	}
}

// SetTopP sets the value for nucleus sampling.
func WithTopP(topP float64) Option {
	return func(p *Params) {
		p.TopP = topP
	}
}

// SetStopWords sets the prompts that will stop predictions.
func WithStopWords(stop ...string) Option {
	return func(p *Params) {
		p.StopWords = stop
	}
}

// SetTokens sets the number of tokens to generate.
func WithTokens(tokens int) Option {
	return func(p *Params) {
		p.Tokens = int32(tokens)
	}
}

// SetTemperature sets the temperature value for text generation.
func WithTemperature(temp float64) Option {
	return func(p *Params) {
		p.Temperature = temp
	}
}

// SetPenalty sets the repetition penalty for text generation.
func WithPenalty(penalty float64) Option {
	return func(p *Params) {
		p.Penalty = penalty
	}
}

// SetRepeat sets the number of times to repeat text generation.
func WithRepeat(repeat int32) Option {
	return func(p *Params) {
		p.Repeat = repeat
	}
}

// SetBatch sets the batch size.
func WithBatch(size int32) Option {
	return func(p *Params) {
		p.Batch = size
	}
}

// SetKeep sets the number of tokens from initial prompt to keep.
func WithNKeep(n int32) Option {
	return func(p *Params) {
		p.NKeep = n
	}
}

func Predict(text string, options ...Option) (Prediction, error) {
	p := &Params{
		StopWords: make([]string, 0),
	}
	for _, opt := range options {
		opt(p)
	}

	var stopWordsSliceEncoded []byte
	if err := codec.Convert(p.StopWords).To(&stopWordsSliceEncoded); err != nil {
		return 0, err
	}

	var pred Prediction

	valByteString := []byte(text)
	if _err := llama_predict(
		&valByteString[0],
		uint32(len(valByteString)),
		p.Seed,
		p.TopK,
		p.TopP,
		(*byte)(unsafe.Pointer(&stopWordsSliceEncoded)),
		uint32(len(stopWordsSliceEncoded)),
		p.Tokens,
		p.Temperature,
		p.Penalty,
		p.Repeat,
		p.Batch,
		p.NKeep,
		&pred,
	); _err != ErrorNone {
		return 0, _err.Error()
	}

	return pred, nil
}

func (p Prediction) NextWithTimeout(timeout time.Duration) (string, error) {
	buf := make([]byte, MaxTokenSize)
	var written uint32
	_err := llama_token(p, &buf[0], uint32(MaxTokenSize), &written, int64(timeout))
	if _err == ErrorWith {
		return "", errors.New(string(buf[:written]))
	}
	return string(buf[:written]), _err.Error()
}

func (p Prediction) Next() (string, error) {
	return p.NextWithTimeout(0)
}

func (p Prediction) Cancel() error {
	return llama_stop(p).Error()
}
