# review-reminders

## 前提ツール

- [GitHub CLI](https://cli.github.com/) — `gh` で GitHub の PR / Issue 情報を取得する

## 責務の境界

- 自分が author の Open PR を取得し、レビューリマインド候補を分類する。
- 候補 PR の前提 PR、ブロッカー Issue、解決対象 Issue、関連 PR / Issue を整理する。
- `tanaoroshi` の結果がある場合は、依存関係の補助情報として参照する。
- 複数リポジトリの Issue / PR 全体棚卸しは `/tanaoroshi` を案内する。
- PR / Issue の詳細取得だけが必要な場合は `/gh-read` を案内する。
- GitHub へのコメント投稿や PR 本文更新が必要な場合は、目的に応じて `/gh-edit`、`/pr-progress`、`/reply-review`、`/diff-comment` を案内する。
- レビュー指摘への対応が必要な場合は `/respond` または `/collect-feedback` を案内する。
