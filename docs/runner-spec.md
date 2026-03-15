---
date    : 2026-03-15         # 文書の初回作成日
author  : 高野順一            # 文書の主要な作成者
title   : runner 仕様
owner   : 責任者              # 文書の現在の責任者 / 最終更新者
updated : YYYY-MM-DD         # 文書の最終更新日
tags    : SE向け, PM向け, 開発者向け, 新人向け
       - 技術解説, 概念, 実践, ノウハウ, 考察, トラブルシューティング
policy  :
  style   : だ、である調、章は見出しレベル2から、
  scope   : 
  tone    : 技術論は正確に
  purpose : 
note    : 
---

**最終更新日: 2026年03月15日**

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

実際のプロセス実行は行わない。

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

`.run` ファイルは次の構造を持つ。

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

---

### 2.3 ヘッダ仕様

ヘッダは次の 3 形式をサポートする。

```text
#<runtime>
#<filename.extension>
#.<extension>
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

---

### 2.6 本文

2 行目以降は **そのままプログラム本文**とする。

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

`.run` 自体は本文に対して独自構文を持たない。

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

stdin 実行や `-c` / `eval` 方式は採用しない。

理由

* 言語依存を避けられる
* 行番号が安定する
* 実行方式が単純になる

---

### 2.9 一時ファイル生成

runner は `.run` の本文を一時ファイルに書き出し、そのファイルを runtime に渡して実行する。

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

#### 2.9.4 一時ファイル配置

一時ファイルは OS の一時ディレクトリに生成する。
runner は実行終了後に一時ファイルを削除する。

例

/tmp/runner_tmp.py

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

---

## 3. `runner.env` 仕様

### 3.1 目的

`runner.env` は **runner の実行環境を定義する設定ファイル**である。

主な用途

* runtime の定義
* 拡張子と runtime の対応付け
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

### 3.9 まとめ

`runner.env` は次の2種類の定義を持つ。

```text
runtime.<name>=<command>
ext.<extension>=<runtime>
```

これにより runner は

```text
.run header
↓
extension/runtime 解決
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

### 4.4 一時ファイル

runner は .run の実行時に一時ファイルを生成する。

一時ファイルは OS の一時ディレクトリに生成する。  
実行終了後に削除する。

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

実装時は以下に注意する。

- OS の一時ディレクトリを利用する
- 同時実行時に衝突しない一意なファイル名を使用する
- 実行終了後に削除する
- 異常終了時も可能な限り削除する

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
