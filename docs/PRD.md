# Beatroom（仮）Product Requirements Document

## 1. 文書の目的

この文書は、Beatmakerが既に接点のあるArtistへ少数のBeatを送り、反応をArtist単位で蓄積するWebアプリ「Beatroom（仮）」のMVP要件を定義する。

詳細な実装方法は各実装フェーズ開始時に確認する。未決定事項を推測で埋めず、本書の「未確定事項」として扱う。

## 2. プロダクト概要

Beatroomは、BeatmakerがArtistごとにprivateなBeat pack（Room）を作り、secret URLで共有し、そのRoomへの反応を次の営業判断に活用するためのツールである。

中心価値は単なるファイル送信ではなく、次の2点にある。

1. Artistが登録不要で、音楽をすぐ聴いて反応できること
2. Beatmakerが「誰に、いつ、何を送り、どう反応されたか」を後から確認できること

将来的には蓄積した履歴から「次に送るBeat」を推薦するが、MVPでは履歴の収集と閲覧までに絞る。

## 3. 対象ユーザー

### 3.1 Primary: Beatmaker / Producer

- 複数のArtistへ継続的にBeatを送る
- Instagram、Discord、Emailなど既存の連絡手段を持つ
- 送付履歴や反応を手作業で管理している、または管理できていない
- 「次に何を送るか」をArtistごとに判断したい

### 3.2 Secondary: Artist

- Beatmakerから送られたリンクをスマートフォン中心に開く
- アカウント登録せずBeatを聴く
- Like、Want this、Comment、Downloadで反応する

ArtistはMVP上のログインユーザーではない。反応は「当該Artist向けRoomから行われた」と記録するが、操作した本人がArtist本人であることは保証しない。

## 4. 解決する課題

- DMやストレージリンクだけでは、Beatが開かれたか、どこまで聴かれたか分からない
- Artistごとの送付履歴と反応が分散し、次の選曲に活かせない
- Beatの送付準備に入力や操作が多いと、営業活動が継続しない
- Artistにアカウント登録を求めると、再生までの離脱が増える

## 5. Product principles

- **Music first:** Artist画面では再生を最優先する
- **Fast sharing:** BeatmakerがArtist選択からURLコピーまで短時間で完了できる
- **Low friction:** 入力を最小限にし、候補選択と自動入力を使う
- **Honest analytics:** secret URLの閲覧者をArtist本人と断定しない
- **History preserving:** 削除操作で過去の送付・反応履歴を壊さない
- **Focused MVP:** 販売、決済、チャット、推薦は既存サービスまたは将来機能に委ねる

## 6. 主要ユースケース

1. Beatmakerが初回登録し、オンボーディングで最初のMP3を1曲アップロードする
2. BeatmakerがBeat LibraryでBeatを検索・絞り込み・並べ替えする
3. BeatmakerがArtistと連絡先を登録する
4. BeatmakerがArtistを選び、1〜5曲を並べたRoomを作成する
5. Beatmakerがsecret URLを既存のDM等で共有する
6. ArtistがログインせずRoomを開き、再生・Like・Want this・Comment・Downloadを行う
7. BeatmakerがDashboard、Room、Artist detailで反応と履歴を確認する
8. Want thisを受けたBeatmakerが、InstagramまたはDiscord等で交渉を続ける

## 7. 機能要件

### 7.1 Beatmaker account

- 複数Beatmakerが利用できるmulti-tenant構成とする
- 各BeatmakerのBeats、Artists、Rooms、analyticsを他ユーザーから分離する
- ログイン方法は未確定

### 7.2 初回オンボーディング

- 初回ログイン後、空のDashboardだけを表示せず、最初のBeat 1曲のアップロードへ案内する
- アップロード完了をオンボーディングの最低完了条件とする
- オンボーディングの具体的な画面数とスキップ可否は未確定

### 7.3 Beat

- MVPの音声形式はMP3のみ
- titleは必須。初期値はファイル名から生成し、編集可能にする
- BPM、key、tagsは任意
- BPM、key、tagsは自由入力を主にせず、候補選択を中心にする
- cover画像はMVP対象外
- Beatをarchivedにしても、過去Roomと反応履歴は保持する

### 7.4 Beat Library

- table形式で表示する
- title、BPM、key、tagsで検索または絞り込みできる
- 「送ったことがある / ない」で絞り込みできる
- 最近作成した順、BPM順で並べ替えできる
- archived Beatは通常一覧から除外し、必要に応じて履歴上で表示する

### 7.5 Artist CRM

- display nameと最低1件のcontactを必須にする
- contact typeはInstagram、Discord、Email、Otherを許可し、複数登録できる
- reference artistや好みを自由に残せるnotesを持つ
- avatarはMVP対象外
- Artist detailの上部で最終送信日を目立たせる
- 履歴はRoom単位の時系列で表示し、折りたためる
- 累計送付曲数を表示する
- Artistをarchivedにしても履歴を保持する

### 7.6 Room作成・管理

- RoomはArtistへの1回の送付単位とする
- 1 Roomにつき1〜5曲を登録できる
- 最大曲数のユーザー設定は将来候補。MVPは5曲固定
- 曲順をdrag & dropで変更できる
- Artistへの短いメッセージを設定できる
- Room名は `Artist名 + 日付` を初期値として自動生成する
- Downloadの許可はRoom全体でON/OFFする
- 作成者はArtistが開く前ならRoom内容を編集できる
- Artistが一度でもRoomを開いた後は、収録Beat、順序、message、Download設定を含むRoom内容を変更できない
- 作成者はRoom URLをいつでもrevokeできる
- 送信済みRoomの削除は物理削除ではなくarchiveを基本とする

### 7.7 Public Room

- 推測困難なsecret token URLでアクセスする
- Artistのログイン、名前入力、passwordは不要
- 有効期限は設けない
- revoked Roomにはアクセスできない
- mobile-firstで設計する
- Playを第一の操作として大きく表示する
- Likeは常時表示する
- Want thisはLikeより強く目立たせる
- CommentはBeat単位の通常テキストとする
- timestamp commentは対象外
- Artist側で自分のCommentを編集・削除できる
- Downloadは補助メニュー内に置き、Roomで許可された場合のみ利用できる
- DownloadはMP3のdemo / writing / recording用途であり、商用利用権を付与しない
- Download時の用途上の注意文はMVPでは表示しない
- 1曲終了時に次の曲を自動再生する
- 波形を表示する
- 波形またはprogress barから再生位置を移動できる（seek）
- 再生速度変更は提供しない

### 7.8 Reactions and events

- Like: 「良い」という軽い肯定反応
- Want this: 「実際に使いたい / 話を進めたい」という強い意思表示
- Want this後、ArtistにはProducerへ通知した旨を表示する
- 交渉はアプリ内で行わず、既存のInstagram、Discord等へ戻す
- Room open、play progress、Like、Want this、Comment、Downloadを記録する
- 明示的なPassは設けない

### 7.9 Analytics and history

- Room opened
- 各Beatの最大再生到達率
- Like
- Want this
- Comment
- Download
- 完了判定は90%以上の実再生とする
- seekで末尾へ移動しただけでは完了扱いにしない
- Artist detailでは、送付RoomとBeatごとの反応を時系列表示する
- Dashboardに総再生数は表示しない
- IP等からArtist本人であると断定しない。URL転送等によるアクセスはsessionとして分けて扱う

### 7.10 Dashboard and notifications

- Recent activity
- 最近作成したRooms
- 最近反応があったArtists
- Want this通知
- MVPの通知はアプリ内のみ
- Email、Discord webhook等は対象外

## 8. UX / visual direction

- dark mode中心
- minimal / music-first
- marketplace感を避ける
- trapカルチャーのニュアンスを適度に含める
- 参照イメージとして日本のorigami、USの2hollisが挙げられているが、具体的なブランド・ビジュアル模倣はしない
- 最終的なdesign tokens、logo、typeface、アクセシビリティ基準は未確定

## 9. 成功指標

### Product validation

- 知人以外のBeatmaker 5人が利用する
- そのうち3人以上がRoomを2回以上作成する

### Portfolio / delivery completion

- 本番環境で主要フローが安定動作する
- 自動テスト、CI/CD、README、設計資料が揃っている
- tenant間データ分離とsecret URLのアクセス制御を検証できる

## 10. 非目標

- Beat販売、決済、license発行
- アプリ内chat
- Artist accountおよび本人確認
- Recommendation / ML / Taste inference
- WAV対応
- Email / Discord通知
- password、期限付き、one-time Room
- timestamp comment
- producer tagの自動付与
- cover画像
- 総再生数を中心にしたvanity metrics

## 11. 未確定事項

実装開始前または該当Phaseの開始時に確認する。

1. 認証方式（email/password、magic link、OAuth等）
2. ~~暫定技術スタックを採用するか~~ → 確定（[ADR 0001](adr/0001-web-api-monorepo.md)）
3. MP3の最大ファイルサイズ、最大再生時間、対応bitrate
4. BPM、key、tagsの候補リストと複数選択ルール
5. Roomが「開かれた」の技術的な判定条件と、誤アクセス時のlock解除手段の要否
6. Like / Want thisを同一sessionから取り消せるか、重複操作をどう扱うか
7. ArtistがCommentを編集・削除するための識別方法
8. sessionの有効期間とプライバシー方針
9. 対応言語（日本語のみ、英語のみ、両対応）
10. デプロイ先、object storage、運用予算の方針は確定（[ADR 0002](adr/0002-deployment-topology.md)、[ADR 0003](adr/0003-object-storage.md)）。CDNの要否は未確定
11. アクセシビリティの達成基準
