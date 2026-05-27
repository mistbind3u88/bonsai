---
name: tanaoroshi
description: 複数GitHubリポジトリのOpen Issue/PRを横断的に棚卸しする。body精査・コメント確認・前提作業の完了状況確認・テーマ別の構造整理・アクション提案を行う。
allowed-tools: Bash(tanaoroshi:*) Bash(gh repo view:*) Bash(gh api user:*) Read
---

# Issue/PR 棚卸し

複数の GitHub リポジトリを横断して Open Issue/PR を精査し、テーマ別に構造化した棚卸しレポートを作成する。

## 入力

引数でリポジトリをスペース区切りで指定する。

```
/tanaoroshi owner/repo1 owner/repo2
```

引数なしの場合は、カレントリポジトリに加えて CLAUDE.md または AGENTS.md の「関連リポジトリ」セクションに記載されたリポジトリも対象とする。関連リポジトリセクションがなければカレントリポジトリのみを対象とする。

## 出力

以下の構成の棚卸しレポートを会話に出力する。

1. テーマ別の構造整理（Open/Closed の依存関係ツリー）
2. 全体の進捗サマリー
3. アクション提案（マージ可能・クローズ可能・方針合意が必要なもの）

## 手順

### Phase 0: 対象リポジトリの決定

引数でリポジトリが指定されていない場合、以下の手順で対象を決定する。

1. `gh repo view --json nameWithOwner` でカレントリポジトリの `owner/repo` を取得する
2. CLAUDE.md または AGENTS.md の「関連リポジトリ」セクションからリポジトリ名を抽出する
3. カレントリポジトリ + 関連リポジトリを対象とする

### 無視リスト

`tanaoroshi` コマンドは、デフォルトでスキル側の `ignore` にリストされた Issue/PR を `summary`・`refs` の出力から除外する。実行リポジトリ側の ignore を追加したい場合は、`--ignore-file` で呼び出し元の cwd から読めるファイルを明示する。形式は `.gitignore` のコメント規則に準拠し、`#` で始まる行はコメントとして扱う。`--ignore-file` で指定するファイルは実在するパスにする。

特定リポジトリの ignore を注入する例:

```bash
tanaoroshi summary --ignore-file .tanaoroshi-ignore .results/tanaoroshi.json
tanaoroshi refs --ignore-file .tanaoroshi-ignore .results/tanaoroshi.json
```

### Phase 1: データ収集

`collect` で各リポジトリの Open Issue/PR を一括取得し、ファイルに保存する。

```bash
tanaoroshi collect <owner/repo> [owner/repo2 ...] > .results/tanaoroshi.json
```

レスポンスはリポジトリごとに `issues` と `prs` を持つ JSON オブジェクト（body 含む全フィールド）。

データが大きい場合は `summary` で body を除いた一覧を取得する。

```bash
tanaoroshi summary .results/tanaoroshi.json
```

### Phase 2: 参照先の解決

`refs` で collect 結果から body 内の参照を自動抽出する（重複排除済み）。

```bash
tanaoroshi refs .results/tanaoroshi.json
```

出力は `{source, ref}` のペア配列。検出する参照パターン:

- `#123` — 同一リポジトリ内の参照
- `owner/repo#123` — クロスリポジトリ参照
- `closes #123`, `fixes #123` — クローズ対象の明示
- `https://github.com/owner/repo/issues/123` — URL 形式の参照
- `https://github.com/owner/repo/pull/123` — URL 形式の参照

抽出した参照のうち Open 一覧にないものを `resolve` でステータス確認する。
引数では `#` の代わりに `:` を区切り文字として使用する（シェルでのコメント解釈を回避するため）。

```bash
tanaoroshi resolve \
  --issues <owner/repo:N> [owner/repo:N ...] \
  --prs <owner/repo:N> [owner/repo:N ...]
```

### Phase 2a: 直近コメントの確認

直近で更新のあった Open Issue/PR について、最新のコメントを確認する。

1. `summary` の出力から `updatedAt` が直近2週間以内の Issue/PR を抽出する
2. 該当する Issue/PR のコメントを `comments` サブコマンドで取得する

```bash
tanaoroshi comments <owner/repo:N> [owner/repo:N ...]
```

- Issue コメントと PR レビューコメントを時系列でマージして返す
- 出力の各コメントには `type` フィールド（`"comment"` = Issue コメント、`"review"` = PR レビューコメント）が付与される
- コメント内容は Phase 3 のテーマ整理で「進捗」「ブロッカー」「方針変更」の判断材料にする

### Phase 2b: リリース PR の分離

テーマ整理の前に、リリース PR を識別して分離する。以下のいずれかに該当する PR はリリース PR として扱う:

- タイトルに `release` / `Release` / `v0.0.0` 等のバージョン番号パターンを含む
- 作成者が bot（`dependabot`, `renovate`, `github-actions`, `release-please` 等）
- ラベルに `release` / `autorelease` を含む

リリース PR は Phase 3 のテーマ整理に含めず、Phase 4 のサマリーで「リリース運用」として別セクションにまとめる。

### Phase 3: テーマ別構造整理

収集したデータを以下の観点で分析し、テーマごとにグループ化する。

#### 3-1. グループ化の基準

- **body 内の相互参照**: Issue/PR が互いに言及し合っている場合は同一テーマ
- **closes/fixes 関係**: PR が Issue を closes している場合は親子関係
- **ブランチ名の共通プレフィックス**: 同一機能の連続 PR
- **タイトルの共通キーワード**: 同一施策の Issue 群
- **クロスリポジトリの対応関係**: 同名ブランチ、同一テーマの Issue が複数リポジトリに存在

#### 3-2. 各テーマの出力形式

テーマごとに以下を整理する。テーマを大見出し、トップレベル Issue または作業グループを小見出しにし、その配下に関係ツリーを置く。

```
## N. テーマ名

**ゴール**: このテーマで達成しようとしていること（1-2文）

### [Issue/OPEN] 親課題タイトル [owner/repo#127](https://github.com/owner/repo/issues/127)

（author: @author / assignee: @assignee / next: 実装）

- 実装: [PR/OPEN] タイトル [owner/repo#125](https://github.com/owner/repo/issues/125)（author: @author / next: レビュー）
  ├─ 前提: [PR/CLOSED] タイトル [owner/repo#123](https://github.com/owner/repo/issues/123)（closed: YYYY-MM-DD）
  ├─ ブロッカー: [Issue/OPEN] タイトル [owner/repo#124](https://github.com/owner/repo/issues/124)（author: @author / next: 方針合意）
  └─ 後続: [Issue/OPEN] タイトル [owner/repo#128](https://github.com/owner/repo/issues/128)（author: @author / next: 優先度判断）
- 関連: [Issue/OPEN] 親課題に直接関連する別Issue [owner/repo#130](https://github.com/owner/repo/issues/130)（author: @author / next: 優先度判断）

### 親 Issue がない作業グループの短い説明

- 実装: [PR/Draft] タイトル [owner/repo#129](https://github.com/owner/repo/issues/129)（author: @author / next: 下書き完成）
  └─ 代替: [PR/OPEN] 同じ目的の別PR [owner/repo#131](https://github.com/owner/repo/issues/131)（author: @author / next: レビュー）
- 関連: [Issue/OPEN] 作業グループに直接関連するIssue [owner/repo#132](https://github.com/owner/repo/issues/132)（author: @author / next: 優先度判断）

**進捗**: 要約と次のアクション
```

- テーマ見出しは `## N. テーマ名` とし、例: `## 1. Google HPA / フィード品質`
- トップレベル Issue は必ず `### [Issue/OPEN] タイトル [owner/repo#N](...)` の小見出しにする
- 親 Issue がない作業単位は、 `### 短い説明` の小見出しにする
- 小見出し直下にトップレベル Issue の後置メタ情報を 1 行で置く。トップレベル Issue の詳細は小見出し側に集約する
- 小見出しに直接関連する PR/Issue は第1階層の通常リスト（`-`）に置く。第1階層では `├─` / `└─` を使わない
- 第1階層の PR/Issue に関連する前提・ブロッカー・後続・代替・関連は第2階層のツリー（`  ├─` / `  └─`）に置く
- 小見出し自体に直接関係する前提・ブロッカー・後続は第1階層の通常リストに置いてよい。特定の PR/Issue にだけ紐づくものは第2階層のツリーに配置する
- 直接関連か間接関連かはツリーの深さで表す。中間見出しはテーマ見出しとトップレベル Issue 見出しに限定する
- ツリー行は `関係ラベル + 種別プレフィックス + タイトル + リンク + 後置メタ情報` の順で書く。小見出しなど関係ラベルを付けない行では `種別プレフィックス + タイトル + リンク + 後置メタ情報` の順にする。`@author` / `assignee` / `next` は行末の括弧内に配置する
- 後置メタ情報は Open では `（author: @author / assignee: @assignee / next: 実装）`、Closed では `（closed: YYYY-MM-DD）` を基本形にする。assignee がない場合は省略する
- ツリーの子要素には関係ラベルを必ず付ける。ラベルは増やしすぎず、関係が一目で分かる以下に絞る
- `実装`: 親 Issue を直接解決する PR、または実装本体の Issue/PR
- `前提`: 親・後続へ進むために完了が必要な Issue/PR。完了状態は `[PR/CLOSED]` / `[Issue/CLOSED]` と `closed` メタ情報で表現し、ラベルとは別に示す
- `ブロッカー`: 方針未決、レビュー停滞、外部依存など、現在進行を止めている Issue/PR
- `後続`: 親または前提完了後に着手する Issue/PR
- `代替`: 同じ目的の別 PR、作り直し PR、競合する解決案
- `関連`: 親子・前後関係までは断定できないが、同じテーマとして把握すべき Issue/PR
- 迷った場合は `関連` を使う。細かい作業種別（docs/refactor/test など）をラベル化せず、タイトルと種別プレフィックスで表現する
- 1 つの根に子要素が多すぎる場合は、重要な 3-6 件に絞り、残りは `└─ その他: N件` として要約する
- 判定の根拠が明確な項目だけを親子化する。根拠が弱い場合は別の作業グループ小見出しに分ける
- PR 概要欄またはコミットメッセージで GitHub が Issue closing keyword として扱う語は、`close` / `closes` / `closed` / `fix` / `fixes` / `fixed` / `resolve` / `resolves` / `resolved`。大文字やコロン付き（例: `CLOSES: #10`）も有効
- 上記 keyword は、PR が default branch を対象にしている場合に GitHub に解釈される。同一リポジトリは `KEYWORD #N`、別リポジトリは `KEYWORD owner/repo#N`、複数 Issue は各 Issue ごとに keyword を付ける
- PR 概要欄の closing keyword は、対象 Issue を閉じる意図が本文上明確で、PR の実装内容・タイトルと対象 Issue の目的が整合している場合だけ親子関係に使う
- PR 概要欄の明示リンクは、単なる参考リンク・関連資料・例示の可能性があるため、リンク周辺の文脈を確認し、`解決対象` / `前提` / `後続` / `関連` のどれに当たるかを判断してから配置する
- 判断できない参照は親子関係にせず `関連` として扱う。次にブランチ名・タイトル類似・コメント文脈で補完する
- クロスリポジトリ関係は同じ作業単位に含め、関係するリポジトリをまたいで整理する
- 1 つのテーマに Open が 10 件を超える場合は、全件列挙よりも相関が分かるツリーを優先し、低優先の独立 Issue は「その他」にまとめる
- 関係ラベルを付けない Issue/PR 行は、行頭に `[Issue/OPEN]` / `[Issue/CLOSED]` / `[PR/OPEN]` / `[PR/Draft]` / `[PR/CLOSED]` のいずれかを必ず付ける（merged PR は `[PR/CLOSED]` として扱う）。関係ラベルを付けるツリー行では、ラベル直後に種別プレフィックスを置く
- Open の PR/Issue はタイトル直後にリンクを置く。担当・次アクションは後置メタ情報の `author` / `assignee` / `next` にまとめる
- Closed は主体注釈を付けない（完了済みのため）
- Open は詳しく（タイトル・残件内容・ブロッカーの有無・`owner/repo#N`）
- Closed は簡潔に（タイトル・完了日・`owner/repo#N`）
- リンク（`owner/repo#N`）はタイトル直後に置き、その後に後置メタ情報を書く

#### 3-3. テーマに属さない個別 Issue

テーマにまとめられない独立した Issue/PR は以下のカテゴリで分類する:

- **インフラ・運用改善**: デプロイ、監視、コスト管理
- **コード品質**: リファクタ、リネーム、設定整理
- **管轄外の可能性**: 他リポジトリや他チームの対応が必要なもの

### Phase 4: 全体サマリーとアクション提案

#### 4-1. 進捗サマリー

テーマごとの Open 数・Approved 未マージ PR 数・Closed 数・進捗状況を表形式で出力する。単なる件数だけでなく、テーマ内で次に詰まっている関係を `詰まり` に明記する。

| テーマ   | Open | Approved | Closed | 詰まり                                                                      | 進捗 |
| -------- | :--: | :------: | :----: | --------------------------------------------------------------------------- | ---- |
| テーマ名 |  N   |    A     |   M    | レビュー依頼先未割当 / 他人 APPROVED 停滞 / Draft待ち / 方針未決 / 長期放置 | 要約 |

- `Approved` には `reviewDecision="APPROVED"` の未マージ PR を数える
- 自分が author の Approved PR は `PR のレビュー・完了対応` の `[merger]`、自分以外が author の Approved PR は `停滞 PR の整理` に配置する
- テーマ内の Open 件数は、未解決 Issue、Draft PR、レビュー中 PR、方針未決 Issue など、承認前の状態にあるものを数える

#### 4-2. アクション提案

アクション提案は主要カテゴリとして「PR のレビュー・完了対応」「完了確認・クローズ判断」「次に着手できるもの」「停滞 PR の整理」を置き、補助カテゴリとして「サイドカー候補」「方針合意が必要なもの」を続ける。各箇条書きには `[主体]` タグ（`[merger]` / `[author]` / `[reviewer]` / `[assignee]` / `[team]` / `[unassigned]` 等）を前置し、その後に `[Issue/…]` または `[PR/…]`、タイトル、リンク、後置メタ情報を続ける。例:

```
PR のレビュー・完了対応:
- [merger] [PR/OPEN] タイトル [owner/repo#N](https://github.com/owner/repo/issues/N)（author: @author / next: マージ可）
- [reviewer] [PR/OPEN] タイトル [owner/repo#N](https://github.com/owner/repo/issues/N)（author: @author / next: レビュー）

次に着手できるもの:
- [author] [PR/Draft] タイトル [owner/repo#N](https://github.com/owner/repo/issues/N)（author: @author / next: 下書き完成）
```

自分が次に動けば前進する項目を「すぐ対応できるもの」として抽出する。自分以外が next 主体の項目は、進捗サマリーまたは「停滞 PR の整理」に配置する。Draft 以外でレビュー・マージ段階に入っている PR は「PR のレビュー・完了対応」、完了後の確認対象となる Open Issue は「完了確認・クローズ判断」、Draft PR と未着手 follow-up は「次に着手できるもの」に配置する。

**すぐ対応できるものの抽出ルール**:

- この節は、主要カテゴリのうち自分が次に動くアクションを含む 3 カテゴリから「自分が次に動けば前進するもの」を抽出するためのルールである。出力見出しはアクション提案カテゴリのいずれかを使う
- 抽出した項目は元のカテゴリ（`PR のレビュー・完了対応` / `完了確認・クローズ判断` / `次に着手できるもの`）に置き、同じ項目は 1 か所にだけ列挙する
- 自分が次に動く項目には、実際の役割に応じて `[merger]` / `[author]` / `[reviewer]` / `[assignee]` などの主体タグを付ける

**PR のレビュー・完了対応**:

- Draft 以外の Open PR で、マージ実行・レビュー・再レビュー・bot 指摘対応・author の再プッシュ・レビュー依頼先の確認で前進するもの
- bot 指摘対応は、優先 5 / 8 の author 対応または補足判定の `author: 再プッシュ` として扱う
- 自分が author の APPROVED PR でマージ可能なものは `[merger]` として「PR のレビュー・完了対応」に配置する
- 自分が review request を受けている PR は `[reviewer]` として「PR のレビュー・完了対応」に配置する
- 自分が author の PR で指摘対応・再プッシュ・レビュー依頼先の確認が必要なものは `[author]` として「PR のレビュー・完了対応」に配置する
- 自分以外が author で、next 主体が author / owner / nobody の PR は、進捗サマリーの「詰まり」または「停滞 PR の整理」に配置する

**完了確認・クローズ判断**:

- この節には、次のアクションが残っている Open Issue だけを配置する。Closed/Merged PR/Issue は、クローズ判断の根拠として扱う
- Closed/Merged PR/Issue は、クローズ判断の根拠として Phase 3 の関係ツリーまたは依存関係サマリーに置く
- `[closer]` 直接対応する PR/Issue が Closed/Merged になり、親 Open Issue を閉じるか確認するだけのもの
- `[closer]` 上位 Issue に統合してクローズ可能な Open Issue
- 一通りリリース済み、動作確認済み、またはサブ Issue がすべて Closed になった親 Open Issue は「完了確認・クローズ判断」に配置する
- この節のアクション提案行には、next が `クローズ判断` / `上位Issueへ統合` のいずれかである Open Issue を配置する

**次に着手できるもの**:

- `[assignee]`/`[unassigned]` 前提の Closed により実装・調査に着手できる Open PR 作成前の Issue
- `[author]` Draft PR のうち、依存がすべて解決済みで下書き完成に進めるもの
- 依存未解決の Draft PR は、Phase 3 の関係ツリーまたは進捗サマリーで前提関係と詰まりを示す
- 進行中 PR に紐づく follow-up は、未着手の別 Issue として切り出されている場合にこの節へ配置する

**停滞 PR の整理**:

- 自分以外が author の Approved 未マージ PR は、この節に配置する。進捗サマリーの `Approved` 列には、author に関わらず全 Approved 未マージ PR を数える
- `[author]` 自分以外が author の APPROVED PR で、APPROVED から 7 日以内に未マージのもの（author にマージ判断を促す。レポートでは件数のみ示し、リストは進捗サマリーの「詰まり」列に集約してよい）
- `[owner]` 自分以外が author の APPROVED PR で、APPROVED から 7 日超未マージのもの（マージ要否を再評価）

**サイドカー候補**:

- サイドカー候補は「ネクストアクション主体の判定ルール」による自動分類とは別枠の、人手で選ぶ推薦である。`[sidecar]` 主体タグと `next: サイドカー調査` は判定表からは生成されず、この節の基準で選定したときにだけ用いる
- 本線テーマのクリティカルパスではないが、別作業として小さく進められる Open Issue を 3-5 件に絞って提案する
- サイドカー候補には、PR 作成前で、本線テーマから独立して進められ、担当者未設定または自分が担当できる Issue を選ぶ
- 軽量性を優先し、調査・再現テスト追加・命名整理・ドキュメント化・小さな設定整理で前進できるものを選ぶ
- 各行は `[sidecar] [Issue/OPEN] タイトル [owner/repo#N](...)（author: @author / next: サイドカー調査） - 推薦理由` の形式にする
- 推薦理由は 1 文で、なぜ本線から独立して進められるか、最初の着手が何かを明記する
- 候補がない場合は `該当なし` と明記する。長期放置 Issue は、軽量性と独立性を満たす場合にだけ候補にする

**方針合意が必要なもの**:

- `[team]` 議論 Issue で方針未決のもの
- `[team]` クロスリポジトリで方針が対立しているもの
- `[owner]` 長期放置（3ヶ月以上）で必要性の再評価が必要なもの

#### 4-3. 依存関係サマリー

アクション提案の前後に、テーマ横断で詰まりやすい依存関係を 5-10 件に絞ってツリー形式で出力する。根の行は `[ブロッカー]` / `[前提]` / `[重複/統合候補]` などの関係タグを行頭に置き、種別プレフィックスをその直後に続ける。形式は以下:

```
**依存関係サマリー**

[ブロッカー] [PR/OPEN] A [owner/repo#1](https://github.com/owner/repo/issues/1)（author: @author / next: レビュー）
└─ blocks: [Issue/OPEN] B [owner/repo#2](https://github.com/owner/repo/issues/2)（author: @author / next: 優先度判断）

[前提] [PR/CLOSED] A [owner/repo#3](https://github.com/owner/repo/issues/3)（closed: YYYY-MM-DD）
└─ enables: [Issue/OPEN] B [owner/repo#4](https://github.com/owner/repo/issues/4)（author: @author / assignee: @assignee / next: 実装/調査）

[重複/統合候補] [Issue/OPEN] A [owner/repo#5](https://github.com/owner/repo/issues/5)（author: @author / next: 優先度判断）
└─ overlaps: [Issue/OPEN] B [owner/repo#6](https://github.com/owner/repo/issues/6)（author: @author / next: 優先度判断）
```

- `blocks` は未完了の前提により後続が進めにくい関係に使う
- `enables` は Closed/Merged により後続が着手可能になった関係に使う
- `overlaps` は同じ目的の Issue/PR が併存し、統合判断が必要な関係に使う
- 関係が推測に留まる場合は `関連` として扱い、確定関係は根拠がある場合にだけ使う

## 出力の共通ルール

- **Issue/PR への言及は、タイトル末尾・本文中・括弧内を問わず、すべて `[owner/repo#N](https://github.com/owner/repo/issues/N)` 形式のクリッカブルリンクで記載する。`#N` や `owner/repo#N` だけの素の表記ではターミナル上でリンクにならない。** 範囲指定（`#256〜#258`）や列挙（`#204, #205`）も個別に完全形式で書く
- **種別プレフィックス** `[Issue/OPEN]` / `[Issue/CLOSED]` / `[PR/OPEN]` / `[PR/Draft]` / `[PR/CLOSED]` を必ず行頭に付ける（merged PR は `[PR/CLOSED]`）。前置ラベルを使う節では、その節のルールに従う
- **Open の PR/Issue** はタイトル直後にリンクを置く。リンクの後に `（author: ... / assignee: ... / next: ...）` として後置する
- Closed はリンクの後に `（closed: YYYY-MM-DD）` を付ける
- リンクはタイトル直後に置く
- テーマ別ツリー、個別 Issue 一覧、サマリー表、アクション提案のいずれでも同じ形式を適用する
- Open を先に、Closed を後に記載する
- PR の状態は Closed に統一する（Merged も Closed として扱う）
- テーマ内では「リポジトリ別」ではなく「作業単位別」に並べる。リポジトリ境界よりも、Issue/PR の相関と次アクションの見通しを優先する
- 親子関係は根拠がある項目だけをツリー化する。根拠が弱い場合は `関連` と明記する
- リポジトリ内のスクリプトやファイルに言及する場合は、絶対パスではなくリポジトリルートからの相対パスで記載する

## ネクストアクション主体の判定ルール

Open の Issue/PR 各行の後置メタ情報に入れる `next` は以下の基準で決める。判定に必要な `isDraft` / `reviewDecision` / `author` / `assignees` / `reviewRequests` / `latestReviews` / ラベルは `summary` の出力に含まれる（`collect` で `reviewRequests,latestReviews` を取得済み）。必要なら `comments` で直近コメントの向き先を補足する。

本人視点の分類が必要な場合は、`gh api user --jq .login` またはユーザーから明示された GitHub login を「自分」として扱う。

### PR

**主判定**（上の行から評価し、最初に一致した行を採用する）:

| 優先 | 条件                                                                                                 | 主体: アクション                                                     |
| ---- | ---------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| 1    | `isDraft=true`                                                                                       | `author: 下書き完成`                                                 |
| 2    | `reviewDecision="APPROVED"` かつ `author.login` が自分                                               | `merger: マージ実行可`                                               |
| 3    | `reviewDecision="APPROVED"` かつ `author.login` が自分以外 かつ APPROVED からの経過日数が 7 日以内   | `author: マージ判断`                                                 |
| 4    | `reviewDecision="APPROVED"` かつ `author.login` が自分以外 かつ APPROVED からの経過日数が 7 日超     | `owner: マージ要否再評価`                                            |
| 5    | `reviewDecision="CHANGES_REQUESTED"` かつ `author.login` が自分                                      | `author: 指摘対応`                                                   |
| 6    | `reviewDecision="CHANGES_REQUESTED"` かつ `author.login` が自分以外                                  | `author: 指摘対応待ち`（自分が動くことはない。「停滞 PR の整理」へ） |
| 7    | (`reviewDecision="REVIEW_REQUIRED"` または空・未設定) かつ `reviewRequests[].login` に自分が含まれる | `reviewer: レビュー`                                                 |
| 8    | (`reviewDecision="REVIEW_REQUIRED"` または空・未設定) かつ `author.login` が自分                     | `author: レビュー待ち`                                               |
| 9    | `reviewDecision="REVIEW_REQUIRED"` または空・未設定（自分が author でもレビュー依頼先でもない）      | `nobody: レビュー依頼待ち`                                           |

**重要**:

- 優先 9 の `nobody: レビュー依頼待ち` は「自分には次の動きがない」状態を意味する。進捗サマリーの「詰まり」列に集約する。
- 優先 4 / 6 / 9 のように **自分が動かない PR** は、レポート上では `[author]`（自分以外）/ `[owner]`（停滞）/ `[nobody]`（依頼待ち）に分類し、「すぐ対応できるもの」の対象は自分が次に動く PR/Issue に限定する。
- `reviewDecision="APPROVED"` の未マージ PR は、進捗サマリーの `Approved` 列に数える。テーマ別ツリーでは `[PR/OPEN]` のまま表示し、レビュー待ち Open PR とは別の状態として扱う。
- APPROVED からの経過日数は `latestReviews[]` の `state="APPROVED"` で `submittedAt` が最新のものを基準に計算する（複数いる場合は最も新しい APPROVED の日時）。

**補足判定**（任意）: `comments` を取得済みの場合に限り、以下で主判定の結果を上書きしてよい。判別できない場合は主判定の結果をそのまま使う。

- 主判定が優先 5 / 7 / 8 のレビュー中 PR で、直近コメントの投稿者が author 以外 → `author: 再プッシュ`
- 主判定が優先 7 のレビュー待ち PR で、直近コメントの投稿者が author → `reviewer: 再レビュー`

### Issue

上の行から評価し、最初に一致した行を採用する。

| 優先 | 条件                                  | 主体: アクション         |
| ---- | ------------------------------------- | ------------------------ |
| 1    | `assignees` あり                      | `assignee: 実装/調査`    |
| 2    | 議論系ラベル（`discussion` 等）を持つ | `team: 方針合意`         |
| 3    | 上記以外                              | `unassigned: 優先度判断` |

`closer: クローズ判断` / `closer: 上位Issueへ統合` は、Phase 4-2 の「完了確認・クローズ判断」のアクション提案行で付与する。Issue 判定表は通常の Open Issue next を担う。Phase 3 の関係ツリーでは主判定の next を維持し、直接対応する PR/Issue の Closed/Merged、上位 Issue への統合可否、サブ Issue の完了状況などを関係ツリーとコメント文脈から判断する。`closer` は Open Issue のクローズ判断アクションに付与し、Closed/Merged PR/Issue は Open Issue のクローズ判断根拠として扱う。

### アクション提案（Phase 4-2）の主体タグ

アクション提案セクションでは行頭に `[主体]` タグを付ける（上記の「主体」を `[…]` で囲んだ形）。代表例:

- `[merger]` — 自分が author の APPROVED 済み PR でマージ実行を待っている
- `[author]` — レビュー待ち、下書き完成、指摘対応、再プッシュ、マージ判断を待っている（自分・他者双方を含む）
- `[reviewer]` — `reviewRequests` に自分が含まれている PR のレビュー・再レビューを待っている
- `[assignee]` — 指名済みの担当者の実装/調査を待っている
- `[team]` — 方針合意を待っている
- `[closer]` — クローズ判断を待っている
- `[owner]` — 長期放置・APPROVED 後の停滞などで必要性／マージ要否の再評価が必要
- `[unassigned]` — 担当未定で優先度判断を待っている
- `[nobody]` — レビュー依頼先が未割当で誰の手番でもない（進捗サマリーに集約し、アクション提案には載せない）
- `[sidecar]` — 本線と独立して小さく進められる候補。通常の「すぐ対応できるもの」とは分けて提案する

## 注意

- Closed の前提作業は「何が完了したか」の事実のみ記載し、詳細な説明は省く
- リリース PR は Phase 2b で分離し、個別テーマに含めない。Phase 4 のサマリーで「リリース運用」として別セクションにまとめる
