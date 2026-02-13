package ws

import (
	"fmt"
	"math/rand"
)

var adjectives = []string{
	"happy", "cosmic", "neon", "fuzzy", "sparkly",
	"groovy", "stellar", "electric", "golden", "mystic",
	"radiant", "vivid", "bold", "swift", "bright", "jolly",
}

var nouns = []string{
	"unicorn", "rainbow", "rocket", "pizza", "penguin",
	"dragon", "taco", "octopus", "phoenix", "panda",
	"koala", "owl", "falcon", "dolphin", "tiger", "wolf",
}

// GenerateChannelKey creates a memorable channel key like "happy-unicorn-42".
func GenerateChannelKey() string {
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	num := rand.Intn(100)
	return fmt.Sprintf("%s-%s-%d", adj, noun, num)
}
