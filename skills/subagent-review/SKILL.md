---
name: subagent-review
description: 同セッションの `general-purpose` サブエージェントへ差分レビューを依頼する。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git diff:*) Bash(git rev-parse:*) Agent Read
---

# subagent-review

変更差分を同セッションの `general-purpose` サブエージェントへ渡してレビューさせる。`/codex-review` / `/claude-review` の別エージェント呼び出しと違い、in-sessionで軽量に走る。

## 手順

### 1. レビュー対象を特定する

`$ARGUMENTS` が指定されていればbase refとして使う。未指定なら `main` を使う。

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

スキル `/review-log read` を実行し、現在のbranchに対応するreview logから既知のレビュー判断を取得する。取得した要約は、テンプレートの `既知の論点・対応済み判断` に含めてサブエージェントに渡す。

同ディレクトリの [subagent-prompt-template.md](subagent-prompt-template.md) を埋めて、サブエージェントに渡すプロンプトを作る。

テンプレートを埋める際は、以下を具体値に置き換える。

- `<title>`: レビュー対象を短く表すタイトル
- `<absolute-repo-path>`: 現在の作業ディレクトリの絶対パス
- `<base-ref-or-sha>`: `$ARGUMENTS` で指定されたbase、または `main`
- `<head-ref-or-sha>`: 現在の `HEAD`
- `<none|staged|unstaged|untracked|mixed>`: 未コミット変更の状態
- `<summary>` / `<change>` / `<decision-or-none>` / `<context-or-none>` / `<known-item-or-none>`: レビューの前提情報
- `<reference-path-url-or-none>` / `<why-this-reference-matters>`: reviewerが明示的に読むべきファイルパス、URL、review-logの関連IDと、その参照理由。該当がなければ `なし`
- `<review-log-summary-or-none>`: `/review-log` から取得した既知のレビュー判断。なければ `なし`

テンプレートの差分欄には、取得した `git diff --stat`、hunk-levelの差分、staged / unstaged / untrackedの内容を貼る。観点と出力形式はテンプレートから削らない。

### 3. サブエージェントにレビューを依頼する

依頼文案を会話中に提示した上で、メインエージェントが新規起動または再利用の判断を済ませた状態で起動する。依頼文案のメタ情報には `work_unit: subagent-review:<base-ref-or-sha>..<head-ref-or-sha>`、`role: diff-review`、`reuse_policy: prefer-reuse` を含める。新規起動を選んでいれば新しいサブエージェントを起動し、再利用を選んでいれば既存サブエージェントを再利用する。`WAIT` / `ASK_USER` / `FIX_PROMPT` に該当する場合は起動せず、その理由と不足前提を報告して停止する。

このスキルはメインエージェントが実行し、サブエージェントには差分レビュー結果の作成だけを依頼する。メインエージェントは指摘の整理、ユーザー判断の確認、`/mark` 連携を担う。

サブエージェントへの依頼文案は、テンプレートの「境界確認」を `/subagent-check` 用の確認に使い、「レビュー対象」以降をサブエージェントへ渡す本文として使う。

#### Codex上でのモデル解決

Codex上で実行している場合は、サブエージェントの実行モデルを `codex-review` skillが管理するモデル設定ソースに合わせる。

1. `codex-review` skillが現在のセッションで利用可能な場合は、そのskillディレクトリを起点に隣接する `config.toml` と `fallback.config.toml` を読む。
2. 許可モデル一覧と指定可能な `reasoning_effort` は、実行時に利用している `spawn_agent` ツール定義を一次情報源として確認する。
3. 写像は名前の一致を優先し、`config.toml` の `model` が `spawn_agent` の許可モデル名と完全一致する場合だけ、その `model` を指定して起動する。`model_reasoning_effort` も `spawn_agent` の許可値と一致する場合だけ指定する。
4. `config.toml` の値が許可モデル外で一致しない場合は、同じskillディレクトリの `fallback.config.toml` を確認し、そちらが一致するならfallbackを使う。
5. primary / fallbackのどちらも一致しない場合は、使えない設定値と `spawn_agent` の許可モデルを示して停止する。
6. `codex-review` skillが利用できない、またはskillディレクトリにモデル設定ソースが見つからない場合は、その旨を示して停止する。

#### 再利用条件

既存サブエージェントを再利用する場合は、まず再利用候補が保持している `model` と `reasoning_effort` が、この手順で解決した目標値と一致しているかを確認する。一致している場合だけ `REUSE` として既存サブエージェントを使う。一致していない場合は、その候補を今回のレビューには再利用せず、新規起動へ切り替える。モデル差異がある候補をそのまま `REUSE` しない。

#### 起動または再利用

`OK` の場合はバックグラウンドで新規起動する（`run_in_background: true`）。Codex上でモデル指定を行う場合も、レビュー内容や依頼境界は変えず、起動パラメータだけを調整する。`REUSE` の場合は、既存サブエージェントに前回レビューから変わった差分、前回指摘への対応、今回確認してほしい範囲だけを追加依頼する。完了通知を受けてから結果を読む。

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

### 5. review logを更新する

指摘がない場合、または全指摘へのユーザー判断が揃った場合は、スキル `/review-log write` を実行し、review source、対象範囲、指摘、判断、理由、対応状態を現在の実行主体側ログへ追記する。

### 6. レビュー完了タグを設置する

指摘がない場合、または全指摘へのユーザー判断が揃った場合は、スキル `/mark review-sub` を実行する。

## 注意

- in-sessionのサブエージェントによる軽量レビュー。別エージェント（codex CLI／claude CLI）による独立視点レビューが必要なときは `/codex-review` / `/claude-review` を使う
- Codex上でモデル指定を行う場合は `codex-review` skillが公開しているモデル設定ソースを一次情報源とし、`spawn_agent` ツール定義にある許可モデル・許可 `reasoning_effort` と完全一致する値だけを使う
- 再利用候補のモデルが今回解決した目標値と一致しない場合は、その候補を今回のレビューに再利用しない
- サブエージェントは必ずバックグラウンドで起動する。同期実行にしない
- 起動前確認を省略しない
- 同じbase/head系統で修正後の再レビューを行う場合は、前回の `diff-review` サブエージェントを再利用する
