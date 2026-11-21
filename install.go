package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// installSelf は実行中のバイナリを dir/qrcode にコピーする。
// dir が PATH に含まれていなければ、追加方法を案内する。
func installSelf(dir string, stdout io.Writer) error {
	src, err := os.Executable()
	if err != nil {
		return fmt.Errorf("実行ファイルの場所を取得できません: %w", err)
	}
	if src, err = filepath.EvalSymlinks(src); err != nil {
		return fmt.Errorf("実行ファイルの場所を取得できません: %w", err)
	}

	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("ホームディレクトリを取得できません: %w", err)
		}
		dir = filepath.Join(home, ".local", "bin")
	}
	if dir, err = filepath.Abs(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("インストール先を作成できません: %w", err)
	}

	name := "qrcode"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	dst := filepath.Join(dir, name)

	if same, _ := sameFile(src, dst); same {
		fmt.Fprintf(stdout, "既にインストール済みです: %s\n", dst)
	} else {
		if err := copyExecutable(src, dst); err != nil {
			return fmt.Errorf("インストールに失敗: %w", err)
		}
		fmt.Fprintf(stdout, "インストールしました: %s\n", dst)
	}

	if !inPath(dir) {
		fmt.Fprintf(stdout, "\n%s は PATH に含まれていません。シェルの設定ファイル (例: ~/.zshrc) に次の行を追加してください:\n\n", dir)
		fmt.Fprintf(stdout, "    export PATH=\"%s:$PATH\"\n", dir)
	}
	return nil
}

// copyExecutable は同じディレクトリに一時ファイルを書いてから rename する。
// 実行中のバイナリを上書きする場合でも安全に置き換えられる。
func copyExecutable(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	tmp, err := os.CreateTemp(filepath.Dir(dst), ".qrcode-install-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name()) // rename 成功後は何もしない

	if _, err := io.Copy(tmp, in); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o755); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), dst)
}

func sameFile(a, b string) (bool, error) {
	ai, err := os.Stat(a)
	if err != nil {
		return false, err
	}
	bi, err := os.Stat(b)
	if err != nil {
		return false, err
	}
	return os.SameFile(ai, bi), nil
}

func inPath(dir string) bool {
	for _, p := range filepath.SplitList(os.Getenv("PATH")) {
		if p != "" && filepath.Clean(p) == filepath.Clean(dir) {
			return true
		}
	}
	return false
}
