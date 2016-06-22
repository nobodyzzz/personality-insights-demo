package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/SlyMarbo/rss"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gin-gonic/gin"
	"github.com/kennygrant/sanitize"
	"github.com/laktek/Stack-on-Go/stackongo"
	"github.com/liviosoares/go-watson-sdk/watson"
	"github.com/liviosoares/go-watson-sdk/watson/personality_insights"
)

func main() {
	var port string
	var p personality_insights.Profile
	var profile *bytes.Buffer

	if appEnv, err := cfenv.Current(); err != nil {
		port = ":8080"
	} else {
		port = fmt.Sprintf(":%d", appEnv.Port)
	}
	r := gin.Default()
	r.LoadHTMLFiles("./static/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/profile/:username/:source", func(c *gin.Context) {
		source := c.Param("source")
		username := c.Param("username")
		client, err := personality_insights.NewClient(watson.Config{})
		if err == nil {
			switch source {
			case "reddit":
				profile, err = RedditProfile(username)
			case "stackoverflow":
				profile, err = StackoverfloProfile(username)
			default:
				err = fmt.Errorf("Unknown source: %s", source)
			}
		}
		if err == nil {
			p, err = client.GetProfile(profile, "text/plain; charset=utf-8", "")
		}
		if err == nil {
			c.JSON(200, p)
		} else {
			c.JSON(400, gin.H{"error": err.Error()})
		}
	})

	r.Static("/static", "./static")
	r.Run(port)
}

type defaultHeaderSet struct{}

func (d defaultHeaderSet) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", "PersonalInsightsDemo/1.0")
	return http.DefaultTransport.RoundTrip(r)
}

func RedditProfile(username string) (*bytes.Buffer, error) {
	time.Sleep(2 * time.Second)
	s := http.DefaultClient.Transport
	http.DefaultClient.Transport = &defaultHeaderSet{}
	feed, err := rss.Fetch(fmt.Sprintf("https://www.reddit.com/user/%s/comments.rss", username))
	http.DefaultClient.Transport = s
	if err != nil {
		return nil, err
	}
	profile := bytes.NewBufferString("")
	for _, item := range feed.Items {
		profile.WriteString(fmt.Sprintf("%s\n", sanitize.HTML(item.Content)))
	}
	return profile, nil
}

func StackoverfloProfile(username string) (*bytes.Buffer, error) {
	session := stackongo.NewSession("stackoverflow")
	params := make(stackongo.Params)
	params.Add("inname", username)
	users, err := session.AllUsers(params)
	if err != nil {
		return nil, err
	}
	if len(users.Items) < 1 {
		return nil, fmt.Errorf("User '%s' not found", username)
	}
	params = make(stackongo.Params)
	params.Add("filter", "withbody")
	profile := bytes.NewBufferString("")
	answers, err := session.AnswersFromUsers([]int{users.Items[0].User_id}, params)
	if err != nil {
		return nil, err
	}

	for _, answer := range answers.Items {
		profile.WriteString(fmt.Sprintf("%s\n", sanitize.HTML(answer.Body)))
	}
	questions, err := session.QuestionsFromUsers([]int{users.Items[0].User_id}, params)
	if err != nil {
		return nil, err
	}
	for _, question := range questions.Items {
		profile.WriteString(fmt.Sprintf("%s\n", sanitize.HTML(question.Body)))
	}
	return profile, nil
}
