---
name: link-agentdoc
description: AGENTS.md / CLAUDE.md に、恒久的に参照すべきドキュメントへの link、delegate、ref を追記する。
allowed-tools: Read Edit Bash(pwd:*) Bash(ls:*) Bash(find:*) Bash(readlink:*) Bash(git rev-parse:*) Bash(git status:*) Bash(git diff:*)
---

# link-agentdoc

`AGENTS.md` / `CLAUDE.md` に、実行者が参照すべきドキュメントへの導線を、役割に応じた形式で追加する。

## 分類オプション

```text
/link-agentdoc link <doc-path> [--target <path-to-AGENTS-or-CLAUDE>]
/link-agentdoc delegate <path-or-pattern> [--target <path-to-AGENTS-or-CLAUDE>]
/link-agentdoc ref <path-or-pattern> [--target <path-to-AGENTS-or-CLAUDE>]
```

| 分類オプション | 追記形式                           | 用途                                                                 |
| -------------- | ---------------------------------- | -------------------------------------------------------------------- |
| `link`         | 見出し + `@...` link               | 参照先を `AGENTS.md` / `CLAUDE.md` のセクション相当として扱う        |
| `delegate`     | 冒頭の役割分担説明 + Markdown link | この文書と委譲先ドキュメントで記載内容を分けて詳細を伝える           |
| `ref`          | 末尾の参考一覧 + Markdown link     | 単に参考として読む探索起点、ディレクトリ、glob、候補群をまとめて示す |

分類判断の詳細は、同じスキルディレクトリの `AGENTS.md` を読む。

## 手順

### 1. 対象ファイルを決める

`--target` が指定されていれば、その path が指す `AGENTS.md` または `CLAUDE.md` を対象にする。

指定がなければ、repo root 配下の `AGENTS.md` / `CLAUDE.md` を列挙し、現在の作業ディレクトリまたは対象スキル・モジュールとの位置関係から、適用範囲が最も近い `AGENTS.md` / `CLAUDE.md` を候補にする。複数候補が同じ範囲にある場合は、次を確認する。

- `CLAUDE.md` が `AGENTS.md` への symlink なら、実体である `AGENTS.md` を編集する
- `AGENTS.md` と `CLAUDE.md` が独立した通常ファイルなら、どちらが対象かをユーザーに確認する

確認には次のコマンドを使う。

```bash
pwd
ls AGENTS.md CLAUDE.md
git rev-parse --show-toplevel
find <repo-root> \( -path "*/.git" -o -path "*/node_modules" \) -prune -o \( -name AGENTS.md -o -name CLAUDE.md \) -print
readlink CLAUDE.md
```

`<repo-root>` には、直前の `git rev-parse --show-toplevel` の出力をリテラル値として入れる。

### 2. 参照先を確認する

`link` の場合は、対象リポジトリ内の単一ファイルとして存在することを確認する。リポジトリ内ファイルは対象 `AGENTS.md` / `CLAUDE.md` からの相対パスで書く。参照先が symlink の場合は `readlink` で実体を確認し、恒久導線としては実体 path を優先する。`@...` を自動展開する環境では参照先を直接読ませ、自動展開しない環境では見出しと `@...` の参照先から読むべきファイルを判断できる導線として扱う。参照先が URL の場合は、分類オプションを `ref` に切り替え、以降は `ref` の追記位置と形式で進める。参照先が対象ファイルへ戻る導線を持つ場合は、相互参照・循環参照になるため、本文へ要点を書くか `delegate` / `ref` に切り替える。

`delegate` / `ref` の場合は、ディレクトリ、glob、候補群、探索起点として意味が通る表記にする。存在確認できるものは確認し、glob や将来作成されるパターンは参照理由で用途を明確にする。`delegate` で委譲先から対象 `AGENTS.md` / `CLAUDE.md` へ戻る導線を置く場合は、相互参照の役割を明記し、同じ判断を往復して読ませる循環参照を作らない。

### 3. 追記位置を選ぶ

`link` は、リンク先をセクション本文の代わりとして扱う。対象ファイルの既存構成に合わせ、同じ粒度の既存セクションと同じ見出しレベルを選び、その直下に `@...` link を置く。独立した主要判断として置く場合は `##`、既存セクション配下の補足判断として置く場合はその配下の見出しレベルを使う。

`delegate` は、対象 `AGENTS.md` / `CLAUDE.md` の導入セクション直下に `## 役割分担` を新設するか、既存の前提セクションへ追記する。本文で担う内容と委譲先ドキュメントに置く詳細を明記し、作業者が読む順序と責務分担を判断できるようにする。

`ref` は、対象ファイルの末尾付近に `## 参考` または既存の参考セクションを置き、そこへまとめて追記する。

### 4. 重複を確認する

対象ファイル本文を `Read` で読み、既存導線を参照先単位で確認する。分類が `link` / `delegate` / `ref` のどれであっても、同じドキュメント・URL・探索起点を指している場合は同一導線として扱う。

重複確認では、次の表記を同じ参照先として正規化する。

- `@path` と Markdown link の `path`
- 対象 `AGENTS.md` / `CLAUDE.md` からの相対パスと repo root からの相対パス
- symlink path と `readlink` で確認した実体 path
- 末尾 slash の有無だけが違うディレクトリ path
- URL の末尾 slash や fragment の有無だけが違い、同じ文書本体を指すもの

既存導線がある場合は、新規追加ではなく次のいずれかで対応する。

- 既存分類が現在の用途に合う場合は、説明文・見出し・参照理由を補足する
- 現在の用途により近い分類へ移す場合は、既存導線を移動・書き換え、古い位置の同じ参照先を削除する
- `link` と `ref` のように同じ参照先が複数分類に分かれている場合は、読み手が最初に判断すべき 1 つの分類へ統合する
- 参照先が広い探索起点と個別ドキュメントの親子関係にある場合は、個別ドキュメントを読む必要が恒常的にあるときだけ `link` または `delegate` として残し、探索起点側の `ref` には役割が重複しない説明を付ける

統合後は、対象ファイル内で同じ正規化参照先が 1 箇所だけに残っていることを確認する。

### 5. 導線を追記する

`link` は、見出しと `@...` link で書く。リンク先の内容をセクション本文として参照させるため、箇条書きの一項目ではなく見出し単位で置く。

```markdown
## <参照先が担う判断・手順の見出し>

@<relative-path>
```

既存セクション配下へ置く場合は、挿入位置に合わせて `###` など同じ階層の見出しレベルを使う。

`delegate` は、対象ファイルと委譲先の役割分担を書いた上で、委譲先を Markdown link で書く。表示名は、委譲先が担う詳細の種類が分かる名前にする。

```markdown
## 役割分担

この AGENTS.md / CLAUDE.md には <この文書が担う内容> を記載する。詳細な <判断基準・手順・対象一覧> は [<表示名>](path-or-url) に記載する。
```

`ref` は、末尾の参考セクションに Markdown link を書く。表示名は、探す対象や使う場面が分かる名前にする。

```markdown
## 参考

- [<表示名>](path-or-url) — <探す対象・使う場面>
```

表示名や理由は、実行者が「いつ読むか」「何を判断するために読むか」を理解できる粒度にする。

### 6. 差分を確認する

```bash
git diff -- <target-agent-doc>
git status -s
```

## 注意

- 参照先本文の作成・大幅更新はこのスキルで行わず、必要な次のスキルとして `/doc-sync` または `/wiki-sync` を案内する
- タスクドキュメントの配置先確認が必要な場合は、次に使うスキルとして `/taskdoc-locate` を案内する
- リポジトリルールへの適合確認が必要な場合は、次に使うスキルとして `/rule-check` を案内する
