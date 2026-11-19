# Beatroom

BeatmakerがArtistごとにprivateなBeat pack（Room）を作り、secret URLで共有して、反応をArtist単位で蓄積するWebアプリです。現在はMVPを開発中です。

Artistは登録不要でRoomを開いて試聴と反応ができ、Beatmakerは「誰に、いつ、何を送り、どう反応されたか」を後から確認できます。

## 構成

| 領域 | 技術 |
|---|---|
| Web | Next.js（App Router）、TypeScript、Tailwind CSS |
| API | Go（`net/http`）、`log/slog` |
| DB | PostgreSQL、golang-migrate |
| 本番（予定） | Vercel、Cloud Run、Neon、Cloudflare R2 |
| CI | GitHub Actions、Dependabot |

Webは表示を担当し、業務logicと認可はGo APIに置きます。判断の理由は[ADR一覧](docs/adr/README.md)にあります。

```text
apps/
  api/   Go API
  web/   Next.js
docs/    要件、設計、ADR
```

## Local開発

### 必要なもの

| Tool | Version |
|---|---|
| Go | 1.27以上 |
| Node.js | 26（`.nvmrc`） |
| pnpm | 12 |
| Docker Engine + Compose | Compose v2以上 |
| GNU Make | 4以上 |
| golangci-lint | v2.15.0（`make check`で使用） |

### 初回setup

```bash
make setup
```

`.env`と`apps/web/.env.local`をexampleからcopyし、Webの依存をinstallします。既存の`.env`は上書きしません。

### 起動

```bash
make db-up
make migrate
make api
```

別のterminalで次を実行します。

```bash
make web
```

http://localhost:3000 を開き、`API: ok`と表示されれば起動できています。

| URL | 内容 |
|---|---|
| http://localhost:3000 | Web |
| http://localhost:8080/healthz | APIのprocess確認 |
| http://localhost:8080/readyz | APIからDBへの接続確認 |

### Command

| Command | 内容 |
|---|---|
| `make help` | command一覧 |
| `make check` | API、Webのformat、lint、型チェック、test |
| `make db-down` | PostgreSQLを停止する。dataは残る |

PostgreSQLのdataを消して作り直す場合は、`docker compose down -v`の後に`make db-up`と`make migrate`を実行します。

## Docs

- [PRD](docs/PRD.md)：要件と未確定事項
- [MVP scope](docs/MVP_SCOPE.md)：scopeとrelease条件
- [User flows](docs/USER_FLOWS.md)：画面とflow
- [Architecture](docs/ARCHITECTURE.md)
- [Data model](docs/DATA_MODEL.md)
- [Implementation plan](docs/IMPLEMENTATION_PLAN.md)：Phase計画と進捗
- [ADR](docs/adr/README.md)：技術判断の記録
