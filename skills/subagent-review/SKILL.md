---
name: subagent-review
description: 同セッションの general-purpose サブエージェントへ差分レビューを依頼する。コミット時の軽量レビュー用途で /check --review=sub から呼び出される。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git diff:*) Bash(git rev-parse:*) Agent Read
---

# subagent-review

変更差分を同セッションの general-purpose サブエージェントへ渡してレビューさせる。`/codex-review` / `/claude-review` の別エージェント呼び出しと違い、in-session で軽量に走る。`/commit` から `/check --review=sub` 経由で呼ばれ、コミット単位のレビューを担う。

## 手順

### 1. レビュー対象を特定する

`$ARGUMENTS` が指定されていれば base ref として使う。未指定なら `main` を使う。

```bash
git status -s
git log --oneline <base>..HEAD
git diff --stat <base>..HEAD
git diff <base>...HEAD
```

未コミット変更がある場合は以下も確認する。

```bash
git diff --staged
git diff
```

### 2. レビューコンテキストを作る

スキル `/review-log read` を実行し、現在の branch に対応する review log から既知のレビュー判断を取得する。取得した要約は、テンプレートの `既知の論点・対応済み判断` に含めてサブエージェントに渡す。

同ディレクトリの [subagent-prompt-template.md](subagent-prompt-template.md) を埋めて、サブエージェントに渡すプロンプトを作る。

テンプレートを埋める際は、以下を具体値に置き換える。

- `<title>`: レビュー対象を短く表すタイトル
- `<absolute-repo-path>`: 現在の作業ディレクトリの絶対パス
- `<base-ref-or-sha>`: `$ARGUMENTS` で指定された base、または `main`
- `<head-ref-or-sha>`: 現在の `HEAD`
- `<none|staged|unstaged|untracked|mixed>`: 未コミット変更の状態
- `<summary>` / `<change>` / `<decision-or-none>` / `<context-or-none>` / `<known-item-or-none>`: レビューの前提情報
- `<review-log-summary-or-none>`: `/review-log` から取得した既知のレビュー判断。なければ `なし`

テンプレートの差分欄には、取得した `git diff --stat`、hunk-level の差分、staged / unstaged / untracked の内容を貼る。観点と出力形式はテンプレートから削らない。

### 3. サブエージェントにレビューを依頼する

依頼文案を会話中に提示した上で、起動前にスキル `/subagent-check` を実行し、既存サブエージェントの完了状態、再利用可否、起動上限の余剰、依頼文の委譲境界を確認する。`subagent-check` が `OK` を返した場合は新規起動し、`REUSE` / `WAIT` / `ASK_USER` / `FIX_PROMPT` の場合はその判定に従う。

このスキルはメインエージェントが実行し、サブエージェントには差分レビュー結果の作成だけを依頼する。メインエージェントは指摘の整理、ユーザー判断の確認、`/mark` 連携を担う。

サブエージェントへの依頼文案は、テンプレートの「境界確認」を `/subagent-check` 用の確認に使い、「レビュー対象」以降をサブエージェントへ渡す本文として使う。

実行はバックグラウンドで行う（`run_in_background: true`）。完了通知を受けてから結果を読む。

#### サブエージェントの出力形式

- 対象ファイルと行番号
- 指摘内容の要約（人間にわかりやすい言葉）
- 深刻度（`要修正` / `検討推奨` / `軽微`）
- 指摘がない場合は「指摘なし」と明記

### 4. 結果を報告する

サブエージェントの完了通知を受けたら出力を読み、ユーザーに報告する。サブエージェントの出力をそのまま転記せず、各指摘を以下の形式で整理する。

- 対象ファイルと行番号
- 指摘内容の要約
- 深刻度（`要修正` / `検討推奨` / `軽微`）

報告後、全ての指摘についてユーザーの判断（修正する / 対応不要）を確認する。ユーザーから全指摘への回答を得るまで次のステップに進まない。

### 5. review log を更新する

指摘がない場合、または全指摘へのユーザー判断が揃った場合は、スキル `/review-log write` を実行し、review source、対象範囲、指摘、判断、理由、対応状態を現在の実行主体側ログへ追記する。

### 6. レビュー完了タグを設置する

指摘がない場合、または全指摘へのユーザー判断が揃った場合は、スキル `/mark review-sub` を実行する。

## 注意

- in-session のサブエージェントによる軽量レビュー。別エージェント（codex CLI／claude CLI）による独立視点レビューが必要なときは `/codex-review` / `/claude-review` を使う
- サブエージェントは必ずバックグラウンドで起動する。同期実行にしない
- 起動前の `/subagent-check` を省略しない
