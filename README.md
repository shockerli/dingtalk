# dingtalk

> 钉钉群机器人SDK


## 安装

```go
go get github.com/shockerli/dingtalk
```


## 使用

### 初始化配置

```go
var robot = dingtalk.NewRobotCustom()
robot.SetWebhook("your_robot_webhook")
robot.SetSecret("your_secret") // 可选
```

### Text

```go
robot.SendText("TEST: Text")
```

### AtAll

```go
robot.SendText("TEST: Text&AtAll", robot.AtAll())
```

### AtMobiles

```go
robot.SendText("TEST: Text&AtMobiles", robot.AtMobiles("19900001111"))
```

### Link

```go
robot.SendLink(
    "TEST: Link",
    "link content",
    "https://github.com/shockerli",
    "https://www.wangbase.com/blogimg/asset/202101/bg2021011601.jpg",
)
```

### Markdown

```go
robot.SendMarkdown("TEST: Markdown", markdown)

robot.SendMarkdown("TEST: Markdown&AtAll", markdown, robot.AtAll())

robot.SendMarkdown("TEST: Markdown&AtMobiles", markdown, robot.AtMobiles("19900001111"))
```

### ActionCard

```go
robot.SendActionCard(
    "TEST: ActionCard&SingleCard",
    "SingleCard content",
    robot.SingleCard("阅读全文", "https://github.com/shockerli"),
)

robot.SendActionCard(
    "TEST: ActionCard&MultiCard",
    "MultiCard content",
    robot.MultiCard("内容不错", "https://github.com/shockerli"),
    robot.MultiCard("不感兴趣", "https://github.com/shockerli"),
)

robot.SendActionCard(
    "TEST: ActionCard&BtnOrientation",
    "BtnOrientation content",
    robot.MultiCard("内容不错", "https://github.com/shockerli"),
    robot.MultiCard("不感兴趣", "https://github.com/shockerli"),
    robot.BtnOrientation("0"),
)

robot.SendActionCard(
    "TEST: ActionCard&Image",
    "![xxx](https://www.wangbase.com/blogimg/asset/202101/bg2021011601.jpg)\n\nSingleCard content with image",
    robot.SingleCard("阅读全文", "https://github.com/shockerli"),
)

robot.SendActionCard(
    "TEST: ActionCard&HideAvatar",
    "24565\n\n![xxx](https://www.wangbase.com/blogimg/asset/202101/bg2021011601.jpg)\n\nSingleCard content with image",
    robot.SingleCard("阅读全文", "https://github.com/shockerli"),
    robot.HideAvatar("1"),
)
```

### FeedCard

```go
robot.SendFeedCard(
    robot.FeedCard("3月15日起，Chromium 不能再调用谷歌 API", "https://bodhi.fedoraproject.org/updates/FEDORA-2021-48866282e5%29", "https://www.wangbase.com/blogimg/asset/202101/bg2021012506.jpg"),
    robot.FeedCard("考古学家在英国发现两枚11世纪北宋时期的中国硬币", "https://www.caitlingreen.org/2020/12/another-medieval-chinese-coin-from-england.html", "https://www.wangbase.com/blogimg/asset/202101/bg2021012208.jpg"),
)
```

### Outgoing

```go
// 获取HTTP请求Body
var contents = getRequestBody()

// 解析Outgoing内容
og, err := robot.ParseOutgoing(contents)
if err != nil {
    // ...
}

// 自定义业务逻辑，生成响应的Text消息内容
var res = doSomeThing(og)

// 发送回复消息
err = robot.SendText(res, robot.WithOutgoing(og))
if err != nil {
    // ...
}
```


## 获取群机器人Token

- 选择自定义机器人

  ![选择自定义机器人](assets/robot-select.png)

- 设置机器人

  ![设置机器人](assets/robot-setting.png)

- 获取Webhook&token

  ![获取Webhook](assets/robot-token.png)



## 测试

1. 打开 `robot_test.go` 文件，修改 `your_robot_webhook` 和 `your_secret`；
2. 运行单元测试 `go test -v *_test.go`；
