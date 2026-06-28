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

| 項目         | 担当スキル                              | タグなし時の処理                                                                                   |
| ------------ | --------------------------------------- | -------------------------------------------------------------------------------------------------- |
| lint, build  | `/static-check`                         | 実行し、PASS の項目を `/mark lint`・`/mark build` で設置                                           |
| test         | `/unit-test`                            | 実行し、PASS なら `/mark test` を設置                                                              |
| doc-check    | `/doc-check`                            | 指摘なし、採用指摘の修正後再チェック通過、または全指摘への判断記録後に `doc-check` 設置            |
| rule-check   | `/rule-check`                           | 指摘なし、採用指摘の修正後再チェック通過、または全指摘への判断記録後に `rule-check` 設置           |
| review-sub   | `/subagent-review`                      | 引数 `--review=sub` 指定時のみ実行。指摘なし、または全指摘へのユーザー判断後に `review-sub` 設置   |
| review-cross | `/codex-review` または `/claude-review` | 引数なし（デフォルト）の場合に実行。指摘なし、または全指摘へのユーザー判断後に `review-cross` 設置 |

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
2. `/doc-check` と `/rule-check` を並列で実行する。
   - どちらもサブエージェントを起動するため、`/check` 側で先に `/subagent-check` を実行し、2 件を並列実行できる実行計画、既存サブエージェントの再利用方針、依頼文の委譲境界を確認する。
   - 実行計画に新規起動が含まれる場合は、新規起動が必要な分だけ余剰枠を確保してから両方を開始する。余剰枠が不足する場合は、既存サブエージェントの結果回収・再利用・完了待ちで計画を満たす。
   - 実行計画を確定できたら、`/doc-check` と `/rule-check` の読み取り検査をサブエージェントへ並列依頼すること、依頼範囲、メインエージェントに残す判断をまとめてユーザーに提示し、承認を得る。この承認を、`/check` から起動する読み取り検査に対する事前承認として扱う。
   - `doc-check` と `rule-check` は、`/check` 側の包括確認を前提として起動し、個別 skill 内で `/subagent-check` を再実行しない。
   - 同じ `/check` 実行内で修正後の再検査を行う場合は、`doc-check` と `rule-check` それぞれの既存サブエージェントを再利用する。役割が違うため、`doc-check` agent と `rule-check` agent を相互に流用しない。
   - 両方の検査結果をそろえてから次へ進む。
3. 4 項目すべてが成功したら、引数で選ばれた review を実行する。引数なしなら cross-review、`--review=sub` なら sub-review。`--review=skip` のときは lint / build / test / doc-check / rule-check のみで終える。review は他チェックを通過したコードに対して行うため、必ず最後に実行する。`--review=sub` で `review-cross` が既に HEAD タグ済みのときは、cross の通過をもって sub-review も通過扱いとする（cross が上位として代替）

いずれかの項目が、コマンド失敗、サブエージェント起動不可、必要な入力不足などで検査自体を完了できない場合は、実行中の項目の完了を待ってから全項目の結果をまとめて報告し、停止する。

検査自体は完了し、`要修正` / `検討推奨` / `軽微` などの指摘が返った場合は、次の操作へ進む前に「2a. 指摘事項の整合性検証」を実行する。採否判断の結果に応じて、修正・必要な再チェック、ユーザー確認、または対応不要としての記録へ分岐する。`doc-check` / `rule-check` は、指摘がない場合、採用した指摘を修正して必要な再チェックを通過した場合、または全指摘を `対応不要` / `ループ懸念` と判断して理由を最終報告に記録した場合に `/check` 側で mark を設置する。`保留` の指摘が残る場合は mark を設置する前にユーザー判断を確認する。`review-sub` / `review-cross` は各 review skill の契約に従い、指摘なし、または全指摘へのユーザー判断が揃った時点で各 review skill が mark を設置する。

### 2a. 指摘事項の整合性検証

`/doc-check`、`/rule-check`、`/subagent-review`、`/codex-review`、`/claude-review` から `要修正`、`検討推奨`、`軽微` などの指摘が返った場合、メインエージェントが各指摘の根拠と既存手順との整合性を検証してから次の操作へ進む。`/doc-check` と `/rule-check` の指摘はメインエージェントが採否を判断する。review 系 skill の指摘はメインエージェントが採否案と整合性を整理し、各 review skill の契約に従ってユーザー判断を確認する。サブエージェントや外部レビューの指摘を、そのまま自動修正しない。

メインエージェントは、各指摘について以下を確認する。

- 指摘が差分、既存ルール、ユーザー指示、実行時の挙動に根拠を持つか
- 指摘の対応が、上位ルール、近傍 `AGENTS.md` / `CLAUDE.md`、既存 skill の責務境界、`allowed-tools`、README の依存図、テンプレートの出力契約と整合するか
- 指摘の対応によって、新しい承認フロー、委譲フロー、テンプレート項目、mark 条件、README 反映などの追加修正が必要になるか
- 指摘が任意の表現改善、好みの差、今回の差分から離れた将来課題、既存から残る別課題のどれに該当するか
- 指摘対応を繰り返すことで、別の指摘を誘発するだけのループになっていないか

判定は以下に分類する。

| 分類         | 意味                                                                 | 次の行動                                                      |
| ------------ | -------------------------------------------------------------------- | ------------------------------------------------------------- |
| `採用`       | 指摘が妥当で、対応しても既存ルール・責務境界・手順と整合する         | 修正し、必要な関連ドキュメント・テンプレート・mark を更新する |
| `保留`       | 指摘は妥当だが、複数の解決方法がありユーザー判断が必要               | 選択肢と影響を示してユーザーに確認する                        |
| `対応不要`   | 実行上の不整合ではない、任意改善、既存からの別課題、またはスコープ外 | 修正せず、理由を記録する                                      |
| `ループ懸念` | 対応すると別の整合指摘を誘発し、今回の目的を超えて修正が循環する     | 修正せず、ループ要因と止める理由を記録する                    |

review 系 skill の `採用` / `対応不要` / `ループ懸念` は、ユーザー判断が揃ってから確定する。ユーザー判断前は採否案として扱う。

`採用` した指摘を修正した場合は、修正範囲に応じて必要なチェックだけを再実行する。再実行後に追加指摘が出た場合も同じ分類を行う。`対応不要` または `ループ懸念` と判断した指摘は、同じ内容で再度指摘されても自動修正せず、既存の判断を引き継ぐ。

矛盾やループを発生させる可能性がある指摘、または対応しないと決めた指摘については、最終報告に判断を明記する。

```text
指摘事項の整合性検証:
  <review-source>: <file-or-topic>
    分類: 採用 | 保留 | 対応不要 | ループ懸念
    判断理由: <差分・ルール・責務境界・スコープに照らした理由>
    次の行動: <修正内容 / ユーザー確認 / 対応しない理由>
```

### 2b. cross-review 失敗時の fallback

引数なし（デフォルト = cross-review）で `/codex-review` または `/claude-review` が CLI 実行不可等で失敗した場合、`/subagent-review` へ fallback して sub-review を走らせる。

- 成功時: `/subagent-review` で指摘がない場合、または全指摘へのユーザー判断が揃った時点で `review-sub` が HEAD タグに設置される（`review-cross` は未通過のまま残す）
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
- lint / build / test は、対応項目の成功またはスキップを確認したら `/mark <type>` でタグを設置する
- doc-check / rule-check は、指摘なし、採用指摘の修正後の再チェック通過、または全指摘への `対応不要` / `ループ懸念` 判断の記録を確認したら `/mark <type>` でタグを設置する
- review-sub / review-cross は、指摘なし、または全指摘へのユーザー判断が揃った時点で `/subagent-review` / `/codex-review` / `/claude-review` がタグを設置する
