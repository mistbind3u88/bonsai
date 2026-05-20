# edit-taskdoc

- タスクドキュメントの作成・更新先は元リポジトリ側の `.claude/docs` に限定する。
- worktree 側の `.claude/docs` にはファイルを作成・更新しない。
- 編集前に `/read-taskdoc` で既存ドキュメントを確認し、重複作成を避ける。
- 削除や恒久ドキュメントへの移管は `/clean-docs` に委譲する。
