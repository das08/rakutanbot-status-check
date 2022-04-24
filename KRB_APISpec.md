# rakutanbot-status-check

## KRB(kuRakutanBot)へのリクエスト
### エンドポイント
- [POST] https://xxx.yyy/callback

### Header
- Content-Type: application/json
- X-Line-Signature: { Bodyを`LINE_CHANNEL_SECRET`でsha256ハッシュしたもの }

### Body
- [mock.json](https://github.com/das08/rakutanbot-status-check/blob/master/request/mock.json)
- `events[0].source.userId`はKRBで設定されているUID

## KRB(kuRakutanBot)のレスポンス


### Flexメッセージを返すとき
```json
{
    "Status": 2000,
    "Flex": [
        {
            "FlexContainer": {
                "type": "bubble",
                "...":0,
            }
        }
    ]
}
```

### Textメッセージを返すとき
- `Status`がエラーコードの場合は必ず`Text`フィールドが存在する
```json
{
    "Status": 2000,
    "Text": "メッセージが送信されました"
}
```

## レスポンスコード
```go
const (
	KRBSuccess            KRBStatus = 2000
	KRBDatabaseError      KRBStatus = 4000
	KRBOmikujiError       KRBStatus = 4000
	KRBGetFavError        KRBStatus = 4003
	KRBInsertFavError     KRBStatus = 4004
	KRBDeleteFavError     KRBStatus = 4005
	KRBGetLecIDError      KRBStatus = 4006
	KRBGetLecNameError    KRBStatus = 4007
	KRBGetUidError        KRBStatus = 4008
	KRBVerifyCodeGenError KRBStatus = 4009
	KRBVerifyCodeDelError KRBStatus = 4010
)
```
