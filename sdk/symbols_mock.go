//go:build !wasi && !wasm

package sdk

func llama_predict(
	text *byte, textLen uint32,
	seed int64,
	topK int64,
	topP float64,
	stopWordsSlicePtr *byte, stopWordsSliceSize uint32,
	tokens int32,
	temperature float64,
	penalty float64,
	repeat int32,
	batch int32,
	nKeep int32,
	predictionId *Prediction,
) Error {
	return 0
}

func llama_token(id Prediction, buf *byte, bufSize uint32, writtenPtr *uint32, ttl int64) Error {
	return 0
}

func llama_stop(id Prediction) Error {
	return 0
}
