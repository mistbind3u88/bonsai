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
| [static-check](./static-check/SKILL.md)         | リポジトリの lint・build の検出と実行           |
| [subagent-check](./subagent-check/SKILL.md)     | サブエージェント起動前の状態確認                |
| [takeover](./takeover/SKILL.md)                 | 前セッションのコンテキスト収集と作業引き継ぎ    |
| [tanaoroshi](./tanaoroshi/SKILL.md)             | 複数リポジトリの Issue/PR 棚卸し                |
| [unit-test](./unit-test/SKILL.md)               | リポジトリのユニットテストの検出と実行          |
| [watch-ci](./watch-ci/SKILL.md)                 | CI 状態の監視と失敗時の確認                     |
| [wiki-sync](./wiki-sync/SKILL.md)               | 開発内容から LLM Wiki への知識同期              |

補助スキルとして [daily-tagging](./.skill/daily-tagging/SKILL.md) も管理しています。

## スキル間の依存関係

各スキルが他のどのスキルへ処理を委譲しているかを示します。層で色分けし、実線は通常フローで実行される委譲、点線は条件付きの委譲または「注意」欄での案内です。

```mermaid
flowchart LR
    classDef l0 fill:#e9f7ef,stroke:#27ae60;
    classDef l1 fill:#fef5e7,stroke:#d68910;
    classDef l2 fill:#eaeded,stroke:#2c3e50,stroke-width:2px;

    subgraph L2G["L2 ワークフロースキル"]
        ship:::l2
        respond:::l2
    end

    subgraph L1G["L1 サービススキル"]
        check:::l1
        commit:::l1
        push:::l1
        catch-up:::l1
        doc-check:::l1
        doc-sync:::l1
        gh-edit:::l1
        clean-docs:::l1
        edit-taskdoc:::l1
        wiki-sync:::l1
        codex-review:::l1
        claude-review:::l1
        start-dev:::l1
        takeover:::l1
    end

    subgraph L0G["L0 単機能スキル"]
        mark:::l0
        subagent-check:::l0
        backup-branch:::l0
        gh-read:::l0
        read-taskdoc:::l0
        pr-progress:::l0
        fixup:::l0
        static-check:::l0
        unit-test:::l0
        tanaoroshi:::l0
        monthly-report:::l0
        link-skills:::l0
        watch-ci:::l0
        reply-review:::l0
        issue-review:::l0
        collect-feedback:::l0
    end

    catch-up --> backup-branch
    catch-up --> pr-progress
    check --> mark
    check --> static-check
    check --> unit-test
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
    gh-edit --> subagent-check
    push --> check
    start-dev --> takeover
    start-dev --> gh-read
    start-dev --> read-taskdoc
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
    start-dev -.-> edit-taskdoc
    push -.-> pr-progress
    ship -.-> pr-progress
    wiki-sync -.-> gh-read
    wiki-sync -.-> read-taskdoc
    wiki-sync -.-> doc-sync

    ship --> commit
    ship --> check
    ship --> push
    ship --> gh-edit
    ship --> watch-ci
    respond --> collect-feedback
    respond --> fixup
    respond --> ship
    respond --> reply-review
```

層の凡例:

| 層                    | 説明                                               |
| --------------------- | -------------------------------------------------- |
| L0 単機能スキル       | 他スキルへ委譲せず、自身の操作だけで責務を完結する |
| L1 サービススキル     | 固有の操作を持ち、責務の範囲内で他スキルを活用する |
| L2 ワークフロースキル | 操作を持たず、複数のスキルを決まった順序で呼び出す |

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
