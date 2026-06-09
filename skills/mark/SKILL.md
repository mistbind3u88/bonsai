---
name: mark
description: チェック通過を示すローカル軽量タグを現在の HEAD に設置・確認・削除する。
allowed-tools: Bash(mark.sh:*)
---

# mark スキル

チェック通過状態をローカル軽量タグで管理する。タグは `git push` では送信されない。

## タグ命名規則

```
mark/<branch>/<type>
```

ブランチ名を含めることで、複数ブランチの開発が干渉しない。

| タグ例                  | 意味                                 |
| ----------------------- | ------------------------------------ |
| `mark/feature-foo/lint` | feature-foo ブランチの lint 通過済み |
| `mark/main/test`        | main ブランチのテスト通過済み        |

## 手順

`mark.sh` を使って操作する。PATH 上に配置されている必要がある。

### タグを設置する

```bash
mark.sh <type>
```

### 状態を確認する

```bash
mark.sh --status
```

### 全タグを削除する

```bash
mark.sh --clean
```

## 注意

- タグはローカル専用。`git push` のデフォルトでは送信されない
- コミットが進むとタグは古いコミットに残るため、再チェック後に再設置が必要
- 以下の処理系が関連処理の成功後に `/mark` を呼び出してタグを設置する:
  - 品質チェック系の呼び出し元: lint/build/test/doc-check/rule-check タグ
  - sub-review 系の呼び出し元: review-sub タグ
  - cross-review 系の呼び出し元: review-cross タグ
