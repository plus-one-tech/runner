# runner test

## 概要

このディレクトリには runner の受け入れテスト（.run ベース）を配置しています。

## 実行方法

```powershell
cd test
..\bin\runner.exe .\all-test.run
```

## 個別実行

```powershell
..\bin\runner.exe .\positive-test.run
..\bin\runner.exe .\negative-test.run
```

## 単体テスト

```bash
go test ./...
```

## 前提

* Windows は pwsh を使用
* bash はテスト対象外（Unix系用）
* 設定は runner.test.env を使用
