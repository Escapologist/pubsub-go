package pubsub

import (
	"fmt"
	"regexp"
	"sort"
	"time"
)

//https://github.com/rhettinger/modernpython/blob/master/pubsub/pubsub.py

type User = string
type Timestamp = time.Time
type HashTag = string

type Post struct {
	Timestamp Timestamp
	User      User
	Text      string
}

func (p *Post) TimeSincePosted() time.Duration {
	now := time.Now()
	return now.Sub(p.Timestamp).Truncate(time.Second)
}

type UserInfo struct {
	Displayname string
	Email       string
	Bio         string
	Photo       string
}

var hashtagPattern = regexp.MustCompile(`[#@]\w+`)
var hashtagIndex = make(map[string][]Post)

type Set[T comparable] map[T]bool

func NewSet[T comparable](items ...T) Set[T] {
	set := Set[T]{}
	for _, i := range items {
		set[i] = true
	}
	return set
}

var userInfo = make(map[string]UserInfo)

var following = make(map[User]Set[User])
var followers = make(map[User]Set[User])

var posts []Post
var userPosts = make(map[User][]Post)

// TODO save posts to a DB
func PostMessage(user User, text string) {
	timestamp := time.Now()
	post := Post{timestamp, user, text}

	posts = append(posts, post)
	if len(userPosts[user]) == 0 {
		userPosts[user] = []Post{}
	}
	userPosts[user] = append(userPosts[user], post)

	for _, h := range hashtagPattern.FindAll([]byte(text), -1) {
		key := string(h)
		if len(hashtagIndex[key]) == 0 {
			hashtagIndex[key] = []Post{}
		}
		hashtagIndex[key] = append(hashtagIndex[key], post)
	}
}

func Follow(user User, followedUser User) {
	if following[user] == nil {
		following[user] = make(Set[User])
	}
	following[user][followedUser] = true

	if followers[followedUser] == nil {
		followers[followedUser] = make(Set[User])
	}
	followers[followedUser][user] = true
}

func PostsByUser(user User, limit ...int) []Post {
	if len(limit) > 0 {
		return userPosts[user][:limit[0]]
	}
	posts := userPosts[user]
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Timestamp.After(posts[j].Timestamp)
	})
	return posts
}

func PostsForUser(user User, limit ...int) []Post {
	followed := following[user]
	relevantPosts := []Post{}
	for k := range followed {
		fmt.Printf("%v\n", k)
		posts := userPosts[k]
		relevantPosts = append(relevantPosts, posts...)
	}
	if len(limit) > 0 {
		return relevantPosts[:limit[0]]
	}
	return relevantPosts
}

func GetFollowers(user User) Set[User] {
	return followers[user]
}

func GetFolloweed(user User) Set[User] {
	return following[user]
}

func AddUserInfo(ui ...UserInfo) {
	for _, u := range ui {
		userInfo[u.Displayname] = u
	}
}

func GetuserInfo(name string) UserInfo {
	return userInfo[name]
}
