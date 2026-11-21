// qrcode-cli: 標準入力から受け取ったテキストをQRコード化してターミナルに表示する。
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

func main() {
	if err := run(os.Stdin, os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "qrcode:", err)
		os.Exit(1)
	}
}

func run(stdin io.Reader, stdout io.Writer, args []string) error {
	fs := flag.NewFlagSet("qrcode", flag.ContinueOnError)
	level := fs.String("level", "M", "誤り訂正レベル (L, M, Q, H)")
	invert := fs.Bool("invert", false, "白黒を反転する (明るい背景のターミナル向け)")
	keepNewline := fs.Bool("keep-newline", false, "末尾の改行を削除しない")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: echo 'text' | qrcode [options]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	recovery, err := parseLevel(*level)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(stdin)
	if err != nil {
		return fmt.Errorf("標準入力の読み込みに失敗: %w", err)
	}
	text := string(data)
	if !*keepNewline {
		text = strings.TrimRight(text, "\r\n")
	}
	if text == "" {
		return errors.New("入力が空です")
	}

	qr, err := qrcode.New(text, recovery)
	if err != nil {
		return fmt.Errorf("QRコードの生成に失敗: %w", err)
	}

	_, err = io.WriteString(stdout, render(qr.Bitmap(), *invert))
	return err
}

func parseLevel(s string) (qrcode.RecoveryLevel, error) {
	switch strings.ToUpper(s) {
	case "L":
		return qrcode.Low, nil
	case "M":
		return qrcode.Medium, nil
	case "Q":
		return qrcode.High, nil
	case "H":
		return qrcode.Highest, nil
	}
	return 0, fmt.Errorf("不正な誤り訂正レベル: %q (L, M, Q, H のいずれか)", s)
}

// render は上下2モジュールを1文字にまとめ、半角ブロック文字で描画する。
// bitmap は true が黒モジュール。暗い背景のターミナルを想定し、
// デフォルトでは黒モジュールを空白、白モジュールをブロックで描く。
func render(bitmap [][]bool, invert bool) string {
	dark := func(y, x int) bool {
		if y >= len(bitmap) {
			return invert // 範囲外はクワイエットゾーン（白）扱い
		}
		return bitmap[y][x] != invert
	}

	var b strings.Builder
	for y := 0; y < len(bitmap); y += 2 {
		for x := range bitmap[y] {
			top, bottom := dark(y, x), dark(y+1, x)
			switch {
			case top && bottom:
				b.WriteRune(' ')
			case top:
				b.WriteRune('▄')
			case bottom:
				b.WriteRune('▀')
			default:
				b.WriteRune('█')
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}
