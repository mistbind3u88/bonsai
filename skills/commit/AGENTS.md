# commit

## 前提ツール

- [git](https://git-scm.com/)
- [gh](https://cli.github.com/) — `brew install gh`

## 責務の境界

- このスキルは、変更内容をレビュー可能な単位へ整理し、適切なコミットメッセージで `git commit` / `git commit --amend` を実行する責務を持つ。
- 既存コミットへの修正追加は `/fixup` へ委譲する。
- コミット後の品質確認は `/check` へ委譲する。
- push、PR 作成・更新、CI 監視はそれぞれ `/push`、`/gh-edit`、`/watch-ci` へ委譲する。
