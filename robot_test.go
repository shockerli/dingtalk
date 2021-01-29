package dingtalk_test

import (
	"testing"

	"github.com/shockerli/dingtalk"
)

var robot *dingtalk.RobotCustom

func init() {
	robot = dingtalk.NewRobotCustom()
	robot.SetAccessToken("your_access_token")
	robot.SetSecret("your_secret")
}

func TestRobotCustom_SendText(t *testing.T) {
	// 根据Secret验证
	if err := robot.SendText("Type: Text, test"); err != nil {
		t.Errorf("SendText() error = %v", err)
	}

	// AtAll
	if err := robot.SendText("Type: Text&AtAll", robot.AtAll()); err != nil {
		t.Errorf("SendText() && AtAll() error = %v", err)
	}

	// AtMobiles
	if err := robot.SendText("Type: Text&AtMobiles", robot.AtMobiles("19900001111")); err != nil {
		t.Errorf("SendText() && AtMobiles() error = %v", err)
	}
}

func TestRobotCustom_SendLink(t *testing.T) {
	// 根据Secret验证
	if err := robot.SendLink("Type: Link, test", "link content", "https://github.com/shockerli/dingtalk", "https://www.wangbase.com/blogimg/asset/202101/bg2021011601.jpg"); err != nil {
		t.Errorf("SendLink() error = %v", err)
	}
}
