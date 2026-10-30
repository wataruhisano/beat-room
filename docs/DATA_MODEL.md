# Beatroom Data Model

## 1. 方針

- 内部IDはUUIDを基本とする
- すべてのprivate business dataをowner userへ辿れるようにする
- 過去履歴を壊す物理削除を避け、`archived_at`を用いる
- Roomへ入れたBeatのtitle等は必要に応じてsnapshotし、後の編集で送付時の表示が変わらないようにする
- raw eventと表示用summaryを分離する
- Artist accountを前提にしない
- timestamp / timezoneはUTC保存、UIでlocal表示を基本とする

## 2. Entity relationship

```text
User
 ├─< Beat
 ├─< Artist ─< ArtistContact
 ├─< Room >─ Artist
 │    ├─< RoomBeat >─ Beat
 │    ├─< RoomSession
 │    └─< Interaction >─ RoomBeat
 └─< Notification

RoomBeat ─< Comment
```

## 3. Tables

### users

| Column | Type | Rules / purpose |
|---|---|---|
| id | uuid | PK |
| auth_subject | text | UNIQUE、Firebase AuthenticationのUID（[ADR 0004](adr/0004-authentication.md)） |
| display_name | text | nullable可否は未確定 |
| created_at | timestamptz | required |
| updated_at | timestamptz | required |

業務tableは`users.id`を参照し、`auth_subject`を持たない。emailをここへ保持するかは、login方式を決めるPhase 1で判断する。

### beats

| Column | Type | Rules / purpose |
|---|---|---|
| id | uuid | PK |
| owner_user_id | uuid | FK users、required |
| title | text | required |
| bpm | smallint | nullable、妥当範囲は未確定 |
| musical_key | text | nullable、許容値を制約 |
| duration_ms | integer | required、serverで取得 |
| storage_key | text | required、object storageの内部key |
| original_filename | text | required |
| mime_type | text | `audio/mpeg` を基本 |
| byte_size | bigint | required |
| upload_status | enum | `pending / ready / failed` |
| archived_at | timestamptz | nullable |
| created_at | timestamptz | required |
| updated_at | timestamptz | required |

`owner_user_id, created_at`、`owner_user_id, bpm`、title検索用indexを検討する。

### tags

| Column | Type | Rules / purpose |
|---|---|---|
| id | uuid | PK |
| name | text | canonical label |
| category | text | nullable。mood等の分類候補 |
| is_active | boolean | 候補として表示するか |

候補を全ユーザー共通にするかuser-definedを許すかは未確定。MVPで固定候補のみならseed dataでよい。

### beat_tags

| Column | Type | Rules / purpose |
|---|---|---|
| beat_id | uuid | FK beats |
| tag_id | uuid | FK tags |

複合PK `(beat_id, tag_id)`。

### artists

| Column | Type | Rules / purpose |
|---|---|---|
| id | uuid | PK |
| owner_user_id | uuid | FK users、required |
| display_name | text | required |
| notes | text | reference artist、好み等 |
| archived_at | timestamptz | nullable |
| created_at | timestamptz | required |
| updated_at | timestamptz | required |

### artist_contacts

| Column | Type | Rules / purpose |
|---|---|---|
| id | uuid | PK |
| artist_id | uuid | FK artists、required |
| type | enum | `instagram / discord / email / other` |
| value | text | required |
| label | text | Other等の補足、nullable |
| is_primary | boolean | 主要連絡先 |
| created_at | timestamptz | required |

applicationまたはtransactionでactive Artistに最低1件のcontactを保証する。

### rooms

| Column | Type | Rules / purpose |
|---|---|---|
| id | uuid | PK |
| owner_user_id | uuid | FK users、required |
| artist_id | uuid | FK artists、required |
| name | text | required、自動生成後に編集可否は未確定 |
| message | text | nullable |
| download_enabled | boolean | required、default false候補 |
| status | enum | `draft / active / revoked / archived` |
| public_token_hash | text | UNIQUE。hash保存を推奨、確定は設計時 |
| locked_at | timestamptz | 最初の有効open時、nullable |
| revoked_at | timestamptz | nullable |
| archived_at | timestamptz | nullable |
| created_at | timestamptz | required |
| updated_at | timestamptz | required |

`owner_user_id, artist_id, created_at`と`public_token_hash`にindex。

### room_beats

| Column | Type | Rules / purpose |
|---|---|---|
| id | uuid | PK |
| room_id | uuid | FK rooms、required |
| beat_id | uuid | FK beats、required |
| position | smallint | 1〜5 |
| title_snapshot | text | 送付時title |
| duration_ms_snapshot | integer | 送付時duration |
| created_at | timestamptz | required |

制約:

- UNIQUE `(room_id, beat_id)`
- UNIQUE `(room_id, position)`
- Roomあたり最大5件はtransaction内で検証
- `locked_at`以後のinsert / update / deleteをapplicationで拒否し、必要ならDB triggerも検討

### room_sessions

| Column | Type | Rules / purpose |
|---|---|---|
| id | uuid | PK |
| room_id | uuid | FK rooms、required |
| anonymous_session_id_hash | text | browser識別子のhash等 |
| started_at | timestamptz | required |
| last_seen_at | timestamptz | required |
| user_agent_family | text | nullable、収集範囲は要確認 |

IP addressやprecise locationはMVPの価値に不要なら保存しない。session期間・cookie方針は未確定。

### interactions

append-only eventを基本とする。

| Column | Type | Rules / purpose |
|---|---|---|
| id | uuid | PK |
| room_id | uuid | FK rooms、required |
| room_beat_id | uuid | FK room_beats、nullable（room_open時はnull） |
| session_id | uuid | FK room_sessions、required |
| type | enum | `room_open / play_started / progress / completed / liked / like_removed / wanted / want_removed / downloaded` |
| occurred_at | timestamptz | required |
| position_ms | integer | progress系のみ |
| listened_delta_ms | integer | 実再生時間の集計用、nullable |
| event_key | text | retry重複防止用idempotency key |
| metadata | jsonb | versioned補助情報。乱用しない |

UNIQUE `(session_id, event_key)`。

progress eventの頻度は未確定。過剰なevent量を避けつつ、最大到達位置と実再生時間を信頼できる方式にする。

### playback_summaries

raw eventから同期または非同期に更新するread model。

| Column | Type | Rules / purpose |
|---|---|---|
| session_id | uuid | FK room_sessions |
| room_beat_id | uuid | FK room_beats |
| max_position_ms | integer | 最大到達位置 |
| listened_ms | integer | seek分を除く実再生累積の上限付き値 |
| completed_at | timestamptz | 90%以上実再生時 |
| updated_at | timestamptz | required |

複合PK `(session_id, room_beat_id)`。

### comments

| Column | Type | Rules / purpose |
|---|---|---|
| id | uuid | PK |
| room_beat_id | uuid | FK room_beats、required |
| session_id | uuid | FK room_sessions、required |
| body | text | required、長さ上限未確定 |
| author_edit_token_hash | text | accountなし編集権の候補設計 |
| edited_at | timestamptz | nullable |
| deleted_at | timestamptz | nullable |
| created_at | timestamptz | required |

Comment編集方式は未確定。tokenをlocal storage等に渡す方式には消失・共有端末リスクがあるため、Phase 3開始時に決める。

### notifications

| Column | Type | Rules / purpose |
|---|---|---|
| id | uuid | PK |
| user_id | uuid | FK users、required |
| type | enum | `like / want_this / comment / download` 等 |
| room_id | uuid | FK rooms |
| room_beat_id | uuid | nullable |
| source_id | uuid | interaction/comment ID |
| read_at | timestamptz | nullable |
| created_at | timestamptz | required |

## 4. Derived fields

以下は原則としてsource of truthにせずqueryまたはread modelで算出する。

- Artist.last_sent_at = latest active/archived Room created/sent timestamp
- Artist.total_beats_sent = distinct RoomBeat count across sent Rooms
- Beat.has_been_sent = exists RoomBeat in non-draft Room
- completion rate = `min(listened_ms / duration_ms_snapshot, 1.0)`
- completed = completion rate >= 0.9

「Roomを作った日時」と「実際に送った日時」は現状区別できない。Share / mark-as-sent操作を設けるかは未確定。

## 5. Invariants

- private queryは常にauthenticated `owner_user_id`でscopeする
- Room.artistと全RoomBeat.beatは同じownerに属する
- archived Beat / Artistを新規Roomへ追加しない
- active Roomには1〜5 RoomBeatが必要
- locked Roomの構成を変更しない
- revoked / archived Roomにpublic accessを許可しない
- `completed`はpositionだけで決めず、実再生時間を使う
- interactionはpublic tokenで解決されたRoom外のRoomBeatを参照できない

## 6. Retention and privacy（未確定）

- event、session、Commentを保持する期間
- account削除時の扱い
- IP addressを処理する場合の保存有無
- cookie / local storageへの同意表示
- uploaded musicの削除要求と履歴snapshotの扱い
- audit logの範囲
