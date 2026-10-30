# ADR 0005: 標準を優先してdependencyを選ぶ

- Status: Accepted
- Date: 2026-10-30

## 背景

Phaseが進むたびにlibraryを選ぶ場面がある。同じ基準で判断でき、理由を説明できる状態にしたい。

## 決定

### 選定基準

標準機能を優先する。libraryを追加するのは次を満たす時とする。

1. 大きな手間を省ける
2. 広く使われ、保守が続いている
3. 後で差し替えられる

業界での普及度は補助的な判断材料として考慮する。

### 運用

- versionはlock fileで固定する
- 更新は自動PRで受け、CIを通してから取り込む
- 新しいdependencyを足す理由をADRまたはPR説明に残す

### 採用するtool

| 用途 | 採用 | 理由 |
|---|---|---|
| Go router | 標準`net/http` | Go 1.22以降のroutingで足りる。middlewareの構成が辛くなればchiへ移れる |
| DB access | sqlc + pgx | SQLを直接書ける。`owner_user_id`を必須引数にしたqueryを生成できる |
| Migration | golang-migrate | 利用者が多い |
| JS package manager | pnpm | 宣言していないpackageを使えない |
| Styling | Tailwind CSS | 普及度を考慮して標準CSS Modulesより優先した。dark themeのtokenも定義しやすい |

## 検討した代替案

- **Gin / Echo**: 高機能で知名度も高いが、独自の書き方が増え、必要な機能に対して大きい
- **GORM**: 手軽だが、実行されるSQLが見えにくい
- **Atlas**: 宣言的で高機能だが、現時点では必要な範囲を超える
- **CSS Modules**: 依存を増やさず、後からの移行先として残る

## 影響

- query層で`owner_user_id`を必須にすることで、IDだけでresourceを取得するcodeを書きにくくなる（IDOR対策）
- public Roomはtoken hashで解決する別系統のqueryを用意する
- **弱点**: 標準機能を選ぶ分、一部のmiddlewareやvalidationは自分で用意する
