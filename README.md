# runner

[English version](./README-en.md)

runner は、スクリプトやタスクを実行するための軽量なコマンドランナーです。

1つのコマンドでソースファイルやタスクをシンプルに実行できます。

## Example

次の例は、`runner` の実行イメージです。

```text
> runner hello.py
[runner] python hello.py
Hello Runner

> runner build
[runner] bash build.sh
Building...
Done.

> runner
[runner] bash runfile.sh
Running default task...
> 
```

runner は、スクリプトの実行やビルドなどのタスクを、1つのコマンドで実行できるツールです。

`runfile.run` に定義されたタスクを、デフォルトの実行対象として扱うこともできます。

実行時には、実際に呼び出されるコマンドを表示してから処理を実行します。

## 概要

`runner` は、スクリプト実行とタスク実行をシンプルかつ予測可能にすることを目的としたツールです。

スクリプトの実行やビルド、テストなどの繰り返し作業を、1つのコマンドでまとめて実行できます。

## 背景

開発では、コードを書いて、ビルドして、テストして……といった同じ作業の繰り返しが大半です。  
それらを簡単に実行できたら便利だと考えました。

最近では C# もコマンドラインからソースを直接実行できるようになり、  
インタプリタのように扱えると面白いのでは、という発想もあります。

たとえば `runner` をエイリアスで `run` にすると、次のように使えます。

```text
>run hello.cs
hello world
```

昔の環境のように、気軽に実行できる感覚を目指しています。

## 設計ポリシー

* **統一された実行** – スクリプトとタスクを同じコマンドで実行
* **透明性** – 実際に実行されるコマンドを常に表示
* **最小設計** – 複雑なDSLや重い設定を持たない

一般的なタスクランナーと異なり、新しいスクリプト言語は導入しません。`.run` ファイルには通常のプログラムコードを書くだけです。

## 基本的な使い方

### スクリプトを実行

```bash
runner hello.py
```

### タスクを実行

```bash
runner build
```

`build.run` が存在する場合：

```text
#bash
dotnet run ./src/hello.cs
```

`.run` ファイルに書かれている内容がそのまま実行されます。

### デフォルトタスクを実行

```bash
runner
```

`runfile.run` が存在する場合、それが実行されます。

## 主なオプション

### 実行内容の確認（dry-run）

```bash
runner --dry-run build.run
runner --dry-run=windows install.run
runner --dry-run=all install.run
```

### 実行せずに検証

```bash
runner --check build.run
```

### 設定ファイルを指定

```bash
runner --env ./runner.env install.run
```

### 実行可能な `.run` ファイルの一覧

```bash
runner --list
```

## `.run` ファイル

`.run` ファイルはシンプルな実行タスクファイルです。

例：

```text
#python
print("Hello Runner")
```

サポートされているヘッダ：

```text
#python
#program.py
#.py
#script
```

それ以降は通常のプログラムコードとして実行されます。

## 設定

`runner.env` により、ランタイムと拡張子を実際のコマンドにマッピングします。

例：

```text
runtime.python=python
runtime.bash=bash
runtime.node=node

ext.py=python
ext.js=node
ext.sh=bash
```

デフォルトではユーザー設定ディレクトリから読み込まれます。
`--env` オプションで明示的に指定することもできます。

## インストールタスク

`install.run` と `runner.env` を用意した場合は、以下のコマンドを実行します。

```powershell
.\bin\runner --env ./runner.env install.run
```

### 補足（Windows）

Windowsでは、実行中の `.exe` ファイルは自分自身を上書きできません。

インストールや更新を行う場合は、必ず別の場所にあるバイナリ（例: `.\bin\runner`）を使用してください。

## 仕様

詳細な仕様はこちら：

**docs/runner-spec.md**

## ステータス

初期リリース（v0.1.0）です。

コア機能は実装済みで利用可能ですが、今後も改善・拡張を予定しています。

## コントリビューション

Pull Request は歓迎しますが、プロジェクトの設計思想に合わない変更は受け付けない場合があります。

## ライセンス

MIT ライセンス
