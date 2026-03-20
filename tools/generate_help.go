//go:generate go run generate_help.go

package main

import (
	"os"
	"os/exec"
)

func main() {
	clis := []struct {
		cmd string
		out string
	}{
		{"../cmd/xair-cli/", "xair-help.md"},
		{"../cmd/x32-cli/", "x32-help.md"},
	}

	for _, cli := range clis {
		helpCmd := exec.Command("go", "run", cli.cmd, "--help")
		out, err := helpCmd.Output()
		if err != nil {
			panic(err)
		}

		// Wrap output in markdown console code block
		wrapped := append([]byte("```console\n"), out...)
		wrapped = append(wrapped, []byte("\n```\n")...)
		os.WriteFile(cli.out, wrapped, 0o644)
	}
}
