# wepo

Webhook(Discord) に POST するやつ

## 使用法

- config.ini を生成する
  - `cp config.example.ini config.ini`
- config.ini 内の `webhook_url` に、Webhook の URL を設定する
  - キー(e.g. `[addr1]`)を追加することで、複数の宛先を設定可能
- 実行時に与えた引数または標準入力の値が post される

```shell
# arg
wepo example

# stdin
cat example.txt | wepo

# other address
wepo -a addr1 example
```
