package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"time"
)

// RobotCustom 群机器人-自定义
// @doc https://developers.dingtalk.com/document/app/custom-robot-access
type RobotCustom struct {
	webhook string // 例: https://oapi.dingtalk.com/robot/send?access_token=xxx
	secret  string // (可选)例: SEC8a9fc6f36f447d7c497f8c8e08accde4c49b4b5a366fa3903f47e250d6746979
}

// NewRobotCustom 实例化
func NewRobotCustom() *RobotCustom {
	return &RobotCustom{}
}

// SetWebhook 设置Token
func (rc *RobotCustom) SetWebhook(t string) *RobotCustom {
	rc.webhook = t
	return rc
}

// SetSecret 设置Secret
func (rc *RobotCustom) SetSecret(s string) *RobotCustom {
	rc.secret = s
	return rc
}

// SendText 发送Text消息
// 示例:
// 	robot.SendText("TEST: Text")
// 	robot.SendText("TEST: Text&AtAll", robot.AtAll())
// 	robot.SendText("TEST: Text&AtMobiles", robot.AtMobiles("19900001111"))
func (rc *RobotCustom) SendText(content string, opts ...RobotOption) error {
	msg := &robotMsg{
		MsgType: msgTypeText,
		Text:    &robotText{Content: content},
	}

	return rc.send(msg, opts...)
}

// SendLink 发送Link消息
// 示例:
// 	robot.SendLink(
//		"TEST: Link",
//		"link content",
//		"https://github.com/shockerli",
//		"https://www.wangbase.com/blogimg/asset/202101/bg2021011601.jpg",
//	)
func (rc *RobotCustom) SendLink(title, text, msgURL, picURL string, opts ...RobotOption) error {
	msg := &robotMsg{
		MsgType: msgTypeLink,
		Link: &robotLink{
			Title:      title,
			Text:       text,
			MessageURL: msgURL,
			PicURL:     picURL,
		},
	}

	return rc.send(msg, opts...)
}

// SendMarkdown 发送Markdown消息
// 使用示例:
// 	robot.SendMarkdown("TEST: Markdown", markdown)
// 	robot.SendMarkdown("TEST: Markdown&AtAll", markdown, robot.AtAll())
// 	robot.SendMarkdown("TEST: Markdown&AtMobiles", markdown, robot.AtMobiles("19900001111"))
func (rc *RobotCustom) SendMarkdown(title, text string, opts ...RobotOption) error {
	msg := &robotMsg{
		MsgType: msgTypeMarkdown,
		Markdown: &robotMarkdown{
			Title: title,
			Text:  text,
		},
	}

	return rc.send(msg, opts...)
}

// SendActionCard 发送ActionCard消息
// 示例:
// 	robot.SendActionCard(
//		"TEST: ActionCard&SingleCard",
//		"SingleCard content",
//		robot.SingleCard("阅读全文", "https://github.com/shockerli"),
//	)
func (rc *RobotCustom) SendActionCard(title, text string, opts ...RobotOption) error {
	msg := &robotMsg{
		MsgType: msgTypeActionCard,
		ActionCard: &robotActionCard{
			Title:          title,
			Text:           text,
			HideAvatar:     "0", // 默认展示
			BtnOrientation: "1", // 默认横向排列
		},
	}

	return rc.send(msg, opts...)
}

// SendFeedCard 发送FeedCard消息
func (rc *RobotCustom) SendFeedCard(opts ...RobotOption) error {
	msg := &robotMsg{
		MsgType: msgTypeFeedCard,
		FeedCard: &robotFeedCard{
			Links: []robotFeedCardLink{},
		},
	}

	return rc.send(msg, opts...)
}

// 发送消息
func (rc *RobotCustom) send(msg *robotMsg, opts ...RobotOption) error {
	// options
	for _, opt := range opts {
		opt(msg)
	}

	v, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	var api = rc.webhook
	var value = make(url.Values)
	var now = time.Now().UnixNano() / 1e6 // 毫秒
	if msg.outgoing.SessionWebhook != "" {
		if msg.outgoing.SessionWebhookExpiredTime < now {
			return fmt.Errorf("SessionWebhookExpiredTime is expired")
		}
		api = msg.outgoing.SessionWebhook
	} else if rc.secret != "" {
		value.Set("timestamp", fmt.Sprintf("%d", now))
		value.Set("sign", rc.sign(now, rc.secret))
		api = rc.webhook + "&" + value.Encode()
	}

	// 请求接口
	data, err := request(api, v)
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
func (*RobotCustom) sign(ts int64, secret string) string {
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
	MsgType    string           `json:"msgtype"` // 消息类型
	At         *robotAt         `json:"at,omitempty"`
	Text       *robotText       `json:"text,omitempty"`
	Link       *robotLink       `json:"link,omitempty"`
	Markdown   *robotMarkdown   `json:"markdown,omitempty"`
	ActionCard *robotActionCard `json:"actionCard,omitempty"`
	FeedCard   *robotFeedCard   `json:"feedCard,omitempty"`
	outgoing   RobotOutgoing
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
// 用于修改msg配置项
type RobotOption func(*robotMsg)

// AtAll 设置是否@所有人
// [NOTICE] 仅针对Text/Markdown类型有效
// 示例:
// 	robot.SendMarkdown("TEST: Markdown&AtAll", markdown, robot.AtAll())
func (rc *RobotCustom) AtAll() RobotOption {
	return func(msg *robotMsg) {
		if msg.MsgType != msgTypeText && msg.MsgType != msgTypeMarkdown {
			return
		}
		if msg.At == nil {
			msg.At = &robotAt{}
		}
		msg.At.IsAtAll = true
	}
}

// AtMobiles 设置@人的手机号
// [NOTICE] 仅针对Text/Markdown类型有效
// 示例:
// 	robot.SendMarkdown("TEST: Markdown&AtMobiles", markdown, robot.AtMobiles("19900001111"))
func (rc *RobotCustom) AtMobiles(m ...string) RobotOption {
	return func(msg *robotMsg) {
		if msg.MsgType != msgTypeText && msg.MsgType != msgTypeMarkdown {
			return
		}
		if msg.At == nil {
			msg.At = &robotAt{}
		}
		msg.At.AtMobiles = m
	}
}

// HideAvatar 隐藏头像(0-显示, 1-隐藏, 默认0)
// 仅针对ActionCard类型
// 示例:
//	robot.SendActionCard(
//		"TEST: ActionCard&HideAvatar",
//		"24565\n\n![xxx](https://www.wangbase.com/blogimg/asset/202101/bg2021011601.jpg)\n\nSingleCard content with image",
//		robot.SingleCard("阅读全文", "https://github.com/shockerli"),
//		robot.HideAvatar("1"),
// 	)
func (rc *RobotCustom) HideAvatar(v string) RobotOption {
	return func(msg *robotMsg) {
		if msg.MsgType != msgTypeActionCard {
			return
		}
		msg.ActionCard.HideAvatar = v
	}
}

// BtnOrientation 按钮排列(0-竖直排列, 1-横向排列, 默认1)
// 仅针对ActionCard类型
// 示例:
// 	robot.SendActionCard(
//		"TEST: ActionCard&BtnOrientation",
//		"BtnOrientation content",
//		robot.MultiCard("内容不错", "https://github.com/shockerli"),
//		robot.MultiCard("不感兴趣", "https://github.com/shockerli"),
//		robot.BtnOrientation("0"),
//	)
func (rc *RobotCustom) BtnOrientation(v string) RobotOption {
	return func(msg *robotMsg) {
		if msg.MsgType != msgTypeActionCard {
			return
		}
		msg.ActionCard.BtnOrientation = v
	}
}

// SingleCard 整体跳转配置
// 仅针对ActionCard类型
// 示例:
// 	robot.SendActionCard(
//		"TEST: ActionCard&SingleCard",
//		"SingleCard content",
//		robot.SingleCard("阅读全文", "https://github.com/shockerli"),
//	)
func (rc *RobotCustom) SingleCard(title, url string) RobotOption {
	return func(msg *robotMsg) {
		if msg.MsgType != msgTypeActionCard {
			return
		}
		msg.ActionCard.SingleTitle = title
		msg.ActionCard.SingleURL = url
	}
}

// MultiCard 添加一个MultiCard项
// 仅针对ActionCard类型
// 示例:¬
// 	robot.SendActionCard(
//		"TEST: ActionCard&MultiCard",
//		"MultiCard content",
//		robot.MultiCard("内容不错", "https://github.com/shockerli"),
//		robot.MultiCard("不感兴趣", "https://github.com/shockerli"),
//	)
func (rc *RobotCustom) MultiCard(title, url string) RobotOption {
	return func(msg *robotMsg) {
		if msg.MsgType != msgTypeActionCard {
			return
		}
		msg.ActionCard.Btns = append(msg.ActionCard.Btns, robotActionCardBtn{
			Title:     title,
			ActionURL: url,
		})
	}
}

// FeedCard 添加一个FeedCard项
// 仅针对FeedCard类型
func (rc *RobotCustom) FeedCard(title, msgURL, picURL string) RobotOption {
	return func(msg *robotMsg) {
		if msg.MsgType != msgTypeFeedCard {
			return
		}
		msg.FeedCard.Links = append(msg.FeedCard.Links, robotFeedCardLink{
			Title:      title,
			MessageURL: msgURL,
			PicURL:     picURL,
		})
	}
}

// WithOutgoing 通过Outgoing的临时消息接口发送
// 示例:
// 	og, err := robot.ParseOutgoing(bytes.NewBufferString(callbackBody))
//	err = robot.SendText("callback", robot.WithOutgoing(og))
func (rc *RobotCustom) WithOutgoing(og RobotOutgoing) RobotOption {
	return func(msg *robotMsg) {
		msg.outgoing = og
	}
}

// ParseOutgoing 解析Outgoing消息体
// 示例:
// 	robot.ParseOutgoing(callbackBody)
func (rc *RobotCustom) ParseOutgoing(r io.Reader) (og RobotOutgoing, err error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, &og)
	return
}

// RobotOutgoing Outgoing回调消息体
// {
//    "conversationId": "ciddz7nmHDaX/7Niz+Gb5VVrw==",
//    "sceneGroupCode": "project",
//    "atUsers": [
//        {
//            "dingtalkId": "$:LWCP_v1:$0sIVIuw1HvQQ5gRAtFWzypo0+T1TgPOP"
//        },
//        {
//            "dingtalkId": "$:LWCP_v1:$I3cyfTzrws4nCbY289cXbKCVcdd1wize"
//        }
//    ],
//    "chatbotUserId": "$:LWCP_v1:$I3cyfTzrws4nCbY289cXbKCVcdd1wize",
//    "msgId": "msgaKcioIqERkELm2T8TlE9CA==",
//    "senderNick": "Jioby",
//    "isAdmin": false,
//    "sessionWebhookExpiredTime": 1612178396066,
//    "createAt": 1612172996026,
//    "conversationType": "2",
//    "senderId": "$:LWCP_v1:$deZJcSfMzexC2YK+oLkk1g==",
//    "conversationTitle": "xxx",
//    "isInAtList": true,
//    "sessionWebhook": "https://oapi.dingtalk.com/robot/sendBySession?session=eb18e18e8669b0a3cd7dff1388fe5e6a",
//    "text": {
//        "content": "  哈哈哈"
//    },
//    "msgtype": "text"
// }
type RobotOutgoing struct {
	// 被@人的信息
	AtUsers []struct {
		DingTalkID string `json:"dingtalkId"` // 加密的人员ID
	} `json:"atUsers"`
	ChatBotUserID             string    `json:"chatbotUserId"`             // 加密的机器人ID
	ConversationID            string    `json:"conversationId"`            // 加密的会话ID
	ConversationTitle         string    `json:"conversationTitle"`         // 会话标题(群聊时才有，即群名)
	ConversationType          string    `json:"conversationType"`          // 1-单聊、2-群聊
	CreateAt                  int64     `json:"createAt"`                  // 消息的时间戳，单位ms
	IsAdmin                   bool      `json:"isAdmin"`                   // 是否为管理员发送的消息
	IsInAtList                bool      `json:"isInAtList"`                //
	MsgID                     string    `json:"msgId"`                     // 加密的消息ID
	MsgType                   string    `json:"msgtype"`                   // 消息类型: 目前只支持Text
	SceneGroupCode            string    `json:"sceneGroupCode"`            // 群组场景类型Code
	SenderID                  string    `json:"senderId"`                  // 加密的发送者ID
	SenderNick                string    `json:"senderNick"`                // 发送者昵称
	SessionWebhook            string    `json:"sessionWebhook"`            // 临时的发送消息接口
	SessionWebhookExpiredTime int64     `json:"sessionWebhookExpiredTime"` // SessionWebhook可用的有效截止时间
	Text                      robotText `json:"text"`                      // Text类型的消息体
}
