# respond

## 前提ツール

- [git](https://git-scm.com/)
- [gh](https://cli.github.com/) — `brew install gh`

## 責務の境界

- respond はレビュー指摘対応ワークフローのオーケストレーションを担い、要対応の指摘に対するコード修正を手順内で行う
- 指摘の収集は `/collect-feedback`、修正の記録は `/fixup`、push は `/ship`、リプライは `/reply-review` へ委譲する
- 実行順序と各ステップの詳細は SKILL.md に従う
