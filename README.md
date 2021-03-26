# wepo

Webhook に JSON を POST するやつ

## 使用法

- 環境変数 WEPO_URL に Webhook の URL を設定
- 実行時に与えた引数または標準入力の値が content に渡される

```shell
# arg
./wepo.sh example

# stdin
cat example.txt | ./wepo.sh
```
