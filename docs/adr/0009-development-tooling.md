# ADR 0009: 開発toolのdefault

- Status: Provisional
- Date: 2026-10-30

## 背景

Phase 0の基盤に必要な項目のうち、ADR 0005の基準から自然に決まるものをまとめて記録する。いずれも後で変更できる。

## 決定

| 項目 | Default |
|---|---|
| Go format / lint | `gofmt`、golangci-lint |
| Go logging | `log/slog`のJSON出力とrequest ID |
| Web | Next.js App Router、TypeScript strict |
| Web lint / format / test | ESLint、Prettier、Vitest |
| E2E | 必要になるPhaseでPlaywrightを導入する |
| Local環境 | PostgreSQLとMinIOはDocker Compose、WebとAPIは手元で直接起動する |
| Command | rootの`Makefile` |
| Environment | `.env.example`をcommitする。本番の秘密情報は各serviceのsecret機能に置く |
| CI | GitHub Actions |
| 依存更新 | Dependabot |

## 検討した代替案

- **WebとAPIもcontainerで起動する**: 環境は揃うが、修正の反映とdebugが遅くなる。必要になれば追加する

## 影響

- ARCHITECTUREのlocal環境の記述をこの内容へ更新する
- 実際に使って問題があれば、Phaseの途中でも見直す
