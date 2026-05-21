# ship

## 前提ツール

- [git](https://git-scm.com/)
- [gh](https://cli.github.com/) — `brew install gh`

## 責務の境界

- ship はワークフローのオーケストレーションを担い、固有の具体操作は持たない。コミット・品質チェック・push・PR 概要欄の更新・CI 監視・進捗コメントの投稿は、`/commit`・`/check`・`/push`・`/gh-edit`・`/watch-ci`・`/pr-progress` へ委譲する
- 実行順序や条件分岐などオーケストレーションの手順は SKILL.md に定義する
