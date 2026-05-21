# takeover

## 前提ツール

- [git](https://git-scm.com/)
- [gh](https://cli.github.com/) — `brew install gh`

## 責務の境界

- 前セッションのコンテキスト（ブランチ状態・Issue/PR・プランファイル・タスクドキュメント・差分概要）の収集と、引き継ぎサマリの提示を担う
- Issue/PR の取得は `/gh-read`、タスクドキュメントの場所解決は `/taskdoc-locate` へ委譲する
- 収集と提示に限定し、コードの変更や作業の実行そのものは行わない
