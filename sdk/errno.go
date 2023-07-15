package sdk

import (
	"errors"
	"io"
)

type Error uint32

var errorStrings = []error{
	/*ErrorNone*/ nil,
	/*ErrorPredictionNotFound:*/ errors.New("prediction was not found"),
	/*ErrorPredictionStopped:*/ errors.New("prediction was stopped by user"),
	/*ErrorPredictionKilled:*/ errors.New("prediction was killed by vm"),
	/*ErrorPredictionTimeout:*/ errors.New("prediction timed out"),
	/*ErrorParsingStopPrompt:*/ errors.New("parsing stop prompt failed"),
	/*ErrorGettingPredictionText:*/ errors.New("failed getting prediction text"),
	/*ErrorReturnId:*/ errors.New("failed to return prediction id"),
	/*ErrorMaximumCapacity:*/ errors.New("reached maximum capacity"),
	/*ErrorEOF:*/ io.EOF,
	/*ErrorWith:*/ errors.New("failed with an error"),
}

func (e Error) String() string {
	return e.Error().Error()
}

func (e Error) Error() error {
	return errorStrings[e]
}

const (
	ErrorNone Error = iota
	ErrorPredictionNotFound
	ErrorPredictionStopped
	ErrorPredictionKilled
	ErrorPredictionTimeout
	ErrorParsingStopPrompt
	ErrorGettingPredictionText
	ErrorReturnId
	ErrorMaximumCapacity
	ErrorEOF
	ErrorWith
)
