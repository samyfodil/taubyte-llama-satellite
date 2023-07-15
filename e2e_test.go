package tests_test

import (
	goContext "context"
	"fmt"
	"io"
	"sync"
	"testing"

	"gotest.tools/v3/assert"
)

func init() {
	if err := initializeAssetPaths(); err != nil {
		panic(err)
	}

	if err := initializeWasm("predict"); err != nil {
		panic(err)
	}

	if err := buildPlugin("cublas"); err != nil {
		panic(err)
	}
}

func TestPredict(t *testing.T) {

	ctx := goContext.Background()
	tvm := newTVM(ctx)

	plugin := plugin(t, ctx)
	defer plugin.Close()

	instance, err := newVM(ctx, tvm)
	assert.NilError(t, err)

	rt, err := instance.Runtime(nil)
	assert.NilError(t, err)

	_, _, err = rt.Attach(plugin)
	assert.NilError(t, err)

	fi := getFunction(t, "predict", rt, plugin)

	go func() {
		reader := io.MultiReader(rt.Stdout(), rt.Stderr())
		p := make([]byte, 1024)
		for {
			n, err := reader.Read(p)
			if err == io.EOF {
				break
			}
			fmt.Print(string(p[:n]))
		}
	}()

	ret := fi.Call(ctx, 0)

	assert.NilError(t, ret.Error())

}

func TestParallelPredict(t *testing.T) {
	threads := 4

	ctx := goContext.Background()
	tvm := newTVM(ctx)

	plugin := plugin(t, ctx)
	defer plugin.Close()

	var wg sync.WaitGroup
	wg.Add(threads)

	for i := 0; i < threads; i++ {
		go func() {
			defer wg.Done()
			instance, err := newVM(ctx, tvm)
			assert.NilError(t, err)

			rt, err := instance.Runtime(nil)
			assert.NilError(t, err)

			_, _, err = rt.Attach(plugin)
			assert.NilError(t, err)

			fi := getFunction(t, "predict", rt, plugin)

			go func() {
				reader := io.MultiReader(rt.Stdout(), rt.Stderr())
				p := make([]byte, 1024)
				for {
					n, err := reader.Read(p)
					if err == io.EOF {
						break
					}
					fmt.Print(string(p[:n]))
				}
			}()

			ret := fi.Call(ctx, 0)

			assert.NilError(t, ret.Error())

		}()
	}

	wg.Wait()
}
