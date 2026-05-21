---
name: link-skills
description: リポジトリの skills/ ディレクトリを Codex または Claude Code から使えるようにリンクする。Codex では ~/.codex/skills に各スキルディレクトリへのジャンクションを作成し、Claude Code では ~/.claude/skills へ skills/ をリンクする。
allowed-tools: Bash(ls:*) Bash(ln:*) Bash(readlink:*) Bash(find:*)
---

# link-skills

このリポジトリの公開スキル（`skills/` 配下）を、Codex または Claude Code から使えるようにリンクする。

## 手順

1. 対象エージェントと OS を確認する

- Codex on Windows の場合: `~/.codex/skills` 配下に各スキルディレクトリへのジャンクションを作成する
- Claude Code on macOS / Linux の場合: `~/.claude/skills` へ `skills/` ディレクトリをリンクする

2. リポジトリ内の公開スキルを確認する

```bash
find skills -name SKILL.md
```

公開スキルは `skills/` 配下にある。`archive/`（退役スキル）と `internal/`（リポジトリ専用の保守スキル）はリンク対象に含めない。

3. Codex on Windows の場合は `~/.codex/skills` の現在の状態を確認する

```bash
ls ~/.codex/skills
```

- 同名エントリが既にある場合はリンク先を確認する
- 想定外の既存ディレクトリやファイルは上書きしない

4. Codex on Windows の場合は各スキルディレクトリへのジャンクションを作成する

```bash
mklink /J %USERPROFILE%\.codex\skills\<skill-name> C:\path\to\dev-skills\skills\<skill-name>
```

- 既に正しいジャンクションがある場合はそのままにする

5. Claude Code on macOS / Linux の場合は `~/.claude/skills` の現在の状態を確認する

```bash
ls -la ~/.claude/skills 2>/dev/null
readlink ~/.claude/skills 2>/dev/null
```

- すでに `skills/` への正しいリンクが存在する場合はそのまま終了する
- `~/.claude/skills` がシンボリックリンクでないディレクトリとして存在する場合は、上書きせず警告を出して終了する

6. Claude Code on macOS / Linux の場合は `~/.claude/skills` へリンクを作成する

```bash
ln -s /path/to/dev-skills/skills ~/.claude/skills
```

7. 結果を確認して報告する

Codex on Windows:

```bash
dir %USERPROFILE%\.codex\skills
```

Claude Code on macOS / Linux:

```bash
ls -la ~/.claude/skills
```
