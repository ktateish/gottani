package gottani_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ktateish/gottani"
)

func TestCombine(t *testing.T) {
	testCases := []string{
		// examples
		"examples/01-simple",
		"examples/02-simple",
		"examples/03-simple",
		"examples/04-thirdparty",
		"examples/05-renaming",
		"examples/06-initializers",
		"examples/07-methods",
		"examples/08-cgo",

		// testdata
		"testdata/issue2",
		"testdata/issue3",
		"testdata/issue4",
		"testdata/issue5",
		"testdata/issue6",
		"testdata/issue9",
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get working dir: %s", err)
	}
	for _, tc := range testCases {
		dir := tc
		t.Run(dir, func(t *testing.T) {
			// reset to cwd
			if err := os.Chdir(cwd); err != nil {
				t.Fatalf("Failed to enter directory: %s: %s", cwd, err)
			}

			// prepare want result
			if err := os.Chdir(dir); err != nil {
				t.Fatalf("Failed to enter directory: %s: %s", dir, err)
			}
			wantSrcPath := "combined.go"
			wantSrc, err := ioutil.ReadFile(wantSrcPath)
			if err != nil {
				abs, err := filepath.Abs(wantSrcPath)
				if err != nil {
					abs = filepath.Join(cwd, dir, wantSrcPath)
				}
				t.Fatalf("Failed to read file: %s: %s", abs, err)
			}

			// do Compbine()
			srcDir := "src"
			gotSrc, err := gottani.Combine(srcDir, "main")
			if err != nil {
				abs, erra := filepath.Abs(srcDir)
				if erra != nil {
					abs = filepath.Join(cwd, dir, srcDir)
				}
				t.Fatalf("Failed to Combine(): %s: %s", abs, err.Error())
			}

			gotResult, err := run(gotSrc)
			if err != nil {
				t.Fatalf("Failed to run combined source: %s", err.Error())
			}

			wantResult, err := run(wantSrc)
			if err != nil {
				t.Fatalf("Failed to run the properly combined source: %s", err.Error())
			}

			// check the exec result
			if !reflect.DeepEqual(gotResult, wantResult) {
				t.Fatalf("Result of running the combined source is wrong: %s", dir)
			}

			// check the combined source code
			if !reflect.DeepEqual(gotSrc, wantSrc) {
				t.Fatalf("Combined source is wrong: %s", dir)
			}
		})
	}
}

func compile(src []byte) (string, error) {
	sf, err := ioutil.TempFile("", "gottani-test-combined-*.go")
	if err != nil {
		return "", fmt.Errorf("creating source file: %w", err)
	}
	defer os.Remove(sf.Name())
	_, err = sf.Write(src)
	if err != nil {
		return "", fmt.Errorf("writing source file: %w", err)
	}
	err = sf.Close()
	if err != nil {
		return "", fmt.Errorf("closing source file: %w", err)
	}

	ef, err := ioutil.TempFile("", "gottani-test-combined-*.exe")
	if err != nil {
		return "", fmt.Errorf("creating exec file: %w", err)
	}
	err = ef.Close()
	if err != nil {
		return "", fmt.Errorf("closing exec file: %w", err)
	}

	cmd := exec.Command("go", "build", "-o", ef.Name(), sf.Name())
	err = cmd.Run()

	return ef.Name(), err
}

func run(src []byte) ([]byte, error) {
	executable, err := compile(src)
	if err != nil {
		return nil, fmt.Errorf("compiling: %w", err)
	}
	defer os.Remove(executable)

	cmd := exec.Command(executable)
	out, err := cmd.Output()
	if err != nil {
		return out, fmt.Errorf("executing: %w", err)
	}
	return out, nil
}
