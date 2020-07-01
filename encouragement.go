package main

import (
	"math/rand"
	"time"
)

var (
	seed = rand.NewSource(time.Now().UnixNano())
	rnd  = rand.New(seed)
)

func getEncouragement() string {
	idx := rnd.Intn(len(data))
	return data[idx]
}

var data = []string{
	"That was a great post!",
	"You have a way with words.",
	"I see the deeper meaning in that post.",
	"You are the Mark Twain of Mattermost",
	"That was an astute comment.",
	"You've written another gem.",
	"I like your style.",
	"Hang in there.",
	"Don't give up.",
	"Keep pushing.",
	"Keep fighting!",
	"Stay strong.",
	"Never give up.",
	"Never say 'die'.",
	"I knew you could do it!",
	"There you go!",
	"Keep up the good work.",
	"Keep it up.",
	"Good job.",
	"I’m so proud of you!",
	"I’ll support you either way.",
	"I’m behind you 100%.",
	"It’s totally up to you.",
	"I trust you to make the hard decisions.",
	"Follow your dreams.",
	"Reach for the stars.",
	"You can do the impossible.",
	"Believe in yourself.",
	"The sky is the limit.",
	"Everything you need to accomplish your goals is already in you.",
}
