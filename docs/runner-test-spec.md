# runner 受け入れテスト一覧

## 1. 基本起動

### AT-001 ファイル指定実行

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

### AT-002 `.run` ファイル明示指定

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

### AT-003 名前指定実行

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

### AT-004 引数なし実行

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

### AT-005 拡張子あり target は直接実行

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

### AT-006 拡張子なし target は `.run` を探索

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

### AT-007 拡張子なし target が見つからない場合

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

### AT-008 dry-run は実行しない

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

### AT-009 dry-run で `.run` を解析する

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

---

## 4. オプション

### AT-010 `--help`

**操作**

```bash
runner --help
```

**期待結果**

* usage を表示する
* 終了コードは 0

---

### AT-011 `--version`

**操作**

```bash
runner --version
```

**期待結果**

* バージョンを表示する
* 終了コードは 0

---

### AT-012 `--list`

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
* 終了コードは 0

---

### AT-013 `--list` は非再帰

**前提**
`sub/deploy.run` が存在する。

**操作**

```bash
runner --list
```

**期待結果**

* `sub/deploy.run` は表示しない

---

### AT-014 無効オプション

**操作**

```bash
runner --check
```

**期待結果**

* エラー終了する
* `unknown option: --check` を表示する
* 非0を返す

---

### AT-015 target 後ろのオプションは不許可

**操作**

```bash
runner hello.py -n
```

**期待結果**

* エラー終了する
* 非0を返す

---

## 5. `.run` ヘッダ

### AT-016 runtime 指定ヘッダ

**前提**
`hello.run` の先頭行が `#python` である。

**操作**

```bash
runner hello.run
```

**期待結果**

* `runtime.python` を参照して runtime を解決する

---

### AT-017 仮想ファイル名指定ヘッダ

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

### AT-018 拡張子指定ヘッダ

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

### AT-019 ヘッダなし

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

### AT-020 不正ヘッダ

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

### AT-021 ヘッダ前空行は不許可

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

### AT-022 ヘッダ後空行は許可

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

---

## 6. `.run` 一時ファイル

### AT-023 `.run` は一時ファイルに展開する

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

### AT-024 一時ファイルは OS 一時ディレクトリに生成

**前提**
`hello.run` が存在する。

**操作**

```bash
runner hello.run
```

**期待結果**

* 一時ファイルは OS の一時ディレクトリに生成される

---

### AT-025 一時ファイルは実行終了後に削除

**前提**
`hello.run` が存在する。

**操作**

```bash
runner hello.run
```

**期待結果**

* 実行終了後に一時ファイルを削除する

---

### AT-026 dry-run では実行しない

**前提**
`hello.run` が存在する。

**操作**

```bash
runner -n hello.run
```

**期待結果**

* 一時ファイル名は決定する
* 実行はしない

---

## 7. runner.env 読み込み

### AT-027 基本読み込み

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

### AT-028 前後空白を無視

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

### AT-029 コメント行を無視

**前提**
`runner.env` にコメント行がある。

**操作**

```bash
runner hello.py
```

**期待結果**

* コメント行は無視する

---

### AT-030 行内コメントはサポートしない

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

* 行内コメントとしては扱わない
* 実装仕様に従って不正入力として扱うか、そのまま値に含めるかをテストで固定する必要がある

---

### AT-031 重複キーは後勝ち

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

### AT-032 key は大小文字を区別

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

## 8. command 分割

### AT-033 command を空白で分割

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

### AT-034 ダブルクォートで1引数化

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

### AT-035 エスケープ `\"`

**前提**
`runner.env` にダブルクォートを含む値がある。

**操作**
適切な `runner.env` を与えて実行する

**期待結果**

* `\"` をダブルクォートとして解釈する

---

### AT-036 エスケープ `\\`

**前提**
`runner.env` にバックスラッシュを含む値がある。

**操作**
適切な `runner.env` を与えて実行する

**期待結果**

* `\\` をバックスラッシュとして解釈する

---

### AT-037 不正クォート

**前提**
`runner.env` に閉じ忘れたクォートがある。

**操作**

```bash
runner hello.py
```

**期待結果**

* エラー終了する
* 非0を返す

---

### AT-038 shell 展開なし

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

### AT-039 extension 未定義

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

### AT-040 runtime 未定義

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

### AT-041 正常終了

**前提**
正常終了するスクリプトがある。

**操作**

```bash
runner hello.py
```

**期待結果**

* 終了コード 0 を返す

---

### AT-042 実行プロセスの非0終了を伝播

**前提**
対象スクリプトが `exit 5` 相当で終了する。

**操作**

```bash
runner fail.py
```

**期待結果**

* runner も終了コード 5 を返す

---

### AT-043 runner 自体のエラー

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

### AT-044 `.run` UTF-8

**前提**
UTF-8 で書かれた `.run` がある。

**操作**

```bash
runner hello.run
```

**期待結果**

* 正常に読み込める

---

### AT-045 `runner.env` UTF-8

**前提**
UTF-8 で書かれた `runner.env` がある。

**操作**

```bash
runner hello.py
```

**期待結果**

* 正常に読み込める

---

### AT-046 UTF-8 BOM を無視

**前提**
UTF-8 BOM 付き `.run` または `runner.env` がある。

**操作**
実行する

**期待結果**

* BOM を無視して正常に読み込める

---

### AT-047 LF を許可

**前提**
LF 改行の `.run` / `runner.env` がある。

**操作**
実行する

**期待結果**

* 正常に読み込める

---

### AT-048 CRLF を許可

**前提**
CRLF 改行の `.run` / `runner.env` がある。

**操作**
実行する

**期待結果**

* 正常に読み込める

---

## 12. `.run` 実行権限

### AT-049 `.run` に実行権限は不要

**前提**
`build.run` に実行権限が付いていない。

**操作**

```bash
runner build
```

**期待結果**

* 正常に実行できる

---

### AT-050 `.run` の直接実行は想定しない

**前提**
`build.run` が存在する。

**操作**
OS から直接 `build.run` を実行しようとする

**期待結果**

* runner の正式な利用方法ではない
* 受け入れ対象外とする

---
