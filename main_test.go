package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBrowserWritesHTML(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())
	var opened string
	orig := openBrowser
	openBrowser = func(path string) error { opened = path; return nil }
	t.Cleanup(func() { openBrowser = orig })

	var out bytes.Buffer
	if err := run(strings.NewReader("<b>hello</b>\n"), &out, []string{"-browser"}); err != nil {
		t.Fatal(err)
	}
	html, err := os.ReadFile(opened)
	if err != nil {
		t.Fatal(err)
	}
	s := string(html)
	if !strings.Contains(s, `src="data:image/png;base64,iVBOR`) {
		t.Errorf("PNG の data URI が埋め込まれていない")
	}
	if !strings.Contains(s, "&lt;b&gt;hello&lt;/b&gt;") {
		t.Errorf("入力テキストがエスケープされて表示されていない")
	}
}

func TestInstall(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "bin")
	t.Setenv("PATH", "/usr/bin")

	var out bytes.Buffer
	if err := run(strings.NewReader(""), &out, []string{"-install", "-install-dir", dir}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(dir, "qrcode"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Errorf("実行権限がない: %v", info.Mode())
	}
	if !strings.Contains(out.String(), "PATH に含まれていません") {
		t.Errorf("PATH の案内が出ていない: %s", out.String())
	}

	// PATH に含まれていれば案内は出ない
	t.Setenv("PATH", dir)
	out.Reset()
	if err := run(strings.NewReader(""), &out, []string{"-install", "-install-dir", dir}); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "PATH に含まれていません") {
		t.Errorf("不要な PATH の案内が出ている: %s", out.String())
	}
}
