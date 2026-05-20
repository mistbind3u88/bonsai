# check

## 前提ツール

- [git](https://git-scm.com/)

## 責務の境界

`check` は品質チェックの検出・実行を直接行わず、`/static-check`（lint・build）、`/unit-test`（test）、`/doc-check`、review スキルへ委譲して結果を集約する。lint・build・test の実行コマンド検出やスコープ判定は各単機能スキルが担う。`check` 自身が行うのは autosquash 後のタグ引き継ぎ（`git diff`）と、各スキルの成否を受けた `/mark` 設置である。
