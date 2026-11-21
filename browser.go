package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"

	qrcode "github.com/skip2/go-qrcode"
)

var pageTmpl = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="ja">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>QR Code</title>
<style>
  body { margin: 0; min-height: 100vh; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 24px; background: #f4f4f5; font-family: system-ui, sans-serif; color: #18181b; }
  img { width: min(80vw, 480px); height: auto; image-rendering: pixelated; background: #fff; border-radius: 12px; box-shadow: 0 2px 12px rgba(0,0,0,.08); }
  pre { max-width: min(90vw, 720px); margin: 0; padding: 12px 16px; white-space: pre-wrap; word-break: break-all; background: #fff; border-radius: 8px; font-size: 14px; }
</style>
</head>
<body>
<img src="{{.Image}}" alt="QR Code">
<pre>{{.Text}}</pre>
</body>
</html>
`))

// openBrowser はテストで差し替えられるようにしている。
var openBrowser = func(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Run()
}

// showInBrowser は QR コード画像を埋め込んだ HTML を一時ファイルに書き出し、ブラウザで開く。
// ブラウザが読み込む前に消えないよう、一時ファイルは削除せずに残す。
func showInBrowser(qr *qrcode.QRCode, text string, stdout io.Writer) error {
	png, err := qr.PNG(-10) // 負の値は 1 モジュールあたりのピクセル数
	if err != nil {
		return fmt.Errorf("画像の生成に失敗: %w", err)
	}

	f, err := os.CreateTemp("", "qrcode-*.html")
	if err != nil {
		return fmt.Errorf("一時ファイルを作成できません: %w", err)
	}
	err = pageTmpl.Execute(f, struct {
		Image template.URL
		Text  string
	}{
		Image: template.URL("data:image/png;base64," + base64.StdEncoding.EncodeToString(png)),
		Text:  text,
	})
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return fmt.Errorf("HTML の書き出しに失敗: %w", err)
	}

	if err := openBrowser(f.Name()); err != nil {
		return fmt.Errorf("ブラウザを起動できません (%s を手動で開いてください): %w", f.Name(), err)
	}
	fmt.Fprintf(stdout, "ブラウザで開きました: %s\n", f.Name())
	return nil
}
