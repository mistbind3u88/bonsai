---
name: check
description: mark で未通過の品質チェック（lint / build / test / doc-check / review）を検出し、対応するスキルを実行・集約する。
allowed-tools: Bash(git diff:*)
---

# check スキル

mark タグが未設置の品質チェックを検出し、対応するスキルへ実行を委譲して結果を集約する。lint・build・test の実行は `/static-check` と `/unit-test` に、ドキュメント整合は `/doc-check` に、レビューは review スキルに委譲する。

## 手順

### 1. チェックタグを確認する

スキル `/mark --status` を実行して、現在の HEAD のチェック通過状況を確認する。

全項目が現在の HEAD にタグ設置済みなら、その旨を報告して終了する。

### 1a. autosquash 後のタグ引き継ぎ

タグが現在の HEAD にないが、タグが付いているコミットが存在する場合、そのコミットと現在の HEAD の間に差分があるか確認する。

```bash
git diff <タグのコミット> HEAD
```

差分がなければ、autosquash でコミットハッシュが変わっただけで内容は同一なので、スキル `/mark <type>` を実行してタグを現在の HEAD に付け替える。差分がある項目のみ再実行する。

### 2. 未通過の項目を実行する

タグが現在の HEAD にない項目を、対応するスキルへ委譲して実行する。

| 項目        | 担当スキル      | タグなし時の処理                                             |
| ----------- | --------------- | ------------------------------------------------------------ |
| lint, build | `/static-check` | 実行し、PASS の項目を `/mark lint`・`/mark build` で設置     |
| test        | `/unit-test`    | 実行し、PASS なら `/mark test` を設置                        |
| doc-check   | `/doc-check`    | 実行し、成功したら `/mark doc-check` を設置                  |
| review      | review スキル   | 現在の実行エージェントと別のエージェントへレビューを依頼する |

`/static-check` は lint と build を、`/unit-test` は test を、それぞれリポジトリに合わせて検出・スコープ判定・実行し、項目別に結果を報告する。`check` はその成否を受けて `/mark` で各タグを設置する。

`/static-check` は lint・build を 1 単位で実行するため、lint・build の両方が現在の HEAD にタグ済みのときだけ起動をスキップする。どちらかが未タグなら `/static-check` を実行する（タグ済みの項目も再実行されるが、HEAD で通過済みのため結果は変わらない）。`/unit-test`・`/doc-check` は対応する項目が現在の HEAD にタグ済みならスキップする。

`/doc-check` の起動・成否判断・`/mark doc-check` はメインエージェントが担う（`/doc-check` 内部の読み取り検査はサブエージェントへ委譲される）。

review では、メインエージェントが `/claude-review` または `/codex-review` を起動し、レビュー結果の扱いと、呼び出した review スキルが `/mark review` を設置したことの確認を担う。

実行順序:

1. `/static-check`・`/unit-test`・`/doc-check` の 3 項目を並列で実行する。互いの結果に依存しないため直列化しない
2. 3 項目すべてが成功したら review を実行する。review は他チェックを通過したコードに対して行うため、必ず最後に実行する

review は、現在作業しているエージェントとは別のエージェントに依頼する。

- Codex 上で実行している場合: スキル `/claude-review` を実行する
- Claude Code 上で実行している場合: スキル `/codex-review` を実行する
- 実行環境を判断できない場合: ユーザーに確認する

いずれかの項目が失敗した場合は、実行中の項目の完了を待ってから全項目の結果をまとめて報告し、停止する。

### 3. 結果サマリーを表示する

全チェック項目の結果を一覧で表示する。

```
チェック結果:
  lint:      OK
  build:     OK
  test:      OK
  doc-check: OK
  review:    OK (<実行した review skill>)
```

## 注意

- `/static-check`・`/unit-test` が項目をスキップと報告した場合（対応言語の変更なし等）も、その項目は通過扱いとし `/mark` でタグを設置する
- lint / build / test / doc-check が成功したら `/mark <type>` でタグを設置する
- review は `/claude-review` または `/codex-review` が完了時に自動でタグを設置する
- `$ARGUMENTS` で `--skip-review` が指定された場合は review をスキップする
