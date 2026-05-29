---
name: diff-comment
description: PR の差分行に、レビュー時だけ必要な意図・背景・判断材料の inline comment を投稿する。
allowed-tools: Bash(gh repo view:*) Bash(gh pr view:*) Bash(gh pr diff:*) Bash(gh api repos/*/pulls/*/comments:*)
---

# diff-comment スキル

PR の diff 上に、レビュー時だけ必要な意図・背景・判断材料を inline comment として投稿する。

## 対象場面

| 場面                                                   | 使い方                                                               |
| ------------------------------------------------------ | -------------------------------------------------------------------- |
| 実装意図をレビュアーへ補足したい                       | 対象行に、なぜその判断をしたかを短く書く                             |
| レビュー時だけ必要な補足を伝えたい                     | PR 上の一時的なレビュー補助として残す                                |
| 差分の読み方やレビューしてほしい観点を対象行に添えたい | 対象行に、確認してほしい観点・判断材料・前提を結び付けて書く         |
| PR 全体の経過や push 後の差分を伝えたい                | 次に使うスキルとして `/pr-progress` を案内する                       |
| 既存レビューコメントへ返信したい                       | 次に使うスキルとして `/reply-review` を案内する                      |
| PR/Issue の本文を更新したい                            | 次に使うスキルとして `/gh-edit` を案内する                           |
| 将来の保守にも必要な説明を残したい                     | ソースコードのコメント、README、設計ドキュメント、LLM Wikiへ反映する |

## 手順

### 1. 対象 PR を特定する

`$ARGUMENTS` から PR 番号、PR URL、対象ファイル、対象行、コメント本文を読み取る。PR の指定がなければ、現在のブランチに紐づく PR を対象にする。

```bash
gh repo view --json owner,name
gh pr view <PR番号またはURL> --json number,url,title,headRefOid,headRefName,baseRefName,isDraft
```

PR 指定がない場合:

```bash
gh pr view --json number,url,title,headRefOid,headRefName,baseRefName,isDraft
```

### 2. コメント対象行を確認する

対象行が PR diff 上でコメント可能な行か確認する。

```bash
gh pr diff <PR番号> --patch
```

確認すること:

- `path` は diff に含まれるファイルである
- `line` は GitHub の PR diff 上で見える行である
- 追加・変更・コンテキスト行へのコメントは `side=RIGHT` を使う
- 削除行へのコメントは `side=LEFT` を使う
- 行番号や side を一意に決められない場合は、投稿候補を提示してユーザーに確認する

### 3. コメント本文を作成する

コメントは対象行に結び付く意図・背景・レビュー観点に絞る。

本文の観点:

- この行の判断理由
- レビュアーに確認してほしい観点
- PR 本文に書くほど広くはないが、この差分の理解に必要な前提
- コードへ恒久的に残す説明か、PR 上の補足で足りる説明かの判断

書き方:

- 1コメント1意図にする
- diff から読み取れる実装説明より、判断や意図を優先して書く
- 対象行から離れた話題は、次に使うスキルとして `/pr-progress` または `/gh-edit` を案内する
- 恒久的に必要な説明はソースコードやドキュメントへ反映する

例:

```markdown
ここは既存 API のレスポンス形状を維持するため、内部では新しい型へ寄せつつ outward-facing なキー名は据え置いています。
```

### 4. 投稿前にユーザー確認する

投稿予定を以下の形式で提示し、ユーザーの承認を得る。

```markdown
投稿予定:

- PR: <URL>
- path: <path>
- line: <line>
- side: RIGHT|LEFT
- body: <本文>
```

ユーザーが既に投稿を明示承認している場合は、この確認を投稿直前の最終提示として扱い、そのまま手順 5 へ進む。

### 5. inline comment を投稿する

Markdown 本文の改行・引用符・バッククォートを安全に扱うため、payload JSON を標準入力から `--input -` で渡す。

単一行コメント:

```bash
gh api repos/<OWNER>/<REPO>/pulls/<PR番号>/comments --method POST --input - <<'JSON'
{
  "body": "<本文>",
  "commit_id": "<headRefOid>",
  "path": "<path>",
  "line": <line>,
  "side": "RIGHT"
}
JSON
```

削除行にコメントする場合は `side="LEFT"` を指定する。

複数行コメントが必要な場合:

```bash
gh api repos/<OWNER>/<REPO>/pulls/<PR番号>/comments --method POST --input - <<'JSON'
{
  "body": "<本文>",
  "commit_id": "<headRefOid>",
  "path": "<path>",
  "start_line": <start_line>,
  "start_side": "RIGHT",
  "line": <end_line>,
  "side": "RIGHT"
}
JSON
```

投稿後、返却された `html_url` を報告する。

## 注意

- コメント対象は PR diff に表示される行に限定する
- 投稿本文は公開 PR に残せる内容で書き、秘匿情報、社内 URL、ローカルパスを除外する
- コメントが複数ある場合は、対象行・本文・投稿順を一覧化してから承認を得る
- draft PR でも、レビュアーに事前共有する意図がある場合は投稿できる。投稿前に draft であることを明示する
