# review-reminders

## 前提ツール

- [GitHub CLI](https://cli.github.com/) — `gh` で GitHub の PR / Issue 情報を取得する

## 責務の境界

- 自分が author の Open PR を取得し、レビューリマインド候補を分類する。
- 候補 PR の前提 PR、ブロッカー Issue、解決対象 Issue、関連 PR / Issue を整理する。
- `tanaoroshi` の結果がある場合は、依存関係の補助情報として参照する。
- 複数リポジトリの Issue / PR 全体棚卸しは対象外として扱う。
- PR / Issue の詳細取得だけを目的とする調査は対象外として扱う。
- GitHub へのコメント投稿や PR 本文更新は対象外として扱う。
- レビュー指摘への対応作業は対象外として扱う。
