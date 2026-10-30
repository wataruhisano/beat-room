# Beatroom User Flows

## 1. 前提

- Beatmakerはログインユーザー
- Artistはログインせず、Artist専用に発行されたRoom URLを使う
- RoomはArtistへの1回の送付単位
- `opened`後のRoom内容は固定する
- 公開Room上の操作をArtist本人の操作とは断定しない

## 2. 画面一覧

### Beatmaker側

| 画面 | 主目的 |
|---|---|
| Sign up / Login | アカウント作成・ログイン |
| Onboarding upload | 最初のBeatを1曲登録 |
| Dashboard | 最近のRoom・Artist反応・Want thisを確認 |
| Beat Library | Beatの検索、絞り込み、並べ替え、追加 |
| Beat detail/edit | metadata編集、archive |
| Artist list | Artist検索、追加、archive |
| Artist detail | 最終送信日、contact、notes、Room履歴を確認 |
| Room composer | Artist・Beat・順序・message・download設定 |
| Room detail | URLコピー、状態、反応、revoke、archive |
| Notifications | アプリ内反応一覧 |

### Artist側

| 画面 | 主目的 |
|---|---|
| Public Room | Beatの再生と反応 |
| Revoked / unavailable | 利用不能状態の説明 |

## 3. Flow A: 初回利用

1. Beatmakerがsign up / loginする
2. 初回利用ならOnboarding uploadへ移動する
3. MP3を選ぶ
4. filenameからtitle候補を表示する
5. Beatmakerがtitleを確認・編集する
6. 必要ならBPM、key、tagsを候補から選ぶ
7. 保存する
8. 成功時、Beat Libraryまたは次の導線へ移動する

### エラー

- MP3以外: 対応形式を明示し、選択し直せる
- upload失敗: 入力値を保持し、再試行できる
- title未入力: 保存せず、該当fieldを示す
- 制限超過: 許容サイズまたは長さを示す（値は未確定）

### 未確定

- 登録直後にArtist追加またはRoom作成まで案内するか
- Onboardingのスキップ可否

## 4. Flow B: 2曲目以降のBeat登録

1. Beat LibraryでAdd Beatを押す
2. MP3を選択する
3. 自動入力されたtitleを確認する
4. 任意metadataを設定する
5. 保存後、tableに追加される

## 5. Flow C: Beatを探す

1. Beat Libraryを開く
2. title search、BPM、key、tags、送付済み/未送付filterを必要に応じて使う
3. 最近作成した順またはBPM順にsortする
4. Beatを選び詳細を開く、またはRoom作成に使う

条件を組み合わせた際のAND/OR仕様は実装前に確定する。

## 6. Flow D: Artist登録

1. Artist listでAdd Artistを押す
2. display nameを入力する
3. Instagram / Discord / Email / Otherから最低1件のcontactを入力する
4. 必要なら複数contact、reference artist、好みのnotesを追加する
5. 保存し、Artist detailへ移動する

### validation

- display name必須
- 有効なcontactが最低1件必要
- contact typeごとの厳密なformat validationは未確定

## 7. Flow E: Room作成と共有

1. Artist detailまたはRoom composerからArtistを選ぶ
2. Beat Libraryから1〜5曲を選ぶ
3. drag & dropで曲順を調整する
4. 必要ならmessageを入力する
5. Room全体のDownload ON/OFFを選ぶ
6. 自動生成されたRoom名（Artist名 + 日付）を確認する
7. previewしてRoomを作成する
8. secret URLをコピーする
9. Instagram / Discord / Email等でArtistへ送る

### validation

- Artistは1人必須
- activeなBeatを1〜5曲選ぶ
- 同一Beatを同じRoomに重複登録できない
- Room作成時点では編集可能

### 失敗時

- 作成失敗時も選択内容を保持し、再試行できる
- upload処理中のBeatは選択不可とする

## 8. Flow F: Roomの編集とlock

1. Beatmakerが未開封Room detailを開く
2. Beat、順序、message、Download設定を編集する
3. Artistがsecret URLを初めて開く
4. systemが最初の有効なopen eventを記録し、Roomの収録内容をlockする
5. 以後、Beat、順序、message、Download設定を含むRoom全体を変更不可にする

## 9. Flow G: ArtistがRoomを聴く

1. Artistがsecret URLを開く
2. systemがtokenの有効性とRoom状態を検証する
3. 有効ならRoomを表示し、sessionとopen eventを記録する
4. Artistが大きなPlay操作で1曲目を再生する
5. 再生中に波形またはprogressからseekできる
6. 曲が終了すると次の曲を自動再生する
7. 最終曲終了時は停止する

### 状態分岐

- token不正: Roomが見つからない旨を表示する
- revoked / archived: 利用できない旨を表示する
- 音声取得失敗: 再試行を提示し、他の曲は可能なら操作できる
- browserのautoplay制限: 最初の再生は必ずArtist操作から始める

## 10. Flow H: Artistの反応

### Like

1. 対象BeatのLikeを押す
2. 即時に視覚feedbackを返す
3. eventを記録する
4. BeatmakerのactivityとArtist履歴へ反映する

### Want this

1. 対象BeatのWant thisを押す
2. 確認が必要かは未確定
3. eventを記録する
4. 「Producerに通知しました」と表示する
5. Beatmakerのアプリ内通知に表示する
6. BeatmakerとArtistは既存の外部連絡手段で続きを行う

### Comment

1. 対象BeatのComment欄を開く
2. テキストを入力して送信する
3. CommentをRoomのArtistに紐づく反応として表示する
4. 同じ閲覧者が後から編集・削除できる

Comment編集者の識別方法は未確定。Artist accountがないため、安全なbrowser token等の設計判断が必要。

### Download

1. menuからDownloadを選ぶ
2. DownloadがOFFなら操作を表示しない
3. MP3を取得し、download eventを記録する

用途上の注意文はMVPでは表示しない。Downloadが商用利用権を付与しないというproduct上の扱いは、利用規約等で別途整備する。

## 11. Flow I: 反応を確認する

1. BeatmakerがDashboardを開く
2. Recent activity、最近のRooms、最近反応のあったArtists、Want this通知を見る
3. activityからRoom detailまたはArtist detailへ移動する
4. Beatごとの最大再生到達率、Like、Want this、Comment、Downloadを確認する
5. Artist detailで過去Roomを時系列に確認する
6. Roomを展開し、各Beatの反応を見る

総再生数を主要KPIとして表示しない。

## 12. Flow J: Revoke / archive

### Room revoke

1. BeatmakerがRoom detailでRevokeを選ぶ
2. 影響（URLが即時利用不可、履歴は保持）を確認する
3. revokeする
4. public URLはunavailable画面を返す

### Beat / Artist archive

1. Beatmakerがarchiveを選ぶ
2. 通常一覧と新規Room候補から除外する
3. 過去Room、event、Comment、Artist履歴は保持する

復元機能をMVPに含めるかは未確定。
