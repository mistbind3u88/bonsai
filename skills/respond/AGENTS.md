# respond

## 前提ツール

- [git](https://git-scm.com/)
- [gh](https://cli.github.com/) — `brew install gh`

## 責務の境界

- `/collect-feedback` → コード修正と `/fixup` → `/ship --skip-gh-edit` → `/reply-review` の順次実行を担う
- 指摘の収集・修正の記録・push・リプライの各操作は対応するスキルが担う。respond 自身はそれらの中身を実装しない
- いずれかのステップが失敗したらワークフロー全体を停止する
