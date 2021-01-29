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

// RobotCustom 群机器人-自定义
// @doc https://developers.dingtalk.com/document/app/custom-robot-access
type RobotCustom struct {
	apiURI      string
	accessToken string
	secret      string
}

// NewRobotCustom 实例化
func NewRobotCustom() *RobotCustom {
	return &RobotCustom{
		apiURI: "https://oapi.dingtalk.com/robot/send",
	}
}

// SendText 发送Text消息
func (r RobotCustom) SendText(content string, opts ...RobotOption) error {
	msg := &robotMsg{
		MsgType: msgTypeText,
		Text:    robotText{Content: content},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return r.send(msg)
}

// SendLink 发送Link消息
func (r RobotCustom) SendLink(title, text, msgURL, picURL string, opts ...RobotOption) error {
	msg := &robotMsg{
		MsgType: msgTypeLink,
		Link: robotLink{
			Title:      title,
			Text:       text,
			MessageURL: msgURL,
			PicURL:     picURL,
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return r.send(msg)
}

// SendMarkdown 发送Markdown消息
func (r RobotCustom) SendMarkdown(title, text string, opts ...RobotOption) error {
	msg := &robotMsg{
		MsgType: msgTypeMarkdown,
		Markdown: robotMarkdown{
			Title: title,
			Text:  text,
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return r.send(msg)
}

// SendActionCard 发送ActionCard消息
func (r RobotCustom) SendActionCard(title, text string, opts ...RobotOption) error {
	msg := &robotMsg{
		MsgType: msgTypeActionCard,
		ActionCard: robotActionCard{
			Title:          title,
			Text:           text,
			HideAvatar:     "0", // 默认展示
			BtnOrientation: "1", // 默认横向排列
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return r.send(msg)
}

// SendFeedCard 发送FeedCard消息
func (r RobotCustom) SendFeedCard(opts ...RobotOption) error {
	msg := &robotMsg{
		MsgType: msgTypeFeedCard,
		FeedCard: robotFeedCard{
			Links: []robotFeedCardLink{},
		},
	}
	for _, opt := range opts {
		opt(msg)
	}
	return r.send(msg)
}

// 发送消息
func (r RobotCustom) send(msg interface{}) error {
	v, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	var value = make(url.Values)
	value.Set("access_token", r.accessToken)
	if r.secret != "" {
		t := time.Now().UnixNano() / 1e6
		value.Set("timestamp", fmt.Sprintf("%d", t))
		value.Set("sign", r.sign(t, r.secret))
	}

	data, err := request(r.apiURI+"?"+value.Encode(), v)
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

// 签名算法
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

// 消息@人的设置
// [NOTICE] 仅针对Text/Link/Markdown类型有效
type robotAt struct {
	AtMobiles []string `json:"atMobiles,omitempty"` // 被@人的手机号
	IsAtAll   bool     `json:"isAtAll,omitempty"`   // 是否@所有人
}

// 消息类型: Text
type robotText struct {
	Content string `json:"content"` // 消息内容
}

// 消息类型: Link
type robotLink struct {
	Title      string `json:"title"`            // 消息标题
	Text       string `json:"text"`             // 消息内容，如果太长只会部分展示
	MessageURL string `json:"messageUrl"`       // 点击消息跳转的UR
	PicURL     string `json:"picUrl,omitempty"` // 图片URL
}

// 消息类型: Markdown
type robotMarkdown struct {
	Title string `json:"title"` // 首屏会话透出的展示内容
	Text  string `json:"text"`  // Markdown格式的消息
}

// 消息类型: ActionCard
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

// 消息类型: FeedCard
type robotFeedCard struct {
	Links []robotFeedCardLink `json:"links"`
}

type robotFeedCardLink struct {
	Title      string `json:"title"`
	MessageURL string `json:"messageURL"` // 点击单条信息到跳转链接
	PicURL     string `json:"picURL"`     // 单条信息后面图片的URL
}

// RobotOption 群机器人-消息配置项
type RobotOption func(*robotMsg)

// AtAll 设置是否@所有人
// [NOTICE] 仅针对Text/Link/Markdown类型有效
func (r RobotCustom) AtAll() RobotOption {
	return func(msg *robotMsg) {
		msg.At.IsAtAll = true
	}
}

// AtMobiles 设置@人的手机号
// [NOTICE] 仅针对Text/Link/Markdown类型有效
func (r RobotCustom) AtMobiles(m ...string) RobotOption {
	return func(msg *robotMsg) {
		msg.At.AtMobiles = m
	}
}

// HideAvatar 隐藏缩略图
// 仅针对ActionCard类型
func (r RobotCustom) HideAvatar(v string) RobotOption {
	return func(msg *robotMsg) {
		msg.ActionCard.HideAvatar = v
	}
}

// BtnOrientation 按钮排列(0: 竖直排列, 1: 横向排列)
// 仅针对ActionCard类型
func (r RobotCustom) BtnOrientation(v string) RobotOption {
	return func(msg *robotMsg) {
		msg.ActionCard.BtnOrientation = v
	}
}

// SingleCard 整体调整配置
// 仅针对ActionCard类型
func (r RobotCustom) SingleCard(title, url string) RobotOption {
	return func(msg *robotMsg) {
		msg.ActionCard.SingleTitle = title
		msg.ActionCard.SingleURL = url
	}
}

// MultiCard 添加一个MultiCard项
// 仅针对ActionCard类型
func (r RobotCustom) MultiCard(title, url string) RobotOption {
	return func(msg *robotMsg) {
		msg.ActionCard.Btns = append(msg.ActionCard.Btns, robotActionCardBtn{
			Title:     title,
			ActionURL: url,
		})
	}
}

// FeedCard 添加一个FeedCard项
// 仅针对FeedCard类型
func (r RobotCustom) FeedCard(title, msgURL, picURL string) RobotOption {
	return func(msg *robotMsg) {
		msg.FeedCard.Links = append(msg.FeedCard.Links, robotFeedCardLink{
			Title:      title,
			MessageURL: msgURL,
			PicURL:     picURL,
		})
	}
}
