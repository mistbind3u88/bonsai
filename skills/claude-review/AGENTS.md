# claude-review

- Codex 上で作業した差分を、Claude CLI による別エージェント視点で確認する。
- レビューではバグ、回帰、仕様漏れ、テスト不足を優先して指摘する。
- Claude CLI が実行できない場合は失敗として呼び出し元へ返す。スキル内部での自動 fallback は持たない（fallback の判断は `/check` 側）。

## 前提ツール

- `claude` CLI（Claude Code）がインストールされ、`claude -p` の非対話実行と認証が利用できること。
- Claude に渡す tools は読み取り・差分確認に必要なものだけに限定すること。Bash は `git diff`、`git log`、`git show`、`git grep`、`rg`、`sed`、`nl`、`cat` などの調査用途に限り、編集・コミット・push・削除には使わせないこと。
