---
date    : 2026-03-15         # 文書の初回作成日
author  : 高野順一            # 文書の主要な作成者
title   : runner 受け入れテスト一覧
owner   : 責任者              # 文書の現在の責任者 / 最終更新者
updated : 2026-03-20         # 文書の最終更新日
tags    : SE向け, PM向け, 開発者向け, 新人向け
       - 技術解説, 概念, 実践, ノウハウ, 考察, トラブルシューティング
policy  :
  style   : だ、である調、章は見出しレベル2から、
  scope   : 
  tone    : 技術論は正確に
  purpose : 
note    : 
---

**最終更新日: 2026年03月20日**

# runner 受け入れテスト一覧

## 1. 基本起動

### AT-010 ファイル指定実行

**前提**
`hello.py` が存在する。

**操作**

```bash
runner hello.py
```

**期待結果**

* `hello.py` を対象として実行する
* 実行前にコマンドを表示する
* 正常終了時は 0 を返す

---

### AT-011 `.run` ファイル明示指定

**前提**
`hello.run` が存在する。

**操作**

```bash
runner hello.run
```

**期待結果**

* `hello.run` を対象として実行する
* `.run` ヘッダを解析する
* 実行前にコマンドを表示する

---

### AT-012 名前指定実行

**前提**
`build.run` が存在する。

**操作**

```bash
runner build
```

**期待結果**

* `build.run` を探索して実行する
* 実行前にコマンドを表示する

---

### AT-013 引数なし実行

**前提**
`runfile.run` が存在する。

**操作**

```bash
runner
```

**期待結果**

* `runfile.run` を探索して実行する
* 実行前にコマンドを表示する

---

## 2. target 解決優先順位

### AT-020 拡張子あり target は直接実行

**前提**
`hello.py` が存在する。

**操作**

```bash
runner hello.py
```

**期待結果**

* `<target>.run` は探索しない
* `hello.py` を直接実行対象とする

---

### AT-021 拡張子なし target は `.run` を探索

**前提**
`build.run` が存在する。

**操作**

```bash
runner build
```

**期待結果**

* `build.run` を探索する
* 見つかった場合はそれを実行する

---

### AT-022 拡張子なし target が見つからない場合

**前提**
`missing.run` が存在しない。

**操作**

```bash
runner missing
```

**期待結果**

* エラー終了する
* 非0を返す
* `target not found: missing.run` を表示する

---

## 3. dry-run

### AT-030 dry-run は実行しない

**前提**
`hello.py` が存在し、実行すると副作用が出る内容になっている。

**操作**

```bash
runner -n hello.py
```

**期待結果**

* 対象解決を行う
* 実行コマンド生成を行う
* 実際のプロセス実行は行わない
* 副作用は発生しない

---

### AT-031 dry-run で `.run` を解析する

**前提**
`build.run` が存在する。

**操作**

```bash
runner -n build
```

**期待結果**

* `build.run` を探索する
* ヘッダ解析を行う
* runtime 解決を行う
* 一時ファイル名を決定する
* 実行はしない

### AT-032 `#script` の dry-run は全 OS を表示

**前提**
`install.run` の先頭が `#script` であり、`@windows` / `@linux` / `@macos` の各ブロックが存在する。

**操作**

```bash
runner -n install.run
```

**期待結果**

* `@windows` / `@linux` / `@macos` を順に処理する
* 各 OS ブロックの runtime ヘッダを解析する
* 各 OS ブロックに対して変数展開を行う
* 展開後の本文を表示する
* 一時ファイルは生成しない
* 実行はしない

---

### AT-033 --dry-run=windows

入力:

runner --dry-run=windows script.run

期待:

* windows ブロックのみ表示

---

### AT-034 --dry-run=linux

入力:

runner --dry-run=linux script.run

期待:

* linux ブロックのみ表示

---

### AT-035 --dry-run=macos

入力:

runner --dry-run=macos script.run

期待:

* macos ブロックのみ表示

---

### AT-036 --dry-run=all

入力:

runner --dry-run=all script.run

期待:

* 全 OS ブロックを順に表示

---

### AT-037 不正 OS

入力:

runner --dry-run=freebsd script.run

期待:

```error
[runner] unknown os: freebsd
```

---

## 4. オプション

### AT-040 `--help`

**操作**

```bash
runner --help
```

**期待結果**

* usage を表示する
* 終了コードは 0

---

### AT-041 `--version`

**操作**

```bash
runner --version
```

**期待結果**

* バージョンを表示する
* 終了コードは 0

---

### AT-042 `--list`

**前提**
カレントディレクトリに `runfile.run`, `build.run`, `test.run` が存在する。

**操作**

```bash
runner --list
```

**期待結果**

* カレントディレクトリのみを探索する
* 再帰探索は行わない
* `.run` タスク一覧を表示する
* `.run` を拡張子なしで表示する
* ファイル名の辞書順（昇順）で表示する
* 終了コードは 0

---

### AT-043 `--list` は非再帰

**前提**
`sub/deploy.run` が存在する。

**操作**

```bash
runner --list
```

**期待結果**

* `sub/deploy.run` は表示しない

---

### AT-044 `--list` 対象なし

**前提**
カレントディレクトリに `.run` ファイルが存在しない。

**操作**

```bash
runner --list
```

**期待結果**

* 何も表示しない
* 終了コードは 0

---

### AT-045 target 後ろのオプションは不許可

**操作**

```bash
runner hello.py -n
```

**期待結果**

* エラー終了する
* 非0を返す

---

### AT-046 --check は .run を解析する

入力:

runner --check build.run

期待:

* ヘッダ解析が行われる
* エラーがなければ成功

---

### AT-047 --check は全 OS を検証する

入力:

runner --check script.run

期待:

* 全 OS ブロックが検証される
* OS に依存しない

---

### AT-048: --check 正常終了

入力:

runner --check valid.run

期待:

* 出力は最小
* 終了コード 0

---

### AT-049: --check と --dry-run 同時指定

入力:

runner --check --dry-run build.run

期待:

```error
[runner] invalid option combination
```

---

### AT-04A 無効オプション

**操作**

```bash
runner --exec
```

**期待結果**

* エラー終了する
* `unknown option: --exec` を表示する
* 非0を返す

---

## 5. `.run` ヘッダ

### AT-050 runtime 指定ヘッダ

**前提**
`hello.run` の先頭行が `#python` である。

**操作**

```bash
runner hello.run
```

**期待結果**

* `runtime.python` を参照して runtime を解決する

---

### AT-051 仮想ファイル名指定ヘッダ

**前提**
`hello.run` の先頭行が `#program.py` である。

**操作**

```bash
runner hello.run
```

**期待結果**

* 拡張子 `py` を抽出する
* `ext.py` を参照する
* `runtime.python` を解決する

---

### AT-052 拡張子指定ヘッダ

**前提**
`hello.run` の先頭行が `#.py` である。

**操作**

```bash
runner hello.run
```

**期待結果**

* `ext.py` を参照する
* `runtime.python` を解決する

---

### AT-053 ヘッダなし

**前提**
`bad.run` の1行目がヘッダでない。

**操作**

```bash
runner bad.run
```

**期待結果**

* エラー終了する
* `invalid .run format` を表示する
* 非0を返す

---

### AT-054 不正ヘッダ

**前提**
`bad.run` の1行目が不正なヘッダである。

**操作**

```bash
runner bad.run
```

**期待結果**

* エラー終了する
* `invalid .run header` を表示する
* 非0を返す

---

### AT-055 ヘッダ前空行は不許可

**前提**
`bad.run` の先頭が空行で、その次に `#python` がある。

**操作**

```bash
runner bad.run
```

**期待結果**

* エラー終了する
* 1行目が必ずヘッダである前提を満たさないため失敗する

---

### AT-056 ヘッダ後空行は許可

**前提**
`ok.run` の内容が以下である。

```text
#python

print("Hello")
```

**操作**

```bash
runner ok.run
```

**期待結果**

* 正常に実行できる

### AT-057 `#script` ヘッダ

**前提**
`install.run` の先頭行が `#script` である。

**操作**

```bash
runner -n install.run
```

**期待結果**

* `#script` として解釈する
* 通常 `.run` とは異なり、OS ブロックモードで解析する

### AT-058 OS ブロック内 runtime ヘッダ必須

**前提**
`install.run` が `#script` で始まり、`@linux` の直後に runtime ヘッダがない。

**操作**

```bash
runner install.run
```

**期待結果**

* エラー終了する
* `runtime header required in os block: linux` を表示する
* 非0を返す

### AT-059 `#script` の不正構造

**前提**
`install.run` が `#script` で始まり、OS ブロック外に本文行がある。

**操作**

```bash
runner install.run
```

**期待結果**

* エラー終了する
* `invalid script block` を表示する
* 非0を返す

---

### AT-05A 同一OSブロック重複

前提:
install.run に同一 OS ブロックが複数存在する（例: @windows が2回）

操作:
runner install.run を実行

期待結果:
エラー終了する
標準エラーに以下を含む

```text
[runner] duplicate os block: windows
```

---

### AT-05B 未知のOSブロック

前提:
install.run に未対応の OS マーカー（例: @freebsd）が含まれる

操作:
runner install.run を実行

期待結果:
エラー終了する
標準エラーに以下を含む

```text
[runner] unknown os block: freebsd
```

---

### AT-05C OSブロックなし

前提:
install.run が #script 形式であるが、@windows / @linux / @macos のいずれも含まない

操作:
runner install.run を実行

期待結果:
エラー終了する
標準エラーに以下を含む

```text
[runner] os block not found: <current-os>
```

---

## 6. `.run` 一時ファイル

### AT-060 `.run` は一時ファイルに展開する

**前提**
`hello.run` が存在する。

**操作**

```bash
runner hello.run
```

**期待結果**

* 本文を一時ファイルに展開する
* 一時ファイルを runtime に渡して実行する

---

### AT-061 一時ファイルは OS 一時ディレクトリに生成

**前提**
`hello.run` が存在する。

**操作**

```bash
runner hello.run
```

**期待結果**

* 一時ファイルは OS の一時ディレクトリに生成される

---

### AT-062 一時ファイルは実行終了後に削除

**前提**
`hello.run` が存在する。

**操作**

```bash
runner hello.run
```

**期待結果**

* 実行終了後に一時ファイルを削除する

---

### AT-063 dry-run では実行しない

**前提**
`hello.run` が存在する。

**操作**

```bash
runner -n hello.run
```

**期待結果**

* 一時ファイル名は決定するが、生成しない
* 実行はしない

### AT-064 `#script` は選択 OS ブロックだけを一時ファイル化する

**前提**
`install.run` が `#script` で始まり、各 OS ブロックを持つ。

**操作**

```bash
runner install.run
```

**期待結果**

* 現在 OS に対応するブロックだけを選択する
* 選択したブロックの本文だけを一時ファイルへ展開する
* `#script` や `@windows` などの構造行は一時ファイルに含めない

---

### AT-065 Windows pwsh 一時ファイル拡張子

前提:
#script または #pwsh 形式の .run ファイルが存在する

操作:
Windows 環境で runner hello.run を実行

期待結果:
pwsh に渡される一時ファイルは .ps1 拡張子を持つ

---

### AT-066 非選択OSブロックは実行されない

前提:
install.run に複数の OS ブロックがあり、それぞれ異なる副作用（例: ファイル生成）を持つ

操作:
runner install.run を実行

期待結果:
現在の OS に対応するブロックのみ実行される
他の OS ブロックは実行されない

---

## 7. runner.env 読み込み

### AT-070 基本読み込み

**前提**
`runner.env` に `runtime.python=python` と `ext.py=python` がある。

**操作**

```bash
runner hello.py
```

**期待結果**

* `ext.py` → `python`
* `runtime.python` → `python`
* の順で解決する

---

### AT-071 前後空白を無視

**前提**
`runner.env` に以下がある。

```text
runtime.python = python
ext.py = python
```

**操作**

```bash
runner hello.py
```

**期待結果**

* 正常に解決できる

---

### AT-072 コメント行を無視

**前提**
`runner.env` にコメント行がある。

**操作**

```bash
runner hello.py
```

**期待結果**

* コメント行は無視する

---

### AT-073 行内コメントはサポートしない

**前提**
`runner.env` に以下がある。

```text
runtime.python=python # invalid comment
```

**操作**

```bash
runner hello.py
```

**期待結果**

* `#` 以降をコメントとしては扱わない
* 値の一部として読み込む
  
---

### AT-074 重複キーは後勝ち

**前提**
`runner.env` に以下がある。

```text
runtime.python=python
runtime.python=python3
```

**操作**

```bash
runner hello.py
```

**期待結果**

* `python3` を採用する

---

### AT-075 key は大小文字を区別

**前提**
`runner.env` に `ext.PY=python` のみがある。

**操作**

```bash
runner hello.py
```

**期待結果**

* `ext.py` は未定義とみなす
* エラー終了する

---

### AT-076: --env 指定ファイルを使用

入力:

runner --env testdata/env/runner.env hello.py

期待:

* 指定された runner.env を使用する

---

### AT-077: --env ファイル不存在

入力:

runner --env notfound.env hello.py

期待:

```error
[runner] file not found: notfound.env
```

---

### AT-078: カレント runner.env は読まない

入力:

(カレントに runner.env が存在)

runner hello.py

期待:

* カレントの runner.env は使用されない

---

## 8. command 分割

### AT-080 command を空白で分割

**前提**
`runner.env` に以下がある。

```text
runtime.python=python -u
```

**操作**

```bash
runner hello.py
```

**期待結果**

* `python` と `-u` を別引数として扱う

---

### AT-081 ダブルクォートで1引数化

**前提**
`runner.env` に以下がある。

```text
runtime.python="C:\Program Files\Python\python.exe" -u
```

**操作**

```bash
runner hello.py
```

**期待結果**

* `"C:\Program Files\Python\python.exe"` を1つの引数として扱う
* `-u` は別引数として扱う

---

### AT-082 エスケープ `\"`

**前提**
`runner.env` に以下がある。

```text
runtime.echo=echo "a\"b"
ext.txt=echo
```

`hello.txt` が存在する。

**操作**

```text
runner hello.txt
```

**期待結果**

* `\"` を `"` として解釈する
* 実行引数に `a"b` が渡される

---

### AT-083 エスケープ `\\`

**前提**
`runner.env` に以下がある。

```text
runtime.echo=echo "C:\\Tools\\Python"
ext.txt=echo
```

`hello.txt` が存在する。

**操作**

```text
runner hello.txt
```

**期待結果**

* `\\` を `\` として解釈する
* 実行引数に `C:\Tools\Python` が渡される

---

### AT-084 不正クォート

**前提**
`runner.env` に閉じ忘れたクォートがある。

**操作**

```bash
runner hello.py
```

**期待結果**

* エラー終了する
* `invalid runtime command` を表示する
* 非0を返す

---

### AT-085 shell 展開なし

**前提**
`runner.env` に以下がある。

```text
runtime.python=$PYTHON
```

**操作**

```bash
runner hello.py
```

**期待結果**

* `$PYTHON` を展開しない
* shell を経由しない

---

## 9. 未定義 runtime / extension

### AT-090 extension 未定義

**前提**
`runner.env` に `ext.py` が存在しない。

**操作**

```bash
runner hello.py
```

**期待結果**

* エラー終了する
* `extension not mapped: .py` を表示する

---

### AT-091 runtime 未定義

**前提**
`runner.env` に `ext.py=python` はあるが `runtime.python` が存在しない。

**操作**

```bash
runner hello.py
```

**期待結果**

* エラー終了する
* `runtime not defined: python` を表示する

---

## 10. 終了コード

### AT-100 正常終了

**前提**
正常終了するスクリプトがある。

**操作**

```bash
runner hello.py
```

**期待結果**

* 終了コード 0 を返す

---

### AT-101 実行プロセスの非0終了を伝播

**前提**
対象スクリプトが `exit 5` 相当で終了する。

**操作**

```bash
runner fail.py
```

**期待結果**

* runner も終了コード 5 を返す

---

### AT-102 runner 自体のエラー

**前提**
存在しない target を指定する。

**操作**

```bash
runner missing
```

**期待結果**

* 非0を返す

---

## 11. 文字コード・改行コード

### AT-110 `.run` UTF-8

**前提**
UTF-8 で書かれた `.run` がある。

**操作**

```bash
runner hello.run
```

**期待結果**

* 正常に読み込める

---

### AT-111 `runner.env` UTF-8

**前提**
UTF-8 で書かれた `runner.env` がある。

**操作**

```bash
runner hello.py
```

**期待結果**

* 正常に読み込める

---

### AT-112 UTF-8 BOM を無視

**前提**
先頭に UTF-8 BOM を含む `hello.run` がある（1行目は `#python`）。

**操作**

```bash
runner hello.run
```

**期待結果**

* BOM を無視してヘッダ解析できる
* 正常に読み込める

---

### AT-113 LF を許可

**前提**
LF 改行の `.run` / `runner.env` がある。

**操作**

```bash
runner hello.run
```

**期待結果**

* 正常に読み込める
* 正常終了する

---

### AT-114 CRLF を許可

**前提**
CRLF 改行で作成した `hello.run` と `runner.env` がある。

**操作**

```bash
runner hello.run
```

**期待結果**

* 正常に読み込める
* 正常終了する

---

## 12. `.run` 実行権限

### AT-120 `.run` に実行権限は不要

**前提**
`build.run` に実行権限が付いていない。

**操作**

```bash
runner build
```

**期待結果**

* 正常に実行できる

---

## 14. 非対象

### N-001 `.run` の直接実行は受け入れ対象外

**前提**
`build.run` が存在する。

**操作**
OS から直接 `build.run` を実行しようとする

**期待結果**

* runner の正式な利用方法ではない
* 受け入れ対象外とする

---
