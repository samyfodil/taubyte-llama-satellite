package tests_test

import (
	goContext "context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/otiai10/copy"
	builder "github.com/taubyte/builder"
	"github.com/taubyte/go-interfaces/services/tns/mocks"
	"github.com/taubyte/go-interfaces/vm"
	"github.com/taubyte/utils/id"
	vmPlugin "github.com/taubyte/vm-orbit/plugin/vm"
	fileBE "github.com/taubyte/vm/backend/file"
	"github.com/taubyte/vm/context"
	loader "github.com/taubyte/vm/loaders/wazero"
	resolver "github.com/taubyte/vm/resolvers/taubyte"
	service "github.com/taubyte/vm/service/wazero"
	source "github.com/taubyte/vm/sources/taubyte"
	"gotest.tools/v3/assert"
)

var (
	wd         string
	assetDir   string
	taubyteDir = ".taubyte"
	goMod      = "go.mod"
	buildDir   string
)

func initializeAssetPaths() (err error) {
	if wd, err = os.Getwd(); err != nil {
		return
	}

	assetDir = path.Join(wd, "fixtures")
	buildDir = path.Join(assetDir, "build")

	return
}

func goExtension(fileName string) string {
	return fileName + ".go"
}

func initializeWasm(fileName string) error {
	wasmPath := path.Join(assetDir, fileName+".wasm")
	goFile := goExtension(fileName)

	wasmStat, err := os.Stat(path.Join(assetDir, "predict.wasm"))
	if err == nil {
		codeStat, err := os.Stat(path.Join(buildDir, goFile))
		if err == nil && codeStat.ModTime().Compare(wasmStat.ModTime()) < 0 {
			// no need to rebuild
			fmt.Println("Skipping WebAssembly build")
			return nil
		}
	}

	tempDir, err := os.MkdirTemp("/tmp", "*")
	if err != nil {
		return fmt.Errorf("creating temp dir failed with: %w", err)
	}

	buildSrcDir := path.Join(tempDir, "wasm/src")
	err = os.MkdirAll(buildSrcDir, 0755)
	if err != nil {
		return fmt.Errorf("creating wasm/src dir failed with: %w", err)
	}

	if err = copy.Copy(path.Join(buildDir, goFile), path.Join(buildSrcDir, goFile)); err != nil {
		return fmt.Errorf("copying `%s` failed with: %w", goFile, err)
	}

	if err = copy.Copy(path.Join(buildDir, goMod), path.Join(buildSrcDir, goMod)); err != nil {
		return fmt.Errorf("copying go.mod failed with: %w", err)
	}

	if err = copy.Copy(path.Join(buildDir, taubyteDir), path.Join(tempDir, taubyteDir)); err != nil {
		return fmt.Errorf("copying taubyteDir failed with: %w", err)
	}

	err = os.MkdirAll(path.Join(tempDir, "sdk"), 0755)
	if err != nil {
		return fmt.Errorf("creating sdk dir failed with: %w", err)
	}
	if err = copy.Copy("sdk", path.Join(tempDir, "sdk/sdk")); err != nil {
		return fmt.Errorf("copying sdk failed with: %w", err)
	}

	if err = copy.Copy("go.mod", path.Join(tempDir, "sdk/go.mod")); err != nil {
		return fmt.Errorf("copying sdk/go.mod failed with: %w", err)
	}

	_builder, err := builder.New(goContext.TODO(), tempDir)
	if err != nil {
		return fmt.Errorf("creating new builder failed with: %w", err)
	}

	out, err := _builder.Build()
	if err != nil {
		return fmt.Errorf("builder.Build failed with: %w", err)
	}

	if err := copy.Copy(path.Join(out.OutDir(), "artifact.wasm"), wasmPath); err != nil {
		return fmt.Errorf("copying wasm build failed with: %w", err)
	}

	return nil
}

func plugin(t *testing.T, ctx goContext.Context) vm.Plugin {
	wd, err := os.Getwd()
	assert.NilError(t, err)

	pluginBinary := path.Join(wd, "plugin", "plugin")
	_plugin, err := vmPlugin.Load(pluginBinary, ctx)
	assert.NilError(t, err)

	return _plugin
}

func getFunction(t *testing.T, wasmFile string, rt vm.Runtime, plugin vm.Plugin) vm.FunctionInstance {

	wasmFile = path.Join(assetDir, wasmFile+".wasm")

	mod, err := rt.Module("/file/" + wasmFile)
	assert.NilError(t, err)

	fi, err := mod.Function("wapredict")
	assert.NilError(t, err)

	return fi
}

func newTVM(ctx goContext.Context) vm.Service {
	tns := mocks.New()
	rslver := resolver.New(tns)
	ldr := loader.New(rslver, fileBE.New())
	src := source.New(ldr)
	return service.New(ctx, src)
}

func newVM(ctx goContext.Context, vmService vm.Service) (vm.Instance, error) {

	mocksConfig := mocks.InjectConfig{
		Branch:      "master",
		Commit:      "head_commit",
		Project:     id.Generate(),
		Application: id.Generate(),
		Cid:         id.Generate(),
	}

	_ctx, err := context.New(
		ctx,
		context.Application(mocksConfig.Application),
		context.Project(mocksConfig.Project),
		context.Resource(mocksConfig.Cid),
		context.Branch(mocksConfig.Branch),
		context.Commit(mocksConfig.Commit),
	)
	if err != nil {
		return nil, err
	}

	instance, err := vmService.New(_ctx, vm.Config{})
	if err != nil {
		return nil, err
	}

	return instance, err
}

func buildPlugin(buildTags ...string) error {
	fmt.Println("(re)Build plugin")
	pluginDir := "./plugin"
	args := []string{"build"}
	if len(buildTags) != 0 {
		args = append(args, "-tags", strings.Join(buildTags, ","))
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = pluginDir

	return cmd.Run()
}
