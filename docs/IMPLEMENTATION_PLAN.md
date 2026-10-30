# Beatroom Implementation Plan

## 1. 実装原則

- Phaseごとに未確定事項を確認してから実装する
- 各Phaseはdemo可能なvertical sliceで終える
- 後続Phaseのための抽象化を先に作りすぎない
- public Room、tenant isolation、履歴保持を後付けにしない
- 完了条件は「codeがある」ではなく、testと主要UIで挙動を確認できること

## 2. Phase 0 — Decisions and foundation

### 目的

開発・test・deployの最小基盤を作り、未確定な技術判断を記録する。

### 開始時に確認

- Next.js + Go + PostgreSQLを採用するか
- monorepo構成
- auth provider
- object storage
- deployment providerと予算
- package / dependency方針

### Tasks

- repository構成とdevelopment commands
- Web / API / DBのlocal起動
- environment variable管理
- migration framework
- lint / format / typecheck / test
- CI pipeline
- health check
- structured logging / error handlingの基礎
- ADR templateと初期ADR
- skeleton UIとdark theme tokens

### Done

- 新規環境でREADMEだけを使いlocal起動できる
- CIが自動実行される
- empty migrationをCIで適用でき、rollback方針がADRにある
- WebからAPI healthを確認できる

## 3. Phase 1 — Account, onboarding, Beat Library, Artists

### 開始時に確認

- 認証UX
- MP3 size / duration上限
- BPM / key / tags候補
- onboardingの終了地点とskip可否
- Artist contact validation

### Tasks

- authenticationとserver-side user identity
- tenant-scoped data access foundation
- private object storage upload
- MP3 validation / duration extraction
- first Beat onboarding
- Beat CRUD / archive
- Beat Library table
- search / filters / sort
- Artist CRUD / archive
- multiple contacts
- notes
- initial responsive navigation

### Tests

- cross-tenant Beat / Artist access denial
- invalid MP3 / upload failure / retry
- required title
- Artist requires display name + one contact
- archive hides new-selection candidates

### Done

- 2つのtest account間でdataが混ざらない
- 初回userが1曲uploadしてLibraryで確認できる
- Artistを登録・編集・archiveできる

## 4. Phase 2 — Room creation and public player

### 開始時に確認

- secret token hash方式と長さ
- Roomのdraft / active切替タイミング
- Room openの技術的な判定条件と、誤アクセス時のlock解除手段の要否
- Room名を編集可能にするか
- waveform libraryまたは実装方式

### Tasks

- Room / RoomBeat schema
- Room composer
- 1〜5曲validation
- drag & drop ordering
- message / Download setting
- generated Room name
- secret URL generation
- public Room response boundary
- anonymous session
- first-open transactional lock
- signed audio access / Range support
- mobile-first player
- waveform / seek / next autoplay
- Room detail / copy link
- revoke / archive

### Tests

- 0曲 / 6曲 / duplicate Beat拒否
- 他tenantのArtist / BeatをRoomへ追加不可
- concurrent first openでもlockが一貫する
- locked Roomの変更拒否
- revoke後のRoom / audio access拒否
- mobile browserでplay / seek / autoplay fallback確認

### Done

- Beatmakerが1〜5曲のURLを作り、別browserでログインなし再生できる
- 初回open後にRoomが固定される
- revokeが即時に効く
- Singapore構成（[ADR 0002](adr/0002-deployment-topology.md)）でRoom表示の遅延を実測し、DBのregionやproviderの見直し要否を判断する

## 5. Phase 3 — Reactions, comments, playback analytics

### 開始時に確認

- Like / Want thisのtoggle / 取消 / 重複ルール
- Comment edit / delete authorization
- Comment文字数上限
- progress送信間隔
- anonymous session期間

### Tasks

- interaction event ingestion
- idempotency
- playback summary
- 90% actual-listening completion rule
- Like UI / event
- Want this UI / confirmation message / event
- Beat Comment create / edit / delete
- Download gating / signed URL / event
- privacy-conscious session metadata
- rate limiting

### Tests

- seek-onlyではcompletedにならない
- retryでeventが重複しない
- Commentを別sessionが編集できない
- Download OFFではaccess不可
- Room外のBeatへeventを送れない
- event logにsecret / signed URLを残さない

### Done

- public Roomの全MVP反応が保存される
- analytics summaryがraw eventと整合する
- Artist本人と断定しない表示になっている

## 6. Phase 4 — Dashboard, Artist history, notifications

### 開始時に確認

- activityの保持・pagination
- unread notificationの定義
- 「送信日」をRoom作成日で代用するか、mark-as-sentを追加するか

### Tasks

- Dashboard Recent activity
- recent Rooms
- recently active Artists
- Want this emphasis
- in-app notifications / read state
- Artist detailの最終送信日
- Room単位の折りたたみtimeline
- 各Beatのcompletion / Like / Want this / Comment / Download表示
- 累計送付曲数
- loading / empty / error states

### Tests

- eventから適切なArtist / Roomへnavigationできる
- archived entityを含む過去履歴が表示される
- total play countを主要表示しない
- notification tenant isolation

### Done

- BeatmakerがDashboardから重要反応を見つけ、Artistの過去履歴へ辿れる
- archived dataを壊さずに履歴を表示できる

## 7. Phase 5 — Hardening and release

### Tasks

- full end-to-end suite
- accessibility review
- supported browser matrixでmanual QA
- performance確認（Room load / audio start）
- security review
- storage lifecycle / orphan upload cleanup
- DB backup / restore確認
- backupからの復元によるrollback rehearsal
- validation開始前のcold start解消
- privacy policy / terms / Download説明の準備
- production deployment
- operational dashboard / alerts
- seed / demo data
- READMEとarchitecture docs更新

### Done

- MVP_SCOPEのrelease gateを満たす
- known issuesとout-of-scopeが明記される
- product validationを開始できる

## 8. Post-MVP backlog

優先順位は利用データとinterviewで決める。

- Artist detailのNext Beats recommendation
- Taste Memory
- audio analysisによるBPM / key / tag候補
- 一定期間後の再送候補
- Room expiration / password
- Email / Discord notifications
- BeatStars等の外部購入導線
- sales pipeline
- WAV / stems
- configurable max Beats
- Artist account / stronger identity

## 9. Claude Codeへの進め方

各Phase開始時に、次の順で作業する。

1. 該当Phaseの「開始時に確認」を質問する
2. 既存docsとcodeの差分を確認する
3. そのPhaseのacceptance criteriaをtest caseへ変換する
4. 小さいvertical sliceに分解する
5. schema migration → API → UI → testsの順序をtaskごとに調整する
6. Phase終了時にdocs、README、未確定事項を更新する

一度に全Phaseを実装しない。未確定事項への回答がない場合は、可逆な仮決定をADRに残し、product behaviorを勝手に拡張しない。
