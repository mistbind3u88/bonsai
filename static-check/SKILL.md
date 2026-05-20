---
name: static-check
description: リポジトリの lint・build を検出して実行し、結果を報告する。コードを実行しない静的な検査をまとめて担う。
allowed-tools: Bash(make lint:*) Bash(make build:*) Bash(make -n:*) Bash(npm run:*) Bash(yarn run:*) Bash(pnpm run:*) Bash(cargo build:*) Bash(cargo clippy:*) Bash(go build:*) Bash(golangci-lint run:*) Bash(gofumpt:*) Bash(git diff:*) Read
---

# static-check スキル

リポジトリに合わせて lint・build を実行し、結果を項目別に報告する。コードを実行しない静的な検査（リント、コンパイル）をまとめて担う。タグの設置は行わない（呼び出し元の責務）。

## 手順

### 1. 実行コマンドを検出する

以下の優先順位で lint・build のコマンドを特定する。ドキュメントに記載があればそれを使い、自動検出には進まない。

#### 優先: リポジトリのドキュメント

リポジトリの `AGENTS.md`、`CLAUDE.md`、`README.md`、CI 定義、またはスキル定義に lint・build の実行コマンドが記載されていればそれを使う。記載で特定できた項目はフォールバックに進まない。

#### フォールバック: プロジェクトルートのファイルから自動検出

ドキュメントに記載がない項目について、プロジェクトルートのファイルから推定する。

| ファイル       | 検出方法                                           |
| -------------- | -------------------------------------------------- |
| `Makefile`     | `make` のターゲット一覧から `lint`・`build` を探す |
| `package.json` | `scripts` フィールドから `lint`・`build` を探す    |
| `Cargo.toml`   | `cargo clippy`・`cargo build` を使う               |
| `go.mod`       | `golangci-lint run`・`go build ./...` を使う       |

ドキュメントにもプロジェクトファイルにも明示がない項目は言語標準のコマンドで実行する（例: Go なら `go build ./...`）。言語標準のコマンドも特定できない項目はスキップとして扱う。

### 2. 変更範囲に応じてスコープを判定する

スコープ判定は無関係な言語の検査を省く最適化である。変更ファイルの基点を決められる場合に行い、決められない場合は省略して検出した lint・build をすべて実行する。

ブランチに上流（`@{upstream}`）が設定されていれば、それを基点に変更ファイルを確認する。

```bash
git diff --name-only @{upstream}...HEAD
```

変更が特定の言語やディレクトリに閉じている場合、無関係な言語の lint・build はスキップする。上流が未設定で基点を決められない場合は、スコープ判定を省いて検出した検査をすべて実行する。

例:

- ドキュメント（Markdown）のみの変更 → Go の build はスキップ。Markdown を対象とする lint は実行する
- Go のみの変更 → Python 関連の検査はスキップ

変更対象の言語に対応する検査ツールがプロジェクトに存在しない場合はスキップとして扱う。

### 3. lint・build を実行する

検出したコマンドを実行する。lint → build の順に実行し、いずれかが失敗したら停止する。

### 4. 結果を報告する

lint・build それぞれを `PASS` / `FAIL` / `SKIP` で報告する。`FAIL` の場合は失敗の要点を添える。

```
static-check 結果:
  lint:  PASS
  build: SKIP（Go の変更なし）
```

## 注意

- このスキルは検出と実行、結果報告のみを行う。mark タグの設置は呼び出し元（`/check` 等）が行う
- 検出できなかった項目は「スキップ」として扱い、ブロッカーにしない
