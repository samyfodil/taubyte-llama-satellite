package sdk

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
