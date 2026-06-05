# subagent-review サブエージェント依頼テンプレート

以下を埋めて、差分レビュー依頼としてサブエージェントへ渡す。

## 境界確認

- 実行中のスキル: `/subagent-review`
- work_unit: `subagent-review:<base-ref-or-sha>..<head-ref-or-sha>`
- role: `diff-review`
- reuse_policy: `prefer-reuse`
- 依頼する役割: 対象差分のコードレビュー
- サブエージェントが返すもの: 修正価値のある指摘、根拠、修正方針
- メインエージェントが担うもの: 指摘の採否確認、修正実行、`mark`、commit、push、ユーザー報告

## レビュー対象

- レビュータイトル: `<title>`
- リポジトリ: `<absolute-repo-path>`
- 比較元: `<base-ref-or-sha>`
- 比較先: `<head-ref-or-sha>`
- 未コミット変更: `<none|staged|unstaged|untracked|mixed>`
- 変更概要: `<summary>`
- 主な変更点:
  - `<change>`
- 設計判断・トレードオフ:
  - `<decision-or-none>`
- 関連 Issue / PR / ユーザー指示:
  - `<context-or-none>`
- 明示的に参照すべき文書・根拠:
  - `<reference-path-url-or-none> — <why-this-reference-matters>`
- 既知の論点・対応済み判断:
  - `<known-item-or-none>`
  - `<review-log-summary-or-none>`

明示参照は、レビュー判断の根拠として優先して読む。review log の既知判断は、重複指摘を防ぐための文脈として扱う。各判断にファイルパス、URL、関連 `R0xx` が含まれる場合は、その参照先を確認したうえで、同じ根拠の指摘は繰り返さず、現在の差分で判断が破綻している場合は新しい根拠を示して指摘する。

## 差分

以下を貼る。

- 差分 stat: `<git-diff-stat>`
- コミット済み差分: `<committed-diff-or-none>`
- staged 差分: `<staged-diff-or-none>`
- unstaged 差分: `<unstaged-diff-or-none>`
- untracked ファイル: `<untracked-files-or-none>`

## レビュー観点

優先して確認する。

- バグ、回帰、仕様漏れ、境界条件漏れ
- テスト不足、検証不足、失敗時の扱い
- データ破壊、セキュリティ、権限、秘匿情報、外部サービス利用のリスク
- 既存仕様、公開契約、運用手順、関連ドキュメントとの矛盾
- ユーザー指示、対象リポジトリのルール、変更スコープとの矛盾
- スキル定義の変更では、対象リポジトリのルールに照らした `allowed-tools` と責務境界の過不足
- 今回の差分で生じた保守性、安全性、運用性の問題

次の内容は、具体的な影響がある場合にだけ指摘する。

- 好みの差
- 根拠の薄いリファクタ提案
- 今回の差分から離れた将来課題
- 動作、運用、読み手の理解に影響しない表記揺れ

## 出力形式

指摘がある場合:

```text
Review findings:
- <file>:<line> [要修正|検討推奨|軽微]
  指摘: <内容>
  根拠: <差分・仕様・ルール上の根拠>
  参照: <確認したファイルパス・URL・関連判断ID>
  修正方針: <どう直すべきか>
```

指摘がない場合:

```text
指摘なし
```

自由記述の所感は不要。
