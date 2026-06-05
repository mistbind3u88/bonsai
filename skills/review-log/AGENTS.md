# review-log

## 前提ツール

- [git](https://git-scm.com/) — branch 名、worktree root、作業状態を確認するために使う。

## 責務の境界

- review 系 skill のフィードバックループで、過去の指摘・判断・対応状態を worktree-local な review log に記録する。
- 次回 reviewer に渡す `次回レビューへの引き継ぎ` の要約を整える。
- 判断の主体を `開発者` と `エージェント` で区別して記録する。
- review の実行、指摘の採否判断、修正、mark、commit、push は担わず、呼び出し元の review skill または `/check` が担う。
