![Project Banner](recording.gif)

# LLAMA Satellite for Taubyte WebAssembly VM

Welcome to the LLAMA Satellite for Taubyte WebAssembly VM project. This tool extends the functionality of the Taubyte WebAssembly Virtual Machine by introducing Large Language Model (LLAMA) capabilities. It's built upon `llama.cpp` and employs `go-llama-cpp` for cgo bindings. 

With this plugin, you can augment your applications with advanced language understanding features.

## Table of Contents

- [File Structure](#file-structure)
- [Example Usage from WebAssembly](#example-usage-from-webassembly)
- [Installation](#installation)
- [Acceleration](#acceleration)
- [Model Setup](#model-setup)
- [Plugin Compilation](#plugin-compilation)
- [WebAssembly](#webassembly)
- [Testing](#testing)

## File Structure
- `plugin/` - Code for the plugin itself
- `sdk/` - Wrapper around the low-level functions exported by the plugin
- `fixtures/build/` - Code to be compiled to webassembly and run on a Taubyte Virtual Machine during testing
- `models/` - Helper tool to download models

## Example Usage from WebAssembly
Using the plugin is straightforward. In just a few lines of code, you can build your own planet-scale ChatGPT clone-API! Here's a simple example:

```go
package lib

import (
	"fmt"
	"io"
	"github.com/samyfodil/taubyte-llama-satellite/sdk"
)

//export wapredict
func wapredict(uint32) uint32 {
	p, err := sdk.Predict(
		"How old is the universe?",
		sdk.WithTopK(90),
		sdk.WithTopP(0.86),
		sdk.WithBatch(5),
	)
	if err != nil {
		panic(err)
	}

	for {
		token, err := p.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Print(token)
	}
	return 0
}
```

## Installation

This project requires some dependencies to function. Use the following commands to clone the submodules locally:

```bash
git clone --recurse-submodules deps/go-llama
```

Then, navigate to the newly cloned directory and make the libbinding:

```bash
cd deps/go-llama
make libbinding.a
```

## Acceleration

You can take advantage of OpenBLAS and CuBLAS for acceleration.

### OpenBLAS

To build and run with OpenBLAS:

```bash
cd deps/go-llama
BUILD_TYPE=openblas make libbinding.a
```

### CuBLAS

To build with CuBLAS:

```bash
cd deps/go-llama
BUILD_TYPE=cublas make libbinding.a
```

## Model Setup

You need to provide the plugin with a model to load. If you do not have one, you can use the tool in the `models` folder:

```bash
cd models
go run .
```

Follow the prompts to select and download the model. Then, specify the path to the model in `plugin/main.go`. For example:

```go
ai, err := New(ctx, "orca-mini", "models/assets/orca-mini-3b.ggmlv3.q4_0.bin")
```

## Plugin Compilation

### Without Acceleration

```bash
cd plugin
go build .
```

### With OpenBLAS

```bash
cd plugin
go build -tags openblas .
```

### With CuBLAS

```bash
cd plugin
go build -tags cublas .
```

## WebAssembly

The WebAssembly code used to test the plugin is in `fixtures/build`. If you modify `predict.go`, the tests will automatically recompile it.

## Testing

Once you have compiled the bindings, you can run the tests:

```bash
go test -v
```


## License

This project is licensed under the BSD 3-Clause License. For more details, see the LICENSE file.
