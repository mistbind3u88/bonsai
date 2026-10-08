---
name: push
description: /check の通過を確認してから git push する。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git rev-parse:*) Bash(git push:*) Bash(gh pr view:*)
---

# push スキル

`/check` が通過していることを確認してから push する。

## 手順

### 1. 前提確認

```bash
git status -s
git log --oneline main..HEAD
git rev-parse --abbrev-ref HEAD
```

- 未コミットの変更がある場合は push せず、先にコミットするようユーザーに伝える
- main / masterにいる場合は、リポジトリの適用ルールとユーザー指示から、そのブランチへの直接pushが許可され、依頼範囲に含まれるかを確認する。両方を確認できれば、ブランチ名だけを理由とする追加確認は省略してよい。許可または依頼範囲が不明な場合は、対象ブランチを示してユーザーに確認する

### 2. チェックを実行する

スキル `/check` を実行する。

`$ARGUMENTS` に `--review=skip` がある場合はスキル `/check --review=skip` を実行する。

全チェックが OK でない場合は push せずに停止する。

### 3. push する

```bash
git push -u origin HEAD
```

push 後、結果を報告する。

### 4. コンフリクト解消時の PR コメント

rebase で main を取り込んでコンフリクトを解消した場合（`--force-with-lease` で push した場合）、スキル `/pr-progress` を実行してコメントを投稿する。push 前に旧 HEAD を記録しておくこと。

## 注意

- `$ARGUMENTS` で `--force` が指定された場合は `git push --force-with-lease` を使う（ユーザーに確認後）
