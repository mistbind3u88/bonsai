---
name: check
description: markで未通過の品質チェック（lint / build / test / doc-check / rule-check / review-sub / review-cross）を検出し、対応するスキルを実行・集約する。
allowed-tools: Bash(git diff:*)
---

# checkスキル

markタグが未設置の品質チェックを検出し、対応するスキルへ実行を委譲して結果を集約する。lint・buildは `/static-check` に、testは `/unit-test` に、ドキュメント整合は `/doc-check` に、ルール適合は `/rule-check` に、reviewは引数で選ぶ2種類（sub・cross）に委譲する。

## 引数

`$ARGUMENTS` で走らせるreviewを選ぶ:

| 引数                 | 走らせるreview                                          | 主な用途                                                     |
| -------------------- | ------------------------------------------------------- | ------------------------------------------------------------ |
| （なし・デフォルト） | cross-review（`/codex-review` または `/claude-review`） | `/push` 経由・手動 `/check`                                  |
| `--review=sub`       | sub-review（`/subagent-review`）                        | `/commit` 各コミット末尾・autosquash / 履歴書き換え完了後    |
| `--review=skip`      | なし                                                    | `/commit` autosquash中間step・手動でレビュー抜きで見たいとき |

**`review-cross` は `review-sub` の上位**として扱う。`review-cross` がHEADタグ済みのときは、`--review=sub` でもcrossの通過をもってsubも通過扱いとする（crossがsubの役割を含めて代替）。代替の方向はcross→subの一方向に限る。

## 手順

### 1. チェックタグを確認する

スキル `/mark --status` を実行して、現在のHEADのチェック通過状況を確認する。

引数に応じて対象項目を判定し、すべて現在のHEADにタグ設置済みなら、その旨を報告して終了する:

- 共通: `lint` / `build` / `test` / `doc-check` / `rule-check` がHEADタグ済み
- `--review=sub`: `review-sub` または `review-cross` がHEADタグ済み（crossはsubの上位として代替）
- 引数なし（=cross）: `review-cross` がHEADタグ済み
- `--review=skip`: lint / build / test / doc-check / rule-checkのみ判定する

### 1a. 履歴書き換え後のタグ引き継ぎ

タグが現在のHEADにないが、タグが付いているコミットが存在する場合、そのコミットと現在のHEADの間に差分があるか確認する。

```bash
git diff <タグのコミット> HEAD
```

差分がなければ、autosquash・`git history reword` などでコミットハッシュが変わっただけで内容は同一なので、スキル `/mark <type>` を実行してタグを現在のHEADに付け替える。差分がある項目のみ再実行する。

### 2. 未通過の項目を実行する

タグが現在のHEADにない項目を、対応するスキルへ委譲して実行する。

| 項目         | 担当スキル                              | タグなし時の処理                                                                                   |
| ------------ | --------------------------------------- | -------------------------------------------------------------------------------------------------- |
| lint, build  | `/static-check`                         | 実行し、PASSの項目を `/mark lint`・`/mark build` で設置                                             |
| test         | `/unit-test`                            | 実行し、PASSなら `/mark test` を設置                                                                |
| doc-check    | `/doc-check`                            | 指摘なし、採用指摘の修正後再チェック通過、または全指摘への判断記録後に `doc-check` 設置            |
| rule-check   | `/rule-check`                           | 指摘なし、採用指摘の修正後再チェック通過、または全指摘への判断記録後に `rule-check` 設置           |
| review-sub   | `/subagent-review`                      | 引数 `--review=sub` 指定時のみ実行。指摘なし、または全指摘へのユーザー判断後に `review-sub` 設置   |
| review-cross | `/codex-review` または `/claude-review` | 引数なし（デフォルト）の場合に実行。指摘なし、または全指摘へのユーザー判断後に `review-cross` 設置 |

`/static-check` はlintとbuildを、`/unit-test` はtestを、それぞれリポジトリに合わせて検出・スコープ判定・実行し、項目別に結果を報告する。`check` はその成否を受けて `/mark` で各タグを設置する。

`/static-check` はlint・buildを1単位で実行するため、lint・buildの両方が現在のHEADにタグ済みのときだけ起動をスキップする。どちらかが未タグなら `/static-check` を実行する。`/unit-test`・`/doc-check`・`/rule-check` は対応する項目が現在のHEADにタグ済みならスキップする。

`/doc-check` の起動・成否判断・`/mark doc-check` はメインエージェントが担う（`/doc-check` 内部の読み取り検査はサブエージェントへ委譲される）。

`/rule-check` の起動・成否判断・`/mark rule-check` はメインエージェントが担う（`/rule-check` 内部の読み取り検査はサブエージェントへ委譲される）。

cross-reviewは現在のエージェントとは別のエージェントに依頼する。

- Codex上で実行している場合: スキル `/claude-review` を実行する
- Claude Code上で実行している場合: スキル `/codex-review` を実行する
- 実行環境を判断できない場合: ユーザーに確認する

実行順序:

1. `/static-check`・`/unit-test` を並列で実行する。互いの結果に依存せず、サブエージェントも使わないため直列化しない
2. `/doc-check` と `/rule-check` を並列で実行する。
   - どちらもサブエージェントを起動するため、`/check` 側で先に `/subagent-check` を実行し、2件を並列実行できる実行計画、既存サブエージェントの再利用方針、依頼文の委譲境界を確認する。
   - 実行計画に新規起動が含まれる場合は、新規起動が必要な分だけ余剰枠を確保してから両方を開始する。余剰枠が不足する場合は、既存サブエージェントの結果回収・再利用・完了待ちで計画を満たす。
   - 実行計画を確定できたら、`/doc-check` と `/rule-check` の読み取り検査をサブエージェントへ並列依頼すること、依頼範囲、メインエージェントに残す判断をまとめてユーザーに提示する。`/check` が明示起動された範囲では、この提示を情報共有として扱い、追加承認は求めずに進める。
   - `doc-check` と `rule-check` は、`/check` 側の包括確認を前提として起動し、個別skill内で `/subagent-check` を再実行しない。
   - 同じ `/check` 実行内で修正後の再検査を行う場合は、`doc-check` と `rule-check` それぞれの既存サブエージェントを再利用する。役割が違うため、`doc-check` agentと `rule-check` agentを相互に流用しない。
   - 両方の検査結果をそろえてから次へ進む。
3. 4項目すべてが成功したら、引数で選ばれたreviewを実行する。引数なしならcross-review、`--review=sub` ならsub-review。`--review=skip` のときはlint / build / test / doc-check / rule-checkのみで終える。reviewは他チェックを通過したコードに対して行うため、必ず最後に実行する。`--review=sub` で `review-cross` が既にHEADタグ済みのときは、crossの通過をもってsub-reviewも通過扱いとする（crossが上位として代替）

いずれかの項目が、コマンド失敗、サブエージェント起動不可、必要な入力不足などで検査自体を完了できない場合は、実行中の項目の完了を待ってから全項目の結果をまとめて報告し、停止する。

検査自体は完了し、`要修正` / `検討推奨` / `軽微` などの指摘が返った場合は、次の操作へ進む前に「2a. 指摘事項の整合性検証」を実行する。採否判断の結果に応じて、修正・必要な再チェック、ユーザー確認、または対応不要としての記録へ分岐する。`doc-check` / `rule-check` は、指摘がない場合、採用した指摘を修正して必要な再チェックを通過した場合、または全指摘を `対応不要` / `ループ懸念` と判断して理由を最終報告に記録した場合に `/check` 側でmarkを設置する。`保留` の指摘が残る場合はmarkを設置する前にユーザー判断を確認する。`review-sub` / `review-cross` は各review skillの契約に従い、指摘なし、または全指摘へのユーザー判断が揃った時点で各review skillがmarkを設置する。

### 2a. 指摘事項の整合性検証

開始条件:

- `/doc-check`、`/rule-check`、`/subagent-review`、`/codex-review`、`/claude-review` から `要修正`、`検討推奨`、`軽微` などの指摘が返った場合、この手順へ進む。
- `/doc-check` と `/rule-check` の指摘はメインエージェントが採否を判断する。
- review系skillの指摘はメインエージェントが採否案と整合性を整理し、各review skillの契約に従ってユーザー判断を確認する。
- サブエージェントや外部レビューの指摘を、そのまま自動修正しない。

確認項目:

メインエージェントは、各指摘について以下を確認する。

- 指摘が差分、既存ルール、ユーザー指示、実行時の挙動に根拠を持つか
- 指摘の対応が、上位ルール、近傍 `AGENTS.md` / `CLAUDE.md`、既存skillの責務境界、`allowed-tools`、READMEの依存図、テンプレートの出力契約と整合するか
- 指摘の対応によって、新しい承認フロー、委譲フロー、テンプレート項目、mark条件、README反映などの追加修正が必要になるか
- 指摘が任意の表現改善、好みの差、今回の差分から離れた将来課題、既存から残る別課題のどれに該当するか
- 指摘対応を繰り返すことで、別の指摘を誘発するだけのループになっていないか

分類:

判定は以下に分類する。

| 分類         | 意味                                                                 | 次の行動                                                     |
| ------------ | -------------------------------------------------------------------- | ------------------------------------------------------------ |
| `採用`       | 指摘が妥当で、対応しても既存ルール・責務境界・手順と整合する         | 修正し、必要な関連ドキュメント・テンプレート・markを更新する |
| `保留`       | 指摘は妥当だが、複数の解決方法がありユーザー判断が必要               | 選択肢と影響を示してユーザーに確認する                       |
| `対応不要`   | 実行上の不整合ではない、任意改善、既存からの別課題、またはスコープ外 | 修正せず、理由を記録する                                     |
| `ループ懸念` | 対応すると別の整合指摘を誘発し、今回の目的を超えて修正が循環する     | 修正せず、ループ要因と止める理由を記録する                   |

分類後の行動:

- review系skillの `採用` / `対応不要` / `ループ懸念` は、ユーザー判断が揃ってから確定する。ユーザー判断前は採否案として扱う。
- `採用` した指摘を修正した場合は、修正範囲に応じて必要なチェックだけを再実行する。再実行後に追加指摘が出た場合も同じ分類を行う。
- `対応不要` または `ループ懸念` と判断した指摘は、同じ内容で再度指摘されても自動修正せず、既存の判断を引き継ぐ。
- 矛盾やループを発生させる可能性がある指摘、または対応しないと決めた指摘については、最終報告に判断を明記する。

```text
指摘事項の整合性検証:
  <review-source>: <file-or-topic>
    分類: 採用 | 保留 | 対応不要 | ループ懸念
    判断理由: <差分・ルール・責務境界・スコープに照らした理由>
    次の行動: <修正内容 / ユーザー確認 / 対応しない理由>
```

### 2b. cross-review失敗時のfallback

引数なし（デフォルト = cross-review）で `/codex-review` または `/claude-review` がCLI実行不可等で失敗した場合、`/subagent-review` へfallbackしてsub-reviewを走らせる。

- 成功時: `/subagent-review` で指摘がない場合、または全指摘へのユーザー判断が揃った時点で `review-sub` がHEADタグに設置される（`review-cross` は未通過のまま残す）
- 結果報告で「cross-reviewが実行できずsub-reviewにfallbackした」旨と、`review-cross` 未通過を明示する
- fallbackの有無は `/check` 内で完結する。呼び出し元（`/push` 等）は通過状況のtagを見て判断する

### 3. 結果サマリーを表示する

実行した全チェック項目の結果を一覧で表示する。

```text
チェック結果:
  lint:          OK
  build:         OK
  test:          OK
  doc-check:     OK
  rule-check:    OK
  review-cross:  OK (codex-review)
```

fallbackした場合は `review-cross` を失敗扱いで示し、代わりに `review-sub` をOKで示す。

## 注意

- `/static-check`・`/unit-test` が項目をスキップと報告した場合（対応言語の変更なし等）も、その項目は通過扱いとし `/mark` でタグを設置する
- lint / build / testは、対応項目の成功またはスキップを確認したら `/mark <type>` でタグを設置する
- doc-check / rule-checkは、指摘なし、採用指摘の修正後の再チェック通過、または全指摘への `対応不要` / `ループ懸念` 判断の記録を確認したら `/mark <type>` でタグを設置する
- review-sub / review-crossは、指摘なし、または全指摘へのユーザー判断が揃った時点で `/subagent-review` / `/codex-review` / `/claude-review` がタグを設置する
