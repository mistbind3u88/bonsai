# codex-review

## 前提ツール

- [git](https://git-scm.com/)
- [codex](https://github.com/openai/codex) — `npm install -g @openai/codex`

## 責務の境界

- codex CLI に読み取り専用の差分レビューを依頼し、結果を整理してユーザー判断を確認する。
- 過去の review 判断の参照と記録は `/review-log`、結果通過時のタグ設置は `/mark review-cross` へ委譲する。
- 指摘への修正、commit、push、fallback 判断は担わず、呼び出し元または関連 skill が担う。
