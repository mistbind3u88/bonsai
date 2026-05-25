---
name: claude-review
description: Claude CLI を使って変更差分のコードレビューを実行する。Codex 上で作業した差分を別エージェント視点で確認するために使う。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git diff:*) Bash(git rev-parse:*) Bash(printf:*) Bash(claude -p:*) Read
---

# claude-review

変更差分を Claude CLI でレビューに出す。Codex 上で作業した差分を、別エージェント視点で確認するために使う。

## 手順

### 1. レビュー対象を特定する

`$ARGUMENTS` が指定されていれば base ref として使う。未指定なら `main` を使う。

```bash
git status -s
git log --oneline <base>..HEAD
git diff --stat <base>..HEAD
git diff <base>...HEAD
```

未コミット変更がある場合は、以下も確認する。

```bash
git diff --staged
git diff
```

### 2. レビューコンテキストを作る

Claude に渡すプロンプトには以下を含める。

- レビュー対象の base ref と HEAD のハッシュ値（例: `レビュー対象: abc1234..def5678`）
- 未コミット変更の有無
- 変更概要
- 主な変更点
- 設計判断やトレードオフ
- 関連 Issue/PR があれば URL
- 既知の指摘事項と対応不要判断があればその理由
- hunk-level の差分

### 3. Claude にレビューを依頼する

`claude -p` を使い、非対話でレビュー結果を出力させる。

```bash
printf '%s' "<レビューコンテキスト + レビュー指示>" | claude -p \
  --permission-mode dontAsk \
  --tools=Read,Bash \
  --allowedTools="Read Bash(git status:*) Bash(git diff:*) Bash(git log:*) Bash(git show:*) Bash(git grep:*) Bash(rg:*) Bash(sed:*) Bash(nl:*) Bash(cat:*)"
```

Claude には、差分を読み、バグ・回帰・仕様漏れ・テスト不足を優先して指摘させる。単なる好み、過度なリファクタ提案、根拠のない推測は避けさせる。

`--tools` と `--allowedTools` は `--tools=Read,Bash`、`--allowedTools="..."` のように `=` 付きで指定する。空白区切りにすると後続のプロンプトやオプションが引数として解釈される場合がある。

`--allowedTools` では、レビューに必要な読み取り系コマンドだけを許可する。Claude は必要に応じて周辺コードを調査できるが、編集・コミット・push・削除などの変更操作はレビュー対象外にする。

レビュー指示には以下を含める。

- 対象ファイルと行番号を示す
- 指摘内容を人間にわかりやすく要約する
- 深刻度を `要修正` / `検討推奨` / `軽微` のいずれかで示す
- 指摘がない場合は「指摘なし」と明記する

### 4. 結果を報告する

Claude の出力をそのまま転記せず、各指摘を以下の形式で整理して報告する。

- 対象ファイルと行番号
- 指摘内容の要約
- 深刻度（`要修正` / `検討推奨` / `軽微`）

報告後、全ての指摘についてユーザーの判断（修正する / 対応不要）を確認する。ユーザーから全指摘への回答を得るまで次のステップに進まない。

### 5. レビュー完了タグを設置する

指摘がない場合、または全指摘へのユーザー判断が揃った場合は、スキル `/mark review-cross` を実行する。

### Claude CLI が実行できない場合

`claude -p` が権限・実行環境の理由で開始または完了できない場合は、失敗を呼び出し元（`/check`）へ報告する。fallback の有無・種別の判断は `/check` 側に委ねる。

## 注意

- Codex 上で作業した差分を Claude に確認させる用途を基本とする
- Claude Code 上で作業している場合は `/codex-review` を使う
- 外部モデルに差分を渡すため、機密情報や送信範囲に注意する
- CLI 実行不可時は失敗を呼び出し元へ報告し、fallback の判断は `/check` 側に委ねる
