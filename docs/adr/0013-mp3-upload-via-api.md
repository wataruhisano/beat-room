# ADR 0013: MP3はAPI経由でuploadし、server側で検証する

- Status: Accepted
- Date: 2026-11-21

## 背景

ADR 0003で、browserからstorageへ直接uploadする方式とAPI経由の方式のどちらにするかをPhase 1に先送りしていた。
同時に、MP3の最大size、最大duration、対応bitrate（PRDの未確定事項）も決める必要がある。
Cloud RunはHTTP/1のrequest bodyを1 requestあたり32MiBまでに制限している（2026-11確認）。

## 決定

- browserはMP3をGo APIへ送る。Goが検証とduration取得をしてからstorageへ保存する
- MP3かどうかは拡張子やMIME typeではなく、file内容の解析で判定する
- 受け付ける条件は次のとおり。値はconfigに置き、変更しやすくする

| 項目 | 値 |
|---|---|
| 最大file size | 20MB |
| 最大duration | 10分 |
| 最小duration | 5秒 |
| bitrate | 制限なし（CBR / VBRとも可） |

- size上限はrequest bodyを読む段階で強制し、上限を超えたら途中で読むのをやめる
- storageへの保存が成功するまで、Beatを利用可能な状態にしない。`upload_status`の要否と状態遷移は、schemaの実装時に確定する

## 検討した代替案

- **storageへの直接upload（presigned URL）**: 大きいfileも扱え、APIを通らない。ただしbucketのCORS、upload途中で残ったobjectの掃除、「upload済みだが未検証」という状態の管理が増える
- **最大25MB**: 32MiBの上限に近く、multipartの余分なbyteやproxyの差で失敗しやすい
- **低bitrateを拒否する**: demoの段階で低bitrateのMP3を送る使い方を止める理由が、MVPにはない

## 影響

- 検証までが1回のrequestで終わり、「invalid MP3 / upload failure / retry」のtestが書きやすい
- bucketのCORSはuploadには不要になる（再生用のsigned URLでは必要になりうるので、Phase 2で確認する）
- 320kbpsのMP3は20MBで約8分なので、一般的なbeatの長さには十分
- **弱点**: upload中はAPIのinstanceがrequestを抱える。file sizeの上限を大きく上げたい場合は、直接upload方式への変更が必要になる
- storageへの保存後にDBへの書き込みが失敗すると、どこからも参照されないobjectが残る。Phase 5のorphan upload cleanupで扱う
