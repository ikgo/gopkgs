package gopkgs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

func Packages111(opts Options) (map[string]Pkg, error) {
	ctx := context.Background()
	out := new(bytes.Buffer)
	//args := []string{"list", "-e", "-json", "-compiled", "-deps=false", "all"}
	args := []string{"list", "-e", "-json", "-deps=false", "all"}
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Env = os.Environ()
	cmd.Dir = opts.WorkDir
	cmd.Stdout = out
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			return nil, fmt.Errorf("couldn't exec 'go list': %s %T", err, err)
		}
		return nil, fmt.Errorf("go list: %s: %s", exitErr, cmd.Stderr)
	}
	if stderr := fmt.Sprint(cmd.Stderr); stderr != "" {
		fmt.Fprintf(os.Stderr, "go list stderr <<%s>>\n", stderr)
	}
	result, err := parsePackageList(out)
	return result, err
}

type jsonPackage struct {
	Dir        string
	ImportPath string
	Name       string
}

func parsePackageList(buf *bytes.Buffer) (map[string]Pkg, error) {
	result := make(map[string]Pkg)
	for dec := json.NewDecoder(buf); dec.More(); {
		pkg := new(Pkg)
		if err := dec.Decode(pkg); err != nil {
			return nil, fmt.Errorf("JSON decoding failed: %v", err)
		}
		if len(pkg.Name) == 0 {
			// bad package
			continue
		}
		if _, found := result[pkg.Dir]; found {
			continue
		}
		// result
		result[pkg.Dir] = *pkg
	}
	return result, nil
}
