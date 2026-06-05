# subagent-review

## 前提ツール

- [git](https://git-scm.com/)

## 責務の境界

- 同セッションの general-purpose サブエージェントへ差分レビューを委譲することに責務を限定する
- 過去の review 判断の参照と記録は `/review-log`、起動前確認は `/subagent-check`、結果通過時のタグ設置は `/mark review-sub` へ委譲する
- レビュー指摘の最終判断（修正する／対応不要）と整理・報告はメインエージェントが担う。サブエージェントの役割は差分レビューの実行に限定する
- 同じ base/head 系統の修正後再レビューでは、`work_unit: subagent-review:<base>..<head>` / `role: diff-review` の既存サブエージェントを再利用する

## サブエージェント依頼

サブエージェントへの依頼は [subagent-prompt-template.md](subagent-prompt-template.md) を埋めて作成する。テンプレートの境界確認、レビュー観点、出力形式、`work_unit`、`role`、`reuse_policy` を削らず、差分レビューに必要な base / head / 変更概要 / hunk-level diff を具体値で埋める。
