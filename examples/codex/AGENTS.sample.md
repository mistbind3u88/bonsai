# Codex Global Settings Sample

このファイルは Codex グローバル設定向けのサンプルです。コピー先では `AGENTS.md` として配置します。このリポジトリ内では magic filename を避けるため `AGENTS.sample.md` として管理します。

利用時は、`<path-to-skills-repo>`、`<owner/repo>`、`<branch>` などのプレースホルダを自分の環境に合わせて置き換えてください。token、API key、private repository URL、社内ホスト名、個人のホームディレクトリは自分のローカル設定側だけで扱います。

## Skill 利用

- セッション開始時に available skill list を確認し、ユーザー要求、リポジトリ workflow、実行中の手順に合う skill を判断する。
- 一致する skill がある場合は、その `SKILL.md` を読んでから適用し、どの skill をなぜ使うかを短く述べる。
- 複数の skill が一致する場合は、タスクを満たす最小セットを選び、作業順序に沿って実行する。
- skill の実行中は、`SKILL.md` に書かれた手順を順に実行する。`とばしてよい`、`省略してよい` などの明示がある手順以外は省かない。
- 実行できない手順がある場合は、その手順を未実施のまま止め、阻害要因と不足している前提条件をユーザーへ報告する。明示的な省略許可がないまま黙って先へ進めない。

## Skill の前提条件

- Skill の実行中に前提条件が満たされていないことが分かった場合は、その前提条件を満たせる既存 skill があるか確認する。
- 既存 skill で解決できる場合は、その skill を実行して前提条件を満たしてから元の skill に戻る。
- 前提条件の解決に使った skill、解決結果、元の skill へ戻った判断をユーザーへ報告する。
- 既存 skill で解決できない場合、または複数の解決方法があり選択が必要な場合は、足りない前提条件と選択肢を具体的に示してユーザーへ確認する。

## 承認が必要な操作

- skill の手順を読んだ時点で、後続にサブエージェント起動、sandbox 外実行、GitHub 投稿、push、force push など承認が必要な操作が含まれることが明らかな場合は、実行前に承認が必要な操作と理由をまとめて提示し、ユーザーの承認を得てから進める。
- 実行途中で初めて判明した承認事項は、その時点で追加確認する。
- ユーザーが skill を明示的に起動し、その `SKILL.md` にサブエージェント起動が必須または標準手順として明記されている場合、その skill 実行範囲に限りサブエージェント起動はユーザー承認済みとして扱う。
- sandbox 外実行、GitHub 投稿、push、force push、破壊的操作など、サブエージェント起動以外の承認事項は個別に確認する。

## サブエージェント

- サブエージェントを起動する前に、既存サブエージェントの完了状態、役割、起動上限の余剰、再利用可否を確認する。
- 完了済みサブエージェントがある場合は結果を回収し、同じ `work_unit` で再利用予定がないものを閉じて枠を空ける。
- 同じまとまった作業内で同じ役割の再検査を行う場合は、既存サブエージェントの再利用を優先する。
- サブエージェントには、検査、文面確認、想定質問への回答などの限定タスクだけを渡す。
- スキルの起動、手順解釈、結果回収、最終判断、mark、ユーザー報告はメインエージェントが担う。
- サブエージェントへの依頼文には、`work_unit`、`role`、`reuse_policy`、返すべき結果、メインエージェントに残す判断を含める。

## Review

- Codex 上で作業した差分の cross-review では、必要に応じて `claude -p` を read-only の外部レビューとして使う。
- `claude -p` に渡す内容は、レビュー対象の差分、変更概要、設計判断、既知の論点、読み取り専用の許可に限定する。
- `claude -p` の実行では、読み取り専用の権限だけを許可する。
- CLI 実行不可や認証未完了で cross-review が完了できない場合は、呼び出し元 workflow の fallback 方針に従う。
- review 系 skill のフィードバックループでは、過去の指摘、判断、対応状態を worktree-local な review log に記録し、次回 reviewer へ明示的に渡す。

## GitHub と token

- token 値は secret として扱う。
- token の存在確認が必要な場合は、存在有無だけを確認する。
- GitHub CLI に token を渡す場合は、プロセス内の環境変数として渡し、ログに値が出ないコマンドだけを使う。
- SAML SSO enforcement などで GitHub app / connector が読めない場合は、`gh` CLI で読み取りを試し、必要な organization authorization をユーザーへ報告する。

## PATH と補助コマンド

- このリポジトリの skill が使う補助コマンドは `<path-to-skills-repo>/tools` に集約される前提で扱う。
- shell 設定では `<path-to-skills-repo>/tools` を PATH に追加する。
- skill の `allowed-tools` には skill 内部パスを直接書かず、PATH に載る bare command 名を使う。

例:

```zsh
path=(<path-to-skills-repo>/tools $path)
export PATH
```

POSIX shell では次の形を使う。

```sh
export PATH="<path-to-skills-repo>/tools:$PATH"
```

## Git 操作

- main / master で開発作業を始める前に、作業ブランチを作成または使用する。
- リポジトリまたはユーザーが default branch での直接作業を明示的に許可している場合は、その指示を優先する。
- コミットメッセージのタイトル（subject）は、リポジトリまたはユーザーが別形式を指定している場合を除き日本語で書く。
- force push は、rebase や autosquash などで履歴が実際に diverge した場合に `--force-with-lease` を使う。
- fast-forward 可能な場合は通常の push を使う。

## ドキュメント

- コード変更とドキュメント変更は 1 つの開発単位として扱う。
- 公開挙動、データ schema、運用、開発者 workflow、テスト、設計意図に影響するコード変更では、関連ドキュメントを同じ実装フェーズで更新する。
- 変更せずに質問へ回答する場合は、read-only の確認に留める。
- 情報収集と相談の開始時は、目的、前提、未確認事項、論点、選択肢を整理する。
