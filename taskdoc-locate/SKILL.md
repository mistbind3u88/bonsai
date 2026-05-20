---
name: taskdoc-locate
description: タスクドキュメントを置く元リポジトリ側の .claude/docs ディレクトリの場所を解決する。worktree 側ではなく元リポジトリ側を返す。
allowed-tools: Bash(git rev-parse:*) Bash(git worktree list:*) Bash(ls:*)
---

# taskdoc-locate スキル

タスクドキュメントが置かれる（あるいは置くべき）元リポジトリ側の `.claude/docs` ディレクトリの絶対パスを解決して返す。worktree から作業している場合も、worktree 側ではなく元リポジトリ側を返す。

タスクドキュメントの読み取り・作成・整理は呼び出し元が行う。このスキルは場所解決だけを担う。

## 手順

### 1. 現在のリポジトリと元リポジトリを特定する

```bash
git rev-parse --show-toplevel
git worktree list --porcelain
```

`git worktree list --porcelain` の先頭に出る `worktree` を元リポジトリとして扱う。現在の作業ディレクトリが linked worktree の場合も、元リポジトリ側を対象にする。

元リポジトリを一意に判断できない場合は、候補パスを提示して停止する。

### 2. `.claude/docs` のパスを返す

元リポジトリ直下の `.claude/docs` を解決対象とし、`ls` で存在を確認する。

```bash
ls -d <元リポジトリ>/.claude/docs
```

- ディレクトリが存在する場合: その絶対パスを報告する
- 存在しない場合: 解決した絶対パス（未作成）と、まだ存在しない旨を報告する

ローカル環境に依存する固定パス（例: `~/Workspace/<repo>`）を前提にしない。

## 注意

- このスキルは場所解決のみを行う。タスクドキュメントの検索・読み取りは呼び出し元が解決したパスに対して直接行う
- タスクドキュメントの記載形式・作成方針はリポジトリの `AGENTS.md`（タスクドキュメント配置の判断基準）に従う
- 削除や恒久ドキュメントへの移管はスキル `/clean-docs` で行う
