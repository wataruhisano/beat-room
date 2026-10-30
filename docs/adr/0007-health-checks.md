# ADR 0007: health checkの役割を分ける

- Status: Accepted
- Date: 2026-10-30

## 背景

依存先を確認するhealth checkを自動実行すると、一時的な障害でAPIの再起動が連鎖する。またNeonのscale to zeroが効かなくなり、無料枠を消費する。

## 決定

| Endpoint | 確認内容 | 用途 |
|---|---|---|
| `/healthz` | processの応答のみ。依存先は確認しない | Cloud Runの自動check、WebからのAPI確認 |
| `/readyz` | DB接続 | deploy直後の確認と手動確認。定期実行しない |

object storageはどちらでも確認しない。障害はuploadと再生のerrorとして検知する。

## 検討した代替案

- **`/readyz`にstorage確認も含める**: storageの一時障害で、loginやLibraryなど無関係な機能まで止まる
- **Cloud Runの自動checkでDBを確認する**: 上記のとおり、再起動の連鎖と費用増につながる

## 影響

- 部分的な障害が全体の障害に広がりにくい
- **弱点**: DB障害を`/healthz`では検知できないため、error logとalertで補う
