package tests_test

import (
	"bytes"
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
	defer instance.Close()

	rt, err := instance.Runtime(nil)
	assert.NilError(t, err)
	defer rt.Close()

	_, _, err = rt.Attach(plugin)
	assert.NilError(t, err)

	fi := getFunction(t, "predict", rt, plugin)

	var buf bytes.Buffer
	go func() {
		reader := io.MultiReader(io.TeeReader(rt.Stdout(), &buf), rt.Stderr())
		p := make([]byte, 1024)
		for {
			n, err := reader.Read(p)
			if err != nil {
				if n > 0 {
					fmt.Print(string(p[:n]))
				}
				return
			}
			fmt.Print(string(p[:n]))
		}
	}()

	ret := fi.Call(ctx, 0)

	assert.NilError(t, ret.Error())

	fmt.Println(
		"\n--\n",
		buf.String(),
	)

}

func TestParallelPredict(t *testing.T) {
	threads := 8

	ctx := goContext.Background()
	tvm := newTVM(ctx)

	plugin := plugin(t, ctx)
	defer plugin.Close()

	var wg sync.WaitGroup
	wg.Add(threads)

	bufs := make([]bytes.Buffer, threads)
	for i := 0; i < threads; i++ {
		go func(i int) {
			defer wg.Done()
			instance, err := newVM(ctx, tvm)
			assert.NilError(t, err)
			defer instance.Close()

			rt, err := instance.Runtime(nil)
			assert.NilError(t, err)
			defer rt.Close()

			_, _, err = rt.Attach(plugin)
			assert.NilError(t, err)

			fi := getFunction(t, "predict", rt, plugin)

			wg.Add(1)
			go func() {
				defer wg.Done()
				reader := io.MultiReader(io.TeeReader(rt.Stdout(), &bufs[i]), rt.Stderr())
				p := make([]byte, 1024)
				for {
					n, err := reader.Read(p)
					if err != nil {
						if n > 0 {
							fmt.Print(string(p[:n]))
						}
						return
					}
					fmt.Print(string(p[:n]))
				}
			}()

			ret := fi.Call(ctx, 1+(i%4))
			assert.NilError(t, ret.Error())
		}(i)
	}

	wg.Wait()

	for i := 0; i < threads; i++ {
		fmt.Println()
		fmt.Println(bufs[i].String())
		fmt.Println()
	}
}
