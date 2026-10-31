# ADR 0010: MinIOをやめ、localのobject storageはPhase 1で決める

- Status: Accepted
- Date: 2026-10-31

## 背景

ADR 0003とADR 0009は、local開発のS3互換storageとしてMinIOを指定していた。その後、MinIOのcommunity repositoryがarchiveされ、「no longer maintained」と明記されていることを確認した。保守されないsoftwareを新規に採用する理由はない。

一方、Phase 0で必要なのはWeb、API、PostgreSQLだけで、object storageはMP3 uploadを実装するPhase 1で初めて使う。

## 決定

- MinIOは使わない。ADR 0003とADR 0009のMinIOに関する記述を置き換える
- Phase 0の`compose.yaml`はPostgreSQLだけを起動する
- localのobject storageはPhase 1開始時に決める

S3互換APIで操作し、DBには`storage_key`のみを保存するというADR 0003の他の決定は変えない。

## 検討した代替案

- **今SeaweedFSなどの代替を入れる**: 使う機能がない段階では適切に評価できず、管理対象だけが増える
- **今開発用のR2 bucketを用意する**: 本番と同じ実装で確認できるが、credential管理がPhase 0から必要になる

どちらもPhase 1の候補として残す。

## 影響

- Phase 0で管理するserviceが減る
- `/readyz`の確認対象はPostgreSQLのままで、ADR 0007と矛盾しない
- **弱点**: S3互換の動作確認はPhase 1まで行えない
- Phase 1開始時の確認事項に、localのobject storageの選定を追加する
