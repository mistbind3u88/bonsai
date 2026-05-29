# wiki-sync

## 前提ツール

- [git](https://git-scm.com/)

## 責務の境界

- 対象リポジトリの wiki schema を最優先し、配置先・命名・ページ構成は対象リポジトリから読み取る。
- wiki には、将来の LLM が判断に使う設計判断・契約・運用知識を残す。
- README、AGENTS.md、Issue/PR、作業メモの内容は、wiki として読む価値のある知識へ蒸留して反映する。
- コード近傍に置くべき局所的な不変条件は、近い階層のドキュメントやコードコメントへの反映を検討する。
- LLM Wiki が定義されているリポジトリで実行し、定義が見つからない場合は初期化タスクとしてユーザーへ確認する。

## サブエージェント依頼

更新後確認のサブエージェント依頼は、役割ごとのテンプレートを埋めて作成する。

- 想定質問への試行回答: [answer-check-prompt-template.md](answer-check-prompt-template.md)
- schema・体裁・関連ドキュメント整合確認: [presentation-check-prompt-template.md](presentation-check-prompt-template.md)
- 追加観点確認: [additional-check-prompt-template.md](additional-check-prompt-template.md)
