# mark

## 前提ツール

- [git](https://git-scm.com/)

## 責務の境界

- このスキルは、現在の HEAD に対する品質チェック通過状態を軽量タグとして設置・確認・削除する責務を持つ。
- lint、build、test、doc-check、rule-check、review の実行判断と実行順序は `/check` が担う。
- 個別の検査実行は `/static-check`、`/unit-test`、`/doc-check`、`/rule-check`、`/subagent-review`、`/claude-review`、`/codex-review` へ委譲する。
- コミットや push の作成は `/commit`、`/push` へ委譲する。
