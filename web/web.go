package web

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	ps "pubsub.com/pubsub/pubsub"
)

type AppHandler struct {
	Repo          ps.PostsRepositoryInterface
	LoggedInUsers map[string]ps.User
}

func NewAppHandler(repo ps.PostsRepositoryInterface, loggedInUsers map[string]ps.User) *AppHandler {
	return &AppHandler{
		Repo:          repo,
		LoggedInUsers: loggedInUsers,
	}
}

// Auth

func (h *AppHandler) CheckUser(u string, pw string) bool {
	user := h.Repo.GetUserInfo(u)
	err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(pw))
	return err == nil
}

func (h *AppHandler) GetLoggedInUser(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		return ""
	}
	return string(h.LoggedInUsers[cookie.Value])
}

// end Auth

func (h *AppHandler) loginRequired(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		u := h.LoggedInUsers[token.Value]
		if u == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		fn(w, r, u)
	}
}

func (h *AppHandler) ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	u := h.GetLoggedInUser(r)
	if u != "" {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	flashMessage := getFlashMessage(w, r)
	data := map[string]string{"Title": "log in", "User": "", "FlashMessage": flashMessage}

	t, err := template.ParseFiles("templates/base.html", "templates/login.html")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	t.ExecuteTemplate(w, "base", data)
}

func (h *AppHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := r.Form.Get("username")
	pw := r.Form.Get("password")
	loggedIn := h.CheckUser(user, pw)

	if !loggedIn {
		setFlashMessage(w, "Wrong username and/or password!")
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	randValues := make([]byte, 32)
	rand.Read(randValues)
	token := base64.URLEncoding.EncodeToString(randValues)

	h.LoggedInUsers[(token)] = ps.User(user)
	fmt.Printf("Logged in users: %v\n", h.LoggedInUsers)

	cookie := http.Cookie{Name: "token", Value: token}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AppHandler) ShowSignupPage(w http.ResponseWriter, r *http.Request) {
	u := h.GetLoggedInUser(r)
	if u != "" {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	flashMessage := getFlashMessage(w, r)

	data := map[string]string{"Title": "Sign up", "User": "", "FlashMessage": flashMessage}
	t, err := template.ParseFiles("templates/base.html", "templates/signup.html")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	t.ExecuteTemplate(w, "base", data)
}

func (h *AppHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	repeatPassword := r.Form.Get("repeatpassword")
	if password != repeatPassword {
		setFlashMessage(w, "Passwords don't match!")
		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	user := h.Repo.GetUserInfo(username)
	if user.Displayname != "" {
		setFlashMessage(w, "Username already exists")
		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	h.Repo.RegisterUser(username, email, "", "", password)
	setFlashMessage(w, "User created successfully!")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *AppHandler) ShowMainPage(w http.ResponseWriter, r *http.Request, loggedInUser string) {
	title := "Posts from people you follow"
	postsForUser := h.Repo.PostsForUser(ps.User(loggedInUser))

	followed := h.Repo.GetFollowed(loggedInUser)
	followedProfiles := map[string]ps.UserInfo{}
	for u := range followed {
		userInfo := h.Repo.GetUserInfo(u)
		followedProfiles[u] = userInfo
	}

	flashMessage := ""
	if cookie, err := r.Cookie("flash"); err == nil {
		flashMessage = cookie.Value
		http.SetCookie(w, &http.Cookie{Name: "flash", Value: "", MaxAge: -1})
	}

	data := map[string]any{
		"Title":            title,
		"LoggedInUser":     string(loggedInUser),
		"Posts":            postsForUser,
		"FollowedProfiles": followedProfiles,
		"FlashMessage":     flashMessage,
	}

	t, _ := template.ParseFiles("templates/base.html", "templates/main.html")
	t.ExecuteTemplate(w, "base", data)
}

func (h *AppHandler) MakePost(w http.ResponseWriter, r *http.Request, user string) {
	r.ParseForm()
	message := r.Form.Get("message")
	h.Repo.PostMessage(user, message)

	http.SetCookie(w, &http.Cookie{
		Name:  "flash",
		Value: "Post created successfully!",
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AppHandler) FollowUser(w http.ResponseWriter, r *http.Request, user string) {
	followee := r.PathValue("followee")
	h.Repo.Follow(user, followee)

	http.Redirect(w, r, fmt.Sprintf("/user/%s", followee), http.StatusFound)
}

func (h *AppHandler) UnfollowUser(w http.ResponseWriter, r *http.Request, user string) {
	followee := r.PathValue("followee")
	h.Repo.Unfollow(user, followee)

	http.Redirect(w, r, fmt.Sprintf("/user/%s", followee), http.StatusFound)
}

func (h *AppHandler) ShowUserPage(w http.ResponseWriter, r *http.Request, loggedInUser string) {
	name := r.PathValue("name")
	userInfo := h.Repo.GetUserInfo(name)
	posts := h.Repo.PostsByUser(name)
	followers := h.Repo.GetFollowers(name)
	following := h.Repo.GetFollowed(name)
	followed := followers[loggedInUser]
	data := map[string]any{
		"Title":         name,
		"LoggedInUser":  loggedInUser,
		"User":          string(name),
		"UserInfo":      userInfo,
		"Posts":         posts,
		"Followers":     followers,
		"Following":     following,
		"Followed":      followed,
		"IsCurrentUser": loggedInUser == name,
	}

	t, err := template.ParseFiles("templates/base.html", "templates/user.html")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	t.ExecuteTemplate(w, "base", data)
}

func (h *AppHandler) ShowFollowers(w http.ResponseWriter, r *http.Request, loggedInUser string) {
	followee := r.PathValue("followee")
	followers := h.Repo.GetFollowers(followee)
	ui := []ps.UserInfo{}

	for u := range followers {
		fmt.Printf("FOLLOWER: %v\n", u)
		ui = append(ui, h.Repo.GetUserInfo((u)))
		fmt.Printf("UI: %v\n", ui)

	}

	data := map[string]any{
		"Title":        "Followers of " + followee,
		"User":         followee,
		"LoggedInUser": loggedInUser,
		"Follows":      ui,
	}

	t, _ := template.ParseFiles("templates/base.html", "templates/follow.html")
	t.ExecuteTemplate(w, "base", data)
}

func (h *AppHandler) ShowFollowed(w http.ResponseWriter, r *http.Request, loggedInUser string) {
	follower := r.PathValue("follower")
	followed := h.Repo.GetFollowed(follower)
	ui := []ps.UserInfo{}

	for u := range followed {
		ui = append(ui, h.Repo.GetUserInfo((u)))
	}

	data := map[string]any{
		"Title":        "Followed by " + follower,
		"User":         follower,
		"LoggedInUser": loggedInUser,
		"Follows":      ui,
	}

	t, _ := template.ParseFiles("templates/base.html", "templates/follow.html")
	t.ExecuteTemplate(w, "base", data)
}

func (h *AppHandler) SearchPosts(w http.ResponseWriter, r *http.Request, user string) {
	query := r.URL.Query().Get("query")
	postsFound := h.Repo.Search(query)
	fmt.Printf("Params: %v\n", r.URL.Query())

	fmt.Printf("Query: %v\n", query)
	fmt.Printf("Posts: %v\n", postsFound)

	followedProfiles := map[string]ps.UserInfo{}
	for _, post := range postsFound {
		userInfo := h.Repo.GetUserInfo(post.User)
		followedProfiles[post.User] = userInfo
	}
	fmt.Printf("followedProfiles: %v\n", followedProfiles)

	data := map[string]any{
		"LoggedInUser":     user,
		"Posts":            postsFound,
		"Query":            query,
		"FollowedProfiles": followedProfiles,
	}

	t, _ := template.ParseFiles("templates/base.html", "templates/search.html")
	t.ExecuteTemplate(w, "base", data)
}

func (h *AppHandler) SearchHashTag(w http.ResponseWriter, r *http.Request, user string) {
	tag := "#" + r.URL.Path[len("/hashtag/"):]
	posts := h.Repo.SearchByTag(tag)

	followed := h.Repo.GetFollowed(user)
	followedProfiles := map[string]ps.UserInfo{}
	for user := range followed {
		userInfo := h.Repo.GetUserInfo(user)
		followedProfiles[user] = userInfo
	}

	data := map[string]any{
		"Posts":            posts,
		"Query":            tag,
		"FollowedProfiles": followedProfiles,
	}

	t, _ := template.ParseFiles("templates/base.html", "templates/search.html")
	t.ExecuteTemplate(w, "base", data)
}

func MakeServer(h *AppHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /login", h.ShowLoginPage)
	mux.HandleFunc("POST /login", h.HandleLogin)
	mux.HandleFunc("GET /signup", h.ShowSignupPage)
	mux.HandleFunc("POST /signup", h.HandleSignup)

	mux.HandleFunc("/", h.loginRequired(h.ShowMainPage))
	mux.HandleFunc("POST /postmessage", h.loginRequired(h.MakePost))
	mux.HandleFunc("GET /user/{name}", h.loginRequired(h.ShowUserPage))
	mux.HandleFunc("POST /user/{followee}/follow", h.loginRequired(h.FollowUser))
	mux.HandleFunc("POST /user/{followee}/unfollow", h.loginRequired(h.UnfollowUser))
	mux.HandleFunc("GET /user/{follower}/following", h.loginRequired(h.ShowFollowed))
	mux.HandleFunc("GET /user/{followee}/followers", h.loginRequired(h.ShowFollowers))
	mux.HandleFunc("GET /search", h.loginRequired(h.SearchPosts))
	mux.HandleFunc("GET /hashtag/{tag}", h.loginRequired(h.SearchHashTag))

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	return mux
}
