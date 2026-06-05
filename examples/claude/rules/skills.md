## Skills

- セッション開始時に利用可能な skill 群を把握し、その後の作業内容に合致する skill があれば積極的に活用する
- ある skill を実行する上での前提条件が満たされていない場合、その前提条件を満たせる別の skill が存在するなら、先にその skill を実行して前提を満たしてから元の skill を実行する。連鎖する前提も同様に遡って解決する（例: push の前提である check 未通過を check の実行で満たす）
- SKILL.mdを記述する時は、各種エージェント依存の記述ではなくSKILL.mdの公式仕様への準拠を心がける
  - https://agentskills.io/specification
- skillを追加・修正する時は、allowed-toolsで許可するコマンドを必要最小限にする
- skill配下に実装しているスクリプトやツールは `allowed-tools` に含めない。Claude Code のパーミッションマッチングは argv[0] が変数展開で決まる形式（例: `Bash($SKILL_DIR/foo.sh:*)`）を runtime-determined として毎回承認を要求するため、スキル相対パスの指定は機能しない。該当スクリプトは PATH 上に配置した上で bare name で呼び出し、許可は settings.json の `permissions.allow` 側で管理する
