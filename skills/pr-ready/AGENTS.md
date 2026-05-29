# pr-ready

## 前提ツール

- [git](https://git-scm.com/)
- [GitHub CLI](https://cli.github.com/) — `brew install gh`

## 責務の境界

- PR がレビュー可能状態かを確認し、draft PR を Ready for Review に切り替える
- 未コミット変更の整理が必要な場合は、次に使うスキルとして `/commit` を報告する
- push と push 前チェックが必要な場合は、次に使うスキルとして `/push` を報告する
- 品質チェックが必要な場合は、次に使うスキルとして `/check` を報告する
- PR/Issue 本文の作成・更新が必要な場合は、次に使うスキルとして `/gh-edit` を報告する
- PR diff 上の意図補足が必要な場合は、次に使うスキルとして `/diff-comment` を報告する
- CI の監視と失敗分析が必要な場合は、次に使うスキルとして `/watch-ci` を報告する

## 判定方針

- このスキルは不足を自動解消せず、不足項目と次に使うスキルを報告して停止する
- `review-cross` が現在 HEAD で通っている状態を Ready for Review の品質前提にする
- ready PR でも、概要欄・diff comment・CI・check の不足確認を行い、レビュー準備の完了を報告する
