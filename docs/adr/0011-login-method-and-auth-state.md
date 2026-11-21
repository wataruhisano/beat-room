# ADR 0011: loginはGoogleのみとし、login状態はbrowserで扱う

- Status: Provisional
- Date: 2026-11-21

## 背景

ADR 0004でFirebase Authenticationの採用は決めたが、有効にするlogin方法と、Next.jsでlogin状態をどう扱うかはPhase 1に先送りしていた。
あわせて、`users`にemailと表示名を保存するか（DATA_MODEL）も、login方式に合わせて決める必要がある。

## 決定

- 有効にするlogin方法はGoogle sign-inのみとする
- login状態はbrowserのFirebase JS SDKが持つ。WebはGo APIを呼ぶたびにID tokenを`Authorization: Bearer`で送る
- ID tokenを検証するのはGo APIだけとする。Next.jsのserver側では認証を判定しない
- `users`にemailは保存しない
- `users.display_name`はnullableとし、初回loginでGoogleの表示名を入れる
- local開発とtestではFirebase Auth Emulatorを使う。Admin SDKは`FIREBASE_AUTH_EMULATOR_HOST`が設定されているとEmulatorへ接続する。この変数は本番で設定しない

## 検討した代替案

- **Google + email/password**: 誰でも登録できる。ただしpassword reset、email確認、同じemailでの重複accountの扱いが増える
- **Google + email link（magic link）**: passwordは不要。ただしメールの到達性と、別端末でlinkを開いた場合の扱いが必要になる
- **Firebase session cookie**: Next.jsのserver renderingでlogin状態を判定でき、初回描画が速い。ただしcookieの発行と失効、CSRF対策を持つことになり、Next.jsも認証に関わる
- **emailを`users`に保存する**: MVPにはemailを使う機能がない（Email通知はout of scope）。不要な個人情報を持たないほうがprivacy policyも単純になる
- **display_nameを必須にする**: onboardingの手順が1つ増える

## 影響

- login画面はGoogleのbuttonだけになり、password関連の画面が不要になる
- 認証の判断がGoの1か所に集まり、ADR 0001の境界と一致する
- **弱点**: Google accountを使いたくないBeatmakerは登録できない。保護された画面は、login状態の確認が終わるまで描画が一瞬遅れる
- 業務tableは`users.id`を参照するので、email系のloginを後から足しても既存dataへの影響はほとんどない。要望が出たら見直す
- Firebase Authの無料枠はSparkでもBlazeでも50K MAU（2026-11時点の公式pricing）。ADR 0004で求めた再確認はこれで済んだ
- Admin SDKとEmulatorの接続は公式docsにNode.jsの例しかないため、Go（firebase-admin-go）での動作は実装時に確認する
