package main

import (
	"flag"
	"io"
	"os"
)

// parseInterspersed はオプションと位置引数が混在していても解析できるようにする。
// 標準の flag.Parse は最初の位置引数で解析を止めるため、`qrcode hello -browser` の
// -browser がテキスト扱いになってしまう。"--" 以降はすべて位置引数として扱う。
func parseInterspersed(fs *flag.FlagSet, args []string) ([]string, error) {
	var positional []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		rest := fs.Args()
		if len(rest) == 0 {
			return positional, nil
		}
		if consumed := len(args) - len(rest); consumed > 0 && args[consumed-1] == "--" {
			return append(positional, rest...), nil
		}
		positional = append(positional, rest[0])
		args = rest[1:]
	}
}

// isTerminal は r が端末につながった標準入力かどうかを返す。
// 何も渡さずに実行したときに入力待ちで止まらないようにするために使う。
func isTerminal(r io.Reader) bool {
	f, ok := r.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
