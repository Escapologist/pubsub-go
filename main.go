package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	. "pubsub.com/pubsub/pubsub"
)

var following = make(map[User]Set[User])
var followers = make(map[User]Set[User])

var posts []Post
var userPosts = make(map[User][]Post)

var hashtagIndex = make(map[string][]Post)

// Auth
var password = "test"
var loggedInUsers = make(map[string]User)

func CheckUser(u string, pw string) bool {
	return u != "" && pw != ""
}

func GetLoggedInUser(r *http.Request) string {
	cookie, err := r.Cookie("token")
	fmt.Printf("User cookie: %v\n", cookie)
	if err != nil {
		return ""
	}
	return string(loggedInUsers[cookie.Value])
}

// end Auth

func ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		user := r.Form.Get("username")
		pw := r.Form.Get("password")
		loggedIn := CheckUser(user, pw)
		fmt.Printf("User %v\n", user)
		fmt.Printf("Password %v\n", pw)
		fmt.Printf("Logged in %v\n", loggedIn)

		if !loggedIn {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		randValues := make([]byte, 32)
		rand.Read(randValues)
		token := base64.URLEncoding.EncodeToString(randValues)

		loggedInUsers[(token)] = User(user)
		fmt.Printf("Logged in users: %v\n", loggedInUsers)

		cookie := http.Cookie{Name: "token", Value: token}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else {
		u := GetLoggedInUser(r)
		if u != "" {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		data := map[string]string{"Title": "log in"}
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(w, data)
	}
}

func loginRequired(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		fmt.Printf("Token: %v\n", token)

		u := loggedInUsers[token.Value]
		fmt.Printf("Logged in User %v\n", u)
		if u == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		fn(w, r, string(u))
	}
}

// On login POST
func CheckCredentials(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := r.Form.Get("user")
	pw := r.Form.Get("password")
	loggedIn := CheckUser(user, pw)
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
		return
	}

	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusContinue)
		return
	}
	loggedInUsers[string(token)] = User(user)

	cookie := http.Cookie{Name: user, Value: string(token)}
	http.SetCookie(w, &cookie)
	ShowMainPage(w, r, user)
}

func ShowMainPage(w http.ResponseWriter, r *http.Request, user string) {
	title := "Posts from people you follow"
	postsForUser := PostsForUser(User(user))
	fmt.Printf("Posts relevant for user: %v\n", postsForUser)

	data := map[string]any{"Title": title, "User": string(user), "Posts": postsForUser}

	t, _ := template.ParseFiles("templates/main.html")
	t.Execute(w, data)
}

func MakePost(w http.ResponseWriter, r *http.Request, user string) {
	r.ParseForm()
	message := r.Form.Get("message")
	fmt.Printf("message %v", message)
	PostMessage(user, message)
	http.Redirect(w, r, "/", http.StatusFound)
}

func ShowUserPage(w http.ResponseWriter, r *http.Request, user string) {
	name := r.PathValue("name")
	userInfo := GetuserInfo(name)
	posts := PostsByUser(name)
	followers := GetFollowers(name)
	following := GetFolloweed(name)
	data := map[string]any{
		"Title":     name,
		"User":      string(name),
		"UserInfo":  userInfo,
		"Posts":     posts,
		"Followers": followers,
		"Following": following,
	}

	t, _ := template.ParseFiles("templates/user.html")
	t.Execute(w, data)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", ShowLoginPage)
	mux.HandleFunc("/", loginRequired(ShowMainPage))
	mux.HandleFunc("POST /postmessage", loginRequired(MakePost))
	mux.HandleFunc("GET /user/{name}", loginRequired(ShowUserPage))
	mux.HandleFunc("GET /user/{name}/following", loginRequired(ShowUserPage))
	mux.HandleFunc("GET /user/{name}/followers", loginRequired(ShowUserPage))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	log.Fatal(s.ListenAndServe())
}

func init() {
	var user1 User = "rich"
	var user2 User = "bob"

	rich := UserInfo{Displayname: string(user1), Email: "rich@hello.com", Bio: "It's me!", Photo: "rich.webp"}
	bob := UserInfo{Displayname: string(user2), Email: "bob@hello.com", Bio: "It's me, Bob!", Photo: "bob.webp"}

	AddUserInfo(rich, bob)

	PostMessage("rich", "Hello world! #hi")
	PostMessage("rich", "nothing to say. #nada #meh")
	PostMessage("kev", "Ugh! #sounfair")
	PostMessage(user2, "Something! #sounfair")
	PostMessage("bob", "nothing to say. #nada #meh")

	fmt.Printf("%v\n", hashtagIndex)
	fmt.Printf("%v\n", userPosts)
	fmt.Printf("Posts by user: %v\n", PostsByUser(user1))
	fmt.Printf("Posts by user: %v\n", PostsByUser(user1, 1))

	Follow(user1, user2)
	Follow(user1, "kate")
	Follow(user2, "perry")
	Follow(user1, user2)
	Follow("bob", "rich")
	Follow("perry", "rich")

	fmt.Printf("Followers: %v\n", followers)
	fmt.Printf("Following: %v\n", following)

	fmt.Printf("Posts by user: %v\n", PostsByUser(user1))
	fmt.Printf("Posts relevant for user: %v\n", PostsForUser(user1))

	Follow(user1, "perry")
	Follow("perry", user1)
	fmt.Printf("Followed: %v\n", GetFolloweed(user1))
	fmt.Printf("Followers: %v\n", GetFollowers(user1))

}

// // def follow(user: User, followed_user: User) -> None:
// //     user, followed_user = intern(user), intern(followed_user)
// //     following[user].add(followed_user)
// //     followers[followed_user].add(user)

// // def posts_by_user(user: User, limit: Optional[int] = None) -> List[Post]:
// //     return list(islice(user_posts[user], limit))

// // def posts_for_user(user: User, limit: Optional[int] = None) -> List[Post]:
// //     relevant = merge(*[user_posts[u] for u in following[user]], reverse=True)
// //     return list(islice(relevant, limit))

// def get_followers(user: User) -> List[User]:
//     return sorted(followers[user])

// def get_followed(user: User) -> List[User]:
//     return sorted(following[user])

// def search(phrase: str, limit: Optional[int] = None) -> List[Post]:
//     if hashtag_pattern.match(phrase):
//         return list(islice(hashtag_index[phrase], limit))
//     return list(islice((post for post in posts if phrase in post.text), limit))

// def hash_password(password: str, salt: Optional[bytes] = None) -> HashAndSalt:
//     pepper = b'alchemists discovered that gold came from earth air fire and water'
//     salt = salt or secrets.token_bytes(16)
//     return hashlib.pbkdf2_hmac('sha512', password.encode(), salt+pepper, 100_000), salt

// def set_user(user: User, displayname: str, email: str, password: str,
//              bio: Optional[str]=None, photo: Optional[str]=None) -> None:
//     user = intern(user)
//     hashed_password = hash_password(password)
//     user_info[user] = UserInfo(displayname, email, hashed_password, bio, photo)

// def check_user(user: User, password: str) -> bool:
//     hashpass, salt = user_info[user].hashed_password
//     target_hash_pass = hash_password(password, salt)[0]
//     sleep(random.expovariate(10))
//     return secrets.compare_digest(hashpass, target_hash_pass)

// def get_user(user: User) -> Optional[UserInfo]:
//     return user_info.get(user)

// time_unit_cuts = [60, 3600, 3600*24]                                           # type: List[int]
// time_units = [(1, 'second'), (60, 'minute'), (3600, 'hour'), (24*3600, 'day')] # type: List[Tuple[int, str]]

// def age(post: Post) -> str:
//     seconds = time() - post.timestamp
//     divisor, unit = time_units[bisect(time_unit_cuts, seconds)]
//     units = seconds // divisor
//     return '%d %s ago' % (units, unit + ('' if units==1 else 's'))

// def save() -> None:
//     with open('pubsub.pickle', 'wb') as f:
//         pickle.dump([posts, user_posts, hashtag_index, following, followers, user_info], f)

// def restore() -> None:
//     global posts, user_posts, hashtag_index, following, followers, user_info
//     with open('pubsub.pickle', 'rb') as f:
//         posts, user_posts, hashtag_index, following, followers, user_info = pickle.load(f)
