package main

import (
	"context"

	"github.com/taubyte/vm-orbit/plugin"
)

func main() {
	ctx := context.Background()
	ai, err := New(ctx, "gpt4all", "/home/samy/Downloads/GPT4All-13B-snoozy.ggmlv3.q2_K.bin")
	if err != nil {
		panic(err)
	}
	plugin.Export("llama", ai)
}
