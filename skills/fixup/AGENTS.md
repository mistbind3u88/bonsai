# fixup

## 前提ツール

- [git](https://git-scm.com/)

## 責務の境界

- このスキルは、既に行った変更を既存コミットへ fixup として反映する責務を持つ
- デフォルトではfixup commitを作成し、明示的に `--history` が指定された場合だけ `git history fixup` を使う
- 大きな履歴整理、複数 fixup の一括整理、autosquash、push、PR 更新は担わない。必要なら `/commit`、`/push`、`/gh-edit` へ委譲する
