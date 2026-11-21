# ADR 0014: BPM / key / tagsは範囲と固定候補で入力する

- Status: Provisional（tagsの候補listはユーザー確認待ち）
- Date: 2026-11-21

## 背景

PRDで、BPM、key、tagsは任意とし、自由入力は主にせず候補から選ぶことになっている。候補listと複数選択のルールは未確定だった。

## 決定

- **BPM**: 整数の数値入力とし、40〜250を許可する
- **key**: 24 key（12音 × major / minor）の固定listから選ぶ。表記はシャープの短縮形に統一し、DBにもこのcanonicalな値だけを保存する
  - major: `C` `C#` `D` `D#` `E` `F` `F#` `G` `G#` `A` `A#` `B`
  - minor: `Cm` `C#m` `Dm` `D#m` `Em` `Fm` `F#m` `Gm` `G#m` `Am` `A#m` `Bm`
- **tags**: 全ユーザー共通の固定候補（seed data）だけを使う。1曲あたり最大5個とする。userが独自tagを作る機能はpost-MVPで検討する
- tagsは`genre`と`mood`の2分類とする。初期候補は次の案とし、確定前にユーザーが確認する

| category | 候補案 |
|---|---|
| genre | trap, drill, boom bap, R&B, lo-fi, afrobeats, pop, hip hop, jersey club, reggaeton, dancehall, house, rage, plugg, phonk |
| mood | dark, chill, emotional, hype, melodic, sad, aggressive, happy, dreamy, bouncy, epic, romantic, mysterious, energetic, smooth |

## 検討した代替案

- **BPMを候補listから選ぶ**: BPMは連続値で、listは長くなるだけになる
- **keyのフラット表記や併記**（`Db` / `C#`）: 同じkeyが2通りに保存され、絞り込みが割れる
- **user-defined tags**: 自由度は高い。ただし表記ゆれで絞り込みの精度が下がり、post-MVPのaudio analysisやTaste Memoryで使いにくくなる

## 影響

- 絞り込みの値がそろい、将来の推薦（ADR 0008）にも使いやすい
- **弱点**: 候補にないgenreやmoodは付けられない。利用者の要望を見て候補を追加する。追加はseed dataの変更で済み、削除は`tags.is_active`をfalseにする（既存Beatのtagは消さない）
- フラット表記に慣れたユーザーには、表示上の工夫（例：`C#m (D♭m)`の併記）を後で検討する
