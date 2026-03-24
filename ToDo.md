# ■ ToDo（現時点整理版）

## 🟢 Phase 1: コア安定（ほぼ完了）

👉 これはもう完了扱いでOK

* `.env` 前提の統一（`runner.test.env`）
* テスト資産の固定化（`test/` + `copyTestFile`）
* `go test ./...` が通る
* `runfile.run` のデフォルト動作
* `--list / --check / --dry-run` 基本動作

---

## 🟡 Phase 2: 仕様と実装のズレ解消（今やる or すぐやる）

### Issue 1: dry-run 出力仕様の整理

#### 概要

仕様書と実装の差がある

#### 現状

* 実装: `[runner] command: ...` のみ
* 仕様書: OS表示 / script本文 / 区切り線あり

#### 方針（おすすめ）

👉 **仕様書を現実に寄せる（今の実装を正とする）**

#### 対応

* 仕様書の dry-run セクションを簡略化
* 「将来拡張」として脚注に逃がす

---

### Issue 2: install 手順の明確化

#### 概要

自己更新の挙動が分かりにくい

#### 現状

* `runner.exe install` → 場合によって失敗
* `./bin/runner.exe install` → 安定

#### 方針

👉 **install は配布バイナリから実行する**

#### 対応

仕様書に明記：

```txt
インストールは ./bin/runner.exe install を使用する
```

---

## 🟡 Phase 3: 機能未実装（今回の2件）

### Issue 5: `--dry-run=all` の実装

#### 概要

全OS分を順次表示

#### 現状

* オプション解析はある
* 実処理は未実装

#### 対応案

* 各OSごとに plan を生成
* 順番（windows → linux → macos）で表示

#### 優先度

👉 低（将来）

---

### Issue 6: `#filename.ext` の扱い整理

#### 概要

headerの意味が未確定

#### 現状

```go
case headerKindFilename:
    // extensionとして処理
```

👉 `#.ext` と同義

---

#### 対応案

### A: extensionの糖衣として固定（おすすめ）

```txt
#foo.py == #.py
```

👉 シンプル・現状維持

---

### B: filenameに意味を持たせる

```txt
#main.py
```

👉 将来的に：

* ファイル名固定
* 実行対象明示

---

#### 優先度

👉 中（設計判断）

---

## 🔵 Phase 4: UXブラッシュアップ（最後にやる）

### Issue 7: install の自己検知リダイレクト

#### 概要

```txt
runner.exe install
```

を検知して

```txt
./bin/runner.exe install
```

へ誘導

---

#### 対応案

* A: メッセージ表示のみ
* B: 自動リダイレクト

---

#### 優先度

👉 低（最後）

---

### Issue 8: install を `.bat` に逃がすか検討

#### 概要

Windowsのexeロック回避

---

#### 選択肢

* A: `.run` で完結
* B: `install.bat` に委譲

---

#### 優先度

👉 低（設計ポリシー次第）

---

## ■ 最終まとめ

今の状態を一言でいうと👇

👉 **コアは完成、仕様の磨きとUXの調整フェーズ**

---

## ■ 優先順位（超重要）

### 今やる

* Issue 1（dry-run仕様を現実に合わせる）
* Issue 2（install手順明記）

---

### 後でやる

* Issue 5（dry-run=all）
* Issue 6（filename仕様）

---

### 最後

* Issue 7（リダイレクト）
* Issue 8（.bat戦略）

---
