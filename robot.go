package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// 群机器人-自定义
// @doc https://developers.dingtalk.com/document/app/custom-robot-access
type RobotCustom struct {
	apiUri      string
	accessToken string
	secret      string
}

func NewRobotCustom() *RobotCustom {
	return &RobotCustom{
		apiUri: "https://oapi.dingtalk.com/robot/send",
	}
}

func (r RobotCustom) send(msg interface{}) error {
	v, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	var value url.Values
	value.Set("access_token", r.accessToken)
	if r.secret != "" {
		t := time.Now().UnixNano() / 1e6
		value.Set("timestamp", fmt.Sprintf("%d", t))
		value.Set("sign", r.sign(t, r.secret))
	}

	data, err := request(r.apiUri+"?"+value.Encode(), v)
	if err != nil {
		return err
	}

	var response struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}
	if response.ErrCode != 0 {
		return fmt.Errorf("群机器人消息发送失败: %v", response.ErrMsg)
	}

	return nil
}

func (RobotCustom) sign(ts int64, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(fmt.Sprintf("%d\n%s", ts, secret)))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// robot message definition

// 机器人消息类型
const (
	msgTypeText       = "text"
	msgTypeLink       = "link"
	msgTypeMarkdown   = "markdown"
	msgTypeActionCard = "actionCard"
	msgTypeFeedCard   = "feedCard"
)

// 机器人消息结构
type robotMsg struct {
	MsgType    string          `json:"msgtype"` // 消息类型
	At         robotAt         `json:"at,omitempty"`
	Text       robotText       `json:"text,omitempty"`
	Link       robotLink       `json:"link"`
	Markdown   robotMarkdown   `json:"markdown"`
	ActionCard robotActionCard `json:"actionCard,omitempty"`
	FeedCard   robotFeedCard   `json:"feedCard,omitempty"`
}

// 类型: Text
type robotText struct {
	Content string `json:"content"` // 消息内容
}

// @人
type robotAt struct {
	AtMobiles []string `json:"atMobiles,omitempty"` // 被@人的手机号
	IsAtAll   bool     `json:"isAtAll,omitempty"`   // 是否@所有人
}

// 类型: Link
type robotLink struct {
	Title      string `json:"title"`            // 消息标题
	Text       string `json:"text"`             // 消息内容，如果太长只会部分展示
	MessageURL string `json:"messageUrl"`       // 点击消息跳转的UR
	PicURL     string `json:"picUrl,omitempty"` // 图片URL
}

// 类型: Markdown
type robotMarkdown struct {
	Title string `json:"title"` // 首屏会话透出的展示内容
	Text  string `json:"text"`  // Markdown格式的消息
}

// 类型: ActionCard
// * 整体跳转
// * 独立跳转
// [NOTICE]设置singleTitle和singleURL后，btns无效
type robotActionCard struct {
	Title          string               `json:"title"`                    // 首屏会话透出的展示内容
	Text           string               `json:"text"`                     // Markdown格式的消息
	SingleTitle    string               `json:"singleTitle,omitempty"`    // 单个按钮的标题
	SingleURL      string               `json:"singleURL,omitempty"`      // 点击singleTitle按钮触发的URL
	HideAvatar     string               `json:"hideAvatar,omitempty"`     // 0：显示图片，1：隐藏图片
	BtnOrientation string               `json:"btnOrientation,omitempty"` // 0：按钮竖直排列，1：按钮横向排列
	Btns           []robotActionCardBtn `json:"btns,omitempty"`           // 独立跳转的按钮列表
}

type robotActionCardBtn struct {
	Title     string `json:"title"`     // 按钮标题
	ActionURL string `json:"actionURL"` // 点击按钮触发的URL
}

// 类型: FeedCard
type robotFeedCard struct {
	Links []robotFeedCardLink `json:"links"`
}

type robotFeedCardLink struct {
	Title      string `json:"title"`
	MessageURL string `json:"messageURL"` // 点击单条信息到跳转链接
	PicURL     string `json:"picURL"`     // 单条信息后面图片的URL
}
