package main

import (
	"encoding/json"
	"fmt"

	got "github.com/MariaTerzieva/gotumblr"
	"github.com/joho/godotenv"
)

var myEnvs map[string]string

func main() {
	myEnvs, err := godotenv.Read()
	if err != nil {
		fmt.Printf("Error loading .env file.")
	}

	client := got.NewTumblrRestClient(
		myEnvs["CONSUMER_KEY"],
		myEnvs["CONSUMER_SECRET"],
		myEnvs["OAUTH_TOKEN"],
		myEnvs["OAUTH_TOKEN_SECRET"],
		"",
		"http://api.tumblr.com")

	codeBlog := "codeblocks.tumblr.com"
	opts := map[string]string{}

	posts := client.Posts(codeBlog, "quote", opts)
	fmt.Println(posts.Total_posts)

	var allThePosts []got.QuotePost

	for _, elem := range posts.Posts {
		post := got.QuotePost{}
		json.Unmarshal(elem, &post)
		allThePosts = append(allThePosts, post)
	}

	fmt.Println(len(allThePosts))
	fmt.Println(allThePosts)
}

