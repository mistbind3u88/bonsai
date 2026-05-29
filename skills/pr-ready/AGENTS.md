# pr-ready

## 前提ツール

- [git](https://git-scm.com/)
- [GitHub CLI](https://cli.github.com/) — `brew install gh`
- `mark.sh` — このリポジトリの `tools/` を PATH に含める

## 責務の境界

- PR がレビュー可能状態かを確認し、draft PR を Ready for Review に切り替える
- 未コミット変更の整理は `/commit` に委譲する
- push と push 前チェックは `/push` に委譲する
- 品質チェックは `/check` に委譲する
- PR/Issue 本文の作成・更新は `/gh-edit` に委譲する
- PR diff 上の意図補足は `/diff-comment` に委譲する
- CI の監視と失敗分析は `/watch-ci` に委譲する

## 判定方針

- このスキルは不足を自動解消せず、不足項目と次に使うスキルを報告して停止する
- `review-cross` が現在 HEAD で通っている状態を Ready for Review の品質前提にする
- ready PR でも、概要欄・diff comment・CI・check の不足確認を行い、レビュー準備の完了を報告する
