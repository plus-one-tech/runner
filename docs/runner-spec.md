---
date    : 2026-03-15         # 文書の初回作成日
author  : 高野順一            # 文書の主要な作成者
title   : runner 仕様
owner   : 責任者              # 文書の現在の責任者 / 最終更新者
updated : 2026-03-16         # 文書の最終更新日
tags    : SE向け, PM向け, 開発者向け, 新人向け
       - 技術解説, 概念, 実践, ノウハウ, 考察, トラブルシューティング
policy  :
  style   : だ、である調、章は見出しレベル2から、
  scope   : 
  tone    : 技術論は正確に
  purpose : 
note    : 
---

**最終更新日: 2026年03月16日**

# runner 仕様

## はじめに

runner はプログラムやスクリプトの実行を簡単かつ再現可能にするためのツールである。

本ツールは以下の原則に基づいて設計されている。

**1. 実行の可視化**
   runner は実行前に実際に実行されるコマンドを表示する。
**2. shell 非依存**
   コマンド実行は shell を経由せず直接プロセスを起動する。
**3. 最小構造**
   `.run` ファイルは最小の構造のみを持つ。
**4. 単一責務**
   runner はプログラム実行の整理に特化したツールである。

---

## 1. runner

### 1.1 目的

`runner` は、**ソースファイルや `.run` ファイルを、適切な実行環境で起動するためのランチャー**である。

目的は次の3つ。

* 実行方法を統一する
* 実行時に、実際に使うコマンドを見せる
* 反復実行を楽にする

---

### 1.2 起動仕様

`runner` は次の形式で起動する。

```bash
runner [options] [target]
```

`target` は省略可能とする。

---

### 1.3 起動パターン

#### 1.3.1 ファイル指定実行

```bash
runner program.ext
```

指定したファイルを実行する。

例

```bash
runner hello.py
runner app.cs
runner main.rs
runner script.run
```

#### 1.3.2 名前指定実行

```bash
runner build
runner test
runner deploy
```

この場合、`<target>.run` を探索して実行する。

例

```bash
runner build
```

↓

```text
build.run
```

#### 1.3.3 引数なし実行

```bash
runner
```

引数なしの場合は、カレントディレクトリの `runfile.run` を探索して実行する。

### 1.3.4 target 解決優先順位

target は次の優先順位で解決する。

1. target に拡張子がある場合は、そのファイルを直接実行する。
2. target に拡張子がない場合は `<target>.run` を探索する。
3. `<target>.run` が存在しない場合はエラーとする。

例:

runner build → build.run
runner hello.py → hello.py

---

### 1.4 解決結果

| コマンド                | 実行対象          |
| ------------------- | ------------- |
| `runner hello.py`   | `hello.py`    |
| `runner script.run` | `script.run`  |
| `runner build`      | `build.run`   |
| `runner`            | `runfile.run` |

---

### 1.5 実行前表示

runner は実行前に、実際に実行するコマンドを表示する。

例

```text
[runner] command: python hello.py
```

`.run` の場合も同様に表示する。

例

```text
[runner] command: python /tmp/runner_tmp.py
```

---

### 1.6 dry-run

実際には実行せず、解決結果だけを表示するモードをサポートする。

#### 1.6.1 指定方法

```bash
runner --dry-run hello.py
runner -n build
runner -n
```

#### 1.6.2 動作

dry-run では次を行う。

* 対象ファイル解決
* `.run` ヘッダ解析
* runtime 解決
* 実行コマンド生成
* 一時ファイル名決定
* 必要な範囲の変数展開

実際のプロセス実行は行わない。

#### 1.6.3 `#script` の dry-run

`.run` の先頭ヘッダが `#script` の場合、dry-run は **全 OS ブロック** を対象とする。

表示順は次のとおり。

1. `@windows`
2. `@linux`
3. `@macos`

各 OS ブロックについて、runner は次を行う。

* OS ブロック抽出
* ブロック内 runtime ヘッダ解析
* `runner.env` の変数展開
* 展開後のスクリプト内容表示

`#script` の dry-run でも、一時ファイルは生成しない。

#### 1.6.4 `#runtime` の dry-run

`#script` 以外の通常 `.run` は、従来どおり **対象 1 本だけ** を dry-run の対象とする。

---

#### 1.6.5 dry-run 出力内容

dry-run では以下の情報を表示する。

* 対象 OS
* 実行コマンド
* 展開後スクリプト

表示内容は実行時に使用される情報と一致すること。

---

### 1.7 オプション

初期版では次の 4 つ。

| 短縮   | 長形式         | 意味           |
| ---- | ----------- | ------------ |
| `-n` | `--dry-run` | 実行せず内容のみ表示   |
| `-h` | `--help`    | ヘルプ表示        |
|      | `--version` | バージョン表示      |
|      | `--list`    | `.run` タスク一覧 |

#### 1.7.1 オプション配置

オプションは **target の前に置く**。

```bash
runner -n hello.py
runner -n build
```

次の形式は **許可しない**。

```bash
runner hello.py -n
```

理由

* 解析が簡単
* CLIツールとして自然

#### 1.7.2 help

```bash
runner --help
```

表示

```text
usage: runner [options] [target]

examples:
  runner hello.py
  runner build
  runner
```

#### 1.7.3 version

```bash
runner --version
```

例

```text
runner 0.1.0
```

#### 1.7.4 list

カレントディレクトリの `.run` を表示する。
探索はカレントディレクトリのみとし、再帰探索は行わない。
`--list` は `.run` を**拡張子なし**で表示する。

```bash
runner --list
```

例

```text
runfile
build
test
deploy
```

#### 1.7.5 無効オプション

未定義オプションはエラー。

例

```bash
runner --check
```

出力

```text
[runner] unknown option: --check
```

---

### 1.8 エラー条件

#### 1.8.1 引数なしで `runfile.run` がない

```text
[runner] runfile.run not found
```

#### 1.8.2 `target.run` がない

```text
[runner] target not found: build.run
```

#### 1.8.3 対象ファイルが存在しない

```text
[runner] file not found: hello.py
```

---

### 1.9 終了コード

runner は実行したプロセスの終了コードをそのまま返す。

* 0 : 正常終了
* 非0 : runtime またはスクリプトのエラー

runner 自体のエラーも非0を返す。

---

### 1.10 設計方針

起動仕様は次を重視する。

#### シンプル

毎回 `run` を書かせない。

#### 反復実行しやすい

`runner` だけで `runfile.run` を実行できる。

#### タスク実行しやすい

`runner build` のように実行できる。

#### 透明性

何を実行するか必ず表示する。

---

## 2. `.run` ファイル仕様

### 2.1 目的

`.run` は、**runner 専用の簡易実行スクリプト**である。

役割は次のとおり。

* 実行するプログラム本文を 1 ファイルにまとめる
* 先頭行で実行方法を指定する
* `runner` から簡単に実行できるようにする

主な用途例

```text
runfile.run
build.run
test.run
deploy.run
```

---

### 2.2 基本構造

`.run` ファイルは次のいずれかの構造を持つ。

#### 2.2.1 通常モード

```text
#<header>
<program body>
```

1 行目は **ヘッダ**、2 行目以降は **本文** とする。

例

```text
#python
print("Hello Runner")
```

```text
#program.py
print("Hello Runner")
```

```text
#.py
print("Hello Runner")
```

#### 2.2.2 `#script` モード

```text
#script

@windows
#pwsh
...

@linux
#bash
...

@macos
#bash
...
```

`#script` モードでは、本文は **OS ブロックの集合**として扱う。

---

### 2.3 ヘッダ仕様

ヘッダは次の 4 形式をサポートする。

```text
#<runtime>
#<filename.extension>
#.<extension>
#script
```

---

### 2.4 ヘッダ解釈順

ヘッダは次の順で解釈する。

#### 2.4.1 runtime 指定

```text
#python
#bash
#pwsh
```

`<runtime>` として解釈する。

runner は `runner.env` から

```text
runtime.<name>
```

を検索する。

例

```text
#python
```

↓

```text
runtime.python
```

#### 2.4.2 仮想ファイル名指定

```text
#program.py
#build.sh
#script.cs
```

`<filename.extension>` として解釈する。

runner は拡張子を抽出し、`runner.env` から

```text
ext.<extension>
```

を検索して runtime を決定する。

例

```text
#program.py
```

↓

```text
ext.py
```

↓

```text
python
```

↓

```text
runtime.python
```

#### 2.4.3 拡張子指定

```text
#.py
#.cs
#.js
```

`.<extension>` として解釈する。

runner は `runner.env` から

```text
ext.<extension>
```

を検索して runtime を決定する。

例

```text
#.py
```

↓

```text
ext.py
```

↓

```text
python
```

↓

```text
runtime.python
```


#### 2.4.4 `#script`

```text
#script
```

`#script` は **OS ごとに異なるスクリプトを 1 ファイルにまとめるモード**として解釈する。

このヘッダ自体は runtime を表さない。
runtime は各 OS ブロック内の先頭ヘッダで決定する。

---

### 2.5 ヘッダの意味

#### 2.5.1 runtime 指定

```text
#python
```

実行環境を直接指定する。

用途

* `#bash`
* `#pwsh`
* `#python`

のように、**実行環境そのものを明示したい場合**に使う。

#### 2.5.2 仮想ファイル名指定

```text
#program.py
```

仮想ファイル名と拡張子を同時に指定する。

用途

* 一時ファイル名を分かりやすくしたい
* エラー表示やログ上のファイル名を自然にしたい

#### 2.5.3 拡張子指定

```text
#.py
```

拡張子のみを指定する。

用途

* runtime は `runner.env` に任せたい
* 仮想ファイル名までは不要

#### 2.5.4 `#script`

```text
#script
```

OS ごとに異なる処理を 1 つの `.run` にまとめたい場合に使う。

用途

* install 処理
* 初期セットアップ
* OS ごとのコピー先やコマンドが異なる処理

---

### 2.6 本文

#### 2.6.1 通常モード

`#script` 以外の `.run` では、2 行目以降を **そのままプログラム本文**とする。

例

```text
#python
print("Hello Runner")
```

```text
#bash
echo "Hello Runner"
```

本文の文法やコメントは **各 runtime / 言語の仕様に従う**。

#### 2.6.2 `#script` モードの本文

`#script` の場合、2 行目以降は **OS ブロックの集合**とする。

OS ブロックは次の 3 種類のみをサポートする。

```text
@windows
@linux
@macos
```

#### 2.6.3 OS ブロックの構造

各 OS ブロックは次の構造を持つ。

```text
@windows
#pwsh
<block body>
```

```text
@linux
#bash
<block body>
```

```text
@macos
#bash
<block body>
```

`@<os>` の次の **最初の非空行** は、必ず runtime 指定ヘッダでなければならない。

許可する形式は次のみとする。

```text
#<runtime>
```

`#program.py` や `#.py` は、OS ブロック内ヘッダとしてはサポートしない。

#### 2.6.4 OS ブロックの範囲

1 つの OS ブロックは、次のいずれかで終わる。

* 次の `@windows` / `@linux` / `@macos`
* ファイル終端

`@end` はサポートしない。

#### 2.6.5 `#script` モードで許可する内容

`#script` モードでは、OS ブロックの外側に置けるのは次のみとする。

* 空行
* 行頭 `#` のコメント

それ以外の本文行は不正とする。

#### 2.6.6 本文の解釈

OS ブロック内の runtime ヘッダより後ろは、その runtime に渡す **生のスクリプト本文**とする。

本文の文法、コメント、if 文、変数、制御構文は **各 runtime / 言語の仕様に従う**。

runner は本文の意味を解釈しない。

---

### 2.7 実行方法

`.run` は次の形で実行できる。

#### 2.7.1 明示指定

```bash
runner hello.run
```

#### 2.7.2 拡張子省略

```bash
runner hello
```

この場合 `hello.run` を探索して実行する。

#### 2.7.3 デフォルト実行

```bash
runner
```

この場合、カレントディレクトリの `runfile.run` を探索して実行する。

---

### 2.8 実行方式

`.run` の本文は、**一時ファイルに展開して実行する**。

`#script` の場合は、現在 OS に対応する OS ブロックだけを選択し、
そのブロックの runtime ヘッダを除いた本文を一時ファイルに展開して実行する。

stdin 実行や `-c` / `eval` 方式は採用しない。

理由

* 言語依存を避けられる
* 行番号が安定する
* 実行方式が単純になる

---

### 2.9 一時ファイル生成

runner は `.run` の実行対象本文を一時ファイルに書き出し、そのファイルを runtime に渡して実行する。

#### 2.9.1 runtime 指定の場合

```text
#python
print("Hello")
```

一時ファイル名は runner が自動生成する。

例

```text
runner_tmp
```

または runtime に応じた内部名。

#### 2.9.2 仮想ファイル名指定の場合

```text
#program.py
print("Hello")
```

一時ファイル名は

```text
program.py
```

を使用する。

#### 2.9.3 拡張子指定の場合

```text
#.py
print("Hello")
```

一時ファイル名は runner が自動生成する。

例

```text
runner_tmp.py
```

#### 2.9.4 `#script` の場合

`#script` では、選択された OS ブロックだけを対象として一時ファイルを生成する。

一時ファイルの内容には次を含めない。

* `#script`
* `@windows` / `@linux` / `@macos`
* OS ブロック内 runtime ヘッダ

一時ファイル名は runner が自動生成する。

#### 2.9.5 一時ファイル配置

一時ファイルは OS の一時ディレクトリに生成する。
runner は実行終了後に一時ファイルを削除する。

例

```text
/tmp/runner_tmp.py
```

一時ファイル名は、同時実行時に衝突しない一意な名前でなければならない。
固定名の再利用は行わない。

#### 2.9.6 dry-run

dry-run では一時ファイル名は決定するが、ファイル自体は生成しない。

---

#### 2.9.7 Windows の一時ファイル拡張子

Windows (pwsh) で実行する場合、一時ファイルは `.ps1` 拡張子を持つ必要がある。
runner は runtime が pwsh の場合、必ず `.ps1` 拡張子で一時ファイルを生成する。

---

### 2.10 runtime 解決手順

runner は次の手順で実行コマンドを決定する。

#### runtime 指定

```text
#python
```

↓

```text
runtime.python
```

↓

```text
python
```

#### 仮想ファイル名指定

```text
#program.py
```

↓

```text
extension = py
```

↓

```text
ext.py
```

↓

```text
python
```

↓

```text
runtime.python
```

↓

```text
python
```

#### 拡張子指定

```text
#.py
```

↓

```text
ext.py
```

↓

```text
python
```

↓

```text
runtime.python
```

↓

```text
python
```

#### `#script`

```text
#script
@linux
#bash
echo "Hello"
```

↓

```text
current os = linux
```

↓

```text
select @linux block
```

↓

```text
runtime = bash
```

↓

```text
runtime.bash
```

↓

```text
bash
```

---

### 2.11 実行コマンド生成

runner は最終的に次の形式でコマンドを生成する。

```text
<runtime-command> <script-file>
```

例

```text
python program.py
bash build.sh
pwsh script.ps1
dotnet run program.cs
```

`runtime.<name>` に引数が含まれている場合も、そのまま先頭コマンドとして扱う。

例

```text
runtime.python=python -u
```

実行

```text
python -u program.py
```

---

### 2.12 エラー条件

次の場合はエラーとする。

#### 2.12.1 ヘッダが存在しない

```text
[runner] invalid .run format
missing header
```

#### 2.12.2 ヘッダが解釈できない

```text
[runner] invalid .run header
```

#### 2.12.3 extension が未定義

```text
[runner] extension not mapped: .py
```

#### 2.12.4 runtime が未定義

```text
[runner] runtime not defined: python
```

#### 2.12.5 `#script` で OS ブロックが存在しない

現在 OS に対応する OS ブロックが存在しない場合はエラーとする。
また、OS ブロックが1つも存在しない場合もエラーとする。

例

```text
[runner] os block not found: linux
```

#### 2.12.6 OS ブロックの runtime ヘッダが存在しない

`@windows` / `@linux` / `@macos` の直後に runtime ヘッダがない場合はエラーとする。

例

```text
[runner] runtime header required in os block: linux
```

#### 2.12.7 `#script` の構造が不正

OS ブロック外に本文行がある、未対応の OS マーカーがある、または OS ブロック内ヘッダが不正な場合はエラーとする。

例

```text
[runner] invalid script block
```

#### 2.12.8 未定義変数

`runner.env` に存在しない変数を `${...}` で参照した場合はエラーとする。

例

```text
[runner] variable not defined: var.install_dir.linux
```

---

#### 2.12.9 OS ブロック重複

同一 OS ブロックが複数存在する場合はエラーとする。

例

[runner] duplicate os block: windows

---

### 2.13 コメントと空行

#### ヘッダ前の空行

許可しない。
1 行目は必ずヘッダとする。

#### ヘッダ後の空行

許可する。

例

```text
#python

print("Hello")
```

#### `#script` モードの空行

`#script` モードでは、OS ブロックの外側・内側ともに空行を許可する。

#### `#script` モードのコメント

`#script` モードでは、OS ブロックの外側にある行頭 `#` の行をコメントとして許可する。
ただし `#script` と OS ブロック内 runtime ヘッダはコメントではない。

OS ブロック内のコメントは、選択された runtime の文法に従う。

### 2.14 処理順

`.run` 実行時の処理順は次のとおり。

1. `.run` を読み込む
2. 1 行目のヘッダを解析する
3. `#script` の場合は OS ブロックを抽出する
4. runtime を決定する
5. 対象本文に対して変数展開を行う
6. 一時ファイル名を決定する
7. dry-run なら表示のみ行う
8. 通常実行なら一時ファイルを生成し、runtime に直接渡して実行する
9. 実行終了後、一時ファイルを削除する

---

## 3. `runner.env` 仕様

### 3.1 目的

`runner.env` は **runner の実行環境を定義する設定ファイル**である。

主な用途

* runtime の定義
* 拡張子と runtime の対応付け
* `.run` 本文で使う変数の定義
* 一部実行設定の定義

---

### 3.2 ファイル形式

`runner.env` は **key=value 形式のテキストファイル**とする。

例

```text
runtime.python=python
runtime.bash=bash
runtime.pwsh=pwsh

ext.py=python
ext.cs=dotnet
ext.js=node
ext.sh=bash
ext.ps1=pwsh

var.install_dir.windows=C:\tools\runner
var.install_dir.linux=/home/user/.local/bin
var.install_dir.macos=/Users/user/.local/bin
```

#### 3.2.1 コメント

コメントは **行頭 `#` のみ有効**とする。

例

```text
# runtime definitions
runtime.python=python
```

以下はコメントではない。

```text
runtime.python=python # invalid comment
```

#### 3.2.2 空行

空行は無視される。

#### 3.2.3 パース規則

runner.env の読み込みは次の規則で行う。

* key と value の前後の空白は無視する
* key は大小文字を区別する
* 同一 key が複数ある場合は **後勝ち** とする
* 行頭 `#` はコメントとして扱う
* 行内コメントはサポートしない
  
---

### 3.2.4 runner.env が存在しない場合

runner.env が必要な処理において、runner.env が存在しない場合はエラーとする。

例

[runner] file not found: runner.env

---

### 3.3 runtime 定義

runtime は次の形式で定義する。

```text
runtime.<name>=<command>
```

例

```text
runtime.python=python
runtime.bash=bash
runtime.pwsh=pwsh
runtime.node=node
runtime.dotnet=dotnet run
```

#### 3.3.1 command

`<command>` は **実行コマンドと任意の引数を含む文字列**とする。

例

```text
runtime.python=python -u
runtime.node=node --enable-source-maps
```

* `<command>` は空白（スペースまたはタブ）で分割する。
* ダブルクォート `"..."` で囲まれた部分は1つの引数として扱う。
* エスケープは `\"` と `\\` のみ有効。
* shell 展開（`$VAR`、`*`、`` `...` `` など）は行わない。
* 不正なクォート（閉じ忘れ）はエラーとする。
* 分割後の最初のトークンを実行ファイル名とする。

runner は次の形式で実行する。

```text
<command> <script-file>
```

例

```text
python -u program.py
```

`<command>` が解釈できない場合、エラーメッセージは次とする。

```text
[runner] invalid runtime command
```

#### 3.3.2 実行方式

runner は runtime コマンドを **shell を経由せず直接実行する**。

実行形式:

<runtime-command> <script-file>

例:

runtime.python=python -u

実行:

python -u program.py

---

### 3.4 拡張子マッピング

拡張子と runtime の対応付けは次の形式で定義する。

```text
ext.<extension>=<runtime>
```

例

```text
ext.py=python
ext.cs=dotnet
ext.js=node
ext.sh=bash
ext.ps1=pwsh
```

---

### 3.5 実行コマンド生成

runtime 解決後、runner は次の形式でコマンドを生成する。

```text
<runtime-command> <script-file>
```

例

```text
python program.py
dotnet run program.cs
bash script.sh
```

---

### 3.6 典型的な runner.env

```text
# runtime definitions
runtime.python=python
runtime.bash=bash
runtime.pwsh=pwsh
runtime.node=node
runtime.dotnet=dotnet run

# extension mapping
ext.py=python
ext.cs=dotnet
ext.js=node
ext.sh=bash
ext.ps1=pwsh
```

---

### 3.7 未定義 runtime

runtime が定義されていない場合、runner はエラーとする。

例

```text
[runner] runtime not defined: python
```

---

### 3.8 未定義 extension

拡張子が `ext.*` に定義されていない場合もエラーとする。

例

```text
[runner] extension not mapped: .go
```

---

### 3.9 変数定義

`.run` 本文で使用する変数は次の形式で定義する。

```text
var.<name>=<value>
```

例

```text
var.install_dir.windows=C:\tools\runner
var.install_dir.linux=/home/user/.local/bin
var.install_dir.macos=/Users/user/.local/bin
```

`var.` を付けない任意キーはサポートしない。

### 3.10 変数展開

`.run` 本文では、`runner.env` の変数を次の形式で参照できる。

```text
${var.<name>}
```

例

```text
cp runner ${var.install_dir.linux}/runner
```

#### 3.10.1 展開対象

runner が展開するのは **`runner.env` に定義された `var.*` のみ** とする。

OS 環境変数は展開しない。

例

* `$HOME`
* `%LOCALAPPDATA%`
* `$env:LOCALAPPDATA`

これらは runner では解釈せず、そのまま本文文字列として扱う。

#### 3.10.2 展開タイミング

変数展開は、runtime 決定後、実行対象本文に対して行う。

`#script` の場合は、各 OS ブロックを個別に処理する。

* 通常実行では、選択された OS ブロックだけを展開する
* dry-run では、各 OS ブロックを順に同じ規則で展開して表示する

#### 3.10.3 展開規則

次の規則を適用する。

* `${var.<name>}` のみサポートする
* 単純置換のみ行う
* デフォルト値はサポートしない
* ネストはサポートしない
* 式評価はサポートしない
* 未定義変数はエラーとする

### 3.11 まとめ

`runner.env` は次の3種類の定義を持つ。

```text
runtime.<name>=<command>
ext.<extension>=<runtime>
var.<name>=<value>
```

これにより runner は

```text
.run header
↓
extension/runtime 解決
↓
必要なら変数展開
↓
runtime command 実行
```

という処理を行う。

---

## 4. 非機能仕様

### 4.1 文字コード

.run および runner.env は UTF-8 で記述する。

UTF-8 BOM が存在する場合は読み込み時に無視する。

### 4.2 改行コード

改行コードは LF または CRLF を許可する。

### 4.3 .run の実行権限

.run ファイルは runner から実行することを前提とする。

実行権限（chmod +x）は不要であり、直接実行することは想定しない。

### 4.4 OS依存のディレクトリ仕様

runner は OS ごとに次の標準ディレクトリを使用する。

#### 実行ファイル

これらのディレクトリはユーザーの PATH に含まれることを想定する。

* Windows: `%LOCALAPPDATA%\\runner\\runner.exe`
* Linux: `~/.local/bin/runner`
* macOS: `~/.local/bin/runner`

#### ユーザー設定

* Windows: `%APPDATA%\\runner\\runner.env`
* Linux: `~/.config/runner/runner.env`
* macOS: `~/Library/Application Support/runner/runner.env`

#### 一時ファイル

runner は .run の実行時に一時ファイルを生成する。
一時ファイルは OS の一時ディレクトリに生成する。

* Windows: `%TEMP%`
* Linux: `/tmp`
* macOS: `/tmp`

実行終了後に削除する。

#### プロジェクト設定

* 全OS共通: `./runner.env`
  
---

## 5. 将来拡張

runner は将来的に次の機能をサポートする可能性がある。

### 5.1 runner.env のスコープ

現在は `runner.env` を単一ファイルとして扱う。

将来的には次のスコープをサポートする可能性がある。

1. プロジェクト設定  
   `./runner.env`

2. ユーザ設定  
   `~/.runner.env`

優先順位

```text

project runner.env
↓
user runner.env
↓
default

```

これにより、プロジェクトごとの runtime 設定を上書きできるようにする。

---

## 付録A 実装留意事項

本付録は runner の実装時の参考情報であり、
仕様の一部ではない。

### A.1 一時ファイル

.run の本文は一時ファイルとして展開して実行する。

`#script` の場合は、選択された OS ブロックの本文だけを一時ファイルへ展開する。

実装時は以下に注意する。

- OS の一時ディレクトリを利用する
- 同時実行時に衝突しない一意なファイル名を使用する
- 実行終了後に削除する
- 異常終了時も可能な限り削除する
- dry-run では生成しない

### A.2 プロセス実行

shell を経由せず直接プロセスを起動すること。

shell を使用すると

- OS 依存
- クォート差
- セキュリティ問題

が発生する可能性がある。

### A.3 Windows 互換性

Windows では実行中ファイルの削除ができない場合がある。
削除タイミングに注意すること。

### A.4 command 分割

runtime の command 分割は仕様に従って行うこと。
shell 展開は行わない。

---

### A.5 runner 再帰呼び出し

.run 内から runner を呼び出す場合、PATH に依存しない実行が望ましい。

実装では runner 自身の絶対パスを使用するか、実行可能な方法を確保すること。

---
