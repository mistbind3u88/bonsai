# rule-check

## 前提ツール

- [git](https://git-scm.com/)

## 責務の境界

- `rule-check` は、変更内容が適用範囲の `AGENTS.md` / `CLAUDE.md` に記載されたルールへ適合しているかを検査して結果を報告する。
- 読み取り検査はサブエージェントへ限定タスクとして委譲し、起動前確認は `/subagent-check` で行う。
- 同じ base/head の再検査では、`work_unit: rule-check:<base>..<head>` / `role: rule-check` の既存サブエージェントを再利用する。
- ドキュメント追従や README / SKILL.md の整合性確認は `/doc-check` が担う。
- 検出した不適合の修正は、変更内容に応じて対象作業の skill、`/doc-sync`、または `/wiki-sync` に委ねる。
- メインエージェントは `/rule-check` の手順解釈、結果回収、最終判断、ユーザー報告を担う。
- `/check` から呼び出された場合、成功後の `/mark rule-check` は `/check` 側のメインエージェントが担う。

## `doc-check` との違い

- `doc-check` は、コード変更やスキル変更に対して関連ドキュメントが追従しているかを確認する。
- `rule-check` は、変更そのものが適用範囲のルールへ適合しているかを確認する。
