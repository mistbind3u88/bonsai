# Codex Global Settings Sample

このディレクトリは、Codex のグローバル設定に入れる運用ルールのサンプルです。

ここにある `AGENTS.sample.md` は、このリポジトリへ適用する設定ではありません。利用する場合は、自分の Codex グローバル設定に必要な項目だけを写し、環境に合わせてパスや承認方針を調整してください。

## 使い方

1. `AGENTS.sample.md` の内容を読み、必要な項目だけを自分の Codex グローバル設定へ移す。
2. `<path-to-skills-repo>` や `<owner/repo>` のようなプレースホルダを自分の環境に置き換える。
3. token、API key、private repository URL、社内ホスト名、個人名は自分のローカル設定側だけで扱う。
4. `PATH` への追加は shell 設定側で行い、リポジトリには `<path-to-skills-repo>` のようなプレースホルダを書く。

コピー先では `AGENTS.sample.md` の内容を Codex が読むグローバル設定（例: `~/.codex/AGENTS.md`）として配置します。このリポジトリでは自動適用を避けるため、サンプル名を `AGENTS.sample.md` にしています。

## 含めているもの

- skill の前提条件を既存 skill で解決する運用
- 後続で承認が必要な操作を実行前にまとめて確認する運用
- サブエージェントの起動前確認と再利用方針
- cross-review で `claude -p` を read-only で使う方針
- GitHub CLI と token の安全な扱い
- `<path-to-skills-repo>/tools` を PATH に追加する考え方
- review log を worktree-local に置いて reviewer へ渡す考え方
- 単一リポジトリ内のドキュメントで責務を分け、どの文書を基準にするかを決める方針

## ローカル設定側で扱うもの

- 実 token、API key、secret
- private repository の具体名
- 社内 URL、社内ホスト名、内部ネットワーク情報
- 個人のホームディレクトリの絶対パス
- このリポジトリへ直接適用する `AGENTS.md` ルール
