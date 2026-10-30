# ADR 0004: 認証はFirebase Authenticationに任せ、認可はGoで実装する

- Status: Accepted
- Date: 2026-10-30

## 背景

BeatmakerのloginにはAuth providerが必要である。認証の事故は影響が大きく、password管理やreset処理を自前で持つ負担も大きい。
auth providerは乗り換えが最も難しい部品である。

## 決定

- 認証（誰か）はFirebase Authenticationに任せる
- Go APIはFirebaseのID tokenをserver側で検証する
- 認可（何にアクセスできるか）はGoで実装する
- DBの`users`に自前のUUIDを発行し、providerのIDは`users.auth_subject`にのみ保存する
- 業務tableは`users.id`を参照し、providerのIDを持たない

具体的なlogin方式（Google / mailなど）はPhase 1で決める。

## 検討した代替案

- **Clerk**: Next.jsとの相性は最も良い。ただしuser dataがprovider側に集中し、料金も増えやすい
- **Supabase Auth**: browserからDBを直接扱う設計が前提で、ADR 0001の境界と噛み合わない
- **Goで自前実装**: 学びは大きいが、securityの責任が重くMVPの範囲を超える
- **Auth.js**: 認証がNext.js側に置かれ、GoがNext.jsを信頼する構造になる
- **`owner_user_id`にproviderのIDを直接入れる**: 乗り換え時に全tableの更新が必要になる

## 影響

- Cloud Runと同じGCPで管理できる
- provider変更の影響は主に`users`とtoken検証の1か所に閉じる
- **弱点**: Next.jsのserver renderingでlogin状態を扱う場合は追加の工夫が必要
- 無料枠はplanによって異なる。2026-10時点の公式情報にも記載の差があるため、Phase 1の利用開始時に再確認する
