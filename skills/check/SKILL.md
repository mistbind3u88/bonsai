---
name: check
description: mark で未通過の品質チェック（lint / build / test / doc-check / rule-check / review-sub / review-cross）を検出し、対応するスキルを実行・集約する。
allowed-tools: Bash(git diff:*)
---

# check スキル

mark タグが未設置の品質チェックを検出し、対応するスキルへ実行を委譲して結果を集約する。lint・build は `/static-check` に、test は `/unit-test` に、ドキュメント整合は `/doc-check` に、ルール適合は `/rule-check` に、review は引数で選ぶ 2 種類（sub・cross）に委譲する。

## 引数

`$ARGUMENTS` で走らせる review を選ぶ:

| 引数                 | 走らせる review                                         | 主な用途                                                       |
| -------------------- | ------------------------------------------------------- | -------------------------------------------------------------- |
| （なし・デフォルト） | cross-review（`/codex-review` または `/claude-review`） | `/push` 経由・手動 `/check`                                    |
| `--review=sub`       | sub-review（`/subagent-review`）                        | `/commit` 各コミット末尾・autosquash 完了後                    |
| `--review=skip`      | なし                                                    | `/commit` autosquash 中間 step・手動でレビュー抜きで見たいとき |

**`review-cross` は `review-sub` の上位**として扱う。`review-cross` が HEAD タグ済みのときは、`--review=sub` でも cross の通過をもって sub も通過扱いとする（cross が sub の役割を含めて代替）。代替の方向は cross → sub の一方向に限る。

## 手順

### 1. チェックタグを確認する

スキル `/mark --status` を実行して、現在の HEAD のチェック通過状況を確認する。

引数に応じて対象項目を判定し、すべて現在の HEAD にタグ設置済みなら、その旨を報告して終了する:

- 共通: `lint` / `build` / `test` / `doc-check` / `rule-check` が HEAD タグ済み
- `--review=sub`: `review-sub` または `review-cross` が HEAD タグ済み（cross は sub の上位として代替）
- 引数なし（=cross）: `review-cross` が HEAD タグ済み
- `--review=skip`: lint / build / test / doc-check / rule-check のみ判定する

### 1a. autosquash 後のタグ引き継ぎ

タグが現在の HEAD にないが、タグが付いているコミットが存在する場合、そのコミットと現在の HEAD の間に差分があるか確認する。

```bash
git diff <タグのコミット> HEAD
```

差分がなければ、autosquash でコミットハッシュが変わっただけで内容は同一なので、スキル `/mark <type>` を実行してタグを現在の HEAD に付け替える。差分がある項目のみ再実行する。

### 2. 未通過の項目を実行する

タグが現在の HEAD にない項目を、対応するスキルへ委譲して実行する。

| 項目         | 担当スキル                              | タグなし時の処理                                               |
| ------------ | --------------------------------------- | -------------------------------------------------------------- |
| lint, build  | `/static-check`                         | 実行し、PASS の項目を `/mark lint`・`/mark build` で設置       |
| test         | `/unit-test`                            | 実行し、PASS なら `/mark test` を設置                          |
| doc-check    | `/doc-check`                            | 実行し、成功したら `/mark doc-check` を設置                    |
| rule-check   | `/rule-check`                           | 実行し、成功したら `/mark rule-check` を設置                   |
| review-sub   | `/subagent-review`                      | 引数 `--review=sub` 指定時のみ実行。完了で `review-sub` 設置   |
| review-cross | `/codex-review` または `/claude-review` | 引数なし（デフォルト）の場合に実行。完了で `review-cross` 設置 |

`/static-check` は lint と build を、`/unit-test` は test を、それぞれリポジトリに合わせて検出・スコープ判定・実行し、項目別に結果を報告する。`check` はその成否を受けて `/mark` で各タグを設置する。

`/static-check` は lint・build を 1 単位で実行するため、lint・build の両方が現在の HEAD にタグ済みのときだけ起動をスキップする。どちらかが未タグなら `/static-check` を実行する。`/unit-test`・`/doc-check`・`/rule-check` は対応する項目が現在の HEAD にタグ済みならスキップする。

`/doc-check` の起動・成否判断・`/mark doc-check` はメインエージェントが担う（`/doc-check` 内部の読み取り検査はサブエージェントへ委譲される）。

`/rule-check` の起動・成否判断・`/mark rule-check` はメインエージェントが担う（`/rule-check` 内部の読み取り検査はサブエージェントへ委譲される）。

cross-review は現在のエージェントとは別のエージェントに依頼する。

- Codex 上で実行している場合: スキル `/claude-review` を実行する
- Claude Code 上で実行している場合: スキル `/codex-review` を実行する
- 実行環境を判断できない場合: ユーザーに確認する

実行順序:

1. `/static-check`・`/unit-test` を並列で実行する。互いの結果に依存せず、サブエージェントも使わないため直列化しない
2. `/doc-check` と `/rule-check` を順に実行する。どちらもサブエージェントを起動するため、各スキルの `/subagent-check` で完了状態・起動枠・再利用可否を確認してから次へ進む
3. 4 項目すべてが成功したら、引数で選ばれた review を実行する。引数なしなら cross-review、`--review=sub` なら sub-review。`--review=skip` のときは lint / build / test / doc-check / rule-check のみで終える。review は他チェックを通過したコードに対して行うため、必ず最後に実行する。`--review=sub` で `review-cross` が既に HEAD タグ済みのときは、cross の通過をもって sub-review も通過扱いとする（cross が上位として代替）

いずれかの項目が失敗した場合は、実行中の項目の完了を待ってから全項目の結果をまとめて報告し、停止する。

### 2a. cross-review 失敗時の fallback

引数なし（デフォルト = cross-review）で `/codex-review` または `/claude-review` が CLI 実行不可等で失敗した場合、`/subagent-review` へ fallback して sub-review を走らせる。

- 成功時: `/subagent-review` 完了に伴い `review-sub` が HEAD タグに設置される（`review-cross` は未通過のまま残す）
- 結果報告で「cross-review が実行できず sub-review に fallback した」旨と、`review-cross` 未通過を明示する
- fallback の有無は `/check` 内で完結する。呼び出し元（`/push` 等）は通過状況の tag を見て判断する

### 3. 結果サマリーを表示する

実行した全チェック項目の結果を一覧で表示する。

```
チェック結果:
  lint:          OK
  build:         OK
  test:          OK
  doc-check:     OK
  rule-check:    OK
  review-cross:  OK (codex-review)
```

fallback した場合は `review-cross` を失敗扱いで示し、代わりに `review-sub` を OK で示す。

## 注意

- `/static-check`・`/unit-test` が項目をスキップと報告した場合（対応言語の変更なし等）も、その項目は通過扱いとし `/mark` でタグを設置する
- lint / build / test / doc-check / rule-check が成功したら `/mark <type>` でタグを設置する
- review-sub / review-cross は `/subagent-review` / `/codex-review` / `/claude-review` が完了時に自動でタグを設置する
