//go:build wasi || wasm

package sdk

//go:wasm-module llama
//export predict
func llama_predict(
	text *byte, textLen uint32, // text to process
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
) Error

//go:wasm-module llama
//export token
func llama_token(id Prediction, buf *byte, bufSize uint32, writtenPtr *uint32, ttl int64) Error

//go:wasm-module llama
//export stop
func llama_stop(id Prediction) Error
