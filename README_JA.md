# Kompresi
バックグラウンドでフォルダ内にある全てのPNGとJPEG画像を無劣化で圧縮するGoアプリケーションです。

## できること

指定されたディレクトリを監視し、検知されたPNGとJPEG画像をロスレス圧縮します。

PNG画像には[Zopfli](https://github.com/google/zopfli)が、JPEG画像には[Guetzli](https://github.com/google/guetzli)が使用されます。

ZopfliとGuetzliは圧縮が**かなり遅い**です。
このアプリはバックグラウンドでの圧縮を想定していて、それ以外の用途には適していません。

## ダウンロード & セットアップ

### 1. ダウンロード

[ここから](https://github.com/tsg0o0/Kompresi/releases)最新版をダウンロードします。

*Goがあればソースコードからビルドすることもできます。*

### 2. 設定をする

**WindowsとmacOSでは、同梱されているKompresiConfigure(KompresiConfigure.exe)を使用して、画面上で簡単に設定と起動ができます。**

ターミナル(コマンドプロンプト)を開いて、kompresi(kompresi.exe)をドラックアンドドロップなどして起動します。
すると、おそらくセットアップガイドが英語で表示されます。

以下の引数を用いて設定をします。引数とは、ファイル名の後ろに空白を空けて続く文字列です。フォルダのパスはフォルダをドラックアンドドロップすれば自動で入力されます。

- `inputDir '読み取りフォルダのパス'` 画像を読み込むフォルダのパスを設定します。
- `outputDir '書き出しフォルダのパス'` 画像を書き出すフォルダのパスを設定します。
- `deleteOrigin 'Yes か No'` 圧縮後にもとの画像を削除します。圧縮後の画像だけが残ります。
- `optimLv '0 - 2'` 圧縮レベルを0から2から選択します。
  - `0`: 速いですが圧縮率は低いです。
  - `1`: 自動 (試験的)
  - `2`: 遅いですが圧縮率は高いです。(あんまり変わらないかも)
- `help` ヘルプを表示(英語)
- `license` ライセンスを表示

### 3. 実行

引数なしで再度実行します。すると、バックグラウンドでの処理が始まります。

<sub> *引数に特定の画像のパスを指定するとその画像だけを圧縮します。これはパフォーマンステスト用です。* </sub>

## サポート

バグを見つけたり質問がある場合は、遠慮なく[Issue](https://github.com/tsg0o0/Kompresi/issues)を立てるか[私に連絡](https://tsg0o0.com/contact/)してください。

## ライセンス

このソフトウェアは[Mozilla Public License 2.0](https://www.mozilla.org/en-US/MPL/2.0/)の下でライセンスされています。

[非公式日本語訳](https://www.mozilla.jp/documents/mpl/2.0/)もありますが、法的に有効な原文は上記の英語版です。

## チップ

気に入ってくれたら[チップを送って](https://tsg0o0.com/tip/)いただけると嬉しいです！
