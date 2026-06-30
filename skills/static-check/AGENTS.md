# static-check

## 前提ツール

- [git](https://git-scm.com/)
- lint・build の実体ツールはリポジトリ依存（prettier, golangci-lint, gofumpt, cargo, make 等）

## 検出の判断基準

- ドキュメント（`AGENTS.md` / `CLAUDE.md` / `README.md` / CI 定義）に明示されたコマンドを最優先する。プロジェクトファイルからの推定はドキュメントで特定できなかった項目にだけ適用する
- lint と build は別項目として扱う。一方だけ検出できた場合は、検出できた項目を実行し、もう一方はスキップとして報告する

## スコープ判定の判断基準

- 変更ファイルの言語・種別に対応しない検査はスキップする。判定は推測で広げず、`git diff --name-only` の実体に基づく
- lint の対象はツールに依存する。Markdown も対象に含む linter（prettier 等）は、ドキュメントのみの変更でも実行する

## allowed-tools の範囲

- lint・build の実体ツールはリポジトリごとに多様なため、`allowed-tools` は make/npm/yarn/pnpm/cargo/go・golangci-lint・gofumpt といった代表的なラッパー・コマンドに限定する
- リポジトリが lint・build を直接バイナリ（`prettier` 等）で定義している場合、その binary は `allowed-tools` に含まれず実行時にユーザー承認を求めることがある。あらゆる linter バイナリを網羅列挙はしない

## 責務の境界

- このスキルは lint・build の実行と結果報告に限定する。mark タグの設置、履歴書き換え後のタグ引き継ぎ、doc-check、rule-check、review の実行判断や後続処理は担わない
