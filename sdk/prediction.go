package sdk

import (
	"errors"
	"time"
	"unsafe"

	"github.com/taubyte/go-sdk/utils/codec"
)

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
