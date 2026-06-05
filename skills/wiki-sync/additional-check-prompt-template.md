# wiki-sync 追加観点確認テンプレート

ユーザーから追加の品質確認が求められた場合に、不足・曖昧さ・追加調査対象を列挙するサブエージェントへ渡す。

前回の結論、期待する答え、メイン側の解釈は、サブエージェントの回答後にメインエージェントが検証材料として使う。

## 境界確認

- 実行中のスキル: `/wiki-sync`
- work_unit: `wiki-sync:<wiki-entry-or-topic>`
- role: `wiki-additional-check`
- reuse_policy: `prefer-reuse`
- 依頼境界: メインエージェントが `/wiki-sync` の追加観点確認を実行中であり、サブエージェントの役割は指定された問いに対する不足・曖昧さ・追加調査対象を列挙すること
- 返すべき結果: 根拠付きの不足・曖昧さ・追加調査対象の一覧と、参照したファイルパス・URL・セクション
- メイン側に残す判断: 調査結果の採否、wiki への反映、ユーザーへの確認事項
- サブエージェントが行うこと: 対象 wiki、関連コード、明示参照、必要な外部調査条件、ユーザーが確認したい問いを読み、不足・曖昧さ・追加調査対象を返す
- サブエージェントがメイン側に残すこと: スキル実行、wiki 編集、commit、push、ユーザーへの最終判断

## 対象

- リポジトリ: `<absolute-repo-path>`
- 対象 wiki: `<absolute-wiki-page-path>`
- 関連コード: `<absolute-related-code-paths-or-none>`
- 関連ドキュメント: `<absolute-related-doc-paths-or-none>`
- 明示的に参照すべき文書・URL:
  - `<reference-path-url-or-none> — <why-this-reference-matters>`
- 外部調査条件: `<external-research-scope-or-none>`

## ユーザーが確認したい問い

- `<question-or-review-point>`

## 出力形式

以下を根拠付きで簡潔に返してください。

- 不足している前提
- 曖昧な判断基準
- 追加で確認すべきコード、ドキュメント、外部情報
- 参照したファイルパス・URL・セクション
- wiki または近傍ドキュメントへ反映すると判断しやすくなる情報
