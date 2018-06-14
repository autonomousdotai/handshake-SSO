package utils

import (
    "fmt"
    "strings"
    "math/rand"
    "time"
    "github.com/ninjadotorg/handshake-dispatcher/config"
)

func GetForwardingEndpoint(t string) (string, bool) {
    conf := config.GetConfig()
    var endpoint string
    found := false

    for n, ep := range conf.GetStringMap("forwarding") {
        if n == t {
            endpoint = ep.(string)
            found = true
            break
        }
    }

    return endpoint, found
}

func GetServicesEndpoint(t string) (string, bool) {
    conf := config.GetConfig()
    var endpoint string
    found := false
    
    for n, ep := range conf.GetStringMap("services") {
        if n == t {
            endpoint = ep.(string)
            found = true
            break
        }
    }

    return endpoint, found
}

func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}

func RandomNinjaName() (string) {
    nameGroup1 := []string{"Sad", "Shady", "Surprised", "Happy", "Ridiculous", "Endangered", "Silly", "Crazy", "Sweet", "Amazing", "Confused", "Empty", "Excited", "Manic", "Furious", "Hopeful", "Joyful", "Mad", "Peaceful", "Perplexed", "Powerful", "Proud", "Soulful", "Celebrated", "Worthy", "Angry", "Rare", "Enraged", "Spiteful", "Calm", "Insane", "Humble", "Empty", "The", "Kooky", "Burdened", "Tragic", "Panicked", "Desperate", "Desolate", "Amiable", "Pleasant", "Chewy", "Slimy", "Gregarious", "Vain", "Mopey", "Lame", "Bottomless", "Zen", "Shiny", "Renewed", "Reborn", "Mental", "Metal", "Ashamed", "Fair", "Wise", "Worldly", "Curious", "Wide", "Narrow", "Short", "Tall", "Curved", "Wonky", "Flightless", "Bored", "Superior", "Inferior", "Pungent", "Delicate", "Polite", "Gentle", "Senior"}

    nameGroup2 := []string{"Mystic", "Silver", "Silent", "Ghost", "Winged", "Golden", "Deadly", "Young", "Lucky", "Sparkling", "Masked", "Hidden", "Glowing", "Fiery", "Feisty", "Unknown", "Notorious", "Reborn", "Refreshing", "Honorable", "Wicked", "Dr.", "Nobel", "No-eyed", "Blind", "Floating", "Magic", "Stone", "Marble", "Jade", "Cold", "Fanged", "Frozen", "Chewie", "Burnt", "Legless", "All Knowing", "All Seeing", "Crystal", "Unstoppable", "Windy", "Slick", "Secret", "Unmistakable", "Hand of", "Grey", "Red", "Black", "Beautiful", "Ultimate", "Undefeated", "Wondrous", "Flying", "Moonlit", "Blazing", "Hot", "Unwieldy", "Dry", "Parched", "Arid", "Crispy", "Delicious", "Underground", "Oily", "Salty", "Sweet", "Umami", "Bitter", "Steely", "Blackened", "Roasted", "Simmered", "Poached"}

    nameGroup3 := []string{"Moon", "Hamster", "Snowflake", "Eye", "Hand", "Duck", "Reflection", "Butterfly", "Fly", "Ant", "Horse", "Donkey", "Weasel", "Dragon", "Wasp", "Assassin", "Biscuit", "Fork", "Spoon", "Lemon", "Zero", "Chaos", "Master", "Shadow", "Foot", "Bald Patch", "Shoe", "TV Dinner", "Mess", "Egg", "Cat", "Viper", "Pigeon", "Sense of purpose", "Whisper", "Rollerblade", "Pancake", "Beanie Baby", "Misled Youth", "Boot", "Sock", "Jockstrap", "Cake", "Loaf of bread", "Serpent", "Echo", "Watchman", "Kid", "Warrior", "Goat", "Fish", "Squirrel", "Mole", "Tiger", "Spider", "Dog", "Storm", "Potato", "Carrot", "Cucumber", "Answer", "Voice", "Star", "Raindrop", "Heart", "Raven", "Reptile", "Spirit", "Ghost", "Binbag", "Oracle", "Sage", "Ponytail", "Tissue", "Dentist", "Monsoon", "Steel", "Proboscis", "Broth"}

    var r1, r2, r3 int
    r1 = random(0, len(nameGroup1))
    r2 = random(0, len(nameGroup2))
    r3 = random(0, len(nameGroup3))
    
    name := fmt.Sprintf("%s %s %s", nameGroup1[r1], nameGroup2[r2], nameGroup3[r3])
    return strings.Replace(name, " ", "", -1)
}
