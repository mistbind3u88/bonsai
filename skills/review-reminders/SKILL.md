---
name: review-reminders
description: 自分が author の Open PR を横断確認し、レビューリマインド候補・依存関係・自分対応が先の項目を整理する。tanaoroshi JSON があれば依存関係の補助情報として参照する。
allowed-tools: Bash(gh api user:*) Bash(gh repo view:*) Bash(gh pr list:*) Bash(gh pr view:*) Read Grep
---

# review-reminders

自分が author の Open PR を確認し、レビュー依頼・リマインド対象を依存関係つきで整理する。

## 入力

引数で対象リポジトリをスペース区切りで指定する。`tanaoroshi collect` の JSON を補助情報として使う場合は、`--tanaoroshi-json <path>` を指定する。

```text
/review-reminders owner/repo1 owner/repo2
/review-reminders --tanaoroshi-json .results/tanaoroshi.json owner/repo1 owner/repo2
```

引数なしの場合は、`tanaoroshi` の Phase 0 と同じ手順で、カレントリポジトリに加えて CLAUDE.md または AGENTS.md の「関連リポジトリ」セクションに記載されたリポジトリを対象にする。関連リポジトリセクションがなければカレントリポジトリのみを対象にする。

`--tanaoroshi-json` が指定された場合は、PR 一覧と Issue/PR 参照関係の補助情報として使う。GitHub 上の最新状態は `gh` で確認する。

## 手順

### 1. 対象リポジトリを解決する

リポジトリ指定がない場合は、`tanaoroshi` の Phase 0 と同じ手順で、カレントリポジトリと関連リポジトリを対象にする。関連リポジトリセクションの確認が必要な場合は、AGENTS.md または CLAUDE.md を `Read` / `Grep` で確認する。

```bash
gh repo view --json nameWithOwner
```

`--tanaoroshi-json` が指定されている場合は、`Read` で対象ファイルを読み込み、PR 一覧、body 内参照、`refs` 相当の Issue/PR 参照関係を依存関係の補助情報として使う。

### 2. 自分の login を確認する

```bash
gh api user --jq .login
```

### 3. 自分が author の Open PR を取得する

各リポジトリで、自分が author の Open PR を取得する。

```bash
gh pr list \
  --repo <owner/repo> \
  --author <login> \
  --state open \
  --json number,url,title,author,createdAt,updatedAt,isDraft,reviewDecision,reviewRequests,latestReviews,headRefName,baseRefName,labels,statusCheckRollup,mergeable
```

PR の詳細確認が必要な場合は `gh pr view` を使う。

```bash
gh pr view <number> \
  --repo <owner/repo> \
  --json number,url,title,body,author,createdAt,updatedAt,isDraft,reviewDecision,reviewRequests,latestReviews,headRefName,baseRefName,labels,statusCheckRollup,mergeable,closingIssuesReferences,comments
```

### 4. 依存関係を確認する

候補 PR ごとに、以下を依存関係として整理する。

| 判定元                    | 見る内容                                                                                 |
| ------------------------- | ---------------------------------------------------------------------------------------- |
| `baseRefName`             | default branch 以外を向いている PR は、base branch を head に持つ前提 PR 候補を探す      |
| PR body                   | `depends on`、`blocked by`、`前提`、`ブロッカー`、`後続`、`関連`、GitHub Issue/PR リンク |
| `closingIssuesReferences` | PR が解決対象にしている Issue                                                            |
| head branch               | 同じ prefix や連番を持つ stacked PR 候補                                                 |
| comments / latestReviews  | reviewer が示した「先に確認するもの」「ブロッカー」「再レビュー待ち」                    |
| tanaoroshi JSON           | すでに抽出済みの Issue/PR 参照関係                                                       |

依存関係は `前提` / `ブロッカー` / `解決対象` / `関連` に分類する。判断根拠が弱いものは `関連` に置く。

### 5. リマインド可否を分類する

PR ごとに次のカテゴリへ分類する。複数カテゴリに該当する場合は、表の上から評価して最初に一致したカテゴリを採用する。

| カテゴリ           | 条件                                                                                                                     |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------ |
| 自分対応が先       | Draft、`CHANGES_REQUESTED`、CI failure、conflict、未返信コメントなど author 側の次アクションがある                       |
| 依存待ち           | 前提 PR、方針確認 Issue、他者対応など、先に進める対象がある                                                              |
| 完了対応           | `APPROVED` で未マージ、または merge 可能性の確認が必要                                                                   |
| reviewer 未割当    | ready な PR だが `reviewRequests` が空、または明確な reviewer が見えない                                                 |
| すぐリマインドする | `isDraft=false`、自分対応の明確なブロッカーがなく、review request または reviewer activity から 3 営業日以上経過している |

待ち日数は営業日換算で表示する。次の優先順位で採用した日時から数え、採用した起点を出力に書く。

| 優先 | 起点                                                           | 使う場面                                                           |
| ---: | -------------------------------------------------------------- | ------------------------------------------------------------------ |
|    1 | 最後の author 以外の review / comment に author が返信した日時 | 再レビュー待ち                                                     |
|    2 | 最新の review request が作られた日時                           | 初回レビュー待ち                                                   |
|    3 | author 以外の最後の comment / review の日時                    | 相手からの確認待ちを受けて author がまだ返していないか判定する補助 |
|    4 | PR の `updatedAt`                                              | 上記が取れない場合の fallback                                      |

### 6. 出力する

出力は、先にアクション別サマリーを示し、その後に依存関係を含む詳細を置く。

```markdown
## レビューリマインド候補

### サマリー

| カテゴリ           | 件数 | 主な次アクション           |
| ------------------ | :--: | -------------------------- |
| すぐリマインドする |  N   | reviewer へレビュー依頼    |
| 依存待ち           |  N   | 前提 PR / Issue を先に確認 |
| reviewer 未割当    |  N   | reviewer を割り当てる      |
| 自分対応が先       |  N   | author 側の対応を完了する  |
| 完了対応           |  N   | merge 可否を確認する       |

### すぐリマインドする

- [PR/OPEN] タイトル [owner/repo#N](https://github.com/owner/repo/pull/N)
  - 待ち: @reviewer / 3営業日（起点: review request）
  - 状態: REVIEW_REQUIRED / CI OK / conflictなし
  - 依存: 前提なし
  - 次: reviewer にレビュー依頼

### 依存待ち

- [PR/OPEN] タイトル [owner/repo#N](https://github.com/owner/repo/pull/N)
  - 待ち: @reviewer / 5営業日（起点: author 返信）
  - 状態: REVIEW_REQUIRED
  - 依存:
    - 前提: [PR/OPEN] タイトル [owner/repo#M](https://github.com/owner/repo/pull/M)
    - 前提: [PR/CLOSED] タイトル [owner/repo#L](https://github.com/owner/repo/pull/L)（closed: YYYY-MM-DD）
    - ブロッカー: [Issue/OPEN] タイトル [owner/repo#K](https://github.com/owner/repo/issues/K)
  - 次: 前提 PR のレビューを促す

### reviewer 未割当

- [PR/OPEN] タイトル [owner/repo#N](https://github.com/owner/repo/pull/N)
  - 経過: ready 化または作成から 2営業日
  - 状態: REVIEW_REQUIRED / reviewer 未割当
  - 依存: 前提なし
  - 次: reviewer を指定して依頼

### 自分対応が先

- [PR/OPEN] タイトル [owner/repo#N](https://github.com/owner/repo/pull/N)
  - 理由: CHANGES_REQUESTED / CI failure / conflict
  - 依存: 関連 [Issue/OPEN] タイトル [owner/repo#K](https://github.com/owner/repo/issues/K)
  - 次: author 側の対応を完了してから再依頼

### 完了対応

- [PR/OPEN] タイトル [owner/repo#N](https://github.com/owner/repo/pull/N)
  - 状態: APPROVED / CI OK
  - 依存: 前提なし
  - 次: merge 可否を確認
```

## 出力ルール

- PR / Issue のリンク、種別プレフィックス、Closed の後置メタ情報は `tanaoroshi` の出力の共通ルールに準拠する。Open PR の `author` / `assignee` / `next` 相当の情報は、リマインド判断に必要な `待ち` / `状態` / `依存` / `次` へ分解して書く。
- PR / Issue は `[owner/repo#N](https://github.com/owner/repo/pull/N)` または `[owner/repo#N](https://github.com/owner/repo/issues/N)` のクリック可能なリンクで書く。
- PR 行は `[PR/OPEN]` / `[PR/Draft]` / `[PR/CLOSED]` を付ける。
- Issue 行は `[Issue/OPEN]` / `[Issue/CLOSED]` を付ける。
- Closed の前提・解決対象は `（closed: YYYY-MM-DD）` を付ける。
- リマインド対象は、誰に何を依頼するかを `次` に明記する。
- リマインドより先に進めるべき依存関係がある場合は、PR 本体ではなく依存先を次アクションにする。
- GitHub への投稿や PR 本文更新が必要な場合は、目的に応じて `/gh-edit`、`/pr-progress`、`/reply-review`、`/diff-comment` を次に使う skill として案内する。
