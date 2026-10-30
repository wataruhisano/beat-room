# Beatroom Architecture

## 1. Status

この文書はMVPを開始するための**暫定architecture案**である。参照会話ではNext.js、Go API、PostgreSQL、Dockerが実装案として挙がったが、ユーザーによる最終確定はされていない。

Phase 0の最初に次を確認する。

- Next.js + Go APIの分離構成を採用するか
- deployment provider
- object storage / CDN
- authentication provider
- 月額予算と想定traffic

これらを確認するまではprovider固有の実装を増やさない。

## 2. Proposed system shape

```text
Beatmaker browser ─┐
                   ├─ Next.js Web ── Go HTTP API ── PostgreSQL
Artist browser ────┘                     │
                                        ├─ Object Storage (private MP3)
                                        └─ Background/async processing (必要になった時のみ)
```

### Web application

- Beatmaker dashboard / CRUD UI
- Public Room player
- mobile-first responsive UI
- audio playback stateとwaveform UI
- serverから許可された短命audio URLを使用

### API

- authenticated Beatmaker operations
- public Room token resolution
- upload orchestration
- event ingestion / idempotency
- tenant authorization
- signed audio access

### PostgreSQL

- users、Beats、Artists、Rooms、RoomBeats
- sessions、interactions、comments、notifications
- migrationによるschema管理

### Object storage

- MP3をprivate bucket/containerに保存
- DBにはprovider URLではなくstorage keyを保存
- upload/downloadは短命signed URLまたはAPI proxyを利用
- permanent public object URLは使用しない

## 3. Suggested repository layout

monorepoを採用する場合の一例。確定ではない。

```text
/
├─ apps/
│  ├─ web/                 # Next.js
│  └─ api/                 # Go
├─ db/
│  ├─ migrations/
│  └─ seeds/
├─ docs/
├─ infra/                  # provider確定後
├─ compose.yaml
└─ README.md
```

package manager、Go router / ORM / query generator、migration toolはPhase 0で選ぶ。

## 4. API boundary

API pathやpayloadの厳密な仕様は実装時に定義する。必要なcapabilityは以下。

### Authenticated capabilities

- current user / onboarding state
- Beat create、upload complete、list/filter/sort、update、archive
- Artist create、list、detail、update、archive、contacts
- Room create、preview、update-before-lock、list、detail、revoke、archive
- Dashboard activity、notifications、mark read
- Room / Artist analytics read

### Public token capabilities

- public Room取得
- audio access取得
- session開始 / 継続
- playback event送信
- Like / Want this操作
- Comment create / edit / delete
- Download accessとevent記録

public endpointはtokenからRoom scopeを決め、clientが送るowner IDや任意のRoom IDを信用しない。

## 5. Critical request flows

### MP3 upload

1. WebがAPIへmetadataとupload開始を要求
2. APIがowner-scoped Beatを`pending`で作成
3. APIが短命upload URLを返す、またはstream uploadを受ける
4. clientがMP3を送信
5. serverがtype / sizeを検証し、durationを取得する
6. Beatを`ready`にする
7. failure時は`failed`とし、retry / cleanup可能にする

direct uploadとAPI proxyの選択はstorage providerとhosting制限を確認して決める。

### Public Room load

1. clientがsecret tokenを送る
2. APIがtoken hashでRoomを解決する
3. active / revoked / archivedを判定する
4. 最初の有効openをtransaction内で記録し、`locked_at`を設定する
5. public表示に必要な最小情報だけを返す
6. sessionを発行する

### Playback progress

1. playerが一定間隔または状態変化時に進捗を送る
2. event keyでretryをdeduplicateする
3. serverがpositionの急増を実再生時間として数えない
4. playback summaryを更新する
5. `listened_ms >= duration * 0.9`でcompletedを一度だけ記録する

正確なsampling間隔とbackground tabの扱いはPhase 3で決める。

## 6. Security model

### Tenant isolation

- authenticated routeはserver-side identityからuser IDを得る
- repository/query layerですべてowner scopeを必須にする
- IDだけでresource取得する汎用関数を避ける
- cross-tenant integration testsをCIに入れる
- PostgreSQL Row Level Securityを使うかはdata access方式決定後に判断する

### Secret URLs

- cryptographically secure random tokenを用いる
- token entropyは最低128 bit相当を候補とする
- DBにはtoken hashのみ保存する案を推奨
- application log、analytics、error reportingにfull tokenを残さない
- `Referrer-Policy`等で外部へのtoken漏洩を抑える
- revokeを即時反映する

### Audio

- storageはprivate
- public Room検証後に短命accessを発行
- URL expiryは再生体験を壊さない長さにする
- Range requestsに対応し、seekを可能にする
- Download OFFでもbrowser再生用byte取得自体を完全には防げない。この制限はproduct上明示する

### Input and abuse

- Commentをsanitize / escapeする
- public event、Comment、Download URL発行へrate limitを設ける
- uploadのMIME、magic bytes、sizeをserver側で検証する
- CSRF / CORSは認証方式に合わせて設定する

## 7. Consistency and concurrency

- Room初回openとlockはtransactionalに行う
- lock後の編集はserverで必ず拒否する
- Room作成時、1〜5曲・同一owner・重複なしをtransactionで検証する
- Like / Want thisの重複ルールは未確定。決定後DB constraintを置く
- public eventはidempotency keyで二重登録を防ぐ
- archive / revokeとpublic accessのraceをtestする

## 8. Observability

MVPに必要な最低限:

- structured application logs
- request / correlation ID
- error tracking
- upload failure rate
- public Room load / audio error
- DB migration status
- health endpoint

full secret token、Comment本文、contact値、signed URLをlogに含めない。

## 9. Testing strategy

- unit: validation、completion計算、state transition
- integration: DB constraints、tenant isolation、Room lock、revoke、event idempotency
- API contract: authenticated / public response boundaries
- end-to-end: onboarding → Room共有 → Artist反応 →履歴確認
- browser: audio playback、Range / seek、autoplay制限、mobile viewport
- security: token推測困難性、cross-tenant access、revoked access

## 10. Deployment environments

- local: containersでWeb / API / PostgreSQL、storage emulatorまたはdevelopment bucket
- preview/staging: migrationとE2E確認
- production: private storage、TLS、backup、monitoring

hosting構成が未確定のため、環境名以上のprovider設計は行わない。

## 11. Architecture decisions to record

Phase 0以降、重要判断は短いADRとして残す。

1. Web / API frameworkとrepository structure
2. Auth provider / session model
3. Object storageとupload方式
4. public token storage方式
5. anonymous Room session方式
6. Comment edit authorization方式
7. playback progress集計方式
8. deployment topology
