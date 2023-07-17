package sdk

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
