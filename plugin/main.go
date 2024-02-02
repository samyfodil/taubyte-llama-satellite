package main

import (
	"context"

	"github.com/taubyte/vm-orbit/plugin"
)

func main() {
	ctx := context.Background()

	ai, err := New(ctx, "orca-mini", "/tb/plugins/model.bin" /* "/home/samy/Documents/taubyte/RD/llama/llama-2-7b-chat/ggml-model-f32-q4_0.bin" /*"/tb/plugins/model.bin"*/)
	if err != nil {
		panic(err)
	}

	plugin.Export("llama", ai)
}
