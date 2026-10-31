# ADR 0003: S3互換APIとstorage keyでMP3を管理する

- Status: Accepted（local開発のMinIOのみ[ADR 0010](0010-local-object-storage.md)でSuperseded）
- Date: 2026-10-30

## 背景

MP3はprivateに保存し、恒久的な公開URLを持たせない。音声配信はegress費用が増えやすく、storage providerの変更余地も残したい。

## 決定

- Cloudflare R2をprivate bucketとして使う
- GoからはS3互換APIで操作する。endpointとcredentialは環境変数で切り替える
- DBには`storage_key`のみを保存し、provider URLは保存しない
- 再生やdownloadのURLは、権限確認の後に短命signed URLとして都度発行する
- local開発ではMinIOを使う

direct uploadとAPI proxyのどちらを使うかはPhase 1で決める。

## 検討した代替案

- **AWS S3**: 実績は最も豊富。ただし転送量に課金されるため、再生回数が費用に直結する
- **DBにURLを保存する**: 実装は単純だが、移行時に全行の書き換えが必要になり、公開URLの存在を前提にしやすい

## 影響

- S3系serviceへの移行は、objectのcopyと設定変更でおおむね済む
- 本番と同じcodeをlocalのMinIOで確認できる
- **弱点**: S3互換はprovider間で完全に同一ではない。使う機能はprivate object、signed URL、Range requestを中心に絞る
