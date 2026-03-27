# ■ ToDo（現時点整理版）

## 🟢 Phase 1: コア安定（完了）

👉 完了扱い

* `.env` 前提の統一（`runner.test.env`）
* テスト資産の固定化（`test/` + `copyTestFile`）
* `go test ./...` が通る
* `runfile.run` のデフォルト動作
* `--list / --check / --dry-run` 基本動作

---

## 🟢 Phase 2: 仕様確定（完了）

👉 実装と仕様が一致している状態

---

### Issue 1: dry-run 出力仕様

#### 現在の仕様（実装）

* 実行コマンドを表示
* script 実行の場合は、シェルに渡すスクリプト本文を表示
* dry-run では temp ファイルは作成しない

#### 状態

👉 完了

---

### Issue 6: `#filename.ext` の扱い

#### 概要

header の意味の整理

#### 現在の仕様（実装）

```
#foo.py == #.py
```

👉 extension の糖衣として扱う

#### 状態

👉 方針確定（現状維持）

※ 将来拡張の余地あり

---

## 🟡 Phase 3: 未実装機能（将来対応）

---

### Issue 5: `--dry-run=all`

#### 概要

全OS分を順次表示

#### 現状

* オプション解析あり
* 実処理なし（単一OSのみ）

#### 対応案

* 各OSごとに plan を生成
* 表示順：windows → linux → macos

#### 優先度

👉 低

---

## 🟡 Phase 4: 運用・仕様補足

---

### Issue 2: install 手順の明確化

#### 概要

自己更新の挙動の整理

#### 現状

* `runner.exe install` → 失敗する場合あり
* `./bin/runner.exe install` → 安定

#### 理由

* Windowsでは実行中のexeは上書きできない

#### 対応

仕様書に明記：

```
インストールは ./bin/runner.exe install を使用する
```

#### 優先度

👉 中（ドキュメント対応）

---

## 🔵 Phase 5: UXブラッシュアップ（最後）

---

### Issue 7: install の自己検知リダイレクト

#### 概要

```
runner.exe install
```

を検知して

```
./bin/runner.exe install
```

へ誘導

#### 対応案

* A: メッセージ表示
* B: 自動リダイレクト

#### 優先度

👉 低

---

### Issue 8: install を `.bat` に逃がすか検討

#### 概要

Windowsのexeロック回避

#### 選択肢

* A: `.run` で完結（現状）
* B: `.bat` に委譲

#### 優先度

👉 低

---

## ■ 最終まとめ

👉 **コア実装は完成、残りは補助機能とUXのみ**

---

## ■ 優先順位

### 今やる

* Issue 2（install手順の明文化）

---

### 将来

* Issue 5（dry-run=all）

---

### 最後

* Issue 7（リダイレクト）
* Issue 8（.bat戦略）

---
