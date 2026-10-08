## Skills

- セッション開始時に利用可能な skill 群を把握し、その後の作業内容に合致する skill があれば積極的に活用する
- ある skill を実行する上での前提条件が満たされていない場合、その前提条件を満たせる別の skill が存在するなら、先にその skill を実行して前提を満たしてから元の skill を実行する。連鎖する前提も同様に遡って解決する（例: push の前提である check 未通過を check の実行で満たす）
- SKILL.mdを記述する時は、各種エージェント依存の記述ではなくSKILL.mdの公式仕様への準拠を心がける
  - https://agentskills.io/specification
- skillを追加・修正する時は、allowed-toolsで許可するコマンドを必要最小限にする
- skill配下に実装しているスクリプトやツールは、このスキル群の共通運用としてPATH上に配置し、bare nameで呼び出す。恒久的な許可はsettings.jsonの `permissions.allow` 側で管理し、スキルの `allowed-tools` と実際の適用範囲を確認する。現在のClaude Codeでは `${CLAUDE_SKILL_DIR}` を本文と `allowed-tools` で展開して内包スクリプトを許可する[方式](https://code.claude.com/docs/en/skills)もある。通常のshell変数 `$SKILL_DIR` とは区別し、PATH方式を技術上の唯一の選択肢とは扱わない
