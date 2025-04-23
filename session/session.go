package session

import (
	ps "pubsub.com/pubsub/pubsub"
)

func LoadRepo() *ps.PostsRepository {
	repo := ps.NewPostsRepository()

	// Register users
	repo.RegisterUser("rich", "rich@hello.com", "It's me!", "rich.png", "topsecret")
	repo.RegisterUser("paul", "paul@hello.com", "It's me, Paul!", "paul.png", "topsecret")
	repo.RegisterUser("perry", "perry@hello.com", "It's me, Perry!", "perry.png", "topsecret")
	repo.RegisterUser("jack", "jack@hello.com", "Drink! Feck! Girls!", "jack.png", "topsecret")
	repo.RegisterUser("diana", "diana@hello.com", "It's me, Diana!", "diana.png", "topsecret")
	repo.RegisterUser("may", "may@hello.com", "It's me, May!", "may.png", "topsecret")
	repo.RegisterUser("rutger", "rutger@hello.com", "I've seen things you people wouldn't believe...", "rutger.png", "topsecret")

	// Post some messages
	repo.PostMessage("rich", "Hello world! #hi")
	repo.PostMessage("rich", "Figured I'd join this festering hellscape of a website. #socialmedia")
	repo.PostMessage("jack", "That would be an ecumenical matter.")
	repo.PostMessage("paul", "One ball, corner pocket. #TheHustler")
	repo.PostMessage("paul", "Maybe I'm not such a high-class piece of property right now. And a 25% slice of something big is better than a 100% slice of nothing.")
	repo.PostMessage("perry", "How much wood could a woodchuck chuck if a woodchuck could chuck wood? #woodchuck #wood")
	repo.PostMessage("perry", "I'm a lizard person. #lizard #pets")
	repo.PostMessage("diana", `I'm a cat person. #cat`)
	repo.PostMessage("may", `I'm a dog person. #dog`)

	// Follow users
	repo.Follow("rich", "paul")
	repo.Follow("paul", "perry")
	repo.Follow("rich", "paul")
	repo.Follow("perry", "rich")
	repo.Follow("rich", "perry")
	repo.Follow("perry", "rich")
	repo.Follow("rich", "diana")
	repo.Follow("diana", "rich")
	repo.Follow("rich", "may")
	repo.Follow("may", "rich")
	repo.Follow("may", "diana")
	repo.Follow("rutger", "rich")
	repo.Follow("rutger", "diana")
	repo.Follow("rutger", "may")
	repo.Follow("rich", "jack")
	repo.Follow("rich", "perry")
	repo.Follow("perry", "rich")

	return repo
}
