# q

## 前提ツール

- [git](https://git-scm.com/)

## 責務の境界

- 会話履歴、明示された前提、現在のリポジトリ状態を整理し、質問へ回答する。
- 回答に必要な範囲で read-only の確認を行う。
- ファイル編集が必要な場合は、目的に応じて `/doc-sync`、`/wiki-sync`、または対象作業の skill を案内する。
- commit / amend / fixup が必要な場合は、次に使う skill として `/commit` または `/fixup` を案内する。
- push / force push が必要な場合は、次に使う skill として `/push` を案内する。
- GitHub Issue / PR の状態確認が必要な場合は、次に使う skill として `/gh-read` を案内する。
- 前セッションや既存タスクの経緯確認が必要な場合は、次に使う skill として `/takeover` を案内する。
- GitHub 投稿が必要な場合は、目的に応じて `/gh-edit`、`/diff-comment`、`/pr-progress`、`/reply-review` を案内する。
- サブエージェント確認が必要な場合は、起動前確認を含む該当 workflow skill、または `/subagent-check` から始めるよう案内する。
- 他 skill の実行が必要な場合は、次に使う skill や確認事項を案内する。

## 回答方針

- 事実、推測、未確認を分けて答える。
- 推測は推測として明示する。
- 質問に対する結論を先に書き、必要な根拠だけを添える。
