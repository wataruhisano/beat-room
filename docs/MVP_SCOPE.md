# Beatroom MVP Scope

## 1. MVP definition

MVPは次のloopが本番環境で成立する状態とする。

> BeatmakerがBeatとArtistを登録する → 1〜5曲のRoomを作る → secret URLを送る → Artistが登録なしで聴いて反応する → BeatmakerがArtist別の履歴として確認する

## 2. In scope

### Foundation

- Beatmaker authentication（方式は未確定）
- multi-tenant data isolation
- responsive web application
- production deployment
- basic CI
- database migration
- error loggingと最低限の運用監視

### Beats

- MP3 upload / playback / download
- title必須、filenameから初期値生成
- 任意のBPM、key、tags
- table型Library
- title / BPM / key / tags検索・絞り込み
- 送付済み / 未送付filter
- 最近作成順 / BPM順sort
- archive

### Artists

- CRUD（削除はarchive）
- display name必須
- 複数contact、最低1件必須
- reference / taste notes
- 最終送信日
- Room単位の折りたたみ履歴
- 累計送付曲数

### Rooms

- 送付単位で作成
- 1〜5 Beats
- drag & drop ordering
- message
- Room名自動生成
- Room全体のDownload ON/OFF
- secret URL
- 初回open後の内容lock
- revoke
- archive

### Artist experience

- login / name input不要
- mobile-first player
- waveform
- seek
- next track autoplay
- Like
- Want this
- Beat単位Comment
- Artist側Comment edit / delete
- 許可された場合のMP3 Download

### Analytics / activity

- Room open
- 最大再生到達率
- 90%実再生によるcompleted判定
- Like / Want this / Comment / Download event
- Artist timeline
- in-app notifications
- DashboardのRecent activity、recent Rooms、recently active Artists、Want this

### Quality bar

- 主要フローのautomated tests
- tenant boundaryとpublic token accessのsecurity tests
- upload / playback / event記録のintegration tests
- README、environment setup、migration、deployment手順
- main branchのCI

## 3. Explicitly out of scope

- Beat購入、payment、license発行、契約管理
- アプリ内chat
- Artist account、identity verification
- AI / ML、自動metadata解析、Taste Memory、Next Beats recommendation
- WAVやstem upload
- cover art
- password / expiration / one-time Room
- Email通知、Discord webhook、Instagram連携
- timestamp Comment
- producer tag付与
- explicit Pass
- 再生速度変更
- aggregate Like率や総再生数を中心とするadvanced analytics
- Room最大曲数のuser setting
- 送信後のsales pipeline（Interested / Negotiating / Licensed / Placed）
- native mobile app

## 4. Acceptance criteria

### End-to-end

- 新規Beatmakerが最初のMP3を登録できる
- Artistとcontactを登録できる
- 1〜5曲を選択・並べ替え、Room URLを発行できる
- URLを知るArtistがログインせず再生できる
- Like、Want this、Comment、Downloadが許可どおり機能する
- Beatmakerが対象ArtistとRoomに紐づく反応を確認できる
- revoke後はpublic Roomへアクセスできず、過去の履歴は残る
- archived Beat / Artistは新規作成候補から除外され、過去履歴には残る

### Analytics correctness

- 再生進捗はsession・Room・Beatに紐づく
- seekで90%地点へ移動しただけではcompletedにならない
- 実再生が曲長の90%以上に達した時のみcompletedとなる
- UIはアクセス者をArtist本人と断定する表現を使わない

### Security

- User AはUser Bのprivate dataを参照・変更できない
- public APIは有効なRoom tokenに含まれる情報だけを返す
- audio objectの保存先URLを恒久的に公開しない
- revoked Roomから新規audio accessを取得できない
- secret tokenを平文でDB保存するかhash保存するかは設計時に決め、推測困難性を検証する

## 5. Scope control rules

- 各Phase開始時に未確定事項だけを確認する
- 新機能の提案は、上記end-to-end loopを完成させるまでbacklogへ置く
- analytics項目はイベントを増やす前に、ユーザー判断に使う表示先を明確にする
- 将来MLに使える可能性だけを理由に、MVPで不要な個人情報やeventを収集しない
- 技術的に難しい場合、ユーザー価値を変えずに簡素化する。例: 高機能waveform editorではなく再生位置が分かる波形

## 6. Release gate

MVP releaseには以下をすべて満たす。

- end-to-end acceptance criteriaが通る
- tenant isolation testが通る
- secret URL revoke testが通る
- supported mobile / desktop browserでuploadとplaybackを確認する
- privacy / terms上、保存するアクセス情報とMP3利用条件を説明できる
- backup / migration / rollbackの最低限手順がある
- known limitationsをREADMEに記載する

## 7. Product validation after release

- 知人以外のBeatmaker 5人が利用
- 3人以上が2回以上Roomを作成
- 数値だけでなく、Room作成に詰まった点とArtist側の離脱理由を短いinterviewで確認する
