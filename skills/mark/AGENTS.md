# mark

## 前提ツール

- [git](https://git-scm.com/)

## 責務の境界

- このスキルは、現在の HEAD に対する品質チェック通過状態を軽量タグとして設置・確認・削除する責務を持つ。
- lint、build、test、doc-check、rule-check、review の実行判断と実行順序は、このスキルの責務外とする。
- 個別の検査実行、コミット、push は行わず、結果が揃った後のタグ操作だけを担う。
