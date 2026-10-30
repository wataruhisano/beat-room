# Beatroom

BeatmakerがArtistごとにprivateなBeat pack（Room）をsecret URLで送り、反応をArtist単位で蓄積するWebアプリのMVP。

## Docs

- [docs/PRD.md](docs/PRD.md) — 要件と未確定事項
- [docs/MVP_SCOPE.md](docs/MVP_SCOPE.md) — scope、acceptance criteria、release gate
- [docs/USER_FLOWS.md](docs/USER_FLOWS.md) — 画面とflow
- [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) — 暫定architecture（未確定）
- [docs/DATA_MODEL.md](docs/DATA_MODEL.md) — schema案とinvariants
- [docs/IMPLEMENTATION_PLAN.md](docs/IMPLEMENTATION_PLAN.md) — Phase計画

## 進め方

- 一度に全Phaseを実装しない。各Phase開始時に、IMPLEMENTATION_PLANの「開始時に確認」を質問してから着手する
- 未確定事項を推測で埋めない。回答がなければ可逆な仮決定をADRに残す
- MVP_SCOPEのout of scopeに入る機能を追加しない
- tenant isolation、public token access、履歴保持（archive、物理削除しない）を後付けにしない
- Phase終了時にdocs、README、未確定事項を更新する
