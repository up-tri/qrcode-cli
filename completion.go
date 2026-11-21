package main

import (
	_ "embed"
	"fmt"
	"io"
)

//go:embed completion.zsh
var zshCompletion string

func writeCompletion(shell string, stdout io.Writer) error {
	switch shell {
	case "zsh":
		_, err := io.WriteString(stdout, zshCompletion)
		return err
	}
	return fmt.Errorf("未対応のシェルです: %q (対応: zsh)", shell)
}
