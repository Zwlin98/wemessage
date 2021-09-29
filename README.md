# WeMessage

企业微信消息推送部分，主要包括

+ 发送应用消息
+ 接收消息和事件 // 暂未完成

从而利用企业微信来完成 `ios` 上类似app `bark` 的功能。

## Usage

### SendMessage

```go
client := wemessage.NewClient(corpID, corpSecret)

textMsg := wemessage.TextMessage{
    ToUser:  "user",
    ToParty: "",
    ToTag:   "",
    MsgType: "text",
    AgentID: 100001,
    Text: struct {
        Content string `json:"content"`
    }{
        Content: "wemessage test message",
    },
    Safe:                   0,
    EnableIdTrans:          0,
    EnableDuplicateCheck:   0,
    DuplicateCheckInterval: 0,
}

_, err := wemessage.SendMessage(client, textMsg)
```
