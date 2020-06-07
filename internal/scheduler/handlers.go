package scheduler

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/PuerkitoBio/goquery"
)

// HandleArchive handles the archiving of entities
func HandleArchive(job Job) {
	fmt.Printf("Scheduler | Archive -> %s\n", job.Data)

	// Get entity from database
	e := job.Scheduler.DB.GetBookmarkByHash(job.Data)
	if e.ID == 0 {
		fmt.Println("No entity found")
		return
	}

	// Fetch content from url
	res, err := job.Scheduler.fetch(e.Bookmark.URL)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create []bytes from body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	job.Archive = Archive{
		Body: body,
	}
	job.Entity = e
	go job.Scheduler.Next("download-sources", job)
}

// HandleDownloadSources handles the download of CSS and JS source files from archived
// entities
func HandleDownloadSources(job Job) {
	fmt.Printf("Scheduler | Download Sources -> %s\n", job.Data)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(job.Archive.Body))
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find("link[rel=\"stylesheet\"]").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			res, err := job.Scheduler.fetch(href, job.Entity.Bookmark.URL)
			if err != nil {
				fmt.Println(err)
				return
			}

			file, err := job.Scheduler.archiveSource(res, job.Data)
			s.SetAttr("href", file)
		}
	})

	html, _ := doc.Html()
	job.Archive.Body = []byte(html)
	go job.Scheduler.Next("save-html", job)
}

func HandleSaveHtml(job Job) {
	fmt.Printf("Scheduler | Saving HTML -> %s\n", job.Data)

	err := job.Scheduler.archiveHtml(job.Archive.Body, job.Data)
	if err != nil {
		fmt.Println(err)
		return
	}

	go job.Scheduler.Next("archived", job)
}

// HandleDownloadMeta handles the download of meta data
func HandleDownloadMeta(job Job) {
	fmt.Printf("Scheduler | Download Metadata -> %s\n", job.Data)

	// Get entity from database
	e := job.Scheduler.DB.GetBookmarkByHash(job.Data)
	if e.ID == 0 {
		fmt.Println("Not found")
		return
	}

	// Fetch content from url
	res, err := job.Scheduler.fetch(e.Bookmark.URL)
	if err != nil {
		fmt.Printf("Cannot download: %s\n", e.Bookmark.URL)
		return
	}

	// Create []bytes from body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return
	}

	title := doc.Find("title").First().Text()
	desc, _ := doc.Find("meta[name=\"description\"]").First().Attr("content")
	favicon, _ := doc.Find("link[rel=\"icon\"]").First().Attr("href")
	shortcut, _ := doc.Find("link[rel=\"shortcut icon\"]").Last().Attr("href")
	ogpTitle, _ := doc.Find("meta[property=\"og:title\"]").Last().Attr("content")
	// ogpImage, _ := doc.Find("meta[property=\"og:image\"]").First().Attr("content")
	ogpDesc, _ := doc.Find("meta[property=\"og:description\"]").First().Attr("content")

	if ogpTitle != "" {
		title = ogpTitle
	}

	if ogpDesc != "" {
		desc = ogpDesc
	}

	if favicon == "" {
		favicon = shortcut
	}

	job.Result = Result{
		Title:       title,
		Description: desc,
		Image:       favicon,
	}
	job.Entity = e

	res.Body.Close()
	go job.Scheduler.Next("download-image", job)
}

// HandleDownloadImage handles the download of images
func HandleDownloadImage(job Job) {
	fmt.Printf("Scheduler | Download Image -> %s\n", job.Data)
	if job.Result.Image == "" {
		fmt.Println("Image URL empty")
		go job.Scheduler.Next("save", job)
		return
	}

	res, err := job.Scheduler.fetch(job.Result.Image, job.Entity.Bookmark.URL)
	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Handle files like this 0b7e02e7.ico?v=JykvN0w9Ye
	file, err := job.Scheduler.saveImage(res, job.Data)
	if err != nil {
		fmt.Println(err)

		go job.Scheduler.Next("save", job)
		return
	}

	job.Result.Image = file
	go job.Scheduler.Next("save", job)
}

// HandleSave handles the saving of entities to the database
func HandleSave(job Job) {
	fmt.Printf("Scheduler | Saving -> %s\n", job.Data)
	e := job.Scheduler.DB.GetBookmarkByHash(job.Data)

	if e.Name == "" {
		e.Name = job.Result.Title
	}

	if e.Bookmark.Description == "" {
		e.Bookmark.Description = job.Result.Description
	}

	if e.Bookmark.ImageURL == "" {
		e.Bookmark.ImageURL = job.Result.Image
	}

	job.Scheduler.DB.SaveEntity(e)
}

func HandleArchived(job Job) {
	fmt.Printf("Scheduler | Archived -> %s\n", job.Data)
	job.Scheduler.DB.Archived(job.Data)
}
