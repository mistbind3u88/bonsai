---
name: rule-check
description: 変更内容が適用範囲の AGENTS.md / CLAUDE.md に記載されたルールへ適合しているかを検査する。
allowed-tools: Bash(git diff:*) Agent Read
---

# rule-check

変更内容が、適用範囲の `AGENTS.md` / `CLAUDE.md` に記載されたルールへ適合しているかをチェックし、結果を報告する。

**このスキルはメインエージェントが実行し、読み取り検査だけを `general-purpose` サブエージェントへ依頼する。** メインエージェントは `/rule-check` の手順解釈、`/subagent-check`、起動、結果受領、最終判断、報告を担う。`/check` から呼び出された場合の `mark` 設置は、呼び出し元のメインエージェントが担う。

## 手順

### 1. 適用ルールを特定する

変更ファイルごとに、適用範囲の `AGENTS.md` / `CLAUDE.md` を特定する。

- リポジトリルートの `AGENTS.md`
- 変更ファイルの親ディレクトリから上位へたどって見つかる `AGENTS.md`
- `CLAUDE.md` が存在する場合は、同じ範囲のエージェント向けルールとして扱う
- `CLAUDE.md -> AGENTS.md` の symlink は同一ルールとして扱う

### 2. サブエージェントを起動する

サブエージェントへの依頼文案を会話中に提示した上で、起動前にスキル `/subagent-check` を実行し、既存サブエージェントの完了状態、再利用可否、起動上限の余剰、依頼文の委譲境界を確認する。依頼文案のメタ情報には `work_unit: rule-check:<base-ref-or-sha>..<head-ref-or-sha>`、`role: rule-check`、`reuse_policy: prefer-reuse` を含める。`subagent-check` が `OK` を返した場合は、サブエージェントを起動する理由、依頼範囲、読み取り検査だけを任せる境界をユーザーに提示する。`/check` などの呼び出し元で同じ起動範囲について事前承認済みの場合は、その承認に基づいて新規起動する。事前承認がない場合は、ユーザーの承認を得てから新規起動する。`REUSE` / `WAIT` / `ASK_USER` / `FIX_PROMPT` の場合はその判定に従う。

`OK` の場合は `general-purpose` サブエージェントを `run_in_background: true` で起動する。`REUSE` の場合は、既存サブエージェントに前回検査から変わった差分、前回指摘への対応、今回確認してほしい適用ルールを追加依頼する。同期実行にしない。コンソールを占有せず、完了通知を受け取ってから結果を読む。

#### サブエージェントへの依頼文案

依頼文案は、同ディレクトリの [subagent-prompt-template.md](subagent-prompt-template.md) を埋めて作成する。テンプレートの「境界確認」は `/subagent-check` 用の確認にも使い、「検査対象」以降はサブエージェントへ渡す本文として使う。

テンプレートを埋める際は、以下を具体値に置き換える。

- `<absolute-repo-path>`: 現在の作業ディレクトリの絶対パス
- `<base-ref-or-sha>`: 比較元。原則 `main`。既定ブランチが異なる場合はそのブランチ、特定コミットを検査する場合はその親または明示された base
- `<head-ref-or-sha>`: 比較先。通常は `HEAD`
- `<changed-files>`: 変更ファイル一覧
- `<applicable-rule-files>`: 手順1で特定した `AGENTS.md` / `CLAUDE.md`
- `<reference-path-url-or-none>` / `<why-this-reference-matters>`: ルール適用や判断根拠として明示的に読むべき文書、URL、review-log の関連 ID と、その参照理由。該当がなければ `なし`
- `<user-provided-scope-or-none>`: ユーザーが明示した確認範囲や変更概要。なければ `なし`

サブエージェントには、テンプレートの観点カテゴリを使って適用ルールから実際の検査観点を抽出させる。呼び出し元リポジトリの `AGENTS.md` / `CLAUDE.md` に存在しないルールを、固定の必須条件として扱わない。ドキュメントの内容が変更に追従しているかは `/doc-check` が検査し、`rule-check` はルールに定義されたドキュメント品質・コード品質・手順・配置基準への適合を検査する。

### 3. 結果を報告する

サブエージェントの完了通知を受けたら出力を読み、ユーザーに報告する。

- `要修正` がない場合: 成功として報告する（`対象外` は成功扱い）
- `要修正` がある場合: 失敗として一覧を報告する

```text
ルール適合チェック:
  skills/example/SKILL.md: allowed-tools が最小権限ルールと合っていない
  README.md: README に詳細手順を書きすぎており、情報配置ルールと合っていない
```

## 検査観点

検査観点は呼び出し元リポジトリの `AGENTS.md` / `CLAUDE.md` に依存する。[subagent-prompt-template.md](subagent-prompt-template.md) の観点カテゴリは、適用ルールを読み解くための分類軸として使う。対象リポジトリに存在しないルールは `対象外` と判定し、存在するルールに対して変更内容が適合しているかを確認する。ドキュメント内容の変更追従は `/doc-check` が検査し、`rule-check` は適用ルールに定義されたドキュメント品質、コード品質、手順、配置基準への適合を検査する。

## 注意

- このスキルはチェックと報告のみを行う。検出した不適合は、変更内容に応じて対象作業の skill、`/doc-sync`、または `/wiki-sync` で修正する。
- `/check` から呼び出された場合、ルール不適合があれば check 全体を失敗にする。
- `/check` から呼び出された場合、成功後の `/mark rule-check` は `/check` 側が実行する。
- サブエージェントは必ずバックグラウンドで起動する。同期実行にしない。
- 要修正対応後の再検査は、呼び出し側が `/rule-check` を再度実行することで行う（このスキル内でループしない）。同じ `work_unit` の再検査では、前回の `rule-check` サブエージェントを再利用する。
