# ADR 0002: PaaSの組み合わせでdeployし、移行できる形を保つ

- Status: Accepted（region選択のみProvisional）
- Date: 2026-10-30

## 背景

MVP段階の費用は最小にしたい。一方で、成長時にAWSやGCPの別serviceへ移れる余地を残したい。
Artistが初めてRoomを開く瞬間はproductで最も重要な場面であり、cold startの扱いを決める必要がある。

## 決定

| 部品 | 採用 |
|---|---|
| Web | Vercel |
| API | Cloud Run（Docker image） |
| Database | Neon（PostgreSQL） |
| Object storage | Cloudflare R2 |

### Region（Provisional）

Neonに東京regionがないため、APIとDBを同じ場所へ揃えて`asia-southeast1`（Cloud Run）+ `aws-ap-southeast-1`（Neon）のSingapore構成とする。
1 requestの中でDB queryが複数回走るため、API〜DB間の距離を優先する。Neonのregionはproject作成後に変更できない。
Vercelのserver処理のregionも、初回deploy時に確認して揃える。

### Cost guardrail

- GCPのbudget alertを初回deploy前に必ず設定する（通知のみで課金は止まらない）
- Cloud Runの最大instance数に低い上限を設定する
- `/readyz`など、DBへ触れる処理を定期実行しない（ADR 0007）

### Cold start

- MVP開発中は許容する
- Beatmaker 5人によるproduct validationを始める前に解消する。これはrelease gateの条件とする
- Cloud Runの最小instance数を1にし、Neonの起動待ちは実測してから対応を判断する

### Portability

- APIはDocker imageとしてbuildし、Cloud Run固有の機能に依存しない
- DBは標準的なPostgreSQLの機能を使う
- storageはADR 0003、authはADR 0004で乗り換えの範囲を限定する

## 検討した代替案

- **全部AWS**: 求人での認知度は高い。ただしNAT GatewayやRDSなど固定費が読みにくく、IAMとVPCの設定量も多い。S3の転送量課金は音声配信と相性が悪い
- **APIをFly.io / Renderに置く**: 手軽だが、GCPの実運用経験につながらない
- **DBをCloud SQLに置く**: GCPに統一できるが、実質的な無料枠がなく常時費用がかかる
- **DBをSupabaseに置く**: browserから直接DBを扱う前提の機能が中心で、ADR 0001の境界と噛み合わない

## 影響

参考: 2026-10時点の公式情報。料金は変わるため、利用開始時に再確認する。

- Neon Free: 100 CU-hours/月、1 GB/project、転送5 GB/月、履歴保持6時間。scale to zeroは無効化できない。上限超過時はcomputeが次の請求期間まで停止する
- Cloud Run（request-based billing）: 200万requests、180,000 vCPU秒、360,000 GiB秒/月が無料
- R2 Standard: 10 GB-month、Class A 100万回、Class B 1,000万回/月が無料。egressは無料
- Vercel Hobby: 個人の非商用利用に限る。上限を超えても追加購入はできず、一定期間利用が止まる

見直す条件:

- 収益化する場合は、Vercelを有料planへ変更するか移行する
- Neon Freeの上限に近づいた場合、または本番dataの保護が必要になった場合はNeonの有料planを検討する（ADR 0006）
- Singapore構成の遅延が体感で問題になる場合はregionを見直す
