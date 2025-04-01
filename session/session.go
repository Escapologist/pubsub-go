package session

import (
	"fmt"

	ps "pubsub.com/pubsub/pubsub"
)

func LoadRepo() ps.PostsRepositoryInterface {
	repo := ps.NewPostsRepository()
	var user1 ps.User = "rich"
	var user2 ps.User = "bob"

	rich := ps.UserInfo{Displayname: string(user1), Email: "rich@hello.com", Bio: "It's me!", Photo: "rich.png"}
	bob := ps.UserInfo{Displayname: string(user2), Email: "bob@hello.com", Bio: "It's me, Bob!", Photo: "bob.webp"}
	perry := ps.UserInfo{Displayname: "perry", Email: "perry@hello.com", Bio: "It's me, Perry!", Photo: "perry.png"}
	kate := ps.UserInfo{Displayname: "kate", Email: "kate@hello.com", Bio: "It's me, Kate!", Photo: "kiki.jpg"}
	kev := ps.UserInfo{Displayname: "kev", Email: "kev@hello.com", Bio: "It's me, Kev!", Photo: "png-transparent-the-simpson-character-cletus-spuckler-groundskeeper-willie-snake-jailbird-mayor-quimby-ralph-wiggum-the-simpsons-movie-miscellaneous-television-vertebrate.png"}
	joan := ps.UserInfo{Displayname: "joan", Email: "joan@hello.com", Bio: "It's me, Joan!", Photo: "joan.jpg"}
	may := ps.UserInfo{Displayname: "may", Email: "may@hello.com", Bio: "It's me, May!", Photo: "may.png"}
	repo.AddUserInfo(rich, bob, perry, kate, kev, joan, may)

	repo.PostMessage("rich", "Hello world! #hi")
	repo.PostMessage("rich", "My friend Jack claims he can communicate with vegetables. I guess you could say... Jack and the beans talk. #funny")
	repo.PostMessage("kev", `“We’re gonna need more chalk." – detective who discovers my body. #murder`)
	repo.PostMessage(user2, "*a jerk tries to punch me but I catch it perfectly in my mouth and swallow him whole like a snake.* #imawesome")
	repo.PostMessage("bob", `me: I just want 2 minutes of privacy in the bathroom. my kid: best I can do is a paleontology lecture.`)
	repo.PostMessage("perry", "[first day as a spy] Wife: what’s your bosses name? Me: I can’t tell you that Wife: why? Me: because I don’t remember, Linda. #spy")
	repo.PostMessage("perry", `hey "nice" manbun haha it fuckin sucks you hipster asshole [he turns around and reveals he is a samurai from the tokugawa shogunate] oh fuck`)
	repo.PostMessage("joan", `I'm a cat person. #cat`)
	repo.PostMessage("may", `I'm a dog person. #dog`)
	fmt.Printf("%v\n", repo.GetHashTags())
	fmt.Printf("Posts by user: %v\n", repo.PostsByUser(user1))
	fmt.Printf("Posts by user: %v\n", repo.PostsByUser(user1, 1))

	repo.Follow(user1, user2)
	repo.Follow(user1, "kate")
	repo.Follow(user2, "perry")
	repo.Follow(user1, user2)
	repo.Follow("bob", "rich")
	repo.Follow("perry", "rich")
	repo.Follow("rich", "perry")
	repo.Follow("perry", "rich")
	repo.Follow("rich", "kate")
	repo.Follow("rich", "joan")
	repo.Follow("joan", "rich")
	repo.Follow("rich", "may")
	repo.Follow("may", "rich")
	repo.Follow("may", "joan")
	repo.Follow("may", "kate")

	fmt.Printf("Posts by user: %v\n", repo.PostsByUser(user1))
	fmt.Printf("Posts relevant for user: %v\n", repo.PostsForUser(user1))

	repo.Follow(user1, "perry")
	repo.Follow("perry", user1)
	fmt.Printf("Followed: %v\n", repo.GetFollowed(user1))
	fmt.Printf("Followers: %v\n", repo.GetFollowers(user1))
	return repo
}
