---
name: fixup
description: 行った変更を、対象コミットへ fixup として反映する。レビュー指摘の対応やコミットの修正漏れの追加に使う。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git diff:*) Bash(git add:*) Bash(git commit:*) Bash(git show:*) Bash(git rev-parse:*) Bash(git history fixup:*)
---

# fixup スキル

既に行った変更を、適切な既存コミットへ fixup として反映する。修正そのもの（ファイル編集）はこのスキルの前にエージェントが行う。

## フラグ

対象とする変更と反映方法を、`--staged` と `--history` の有無で指定する。

| フラグ      | 意味                                                                           |
| ----------- | ------------------------------------------------------------------------------ |
| (なし)      | fixup 対象ファイルの変更（unstaged・新規ファイルを含む）                       |
| `--staged`  | ステージ済みの変更を対象にする（呼び出し元が対象ファイルを事前にステージ済み） |
| `--history` | `git commit --fixup` を作らず、`git history fixup` で履歴を書き換える          |

デフォルトのfixup commitモードでは、fixupは対象ファイル単位でコミットする。フラグなしモードは対象ファイルを `git add` でステージし、pathspec 指定でコミットすることで、無関係なステージ済みファイルを巻き込まない。対象ファイル内にfixupと無関係な変更を含めないこと、`--staged` 利用時に対象以外をステージしないことは呼び出し元が担保する。

`--history` は experimental な履歴編集として扱う。明示指定がある場合だけ使い、指定がなければ従来どおり fixup commit を作る。

## 手順

### 1. フラグを確認する

`$ARGUMENTS` から `--staged` と `--history` の有無を確認し、対象とする変更の範囲と反映方法を決める。

- `--history` がある場合は履歴書き換えモードとして扱う
- `--history` がない場合は fixup commit モードとして扱う

### 1a. `--history` の前提を確認する

`--history` を使う場合は、次を確認する。

- `--staged` とは併用しない。`git history fixup` に対する index / pathspec の扱いをこのスキルではまだ標準化しない
- ブランチ内にマージコミットがある場合は `git history fixup` を使わず、fixup commit + autosquash に戻す
- `git history` は experimental であることを前提として扱う
- ref の更新範囲は既定値に任せず、`--update-refs=head` を明示する
- まず `--dry-run` で確認し、成功した場合だけ本実行する

前提を満たせない場合は、その時点で停止し、未実施の手順、阻害要因、不足前提を報告する。

### 2. 対象の変更を確認する

- `--staged` あり:

```bash
git diff --staged --name-only
```

ステージ済み変更がない場合はエラーで終了する。

- `--staged` なし:

```bash
git status -s
```

staged・unstaged を含む変更ファイルを把握し、その中から fixup の対象ファイルを特定する。変更がない場合はエラーで終了する。作業ツリーに fixup と無関係な変更が混在し対象を一意に絞れない場合は、どのファイルを対象とするかユーザーに確認する。

### 3. 対象コミットを特定する

```bash
git log --oneline main..HEAD
```

- 変更ファイルと各コミットの変更ファイルを照合して fixup 先のコミットを推定する
- 変更内容が明らかに特定のコミットに属する場合はそのコミットを対象にする
- 同じ対象コミットに対する既存の `fixup!` commit が既にある場合は、その fixup commit 自体ではなく元の対象コミットを fixup 先として扱う
- `--history` 指定時に、同じ対象コミットに対する fixup が既に複数ある、または今回の fixup 追加後に複数になる場合は、その時点で停止し、`git history fixup` を使わない理由と `/commit` のautosquash手順へ進む候補を報告する
- 対象が不明確な場合はユーザーに確認する

### 4. fixup を反映する

- fixup commit モード:
  - `--staged` なし: 対象ファイルを `git add` でステージし（新規ファイルもこれで追跡対象になる）、pathspec を指定してコミットする。pathspec により、無関係なステージ済みファイルは含まれない。

```bash
git add <fixup 対象ファイル>
git commit --fixup=<対象コミットのSHA> -- <fixup 対象ファイル>
```

- `--staged` あり: ステージ済みのインデックスをコミットする。

```bash
git commit --fixup=<対象コミットのSHA>
```

- `--history` モード:
  - まず dry-run で成否を確認する
  - 問題なければ本実行する
  - ref 更新は `HEAD` のみに限定する

```bash
git history fixup <対象コミットのSHA> --dry-run --update-refs=head
git history fixup <対象コミットのSHA> --update-refs=head
```

fixup commitモードでは、`fixup!` プレフィックスは常に 1 つだけにする。対象コミットが既に `fixup!` 付きでも積み重ねない。

### 5. 確認

```bash
git log --oneline -3
git status -s
```

`--history` を使った場合は、対象コミットが書き換わった結果としてハッシュが変わる前提で確認する。

## 注意

- 修正そのもの（ファイル編集）はこのスキルでは行わない。変更を済ませてから `/fixup` を呼ぶ
- `--history` は merge commit を含む履歴整理や複数 fixup の一括整理には使わず、その場合は `/commit` の autosquash 手順を使う
- 既にfixup commitがある対象へさらに修正を足す場合は、元コミットに対するfixup commitを追加し、整理は `/commit` のautosquash手順に寄せる
