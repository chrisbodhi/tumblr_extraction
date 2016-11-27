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
	// Empty because not trying to filter anything
	// Not a huge blog dump that's going to happen here
	opts := map[string]string{}

	posts := client.Posts(codeBlog, "", opts)
	fmt.Println(posts.Total_posts)

	type HugoPost struct {
		Title      string
		Date       string
		Tags       []string
		Categories []string
		Content    string
	}

	var hugoPosts []HugoPost

	for _, elem := range posts.Posts {
		postBase := got.BasePost{}
		json.Unmarshal(elem, &postBase)

		switch postType := postBase.PostType; postType {
		case "link":
			linkPost := got.LinkPost{}
			json.Unmarshal(elem, &linkPost)

			hugoLink := HugoPost{}

			hugoLink.Title = linkPost.Title
			hugoLink.Date = linkPost.BasePost.Date
			hugoLink.Tags = linkPost.BasePost.Tags
			hugoLink.Categories = append(hugoLink.Categories, "imported from tumblr", "link")
			hugoLink.Content = fmt.Sprintf("[%s](%s)", linkPost.Description, linkPost.Url)

			hugoPosts = append(hugoPosts, hugoLink)
		case "photo":
			photoPost := got.PhotoPost{}
			json.Unmarshal(elem, &photoPost)

			hugoPhoto := HugoPost{}

			hugoPhoto.Title = fmt.Sprintf("Photo for %s", photoPost.BasePost.Date) 
			hugoPhoto.Date = photoPost.BasePost.Date
			hugoPhoto.Tags = photoPost.BasePost.Tags
			hugoPhoto.Categories = append(hugoPhoto.Categories, "imported from tumblr", "photo")
			hugoPhoto.Content = fmt.Sprintf("![](%s) <br /> %s", photoPost.Photos[0].Alt_sizes[0].Url, photoPost.Caption)

			hugoPosts = append(hugoPosts, hugoPhoto)
		case "quote":
			quotePost := got.QuotePost{}
			json.Unmarshal(elem, &quotePost)

			fmt.Printf("%+v", quotePost)

		case "text":
			textPost := got.TextPost{}
			json.Unmarshal(elem, &textPost)

		default:
			fmt.Println("\n%v\n", postType)
		}
	}

//	fmt.Printf("%+v", hugoPosts)
//	fmt.Println(len(hugoPosts))
}
