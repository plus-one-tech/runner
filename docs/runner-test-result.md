# runner テスト結果（逆引き表）

| AT     | 概要   | 判定 | 該当テスト  | 備考  |
| ------ | ------------------ | -: | --------------------------------------------------------------------- | ---------- |
| AT-010 | ファイル指定実行    |  ○ | `TestSpec_AT010_FileExecution` | 直接あり |
| AT-011 | `.run` 明示指定実行      |  △ | `TestDryRunDoesNotExecute`, `TestCheckScriptAllOS`, P01, P02, P03     | 間接のみ |
| AT-012 | 名前指定実行      |  ○ | `TestRunNamedTask`, `TestDryRunDoesNotCreateTempFile`, P06     | 直接あり|
| AT-013 | 引数なし実行      |  ○ | `TestBOMAndCRLFRunFile`, P07| 直接あり|
| AT-020 | 拡張子あり直接実行   |  △ | `TestTargetOptionAfterIsRejected`, `TestExtensionNotMapped`, P01      |     |
| AT-021 | `.run` 自動探索 |  ○ | `TestRunNamedTask`, P06     |     |
| AT-022 | target未存在   |  ○ | N02    |     |
| AT-030 | dry-run 実行しない      |  ○ | `TestDryRunDoesNotExecute`  |     |
| AT-031 | dry-run 解析  |  ○ | `TestDryRunPrintsCommand`, `TestDryRunDoesNotCreateTempFile`   |     |
| AT-032 | script dry-run all |  × | -      | 未実装 |
| AT-033 | dry-run windows    |  × | -      |     |
| AT-034 | dry-run linux      |  × | -      |     |
| AT-035 | dry-run macos      |  × | -      |     |
| AT-036 | dry-run all |  × | -      | 未実装 |
| AT-037 | 不正OS |  × | -      |     |
| AT-040 | --help      |  ○ | `TestSpec_AT040_Help`| 直接あり |
| AT-041 | --version   |  ○ | `TestSpec_AT041_Version`    | 直接あり |
| AT-042 | --list      |  ○ | `TestListShowsRunWithoutExtension`, P05   |     |
| AT-043 | list 非再帰    |  ○ | 同上     |     |
| AT-044 | list 空      |  ○ | `TestSpec_AT044_ListEmpty`  | 直接あり |
| AT-045 | target後オプション拒否     |  ○ | `TestTargetOptionAfterIsRejected`  |     |
| AT-046 | --check 基本  |  ○ | `TestCheckScriptAllOS`, P04 |     |
| AT-047 | --check 全OS |  ○ | 同上     |     |
| AT-048 | --check 正常終了|  ○ | 同上     |     |
| AT-049 | check + dry-run    |  ○ | `TestSpec_AT049_CheckWithDryRun`   | 直接あり |
| AT-04A | 無効オプション     |  △ | `TestTargetOptionAfterIsRejected`, N04    |     |
| AT-050 | runtimeヘッダ  |  ○ | `TestResolveNormalHeaderRuntime`   |     |
| AT-051 | 仮想ファイル名     |  × | -      |     |
| AT-052 | extヘッダ      |  ○ | `TestResolveNormalHeaderExt`|     |
| AT-053 | ヘッダなし|  × | -      |     |
| AT-054 | 不正ヘッダ|  × | -      |     |
| AT-055 | ヘッダ前空行禁止    |  × | -      |     |
| AT-056 | ヘッダ後空行許可    |  × | -      |     |
| AT-057 | #script     |  △ | `TestCheckScriptAllOS`, P03 |     |
| AT-058 | script内runtime必須   |  × | -      |     |
| AT-059 | script構造不正  |  × | -      |     |
| AT-05A | OS重複 |  × | -      |     |
| AT-05B | 未知OS |  × | -      |     |
| AT-05C | OSブロックなし    |  ○ | N03    |     |
| AT-060 | temp生成      |  △ | `TestRunNamedTask`, P系全般    |     |
| AT-061 | temp場所      |  ? | -      | 未確認 |
| AT-062 | temp削除      |  ? | -      | 未確認 |
| AT-063 | dry-run tempなし     |  ○ | `TestDryRunDoesNotCreateTempFile`  |     |
| AT-064 | script temp選択      |  × | -      |     |
| AT-065 | pwsh拡張子     |  ○ | `TestTempExtForRuntime`     |     |
| AT-066 | 非選択OS未実行    |  △ | N03    |     |
| AT-070 | env基本|  △ |  | 間接のみ    |
| AT-071 | env空白|  × | -      |     |
| AT-072 | envコメント     |  × | -      |     |
| AT-073 | 行内コメント      |  × | -      |     |
| AT-074 | 重複キー |  × | -      |     |
| AT-075 | 大小区別 |  × | -      |     |
| AT-076 | --env指定     |  ○ | `TestSpec_AT076_EnvSpecified`      | 直接あり |
| AT-077 | env未存在      |  × | -      |     |
| AT-078 | カレントenv無視   |  × | -      |     |
| AT-080 | command分割   |  ○ | `TestCommandQuoteSplit`     |     |
| AT-081 | quote|  ○ | 同上     |     |
| AT-082 | エスケープ"      |  × | -      |     |
| AT-083 | エスケープ\      |  × | -      |     |
| AT-084 | 不正quote     |  ○ | `TestCommandInvalidQuote`   |     |
| AT-085 | shell展開なし   |  × | -      |     |
| AT-090 | ext未定義      |  ○ | `TestExtensionNotMapped`    |     |
| AT-091 | runtime未定義  |  ○ | `TestSpec_AT091_RuntimeNotDefined` | 直接あり |
| AT-100 | 正常終了 |  ○ | P系全般   |     |
| AT-101 | 異常終了伝播      |  ○ | N01〜N04|     |
| AT-102 | runnerエラー   |  ○ | 各エラー系  |     |
| AT-110 | UTF-8|  ○ | `TestBOMAndCRLFRunFile`     |     |
| AT-111 | env UTF-8   |  × | -      |     |
| AT-112 | BOM無視|  ○ | 同上     |     |
| AT-113 | LF許可 |  △ | 多くのテスト |     |
| AT-114 | CRLF許可      |  ○ | 同上     |     |
| AT-120 | 実行権限不要      |  ○ | P系全般   |     |

---

# 🎯 この表から見えること

## 👍 強い領域

* 基本実行（010〜013）
* `.run` 解決
* dry-run 基本
* list / check
* command 分割
* 終了コード
* BOM/CRLF

---

## ⚠️ 弱い領域

* `--dry-run=all` 系（未実装）
* `.run` ヘッダ異常系
* env細かい仕様
* tempファイル管理
* help/version

---

## 🔥 次にやるべき

👉 △しかないところを○にする

特に👇

* AT-010 / 011（1本ずつ追加で締まる）
* AT-060 / 061 / 062（tempまわり）
* AT-070〜076（env）

---
