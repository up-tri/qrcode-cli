# qrcode-cli

標準入力から受け取ったテキストを QR コードに変換し、ターミナルに表示する CLI です。

```
$ echo "https://example.com" | qrcode
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

`go install` で入るコマンド名は `qrcode-cli` です。`qrcode` という名前で使いたい場合は、リポジトリを clone してビルドしてください。

```sh
git clone https://github.com/up-tri/qrcode-cli.git
cd qrcode-cli
go build -o qrcode .
```

## 使い方

```sh
# 文字列を QR コードにする
echo "https://example.com" | qrcode

# ファイルの内容を QR コードにする
qrcode < note.txt

# クリップボードの内容を QR コードにする (macOS)
pbpaste | qrcode

# 誤り訂正レベルを最高にする
echo "hello" | qrcode -level H

# 明るい背景のターミナルで表示する
echo "hello" | qrcode -invert
```

## オプション

| オプション | デフォルト | 説明 |
| --- | --- | --- |
| `-level` | `M` | 誤り訂正レベル。`L` (約7%) / `M` (約15%) / `Q` (約25%) / `H` (約30%) |
| `-invert` | `false` | 白黒を反転します。明るい背景のターミナル向けです |
| `-keep-newline` | `false` | 入力末尾の改行を残します。デフォルトでは `echo` などで付く改行を取り除きます |

## 仕様

- 上下 2 マスを半角ブロック文字 (`▀` `▄` `█`) 1 文字にまとめて描くため、縦横の比率を保ったままコンパクトに表示されます。
- 読み取りに必要な余白 (クワイエットゾーン) を周囲に含めて出力します。
- 入力が空のとき、誤り訂正レベルに不正な値を渡したとき、入力が QR コードの容量を超えたときは、標準エラー出力にメッセージを出して終了コード 1 で終了します。

## ライセンス

MIT
