# ship

## 前提ツール

- [git](https://git-scm.com/)
- [gh](https://cli.github.com/) — `brew install gh`

## 責務の境界

- `/commit` → `/check` → `/push` → `/gh-edit` → `/watch-ci` の順次実行と、PR 有無の確認・autosquash 時の分岐という薄い制御を担う
- コミット・チェック・push・PR 概要欄の更新・CI 監視の各操作は対応するスキルが担う。ship 自身はそれらの中身を実装しない
- いずれかのステップが失敗したらワークフロー全体を停止する
