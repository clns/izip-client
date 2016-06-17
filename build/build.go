package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/clns/izip-client/cmd"
)

type output struct {
	GOOS   string
	GOARCH string
	File   string
}

var outputs = []*output{
	&output{"windows", "amd64", cmd.NAME + "-Windows-x86_64.exe"},
	&output{"linux", "amd64", cmd.NAME + "-Linux-x86_64"},
	&output{"darwin", "amd64", cmd.NAME + "-Darwin-x86_64"},
}

func main() {
	for _, o := range outputs {
		vars := []string{"GOOS=" + o.GOOS, "GOARCH=" + o.GOARCH}
		fmt.Fprintf(os.Stdout, "%s go build -o build/%s main.go ...", strings.Join(vars, " "), o.File)
		cmd := exec.Command("go", "build", "-o", filepath.Join("build", o.File), "main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		env := os.Environ()
		env = append(env, vars...)
		cmd.Env = env
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Fprintln(os.Stdout, "done")
	}
}
