package dingtalk_test

import (
	"bytes"
	"testing"

	"github.com/shockerli/dingtalk"
)

var robot *dingtalk.RobotCustom

func init() {
	robot = dingtalk.NewRobotCustom()
	robot.SetWebhook("your_robot_webhook")
	robot.SetSecret("your_secret")
}

func TestRobotCustom_SendText(t *testing.T) {
	// 根据Secret验证
	if err := robot.SendText("TEST: Text"); err != nil {
		t.Errorf("SendText() error = %v", err)
	}

	// AtAll
	if err := robot.SendText("TEST: Text&AtAll", robot.AtAll()); err != nil {
		t.Errorf("SendText() && AtAll() error = %v", err)
	}

	// AtMobiles
	if err := robot.SendText("TEST: Text&AtMobiles", robot.AtMobiles("19900001111")); err != nil {
		t.Errorf("SendText() && AtMobiles() error = %v", err)
	}
}

func TestRobotCustom_SendLink(t *testing.T) {
	if err := robot.SendLink("TEST: Link", "link content", "https://github.com/shockerli", "https://www.wangbase.com/blogimg/asset/202101/bg2021011601.jpg"); err != nil {
		t.Errorf("SendLink() error = %v", err)
	}
}

func TestRobotCustom_SendMarkdown(t *testing.T) {
	var markdown = `
## 安装
> go get github.com/shockerli/dingtalk

## 使用
> var robot = dingtalk.NewRobotCustom()
> robot.SetAccessToken("your_access_token")
> robot.SetSecret("your_secret") // 可选
`

	if err := robot.SendMarkdown("TEST: Markdown", markdown); err != nil {
		t.Errorf("SendMarkdown() error = %v", err)
	}

	if err := robot.SendMarkdown("TEST: Markdown&AtAll", markdown, robot.AtAll()); err != nil {
		t.Errorf("SendMarkdown() && AtAll() error = %v", err)
	}

	if err := robot.SendMarkdown("TEST: Markdown&AtMobiles", markdown, robot.AtMobiles("19900001111")); err != nil {
		t.Errorf("SendMarkdown() && AtMobiles() error = %v", err)
	}
}

func TestRobotCustom_SendActionCard(t *testing.T) {
	if err := robot.SendActionCard("TEST: ActionCard&SingleCard", "SingleCard content", robot.SingleCard("阅读全文", "https://github.com/shockerli")); err != nil {
		t.Errorf("SendActionCard() && SingleCard() error = %v", err)
	}

	if err := robot.SendActionCard(
		"TEST: ActionCard&MultiCard",
		"MultiCard content",
		robot.MultiCard("内容不错", "https://github.com/shockerli"),
		robot.MultiCard("不感兴趣", "https://github.com/shockerli"),
	); err != nil {
		t.Errorf("SendActionCard() && MultiCard() error = %v", err)
	}

	if err := robot.SendActionCard(
		"TEST: ActionCard&BtnOrientation",
		"BtnOrientation content",
		robot.MultiCard("内容不错", "https://github.com/shockerli"),
		robot.MultiCard("不感兴趣", "https://github.com/shockerli"),
		robot.BtnOrientation("0"),
	); err != nil {
		t.Errorf("SendActionCard() && BtnOrientation() error = %v", err)
	}

	if err := robot.SendActionCard(
		"TEST: ActionCard&Image",
		"![xxx](https://www.wangbase.com/blogimg/asset/202101/bg2021011601.jpg)\n\nSingleCard content with image",
		robot.SingleCard("阅读全文", "https://github.com/shockerli"),
	); err != nil {
		t.Errorf("SendActionCard() && SingleCard() error = %v", err)
	}

	if err := robot.SendActionCard(
		"TEST: ActionCard&HideAvatar",
		"24565\n\n![xxx](https://www.wangbase.com/blogimg/asset/202101/bg2021011601.jpg)\n\nSingleCard content with image",
		robot.SingleCard("阅读全文", "https://github.com/shockerli"),
		robot.HideAvatar("1"),
	); err != nil {
		t.Errorf("SendActionCard() && HideAvatar() error = %v", err)
	}
}

func TestRobotCustom_Outgoing(t *testing.T) {
	var callbackBody = `{"conversationId":"ciddz7nmHDaX/7Niz+Gb5VVrw==","sceneGroupCode":"project","atUsers":[{"dingtalkId":"$:LWCP_v1:$0sIVIuw1HvQQ5gRAtFWzypo0+T1TgPOP"},{"dingtalkId":"$:LWCP_v1:$I3cyfTzrws4nCbY289cXbKCVcdd1wize"}],"chatbotUserId":"$:LWCP_v1:$I3cyfTzrws4nCbY289cXbKCVcdd1wize","msgId":"msgaKcioIqERkELm2T8TlE9CA==","senderNick":"Jioby","isAdmin":false,"sessionWebhookExpiredTime":1612178396066,"createAt":1612172996026,"conversationType":"2","senderId":"$:LWCP_v1:$deZJcSfMzexC2YK+oLkk1g==","conversationTitle":"xxx","isInAtList":true,"sessionWebhook":"https://oapi.dingtalk.com/robot/sendBySession?session=eb18e18e8669b0a3cd7dff1388fe5e6a","text":{"content":"  哈哈哈"},"msgtype":"text"}`

	og, err := robot.ParseOutgoing(bytes.NewBufferString(callbackBody))
	if err != nil {
		t.Errorf("ParseOutgoing() error = %v", err)
	}
	err = robot.SendText("callback", robot.WithOutgoing(og))
	if err != nil {
		t.Errorf("WithOutgoing() error= %v", err)
	}
}
