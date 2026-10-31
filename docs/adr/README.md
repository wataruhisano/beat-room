# Architecture Decision Records

重要な判断を、理由と捨てた選択肢とともに残す。新しいADRは`0000-template.md`をcopyして連番を振る。
一度Acceptedにした内容は書き換えず、変更時は新しいADRを作って古いADRをSupersededにする。

Status: **Accepted** = 確定 / **Provisional** = 可逆な仮決定で、後のPhaseで見直してよい

| ADR | 判断 | Status |
|---|---|---|
| [0001](0001-web-api-monorepo.md) | Next.js + Go API + PostgreSQLのmonorepo | Accepted |
| [0002](0002-deployment-topology.md) | Vercel / Cloud Run / Neon / R2とcost guardrail | Accepted（regionのみProvisional） |
| [0003](0003-object-storage.md) | S3互換APIとstorage key | Accepted（MinIOのみADR 0010でSuperseded） |
| [0004](0004-authentication.md) | Firebase Authenticationと自前users table | Accepted |
| [0005](0005-dependency-policy.md) | dependency選定基準と主要tool | Accepted |
| [0006](0006-migration-strategy.md) | forward-only migration | Accepted |
| [0007](0007-health-checks.md) | `/healthz`と`/readyz`の分担 | Accepted |
| [0008](0008-event-collection.md) | 将来の推薦に向けたevent収集の範囲 | Accepted |
| [0009](0009-development-tooling.md) | 開発toolのdefault | Provisional（MinIOのみADR 0010でSuperseded） |
| [0010](0010-local-object-storage.md) | MinIOをやめ、localのobject storageはPhase 1で決める | Accepted |
