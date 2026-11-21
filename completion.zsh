#compdef qrcode qrcode-cli

# qrcode の zsh 補完。`qrcode -completion zsh` で出力される。
# fpath 上に _qrcode として置いても、.zshrc で eval しても動く。

_qrcode() {
  _arguments \
    '(- *)'{-h,-help}'[ヘルプを表示する]' \
    '-level[誤り訂正レベル]:level:((L\:"約7%を復元できる" M\:"約15%を復元できる (デフォルト)" Q\:"約25%を復元できる" H\:"約30%を復元できる"))' \
    '-invert[白黒を反転する (明るい背景のターミナル向け)]' \
    '-keep-newline[末尾の改行を削除しない]' \
    '-browser[ターミナルではなくブラウザで表示する]' \
    '-install[このコマンドを PATH の通ったディレクトリにインストールする]' \
    '-install-dir[-install のコピー先]:directory:_directories' \
    '(- *)-completion[シェル補完スクリプトを出力する]:shell:(zsh)' \
    '*:text:_nothing'
}

if [[ $funcstack[1] == _qrcode ]]; then
  # fpath から autoload された場合
  _qrcode "$@"
else
  # eval "$(qrcode -completion zsh)" で読み込まれた場合
  compdef _qrcode qrcode qrcode-cli
fi
