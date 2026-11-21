# qrcode-cli

引数または標準入力で受け取ったテキストを QR コードに変換し、ターミナルやブラウザに表示する CLI です。

```
$ qrcode https://example.com
█████████████████████████████████
█████████████████████████████████
████ ▄▄▄▄▄ ███▀▄▄ ▀▀ █ ▄▄▄▄▄ ████
████ █   █ █ ▄█▄█▀ ▀██ █   █ ████
████ █▄▄▄█ █ ▄▄▀█▄▀▀▄█ █▄▄▄█ ████
...
```

## インストール

Go 1.25 以降が必要です。

```sh
go install github.com/up-tri/qrcode-cli@latest
```

`go install` で入るコマンド名は `qrcode-cli` です。`qrcode` という名前で使いたい場合は、`qrcode-cli -install` を実行するか、リポジトリを clone してビルドしてください。

```sh
git clone https://github.com/up-tri/qrcode-cli.git
cd qrcode-cli
go build -o qrcode .
./qrcode -install
```

`-install` は実行中のバイナリを `~/.local/bin/qrcode` にコピーします。コピー先が PATH に含まれていない場合は、シェルの設定ファイルに追加する行を表示します。コピー先は `-install-dir` で変えられます。

```sh
./qrcode -install -install-dir /usr/local/bin
```

### zsh 補完

`qrcode -completion zsh` で補完スクリプトを出力します。次のどちらかの方法で読み込んでください。

**方法 1: `.zshrc` で読み込む** (手軽)

`~/.zshrc` の `compinit` より後ろに次の行を追加します。

```sh
eval "$(qrcode -completion zsh)"
```

**方法 2: fpath にファイルを置く** (シェルの起動が速い)

```sh
mkdir -p ~/.zfunc
qrcode -completion zsh > ~/.zfunc/_qrcode
```

`~/.zshrc` の `compinit` より前に `fpath=(~/.zfunc $fpath)` を追加し、シェルを開き直してください。この方法では、qrcode を更新したときに同じコマンドでファイルを作り直す必要があります。

オプション名、`-level` の値 (`L` `M` `Q` `H`)、`-install-dir` のディレクトリが補完されます。`go install` で入れた `qrcode-cli` コマンドも補完の対象です。

## 使い方

```sh
# 文字列を QR コードにする
qrcode https://example.com

# 標準入力から渡す
echo "https://example.com" | qrcode

# ファイルの内容を QR コードにする
qrcode < note.txt

# クリップボードの内容を QR コードにする (macOS)
pbpaste | qrcode

# 誤り訂正レベルを最高にする
qrcode -level H hello

# 明るい背景のターミナルで表示する
qrcode -invert hello

# ブラウザで表示する (オプションはテキストの後ろにも書ける)
qrcode https://example.com -browser

# "-" で始まるテキストは -- の後ろに書く
qrcode -- -hello
```

## オプション

| オプション | デフォルト | 説明 |
| --- | --- | --- |
| `-level` | `M` | 誤り訂正レベル。`L` (約7%) / `M` (約15%) / `Q` (約25%) / `H` (約30%) |
| `-invert` | `false` | 白黒を反転します。明るい背景のターミナル向けです |
| `-browser` | `false` | ターミナルではなくブラウザで表示します |
| `-install` | `false` | このコマンドを `qrcode` という名前で `-install-dir` にコピーし、PATH から使えるようにします。標準入力は読みません |
| `-install-dir` | `~/.local/bin` | `-install` のコピー先 |
| `-completion` | | シェル補完スクリプトを出力します。対応シェルは `zsh` のみです |
| `-keep-newline` | `false` | 入力末尾の改行を残します。デフォルトでは `echo` などで付く改行を取り除きます |

## 仕様

- テキストは引数か標準入力で渡します。引数が複数あるときは半角スペースでつなげます。引数と標準入力の両方があるときは引数を使います。
- 引数がなく、標準入力も端末のまま (パイプやリダイレクトがない) のときは、入力待ちにならずにヘルプを表示して終了します。
- 上下 2 マスを半角ブロック文字 (`▀` `▄` `█`) 1 文字にまとめて描くため、縦横の比率を保ったままコンパクトに表示されます。
- 読み取りに必要な余白 (クワイエットゾーン) を周囲に含めて出力します。
- `-browser` を付けると、QR コード画像と元のテキストを埋め込んだ HTML を一時ディレクトリに書き出し、既定のブラウザで開きます (macOS は `open`、Linux は `xdg-open`、Windows は `rundll32`)。ブラウザが読み込む前に消えないよう、HTML ファイルは削除せずに残します。
- 入力が空のとき、誤り訂正レベルに不正な値を渡したとき、入力が QR コードの容量を超えたときは、標準エラー出力にメッセージを出して終了コード 1 で終了します。

## ライセンス

MIT
