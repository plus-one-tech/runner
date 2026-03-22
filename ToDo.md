# 🧾 runner 仕様差・改善タスク一覧

## 🟢 Issue 1: 通常 `#pwsh` に `.ps1` 拡張子を付与

### 概要

通常モードの `#pwsh` 実行時に temp ファイルに `.ps1` が付与されていない。

### 現状

* `#script` → `.ps1` 付与される
* 通常 `#pwsh` → 付与されない

### 問題

* PowerShell 実行時に拡張子なしは不自然
* 一貫性がない

### 対応

`resolveNormalHeader()` で `tempExtForRuntime()` を使用する

---

## 🟢 Issue 2: `var.*` 変数展開の実装

### 概要

`runner.env` の `var.*` が読み込まれているが、未使用

### 現状

* `cfg.vars` に格納されている
* `.run` 内では展開されない

### 問題

仕様未達

### 対応

* `${var.xxx}` を本文展開
* 未定義変数はエラー

---

## 🟢 Issue 3: `--check` を全 OS ブロック対象にする

### 概要

現在は1OSのみ検証

### 現状

* 現在OS または dry-run OS のみチェック

### 問題

* script内の他OSブロックの不備を検出できない

### 対応

* `#script` の全OSブロックをループ
* runtime解決まで検証

---

## 🟡 Issue 4: `--dry-run` 出力の拡張

### 概要

現在は command のみ表示

### 仕様

```
[runner] os: windows
--- script ---
...
--- end ---
```

### 現状

```
[runner] command: ...
```

### 対応

* スクリプト本文表示
* OS表示
* `--dry-run=all` 対応

---

## 🟡 Issue 5: `--dry-run=all` の実装

### 概要

全OS分を順次表示

### 現状

未対応（オプション解析のみ）

---

## 🟡 Issue 6: `#filename.ext` の扱い整理

### 概要

現在は extension と同義

### 現状

```go
case headerKindFilename:
    // extensionとして処理
```

### 対応案（どちらか）

* A: extensionの糖衣として明示（現状維持）
* B: filenameとして意味を持たせる

---

## 🔵 Issue 7: Windows bash path 補正（低優先）

### 概要

Windows + bash の temp path 問題

### 方針

* 現在は非推奨なので後回し可

---

# 🎯 優先順位まとめ

## 最優先（すぐやる）

1. `.ps1` 拡張子
2. var 展開

## 次

3. `--check` 強化

## 余裕があれば

4. dry-run 系

---

# 🚀 次にやるなら

👉 **Issue 1（.ps1付与）から着手**

理由：

* 影響範囲が小さい
* すぐ終わる
* 今回触った領域

---

# ひとこと

ここまで整理できてると、

👉 「思いつき開発」じゃなくて
👉 **設計駆動開発**

に完全に移行してます 😏
