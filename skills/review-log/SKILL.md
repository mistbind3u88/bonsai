---
name: review-log
description: review 系 skill のフィードバックループで、過去の指摘・判断・対応状態を worktree-local な review log に記録し、次回 reviewer へ渡す要約を整える。
allowed-tools: Bash(git status:*) Bash(git rev-parse:*) Bash(mkdir -p:*) Read Write
---

# review-log

review 系 skill のフィードバックループで、過去の指摘・判断・対応状態を worktree-local な一時ドキュメントとして記録する。次回 review 時は、同じ指摘の繰り返しを防ぎつつ、判断の前提が崩れた場合の再指摘を可能にする要約だけを reviewer に渡す。

## 配置

現在の実行主体に応じて、次のいずれかへ書き出す。

| 実行主体    | 書き出し先                        |
| ----------- | --------------------------------- |
| Codex       | `.codex/logs/review/<branch>.md`  |
| Claude Code | `.claude/logs/review/<branch>.md` |
| 判断不能    | `.codex/logs/review/<branch>.md`  |

参照時は `.codex/logs/review/<branch>.md` と `.claude/logs/review/<branch>.md` の両方を候補にする。書き込みは現在の実行主体側だけに行う。

実行主体は、このスキルを実行しているエージェントで判断する。Codex 上で実行している場合は `Codex`、Claude Code 上で実行している場合は `Claude Code` とする。会話・実行環境から判断できない場合は `判断不能` として `.codex/logs/review/<branch>.md` に書き込む。

`<branch>` は現在の branch 名を使う。ファイル名では、英数字、`-`、`_` 以外の文字を `__` に置換する。detached HEAD の場合は短縮コミットハッシュを使う。

## 実行モード

`$ARGUMENTS` で実行モードを指定する。

| モード  | 実行する手順                                |
| ------- | ------------------------------------------- |
| `read`  | 対象ファイルを決め、reviewer 向け要約を読む |
| `write` | 対象ファイルを決め、review log を記録する   |

review 前は `/review-log read` を実行する。review 後に指摘がない場合、または全指摘への判断が揃った場合は `/review-log write` を実行する。

## 手順

### 1. 対象ファイルを決める

以下を確認する。

```bash
git rev-parse --abbrev-ref HEAD
git rev-parse --show-toplevel
git status -s
```

`git rev-parse --abbrev-ref HEAD` が `HEAD` を返した場合は detached HEAD と判断し、短縮コミットハッシュを取得する。

```bash
git rev-parse --short HEAD
```

ログは対象リポジトリの worktree 内に作る。作業ツリーの変更有無は記録対象の文脈として扱い、ログ作成のブロッカーにはしない。

`git rev-parse --show-toplevel` の出力を worktree root とし、以降の `mkdir`、`Read`、`Write` は worktree root を起点にした絶対パスで行う。

`write` モードで書き込み先のディレクトリがなければ、現在の実行主体に対応する方だけを作成する。

```bash
mkdir -p "<worktree-root>/.codex/logs/review"
```

Claude Code 上で実行している場合は `<worktree-root>/.claude/logs/review` を作成する。

### 2. review 前に要約を読む

review skill から「参照」目的で使う場合、候補ファイルのうち存在するものだけを読み、`## 次回レビューへの引き継ぎ` セクションだけを抽出する。存在しない候補ファイルは初回 review の正常状態として扱い、エラーにしない。

内容があるログは、片側だけの場合も両方の場合もソースを示して連結する。同じ判断が重複している場合は 1 件に統合する。

```text
[codex]
<.codex/logs/review/<branch>.md の次回レビューへの引き継ぎ>

[claude]
<.claude/logs/review/<branch>.md の次回レビューへの引き継ぎ>
```

reviewer に渡す要約は次の形式にする。

```text
既知のレビュー判断:
以下は過去 review で指摘され、判断済みの事項です。
同じ根拠の指摘は繰り返さず、現在の差分で判断が破綻している場合だけ新しい根拠を示して指摘してください。

<review-handoff-summary>
```

両方の `次回レビューへの引き継ぎ` が空、または候補ログが存在しない場合は `既知のレビュー判断: なし` とする。

### 3. review 後に判断を記録する

review 結果とユーザーまたはメインエージェントの判断が揃ったら、現在の実行主体側のログへ追記する。

既存ログがある場合は、対象ファイル全体を `Read` で取得し、`次回レビューへの引き継ぎ` と `判断履歴` の既存内容を保ったまま更新後の全文を `Write` で書き戻す。`判断履歴` の追加行だけを `Write` して既存内容を置き換える操作は行わない。

既存ログの `判断履歴` が旧形式（`References` / `Related` 列を持たない形式）の場合は、追記前に新形式へ移行する。ヘッダへ `References` と `Related` を追加し、既存行の `Status` の前に `References: なし` と `Related: なし` に相当する空欄または `なし` を補う。既存行の `Decision` は当時の判断記録として保持し、決定主体が明示されていない値を推測で補完しない。移行後に新しい行を追加し、`Rationale` / `References` / `Related` の列意味が分かれる状態にする。

記録には以下を含める。

- review source（`subagent-review` / `claude-review` / `codex-review` など）
- 対象範囲（base ref、HEAD、未コミット変更の有無）
- 指摘の場所、深刻度、内容、根拠
- 判断（`採用` / `対応不要` / `保留` / `ループ懸念`）
- **決定主体**: `ユーザー裁定` / `エージェント判断` / `前回判断の引き継ぎ` のいずれか。Decision セルに括弧書きで添える（例: `対応不要（ユーザー裁定）`）。引き継ぎ時に再指摘の可否を決める強さの指標になるため必須。
- 判断理由。`対応不要` / `ループ懸念` の場合は、**検証方法と根拠を次回 reviewer が再検証できる形で残す**（コード追跡で確認した箇所、Web で参照した公式ドキュメントの URL、実行結果など）。「設計どおり」だけで終えず、確認した根拠を書く。
- 対応状態（`fixed` / `no-action` / `pending`）
- 対応コミットや確認コマンドがあればその要約
- **参照先**: reviewer が明示的に読むべきファイルパス、URL、セクション名、コマンド結果の要約。
- **関連 ID**: 同じ論点が過去の指摘で既出の場合、その `R0xx` を関連として記す。複数 review にまたがって繰り返される論点の系譜を 1 本に保つ。

### 4. reviewer 向け要約を更新する

`## 次回レビューへの引き継ぎ` は、次回 reviewer に渡す必要がある内容だけに圧縮する。

含める内容:

- 対応不要またはループ懸念と判断した指摘と理由。各項目に**決定の強さ**を明示する:
  - `ユーザー裁定` の項目は「設計意図として確定。既存判断として尊重」と書く。reviewer は新しい根拠がある別論点だけを扱う。
  - `エージェント判断` の項目は「新しい証拠・前提の崩れがある場合のみ再指摘可」と書き、判断の根拠（検証方法）を 1 行添えて reviewer が崩しにかかれるようにする。
- 修正済みだが同じ観点で再指摘されやすい判断
- reviewer が前提として読むべき、今回の差分に近い制約

**同一論点が複数 review で繰り返し指摘された場合は、引き継ぎでは 1 項目に統合する。** これまで検討済みの角度（mechanism）を列挙し、reviewer が未検討の新しい角度だけを評価できるようにする。系譜は `判断履歴` の関連 ID で追える形にする。

要約から外し、`判断履歴` に残す内容:

- 既に修正され、再指摘防止に不要な細部
- reviewer の判断に影響しない実行ログ
- 将来課題だけで今回の差分判断に使わない事項

## ログ形式

新規作成時は次の形式で作る。

```markdown
# Review Feedback Log

- Repository: <absolute-repo-path>
- Branch: <branch>
- Last updated: <YYYY-MM-DD>

## 次回レビューへの引き継ぎ

- なし

## 判断履歴

| ID  | Date | Source | Range | Location | Severity | Finding | Decision | Rationale | References | Related | Status |
| --- | ---- | ------ | ----- | -------- | -------- | ------- | -------- | --------- | ---------- | ------- | ------ |
```

`Last updated` は、実行時の会話コンテキストで示される現在日付を `YYYY-MM-DD` 形式で記録する。

新規行の `Decision` 列には決定主体を括弧書きで含める（例: `対応不要（ユーザー裁定）`、`対応不要（エージェント判断）`、`採用`）。旧形式から移行した既存行は、追記時点で確認できる当時の判断値を保持する。`Rationale` 列には判断理由を、`References` 列には検証に使ったファイルパス・URL・セクション・コマンド結果の要約を、`Related` 列には関連 `R0xx` を含める。

追記時は `判断履歴` に行を追加し、必要に応じて `次回レビューへの引き継ぎ` と `Last updated` を更新する。

`ID` は `R001` から始まる連番を使う。既存ログに `R001` から `R007` までがある場合、次の行は `R008` とする。

## 注意

- ログは worktree-local な一時資料として扱う。
- `.codex/logs/review/` と `.claude/logs/review/` は Git 管理対象外の一時資料として扱う。
- 恒久的に残す価値がある判断は、対象リポジトリのルールに従って taskdoc、AGENTS、README、Wiki などへ移す。
- reviewer にはログ全体ではなく `次回レビューへの引き継ぎ` の要約を渡す。
