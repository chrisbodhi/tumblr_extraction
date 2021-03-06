package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	str "strings"
	"text/template"

	got "github.com/MariaTerzieva/gotumblr"
	"github.com/joho/godotenv"
)

var myEnvs map[string]string

type HugoPost struct {
	Title      string
	Date       string
	Tags       []string
	Categories []string
	Content    string
}

var hugoPosts []HugoPost

func createHugoFile(post HugoPost) {
	temp, err := template.ParseFiles("./post.tmpl")
	if err != nil {
		fmt.Println(err)
		return
	}

	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	lower_title := str.ToLower(post.Title)
	title := str.Trim(reg.ReplaceAllString(lower_title, "-"), "-")
	if len(title) > 100 {
		title = title[0:99]
	}

	file, err := os.Create(title + ".md")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}

	err = temp.Execute(file, post)
	if err != nil {
		fmt.Println("execute: ", err)
		return
	}

	file.Close()
}

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

	blog := myEnvs["TUMBLR_SITE"]

	opts := map[string]string{}

	posts := client.Posts(blog, "", opts)
	total := posts.Total_posts
	timesToReq := int(total/20) + 1

	for i := 0; i <= timesToReq; i++ {

		opts["offset"] = fmt.Sprintf("%d", i*20)
		posts := client.Posts(blog, "", opts)

		for _, elem := range posts.Posts {
			postBase := got.BasePost{}
			json.Unmarshal(elem, &postBase)

			switch postType := postBase.PostType; postType {
			case "audio":
				audioPost := got.AudioPost{}
				json.Unmarshal(elem, &audioPost)

				hugoAudio := HugoPost{}

				dateTime := str.Replace(audioPost.BasePost.Date, " GMT", "", 1)
				date := str.Replace(dateTime, " ", "T", 1)

				hugoAudio.Title = fmt.Sprintf("Video for %s", str.Split(audioPost.BasePost.Date, " ")[0])
				hugoAudio.Date = date
				hugoAudio.Tags = audioPost.BasePost.Tags
				hugoAudio.Categories = append(hugoAudio.Categories, "imported from tumblr", "audio")
				hugoAudio.Content = fmt.Sprintf("%s<br /><br />%s", audioPost.Player, audioPost.Caption)

				hugoPosts = append(hugoPosts, hugoAudio)
			case "link":
				linkPost := got.LinkPost{}
				json.Unmarshal(elem, &linkPost)

				hugoLink := HugoPost{}

				dateTime := str.Replace(linkPost.BasePost.Date, " GMT", "", 1)
				date := str.Replace(dateTime, " ", "T", 1)

				hugoLink.Title = linkPost.Title
				hugoLink.Date = date
				hugoLink.Tags = linkPost.BasePost.Tags
				hugoLink.Categories = append(hugoLink.Categories, "imported from tumblr", "link")
				hugoLink.Content = fmt.Sprintf("%s<br /><br />[goto](%s)", linkPost.Description, linkPost.Url)

				hugoPosts = append(hugoPosts, hugoLink)
			case "photo":
				photoPost := got.PhotoPost{}
				json.Unmarshal(elem, &photoPost)

				hugoPhoto := HugoPost{}

				dateTime := str.Replace(photoPost.BasePost.Date, " GMT", "", 1)
				date := str.Replace(dateTime, " ", "T", 1)

				hugoPhoto.Title = fmt.Sprintf("Photo for %s", str.Split(photoPost.BasePost.Date, " ")[0])
				hugoPhoto.Date = date
				hugoPhoto.Tags = photoPost.BasePost.Tags
				hugoPhoto.Categories = append(hugoPhoto.Categories, "imported from tumblr", "photo")
				hugoPhoto.Content = fmt.Sprintf("![%s](%s) <br /> %s", photoPost.BasePost.Post_url, photoPost.Photos[0].Alt_sizes[0].Url, photoPost.Caption)

				hugoPosts = append(hugoPosts, hugoPhoto)
			case "quote":
				quotePost := got.QuotePost{}
				json.Unmarshal(elem, &quotePost)

				hugoQuote := HugoPost{}

				dateTime := str.Replace(quotePost.BasePost.Date, " GMT", "", 1)
				date := str.Replace(dateTime, " ", "T", 1)

				hugoQuote.Title = fmt.Sprintf("Quote for %s", str.Split(quotePost.BasePost.Date, " ")[0])
				hugoQuote.Date = date
				hugoQuote.Tags = quotePost.BasePost.Tags
				hugoQuote.Categories = append(hugoQuote.Categories, "imported from tumblr", "quote")
				hugoQuote.Content = fmt.Sprintf("%s<br /><br />%s", quotePost.Text, quotePost.Source)

				hugoPosts = append(hugoPosts, hugoQuote)
			case "text":
				textPost := got.TextPost{}
				json.Unmarshal(elem, &textPost)

				hugoText := HugoPost{}

				dateTime := str.Replace(textPost.BasePost.Date, " GMT", "", 1)
				date := str.Replace(dateTime, " ", "T", 1)

				hugoText.Title = str.Replace(textPost.Title, "\"", "'", -1)
				hugoText.Date = date
				hugoText.Tags = textPost.BasePost.Tags
				hugoText.Categories = append(hugoText.Categories, "imported from tumblr", "text")
				hugoText.Content = fmt.Sprintf("%s", textPost.Body)

				hugoPosts = append(hugoPosts, hugoText)
			case "video":
				videoPost := got.VideoPost{}
				json.Unmarshal(elem, &videoPost)

				hugoVideo := HugoPost{}

				dateTime := str.Replace(videoPost.BasePost.Date, " GMT", "", 1)
				date := str.Replace(dateTime, " ", "T", 1)

				hugoVideo.Title = fmt.Sprintf("Video for %s", str.Split(videoPost.BasePost.Date, " ")[0])
				hugoVideo.Date = date
				hugoVideo.Tags = videoPost.BasePost.Tags
				hugoVideo.Categories = append(hugoVideo.Categories, "imported from tumblr", "video")
				hugoVideo.Content = fmt.Sprintf("%s<br /><br />%s", videoPost.Player[0].Embed_code, videoPost.Caption)

				hugoPosts = append(hugoPosts, hugoVideo)
			default:
				fmt.Printf("\n%s\n", postType)
			}
		}
	}

	//	fmt.Printf("%+v", hugoPosts)
	fmt.Println(len(hugoPosts))

	for _, post := range hugoPosts {
		createHugoFile(post)
	}
}
