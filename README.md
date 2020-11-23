# order-book-searcher

oanda のオーダーブックを分析するためのコマンドラインツールです。過去のオーダー情報から指定した条件に合致する日時のデータを検索し、CSV 出力します。
オーダーブックの情報取得には oanda API v20 を使用します。詳細はドキュメントを確認してください。

## Usage

### argument

下記のオプションを使用して検索条件を指定します。

| 引数名 | 詳細 |
| --- | --- |
| oanda-key (必須)| oanda の api key を指定します。|
| period (必須)| 集計期間を指定します |
| instrument (必須)| 通貨を指定します |
| stop-order | 逆指値注文の比率を指定します。複数指定した場合はその数値が連続した価格帯が存在している箇所を検索します。 |
| limit-order | 指値注文の比率の下限を指定します。複数指定した場合はその数値が連続した価格帯が存在している箇所を検索します。 |
| losing-position | 損失が出ているポジションの下限比率を指定します。複数指定した場合はその数値が連続した価格帯が存在している箇所を検索します。 |
| profiting-position | 利益が出ているポジションの下限比率を指定します。複数指定した場合はその数値が連続した価格帯が存在している箇所を検索します。 |
| jp | Excel 最適化を行います |
| loc | date-time カラムの time location を指定します。 UTC, JST, MT4 が選択可能です。 |

ex:
```
go run . -oanda-key xxxxxxx -period 2020/10/01-2020/10/04 -instrument EUR_GBP -stop-order 0.5-1.0 -jp -loc MT4

```

### output

| ヘッダー | 詳細 |
| date-time | order book の日時 |
| price | 当時の価格 |
| price-range-{:i} | ヒットした価格帯 |
| short-order-{:i} | ヒットした価格帯の売り注文比率 | 
| long-order-{:i} | ヒットした価格帯の買いり注文比率 |
| short-position-{:i} | ヒットした価格帯の売りポジション比率 |
| long-position-{:i} | ヒットした価格帯の買いポジション比率 |

連続した価格帯での検索を行った場合には、現在価格に近い方から番号付けされ、 {:i} と置き換えられます。





