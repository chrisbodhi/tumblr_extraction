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

	type HugoPost struct {
		Title      string
		Date       string
		Tags       []string
		Categories []string
		Content    string
	}

	var hugoPosts []HugoPost

	client := got.NewTumblrRestClient(
		myEnvs["CONSUMER_KEY"],
		myEnvs["CONSUMER_SECRET"],
		myEnvs["OAUTH_TOKEN"],
		myEnvs["OAUTH_TOKEN_SECRET"],
		"",
		"http://api.tumblr.com")

	codeBlog := "codeblocks.tumblr.com"

	opts := map[string]string{}

	posts := client.Posts(codeBlog, "", opts)
	total := posts.Total_posts
	timesToReq := int(total / 20) + 1

	for i := 0; i <= timesToReq; i++ {

		opts["offset"] = fmt.Sprintf("%d", i * 20)
		posts := client.Posts(codeBlog, "", opts)

		for _, elem := range posts.Posts {
			postBase := got.BasePost{}
			json.Unmarshal(elem, &postBase)

			switch postType := postBase.PostType; postType {
			case "audio":
				audioPost := got.AudioPost{}
				json.Unmarshal(elem, &audioPost)

				hugoAudio := HugoPost{}

				hugoAudio.Title = audioPost.Caption
				hugoAudio.Date = audioPost.BasePost.Date
				hugoAudio.Tags = audioPost.BasePost.Tags
				hugoAudio.Categories = append(hugoAudio.Categories, "imported from tumblr", "audio")
				hugoAudio.Content = audioPost.Player

				hugoPosts = append(hugoPosts, hugoAudio)
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
				hugoPhoto.Content = fmt.Sprintf("![%s](%s) <br /> %s", photoPost.BasePost.Post_url, photoPost.Photos[0].Alt_sizes[0].Url, photoPost.Caption)

				hugoPosts = append(hugoPosts, hugoPhoto)
			case "quote":
				quotePost := got.QuotePost{}
				json.Unmarshal(elem, &quotePost)

				hugoQuote := HugoPost{}

				hugoQuote.Title = fmt.Sprintf("Quote for %s", quotePost.BasePost.Date)
				hugoQuote.Date = quotePost.BasePost.Date
				hugoQuote.Tags = quotePost.BasePost.Tags
				hugoQuote.Categories = append(hugoQuote.Categories, "imported from tumblr", "quote")
				hugoQuote.Content = fmt.Sprintf("%s", quotePost.Text)

				hugoPosts = append(hugoPosts, hugoQuote)
			case "text":
				textPost := got.TextPost{}
				json.Unmarshal(elem, &textPost)

				hugoText := HugoPost{}

				hugoText.Title = textPost.Title
				hugoText.Date = textPost.BasePost.Date
				hugoText.Tags = textPost.BasePost.Tags
				hugoText.Categories = append(hugoText.Categories, "imported from tumblr", "text")
				hugoText.Content = fmt.Sprintf("%s", textPost.Body)

				hugoPosts = append(hugoPosts, hugoText)
			case "video":
				videoPost := got.VideoPost{}
				json.Unmarshal(elem, &videoPost)

				hugoVideo := HugoPost{}

				hugoVideo.Title = videoPost.Caption
				hugoVideo.Date = videoPost.BasePost.Date
				hugoVideo.Tags = videoPost.BasePost.Tags
				hugoVideo.Categories = append(hugoVideo.Categories, "imported from tumblr", "video")
				hugoVideo.Content = videoPost.Player[0].Embed_code

				hugoPosts = append(hugoPosts, hugoVideo)
			default:
				fmt.Printf("\n%s\n", postType)
			}
		}
	}

//	fmt.Printf("%+v", hugoPosts)
	fmt.Println(len(hugoPosts))
}
