package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"os"

	"github.com/acarl005/stripansi"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var (
	client = loadTwitterClient()
	allowTweet = true
)

// Credentials stores Twitter creds
type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func sendCurrentUsersTwitterMessage() {
	if offlineTwitter {
		return
	}
	// TODO: count all users in all rooms
	if len(mainRoom.users) == 0 {
		return
	}
	if !allowTweet {
		return
	}
	allowTweet = false
	usersSnapshot := append(make([]*user, 0, len(mainRoom.users)), mainRoom.users...)
	areUsersEqual := func(a []*user, b []*user) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if b[i] != a[i] {
				return false
			}
		}
		return true
	}
	go func() {
		time.Sleep(time.Second * 60)
		allowTweet = true
		if !areUsersEqual(mainRoom.users, usersSnapshot) {
			return
		}
		l.Println("Sending twitter update")
		names := make([]string, 0, len(mainRoom.users))
		for _, us := range mainRoom.users {
			names = append(names, us.name)
		}
		t, _, err := client.Statuses.Update("People on Devzat rn: "+stripansi.Strip(fmt.Sprint(names))+"\nJoin em with \"ssh devzat.hackclub.com\"\nUptime: "+printPrettyDuration(time.Since(startupTime)), nil)
		if err != nil {
			if !strings.Contains(err.Error(), "twitter: 187 Status is a duplicate.") {
				mainRoom.broadcast(devbot, "err: "+err.Error())
			}
			l.Println("Got twitter err", err)
			return
		}
		mainRoom.broadcast(devbot, "https\\://twitter.com/"+t.User.ScreenName+"/status/"+t.IDStr)
	}()
}

func loadTwitterClient() *twitter.Client {
	d, err := ioutil.ReadFile("twitter-creds.json")

	if os.IsNotExist(err) {
		offlineTwitter = true
		l.Println("Did not find twitter-creds.json. Enabling offline mode.")
	} else {
		panic(err)
	}

	if offlineTwitter {
		return nil
	}

	twitterCreds := new(Credentials)
	err = json.Unmarshal(d, twitterCreds)
	if err != nil {
		panic(err)
	}
	config := oauth1.NewConfig(twitterCreds.ConsumerKey, twitterCreds.ConsumerSecret)
	token := oauth1.NewToken(twitterCreds.AccessToken, twitterCreds.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	t := twitter.NewClient(httpClient)
	return t
}
