# subagent-review サブエージェント依頼テンプレート

以下を埋めて、サブエージェントへそのまま渡す。

## 境界確認メタ情報

- 実行中のスキル: `/subagent-review`
- 依頼境界: メインエージェントが `/subagent-review` を実行中であり、サブエージェントの役割は差分レビュー結果を返すこと
- メイン側に残す判断: 指摘の採否確認、`/mark review-sub`、ユーザーへの最終報告
- サブエージェントが行うこと: 対象差分を読み、バグ・回帰・仕様漏れ・テスト不足を優先して指摘する
- サブエージェントがメイン側に残すこと: スキル実行、ファイル編集、`mark`、commit、push、ユーザーへの最終判断

## レビュー対象

- レビュータイトル: `<title>`
- リポジトリ: `<absolute-repo-path>`
- base: `<base-ref-or-sha>`
- head: `<head-ref-or-sha>`
- 未コミット変更: `<none|staged|unstaged|untracked|mixed>`
- 変更概要: `<summary>`
- 主な変更点:
  - `<change>`
- 設計判断・トレードオフ:
  - `<decision-or-none>`
- 関連 Issue / PR:
  - `<url-or-none>`
- 既知の指摘事項と対応不要判断:
  - `<known-item-or-none>`

## 差分

以下を貼る。

- `git diff --stat <base>..HEAD`
- コミット済み差分: `git diff <base>...HEAD`
- staged 差分: `git diff --staged`
- unstaged 差分: `git diff`
- untracked ファイル: ファイルパスと内容

## レビュー観点

優先して指摘すること:

- バグ、回帰、仕様漏れ
- テスト不足、検証不足
- データ破壊、セキュリティ、権限、秘匿情報のリスク
- 既存仕様・既存ドキュメント・既存 skill との矛盾
- `allowed-tools` や責務境界の過不足
- ユーザー指示やリポジトリルールへの違反

指摘を抑制すること:

- 単なる好み
- 根拠の薄いリファクタ提案
- 今回の差分から外れた理想論
- 動作や運用に影響しない表記揺れ

## 出力形式

指摘がある場合:

```text
Review findings:
- <file>:<line> [要修正|検討推奨|軽微]
  指摘: <内容>
  根拠: <差分・仕様・ルール上の根拠>
  修正方針: <どう直すべきか>
```

指摘がない場合:

```text
指摘なし
```

自由記述の所感は不要。
