# codex-review

## 前提ツール

- [git](https://git-scm.com/)
- [codex](https://github.com/openai/codex) — `npm install -g @openai/codex`

## 責務の境界

- 入力対象は、base/headの差分、未コミット変更、レビュー観点、既知のレビュー判断、レビュー用モデル設定。
- 出力対象は、codex CLIによる読み取り専用の差分レビュー結果と、ユーザー判断を確認するための指摘一覧。
- 停止条件は、対象範囲・認証・実行権限を確認できない場合、CLIレビューが完了しない場合、または指摘へのユーザー判断が揃っていない場合。
- primary model が busy / capacity / rate-limit 系で開始できない場合は、`fallback.config.toml` の設定で 1 回だけ再実行する。
- 過去の review 判断の参照と記録は `/review-log`、結果通過時のタグ設置は `/mark review-cross` へ委譲する。
- 指摘への修正、commit、push、cross-review 全体の fallback 判断は担わず、呼び出し元または外側の workflow が担う。
