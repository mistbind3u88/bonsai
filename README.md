# dev-skills

開発で使っている、比較的汎用的なSKILLをまとめて育てていくリポジトリです。

## SKILL is 何？

SKILL.mdは、コーディングエージェントに「この流れで進めてね」と伝えるための手順書です。

コーディングエージェントの活用が広がり始めたごく初期には、各エージェントごとに形式がばらばらでみんな手探りで運用していました。最近は `AGENTS.md + SKILL.md` という形に寄ってきていて、だいぶ扱いやすくなってきた印象です。

### 参考リンク

- [Agent Skills公式](https://agentskills.io/home)
- [AGENTS.md公式](https://agents.md/)

## AGENTS.mdとSKILL.mdの違い

ざっくり言うと、`AGENTS.md` は「このディレクトリで守るルール」、`SKILL.md` は「特定タスクの進め方」です。

```mermaid
flowchart TD
    A["リポジトリ全体の方針"] --> B["AGENTS.md"]
    B --> C["スコープ: ディレクトリ単位"]
    B --> D["内容: 制約・判断基準・安全ルール"]
    E["個別タスクを実行"] --> F["SKILL.md"]
    F --> G["スコープ: スキル単位"]
    F --> H["内容: 手順・コマンド・期待成果"]
```

| 項目         | AGENTS.md                            | SKILL.md                    |
| ------------ | ------------------------------------ | --------------------------- |
| 主な役割     | ガードレールを定義                   | 実行手順を定義              |
| 読まれる場面 | そのディレクトリ配下で作業するとき   | そのスキルを呼び出したとき  |
| 置き場所     | リポジトリルートや各サブディレクトリ | 各スキルディレクトリ        |
| 例           | 「秘匿情報を含めない」               | 「Issueを作成してPRを作る」 |

普段の開発はClaude Codeをメインに、Codexをサブで使っています。

Claude CodeではSKILLをカスタムコマンドのように読み込ませる使い方もできますが、このリポジトリのSKILL.mdは、特定エージェントに依存しないプレーンな形で書くようにしています。

## Skills

| スキル                                          | 概要                                            |
| ----------------------------------------------- | ----------------------------------------------- |
| [backup-branch](./backup-branch/SKILL.md)       | autosquash や大きな履歴編集前の退避ブランチ作成 |
| [catch-up](./catch-up/SKILL.md)                 | main 取り込みと rebase 後の確認                 |
| [check](./check/SKILL.md)                       | 品質チェックの一連実行                          |
| [claude-review](./claude-review/SKILL.md)       | Claude CLIによるコードレビュー                  |
| [clean-docs](./clean-docs/SKILL.md)             | `.claude/docs` のタスクドキュメント整理         |
| [collect-feedback](./collect-feedback/SKILL.md) | 変更内容に対するフィードバック収集と整理        |
| [codex-review](./codex-review/SKILL.md)         | codex CLIによるコードレビュー                   |
| [commit](./commit/SKILL.md)                     | gitコミット（段階的コミット、fixup、amend）     |
| [doc-check](./doc-check/SKILL.md)               | ドキュメント整合性の確認                        |
| [doc-sync](./doc-sync/SKILL.md)                 | ドキュメント整合性の修正                        |
| [edit-taskdoc](./edit-taskdoc/SKILL.md)         | 元リポジトリ側タスクドキュメントの作成・更新    |
| [fixup](./fixup/SKILL.md)                       | 既存コミットへの fixup 追加                     |
| [gh-edit](./gh-edit/SKILL.md)                   | GitHub PR/Issueの作成・更新                     |
| [gh-read](./gh-read/SKILL.md)                   | GitHub Issue/PR の参照と要約                    |
| [issue-review](./issue-review/SKILL.md)         | Issue とコードベースの照合・有効性判定          |
| [link-skills](./link-skills/SKILL.md)           | Codex / Claude 向けスキルリンク作成             |
| [mark](./mark/SKILL.md)                         | チェック済み状態のタグ設置                      |
| [monthly-report](./monthly-report/SKILL.md)     | GitHub 活動データからの月次報告作成             |
| [pr-progress](./pr-progress/SKILL.md)           | PR 進捗コメントの整形・更新                     |
| [push](./push/SKILL.md)                         | push 前確認と push 実行                         |
| [read-taskdoc](./read-taskdoc/SKILL.md)         | 元リポジトリ側タスクドキュメントの参照          |
| [reply-review](./reply-review/SKILL.md)         | レビューコメントへの返信支援                    |
| [respond](./respond/SKILL.md)                   | 指摘対応から返信までのワークフロー              |
| [ship](./ship/SKILL.md)                         | check から PR 更新までの出荷フロー              |
| [start-dev](./start-dev/SKILL.md)               | 作業開始時のブランチ準備と情報収集              |
| [subagent-check](./subagent-check/SKILL.md)     | サブエージェント起動前の状態確認                |
| [takeover](./takeover/SKILL.md)                 | 前セッションのコンテキスト収集と作業引き継ぎ    |
| [tanaoroshi](./tanaoroshi/SKILL.md)             | 複数リポジトリの Issue/PR 棚卸し                |
| [watch-ci](./watch-ci/SKILL.md)                 | CI 状態の監視と失敗時の確認                     |
| [wiki-sync](./wiki-sync/SKILL.md)               | 開発内容から LLM Wiki への知識同期              |

補助スキルとして [daily-tagging](./.skill/daily-tagging/SKILL.md) も管理しています。

## スキル間の依存関係

各スキルが他のどのスキルへ処理を委譲しているかを示します。役割で色分けし、実線は本文で「実行」「委譲」と明示している呼び出し、点線は「注意」欄での委譲案内です。

```mermaid
flowchart LR
    classDef utility fill:#e6f3ff,stroke:#1a73e8;
    classDef read fill:#e9f7ef,stroke:#27ae60;
    classDef edit fill:#fef5e7,stroke:#d68910;
    classDef checkc fill:#fdedec,stroke:#c0392b;
    classDef history fill:#f5eef8,stroke:#8e44ad;
    classDef bootstrap fill:#fdf2e9,stroke:#e67e22;
    classDef workflow fill:#eaeded,stroke:#2c3e50,stroke-width:2px;
    classDef atomic fill:#ffffff,stroke:#7f8c8d;

    mark:::utility
    subagent-check:::utility
    backup-branch:::utility

    tanaoroshi:::atomic
    monthly-report:::atomic
    link-skills:::atomic
    watch-ci:::atomic
    reply-review:::atomic
    issue-review:::atomic
    collect-feedback:::atomic

    gh-read:::read
    read-taskdoc:::read

    commit:::edit
    fixup:::edit
    gh-edit:::edit
    edit-taskdoc:::edit
    doc-sync:::edit
    wiki-sync:::edit
    clean-docs:::edit

    check:::checkc
    doc-check:::checkc
    codex-review:::checkc
    claude-review:::checkc

    catch-up:::history
    push:::history
    pr-progress:::history

    start-dev:::bootstrap
    takeover:::bootstrap

    ship:::workflow
    respond:::workflow

    catch-up --> backup-branch
    catch-up --> pr-progress
    check --> mark
    check --> doc-check
    check --> claude-review
    check --> codex-review
    claude-review --> mark
    claude-review --> subagent-check
    codex-review --> mark
    codex-review --> subagent-check
    clean-docs --> read-taskdoc
    clean-docs --> wiki-sync
    commit --> check
    commit --> fixup
    commit --> catch-up
    commit --> backup-branch
    commit --> mark
    doc-check --> subagent-check
    doc-sync --> doc-check
    edit-taskdoc --> read-taskdoc
    fixup --> check
    fixup --> mark
    gh-edit --> subagent-check
    push --> check
    push --> pr-progress
    start-dev --> takeover
    start-dev --> gh-read
    start-dev --> read-taskdoc
    start-dev --> edit-taskdoc
    takeover --> gh-read
    takeover --> read-taskdoc
    wiki-sync --> subagent-check

    doc-check -.-> doc-sync
    doc-check -.-> wiki-sync
    edit-taskdoc -.-> clean-docs
    read-taskdoc -.-> edit-taskdoc
    read-taskdoc -.-> clean-docs
    gh-read -.-> gh-edit
    gh-read -.-> collect-feedback
    start-dev -.-> gh-edit
    wiki-sync -.-> gh-read
    wiki-sync -.-> read-taskdoc
    wiki-sync -.-> doc-sync

    ship --> commit
    ship --> check
    ship --> push
    ship --> gh-edit
    ship --> watch-ci
    ship --> pr-progress
    respond --> collect-feedback
    respond --> fixup
    respond --> ship
    respond --> reply-review
```

役割の凡例:

| 役割      | 説明                                   |
| --------- | -------------------------------------- |
| utility   | 他スキルから共通基盤として使われる     |
| read      | 情報を参照・取得する                   |
| edit      | ファイルや GitHub などの状態を変更する |
| check     | 検証・レビューを行い結果を返す         |
| history   | git 履歴・PR の進行に対する操作        |
| bootstrap | 作業開始・引き継ぎの準備               |
| workflow  | L0/L1 スキルを順に呼び出す複合フロー   |
| atomic    | 他スキルへ委譲しない独立した本体機能   |

## Setup

エージェントによっては、AGENTS.mdやSKILL.mdをリポジトリに置いただけでは、デフォルトで読んでくれないことがあります。

たとえば Claude Code では `CLAUDE.md` と `.claude/skills`、Codex では `AGENTS.md` と `~/.codex/skills` を使います。こういうときはリンクでつなぐのが手軽です。

このリポジトリで管理する補助スクリプトは `.tools` に集約しています。スキルから `mark.sh` や `tanaoroshi` などを使うため、`.tools` を PATH に追加してください。

```bash
export PATH="/path/to/dev-skills/.tools:$PATH"
```

### Windows で Codex を使う場合

`~/.codex/skills` 配下に、各スキルディレクトリへのジャンクションを作成します。

例:

```bash
mklink /J %USERPROFILE%\.codex\skills\commit C:\path\to\dev-skills\commit
```

### macOS / Linux で Claude Code を使う場合

リポジトリ全体を `.claude/skills` へリンクする運用ができます。

例:

```bash
ln -s /path/to/dev-skills ~/.claude/skills
```
