# ADR 0012: localとCIのobject storageはSeaweedFSとし、R2の固有挙動は節目で確認する

- Status: Accepted
- Date: 2026-11-21

## 背景

ADR 0010でMinIOをやめ、localのobject storageはPhase 1で決めることにした。Phase 1でMP3 uploadを実装するため、ここで選ぶ。
本番はCloudflare R2で、GoからはS3互換APIで操作する（ADR 0003）。CIでもsecretなしでstorageを使ったtestを動かしたい。

## 決定

- localとCIでは、SeaweedFSのS3 gatewayをDocker Composeで起動する（image `chrislusf/seaweedfs`、single containerの`weed mini`相当）
- signed URL、Range request、CORSのようにproviderごとの違いが出やすい挙動は、開発用のR2 bucketで節目ごとに手動確認する。最初の確認はPhase 2のaudio配信の実装時に行う
- 開発用R2のcredentialは各自の`.env`に置き、CIには渡さない

## 検討した代替案

- **Garage**: presigned URLとCORSには対応している。ただし2026年の開発状況を公式情報で確認できず、互換性の表も更新が止まっている
- **開発用R2 bucketだけを使う**: 本番と同じ挙動で確認できる。ただし各環境とCIでcredentialの管理が必要になり、CIがnetworkとsecretに依存する
- **storageをmockする**: 速いが、S3互換APIとの結合部分をtestできない

## 影響

- offlineで開発でき、CIでもstorageを含むtestを動かせる
- SeaweedFSは2026-11-16に4.54がreleaseされるなど保守が続いている。Apache-2.0。S3 APIのpresigned URL、Range GET、bucket CORSに対応している（公式wiki、2026-11確認）
- `/readyz`にstorageを含めるかは、uploadの実装時にADR 0007の基準で判断する
- **弱点**: R2とSeaweedFSの細かい挙動の違いは、CIでは検出できない。手動確認を省くと本番で初めて気づく
