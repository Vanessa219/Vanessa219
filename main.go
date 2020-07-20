package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/88250/gulu"
	"github.com/parnurzeal/gorequest"
)

var logger = gulu.Log.NewLogger(os.Stdout)

const (
	githubUserName = "Vanessa219"
	hacpaiUserName = "Vanessa"
)

func main() {
	result := map[string]interface{}{}
	response, data, errors := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Get("https://hacpai.com/api/v2/user/"+hacpaiUserName+"/events?size=8").Timeout(7*time.Second).
		Set("User-Agent", "Profile Bot; +https://github.com/"+githubUserName+"/"+githubUserName).EndStruct(&result)
	if nil != errors || http.StatusOK != response.StatusCode {
		logger.Fatalf("fetch events failed: %+v, %s", errors, data)
	}
	if 0 != result["code"].(float64) {
		logger.Fatalf("fetch events failed: %s", data)
	}

	buf := &bytes.Buffer{}
	buf.WriteString("\n\n")
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	updated := time.Now().In(cstSh).Format("2006-01-02 15:04:05")
	buf.WriteString("我的近期动态（点个 [Star](https://github.com/" + githubUserName + "/" + githubUserName + ") 将触发自动刷新，最近更新时间：`" + updated + "`）：\n\n")
	for _, event := range result["data"].([]interface{}) {
		evt := event.(map[string]interface{})
		operation := evt["operation"].(string)
		title := evt["title"].(string)
		typ := evt["type"].(string)
		var emoji string
		switch typ {
		case "article":
			emoji = "📝"
		case "comment":
			emoji = "💬"
		case "comment2":
			emoji = "🗣"
		case "breezemoon":
			emoji = "🌙"
			title = operation
		case "vote-article":
			emoji = "👍📝"
		case "vote-comment":
			emoji = "👍💬"
		case "vote-comment2":
			emoji = "👍🗣"
		case "vote-breezemoon":
			emoji = "👍🌙"
			title = operation
		case "reward-article":
			emoji = "💰📝"
		case "thank-article":
			emoji = "💗📝"
		case "thank-comment":
			emoji = "💗💬"
		case "accept-comment":
			emoji = "✅💬"
		case "thank-comment2":
			emoji = "💗🗣"
		case "thank-breezemoon":
			emoji = "💗🌙"
			title = operation
		case "follow-user":
			emoji = "⭐️👨‍💻"
		case "follow-tag":
			emoji = "⭐️🏷️"
		case "collect-article":
			emoji = "🗃📝"
		}

		url := evt["url"].(string)
		content := evt["content"].(string)
		buf.WriteString("* " + emoji + " [" + title + "](" + url + ")\n\n" + "  > " + content + "\n")
	}
	buf.WriteString("\n📝 帖子 ｜ 💬 评论 ｜ 🗣 回帖 ｜ 🌙 清月 ｜ 👨‍💻 用户 ｜ 🏷️ 标签  \n⭐️ 关注 ｜ 👍 赞同 ｜ 💗 感谢 ｜ 💰 打赏 ｜ 🗃 收藏")
	buf.WriteString("\n\n")

	fmt.Println(buf.String())

	readme, err := ioutil.ReadFile("README.md")
	if nil != err {
		logger.Fatalf("read README.md failed: %s", data)
	}

	startFlag := []byte("<!--events start -->")
	beforeStart := readme[:bytes.Index(readme, startFlag)+len(startFlag)]
	newBeforeStart := make([]byte, len(beforeStart))
	copy(newBeforeStart, beforeStart)
	endFlag := []byte("<!--events end -->")
	afterEnd := readme[bytes.Index(readme, endFlag):]
	newAfterEnd := make([]byte, len(afterEnd))
	copy(newAfterEnd, afterEnd)
	newReadme := append(newBeforeStart, buf.Bytes()...)
	newReadme = append(newReadme, newAfterEnd...)
	if err := ioutil.WriteFile("README.md", newReadme, 0644); nil != err {
		logger.Fatalf("write README.md failed: %s", data)
	}
}
