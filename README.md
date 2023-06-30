# Blog
[link to blog](https://crazcalm.github.io/blog/)

## Setup
- Clone repo
- Install Theme
- Use Hugo Binary in Repo (bin/hugo-v50/hugo)

### Install Theme
The current theme is [Hugo-Geo](https://github.com/alexurquhart/hugo-geo). You can clone this repo under the `themes` directory.

### Hugo Binary 

Hugo has changed enough that the current version does not build this site as is. As such, I am using [hugo version 0.50](https://github.com/gohugoio/hugo/releases/tag/v0.50) to build this blog.

## Hugo commands:
### Draft a New Post:
```
hugo new posts/<name_of_post>.md
```

### Build blog with Drafts:
```
hugo server -D
```

### Publish Draft:
Open the draft post and set `draft = false`

Then run the below to build the site, which create all the needed resources for the blog in the `public` directory.

```
hugo
```

### Publish to Github Pages
```
bash deploy.sh
```

## Gotchas:
### Blog Date

The blogs are sorted by date. If there is no date in your post file, then the date will default to something really old and you will have paginate to the last page of the blog to find your post...

Make sure to add a date to you post.