---
name: pr-ready
description: PR のレビュー準備が完了していることを確認し、draft PR を Ready for Review に切り替える。
allowed-tools: Bash(git status:*) Bash(git rev-parse:*) Bash(gh pr view:*) Bash(gh pr diff:*) Bash(gh pr ready:*)
---

# pr-ready スキル

PR がレビュアーへ渡せる状態かを最終確認し、draft の場合は Ready for Review に切り替える。

## 役割

`pr-ready` はレビュー準備の最終ゲートである。前提が揃っていない場合は修正や他スキルの実行へ進まず、不足項目と次に使うべきスキルを報告して停止する。

Ready 化の対象は PR の状態で分ける。

| PR 状態         | 扱い                                                             |
| --------------- | ---------------------------------------------------------------- |
| draft PR        | すべての前提と最終確認が揃った場合に `gh pr ready` を実行する    |
| ready PR        | `gh pr ready` は実行せず、レビュー準備の不足確認と完了報告を行う |
| closed / merged | レビュー準備の対象外として停止する                               |

## 手順

### 1. 対象 PR を特定する

`$ARGUMENTS` から PR 番号または PR URL を受け取る。未指定の場合は現在のブランチに紐づく PR を対象にする。

```bash
gh pr view <PR番号またはURL> --json number,url,title,state,isDraft,headRefOid,headRefName,baseRefName,body,reviewDecision,statusCheckRollup,changedFiles,files
```

未指定の場合:

```bash
gh pr view --json number,url,title,state,isDraft,headRefOid,headRefName,baseRefName,body,reviewDecision,statusCheckRollup,changedFiles,files
```

`state` が `OPEN` 以外の場合は、対象外として停止する。

### 2. ローカル状態の前提を確認する

作業ツリーと HEAD が PR head と一致していることを確認する。

```bash
git status --short
git rev-parse HEAD
```

前提:

- `git status --short` が空である
- `git rev-parse HEAD` が PR の `headRefOid` と一致している

前提が揃わない場合:

| 不足項目                        | 報告する次アクション                                     |
| ------------------------------- | -------------------------------------------------------- |
| 未コミット変更がある            | `/commit` で変更を整理する                               |
| ローカル HEAD と PR head が違う | `/push` で現在 HEAD を反映する、または対象 PR を確認する |
| 対象 PR を特定できない          | `/gh-read` で対象 PR を確認する                          |

このステップで不足があれば停止する。

### 3. 品質チェックの前提を確認する

`/mark --status` で現在 HEAD のチェック通過状態を確認する。

```bash
/mark --status
```

前提:

- `lint` が現在の HEAD
- `build` が現在の HEAD
- `test` が現在の HEAD
- `doc-check` が現在の HEAD
- `review-cross` が現在の HEAD

上記は `/check` が要求する現在 HEAD の品質項目に、Ready for Review 前提として `review-cross` を加えたものとする。`review-sub` は `review-cross` の代替にしない。Ready for Review に出す直前は別エージェント視点の cross-review まで通った状態を要求する。将来 `/check` の required 項目が増えた場合は、`pr-ready` でも同じ項目を現在 HEAD で要求する。

前提が揃わない場合は、`/check` を実行するよう報告して停止する。

### 4. CI 状態を確認する

PR の `statusCheckRollup` を確認する。

前提:

- 必須の CI / status check が成功している
- pending / queued / in_progress の check が残っていない
- failed / cancelled / timed_out の check がない
- CI 状態を取得できている

前提が揃わない場合:

| 状態                  | 報告する次アクション             |
| --------------------- | -------------------------------- |
| pending / in_progress | `/watch-ci` で完了を確認する     |
| failed / cancelled    | `/watch-ci` で失敗内容を確認する |
| CI 状態を確認できない | `/watch-ci` で状態を確認する     |
| CI が設定されていない | CI がないことを明示して次へ進む  |

CI が pending / in_progress / failed / cancelled / timed_out、または取得不可の場合は Ready 化へ進まず停止する。

### 5. PR 概要欄の最終確認を行う

PR の `title` と `body` を読み、レビュー準備として必要な情報が揃っているか確認する。

確認すること:

- PR の目的・背景・レビュアーが判断すべき観点が本文にある
- `## やったこと` などの本文が実装の逐語説明ではなく、変更の意味と影響を説明している
- `## やってないこと` は、この PR と隣接する明確な除外判断だけを扱っている
- References / Closes / ref は PR の主題と本文中で説明できる関係にある
- 確認したこと、CI、レビュー時に見てほしい点が必要な粒度で書かれている

不足があれば、`/gh-edit` で概要欄を更新するよう報告して停止する。

### 6. diff comment の必要性を確認する

PR diff を確認し、特定行に紐づけて補足すべき意図が残っているか確認する。

```bash
gh pr diff <PR番号> --patch
```

補足先は内容で分ける。

| 補足内容                               | 扱い                                                          |
| -------------------------------------- | ------------------------------------------------------------- |
| PR 全体の目的・背景・スコープ          | `/gh-edit` で概要欄へ反映する                                 |
| 特定行の判断理由・互換性・レビュー観点 | `/diff-comment` で対象行へ投稿する                            |
| 将来の保守にも必要な説明               | コードコメント、README、設計ドキュメント、LLM Wiki へ反映する |

`/diff-comment` が必要な場合は、対象候補を `path` / `line` / `補足したい意図` の形式で列挙し、`/diff-comment` を実行するよう報告して停止する。

### 7. Ready for Review にする

ここまでの前提と最終確認がすべて揃っている場合だけ実行する。

draft PR の場合:

```bash
gh pr ready <PR番号>
```

ready PR の場合は `gh pr ready` を実行せず、レビュー準備が完了していることを報告する。

### 8. 結果を報告する

以下を報告する。

- PR URL
- draft から ready に切り替えたか、既に ready だったか
- 現在 HEAD と PR head の SHA
- `/mark --status` の required 項目が現在 HEAD で揃っていること
- CI の状態
- 概要欄と diff comment の最終確認結果

## 停止時の出力

不足がある場合は、次の形式で報告する。

```markdown
pr-ready は停止しました。

- 不足項目: <不足している前提またはレビュー準備>
- 理由: <Ready for Review に進めない理由>
- 次に使う skill: /<skill-name>
```

複数の不足がある場合は、最初に解決すべき順に並べる。
