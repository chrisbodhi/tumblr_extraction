# E-Tumblr-L
## Extract Transform Load &mdash; from Tumblr to Hugo

A method of extracting data directly from Tumblr, formatting it, and saving it as Markdown files, suitable for rendering on a [Hugo](https://gohugo.io/) site.

### To Use

- Get Go installed, including the GOPATH and GOROOT set up.
- Clone this repo, however you prefer to clone open source Go projects.
- Copy the `dotenv.example` file to `.env` in the project root.
- Fill in the `.env` with the corresponding sekrit keys from Tumblr's developer console.
- While you're in the `.env` file, add the blog whose content you want to ETL.
- Inspect the `post.tmpl` file to ensure that the layout is how you want it to be on your blog.
- Run `go get github.com/MariaTerzieva/gotumblr` for interacting with the Tumblr API.
- Run `go get github.com/joho/godotenv` for loading a config stored in a `.env` file.
- Run `go run main.go` from the project root, then wait.
- Check out all of the Markdown files you now have!

### Very Important

The code in `main.go` is ugly af and whatever the opposite of DRY is. _I'm so sorry._ Maybe one day, I'll refactor it and/or add goroutines whenever I'm ready to learn about them. Pull requests and kind, constructive comments &amp; suggestions are always welcome, of course.

