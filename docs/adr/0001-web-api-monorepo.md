# ADR 0001: Next.js + Go API + PostgreSQLのmonorepoを採用する

- Status: Accepted
- Date: 2026-10-30

## 背景

BeatroomはBeatmakerの管理画面と、Artist向けのpublic Roomを持つ。tenant isolationとpublic token accessを後付けにしない必要がある。
開発者は1人で、実運用に使えるproductであることに加え、実務のWeb開発で一般的な構成に近いことも重視する。

## 決定

- WebはNext.js（TypeScript）、APIはGo、DBはPostgreSQLとする
- Next.jsの担当は表示に限る。DB accessと認可はすべてGo APIに置く
- 1つのrepositoryに`apps/web`と`apps/api`を置く（monorepo）

## 検討した代替案

- **Next.jsのみ（TS full-stack）**: 最速で部品も少ない。ただしDB accessと権限checkの境界が画面側のcodeと混ざりやすい
- **TS frontend + TS backend（Honoなど）**: 1言語で分離できる。Beatroomの処理はGoでも自然に書けるため、言語統一より学習中のGoを活かす方を選んだ
- **polyrepo**: チームが分かれた組織向け。1人開発ではAPIと画面の同時変更とCI管理の手間が増える

## 影響

- tenant isolationとpublic tokenの確認箇所をGoに集約できる
- 画面以外のclientが増えても同じAPIを使える
- 将来の推薦処理は、Python batchがPostgreSQLを読む形で追加できる
- **弱点**: 現在の規模では分離は過剰でもある。2言語の運用、CORS、認証tokenの受け渡しを扱う必要がある
- Next.jsのserver機能へロジックを置き始めると分離の意味が失われるため、reviewで確認する
