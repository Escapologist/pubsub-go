package pubsub

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"sort"
	"strings"
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

type Set[T comparable] map[T]bool

func NewSet[T comparable](items ...T) Set[T] {
	set := Set[T]{}
	for _, i := range items {
		set[i] = true
	}
	return set
}

type PostsRepositoryInterface interface {
	PostMessage(user User, text string)
	PostsByUser(user User, limit ...int) []Post
	PostsForUser(user User, limit ...int) []Post
	Follow(user User, followedUser User)
	GetFollowers(user User) Set[User]
	GetFollowed(user User) Set[User]
	AddUserInfo(ui ...UserInfo)
	GetuserInfo(name string) UserInfo
	Search(query string) []Post
	GetHashTags() []string
	SearchByTag(tag string) []Post
}

type PostsRepository struct {
	Posts     []Post
	UserPosts map[User][]Post
	UserInfo  map[string]UserInfo

	Following map[User]Set[User]
	Followers map[User]Set[User]

	HashtagIndex map[string][]Post
}

func NewPostsRepository() *PostsRepository {
	return &PostsRepository{
		Posts:        []Post{},
		UserPosts:    make(map[User][]Post),
		UserInfo:     make(map[string]UserInfo),
		Following:    make(map[User]Set[User]),
		Followers:    make(map[User]Set[User]),
		HashtagIndex: make(map[string][]Post),
	}
}

// TODO save posts to a DB
func (r *PostsRepository) PostMessage(user User, text string) {
	timestamp := time.Now()
	post := Post{timestamp, user, text}

	r.Posts = append(r.Posts, post)
	if len(r.UserPosts[user]) == 0 {
		r.UserPosts[user] = []Post{}
	}
	r.UserPosts[user] = append(r.UserPosts[user], post)

	for _, h := range hashtagPattern.FindAll([]byte(text), -1) {
		key := string(h)
		if len(r.HashtagIndex[key]) == 0 {
			r.HashtagIndex[key] = []Post{}
		}
		r.HashtagIndex[key] = append(r.HashtagIndex[key], post)
	}
}

func (r *PostsRepository) Follow(user User, followedUser User) {
	if r.Following[user] == nil {
		r.Following[user] = make(Set[User])
	}
	r.Following[user][followedUser] = true

	if r.Followers[followedUser] == nil {
		r.Followers[followedUser] = make(Set[User])
	}
	r.Followers[followedUser][user] = true
}

func (r *PostsRepository) PostsByUser(user User, limit ...int) []Post {
	if len(limit) > 0 {
		return r.UserPosts[user][:limit[0]]
	}
	posts := r.UserPosts[user]
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Timestamp.After(posts[j].Timestamp)
	})
	return posts
}

func (r *PostsRepository) PostsForUser(user User, limit ...int) []Post {
	followed := r.Following[user]
	relevantPosts := []Post{}
	for k := range followed {
		fmt.Printf("%v\n", k)
		posts := r.UserPosts[k]
		relevantPosts = append(relevantPosts, posts...)
	}
	if len(limit) > 0 {
		return relevantPosts[:limit[0]]
	}
	return relevantPosts
}

func (r *PostsRepository) GetFollowers(user User) Set[User] {
	return r.Followers[user]
}

func (r *PostsRepository) GetFollowed(user User) Set[User] {
	return r.Following[user]
}

func (r *PostsRepository) AddUserInfo(ui ...UserInfo) {
	for _, u := range ui {
		r.UserInfo[u.Displayname] = u
	}
}

func (r *PostsRepository) GetuserInfo(name string) UserInfo {
	return r.UserInfo[name]
}

func (r *PostsRepository) Search(query string) []Post {
	res := []Post{}
	for _, p := range r.Posts {
		if strings.Contains(p.Text, query) {
			res = append(res, p)
		}
	}
	return res
}

func (r *PostsRepository) SearchByTag(tag string) []Post {
	return r.HashtagIndex[tag]
}

func (r *PostsRepository) GetHashTags() []string {
	return slices.Sorted(maps.Keys(r.HashtagIndex))
}
