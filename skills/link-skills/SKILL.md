---
name: link-skills
description: 公開スキルをCodexまたはClaude Codeへ登録するためのリンクを作成する。利用先と既存リンクを確認してセットアップするときに使う。
allowed-tools: Bash(ls:*) Bash(ln:*) Bash(readlink:*) Bash(find:*) Bash(mkdir -p:*) Bash(cmd /c mklink /J:*) Bash(cmd /c dir:*)
---

# link-skills

このリポジトリの公開スキル（`skills/` 配下）を、Codex または Claude Code から使えるようにリンクする。

## 手順

1. 対象エージェントと OS を確認する

- Codex on Windows の場合: `~/.agents/skills` 配下に各スキルディレクトリへのジャンクションを作成する
- Claude Code on macOS / Linux の場合: `~/.claude/skills` へ `skills/` ディレクトリをリンクする

Codexのユーザー用配置先は `~/.agents/skills`、リポジトリ用配置先は `.agents/skills`。この手順はユーザー用の登録を扱う。配置先は[公式仕様](https://learn.chatgpt.com/docs/build-skills)で確認する。対象外のOSやリポジトリ用の登録を求められた場合は、対象環境と配置先を示して手順の確認を求める。

2. リポジトリ内の公開スキルを確認する

```bash
find skills -name SKILL.md
```

公開スキルは `skills/` 配下にある。`archive/`（退役スキル）と `internal/`（リポジトリ専用の保守スキル）はリンク対象に含めない。

3. Codex on Windows の場合は `~/.agents/skills` の現在の状態を確認する

```bash
ls ~/.agents/skills
```

- 同名エントリが既にある場合はリンク先を確認する
- 想定外の既存ディレクトリやファイルは上書きしない
- 従来の `~/.codex/skills` や対象リポジトリの `.agents/skills` に同じスキルがある場合は、リンク先と利用中の配置を確認する。同名スキルは自動統合されないため、登録を重複させる前にユーザーへ配置の選択を確認する。既存の登録はこの手順で削除しない
- 登録先ディレクトリが存在しない場合は、対象パスを確認して作成する。作成に必要な権限がない場合は、不足権限を示して停止する

4. Codex on Windows の場合は各スキルディレクトリへのジャンクションを作成する

次はWindowsのコマンドプロンプト用の例。PowerShellから実行する場合は `cmd /c mklink /J` で呼び出す。

```bat
mklink /J "%USERPROFILE%\.agents\skills\<skill-name>" "C:\path\to\bonsai\skills\<skill-name>"
```

- 既に正しいジャンクションがある場合は作成だけを省略し、手順7で結果を確認する

5. Claude Code on macOS / Linux の場合は `~/.claude/skills` の現在の状態を確認する

```bash
ls -la ~/.claude/skills 2>/dev/null
readlink ~/.claude/skills 2>/dev/null
```

- すでに `skills/` への正しいリンクが存在する場合は手順6の作成を省略し、手順7で結果を確認する
- `~/.claude/skills` がシンボリックリンクでないディレクトリとして存在する場合は、上書きせず警告を出して終了する

6. Claude Code on macOS / Linux の場合は `~/.claude/skills` へリンクを作成する

```bash
ln -s /path/to/bonsai/skills ~/.claude/skills
```

7. 結果を確認して報告する

Codex on Windows:

```bat
dir "%USERPROFILE%\.agents\skills"
```

Claude Code on macOS / Linux:

```bash
ls -la ~/.claude/skills
```

リンク先が意図した公開スキルを指すことと、対象エージェントのスキル一覧に登録が現れることを確認する。一覧を確認できない場合は、リンク作成の結果と読込未確認を分けて報告し、登録完了とは扱わない。
