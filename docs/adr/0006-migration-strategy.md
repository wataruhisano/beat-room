# ADR 0006: forward-only migrationを採用する

- Status: Accepted
- Date: 2026-10-30

## 背景

column削除やdata変換は、down migrationを実行しても失われたdataを戻せない。downを必須にしてもtestが伴わなければ、戻せるという誤った安心につながる。

## 決定

- 本番のmigrationは前進のみとする。問題が起きた場合は修正migrationを追加する
- 破壊的な変更はexpand / contract（追加、移行、後日削除）で段階的に行う
- 本番の復旧はbackupからの復元と修正migrationで行う
- down fileはlocal開発用に任意で書いてよい
- CIは次を確認する
  1. 空のDBへ全migrationを適用できる
  2. sqlcの生成codeとschemaが一致している

## 検討した代替案

- **downを必須にし、CIでup → down → upを確認する**: 現実の障害時に使えない場面が多く、維持費に対する効果が小さい
- **forward-onlyに加え、書かれたdownのみCIで確認する**: 有効だが、当面は必須にしない

## 影響

- IMPLEMENTATION_PLANのPhase 0 Doneを、rollback方針の記録へ変更する
- Phase 5のrollback rehearsalは、backupからの復元訓練として実施する
- **注意**: Neon Freeの履歴保持は6時間で、release gateのbackup条件には足りない。validation開始前に有料planまたは定期`pg_dump`を用意する
