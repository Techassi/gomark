package scheduler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"

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
	job.Work = "download-sources"
	go job.Scheduler.Schedule(job)
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

	root, err := url.Parse(job.Entity.Bookmark.URL)
	doc.Find("link[rel=\"stylesheet\"]").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			u, err := url.Parse(href)
			if !u.IsAbs() {
				u.Scheme = root.Scheme
				u.Host = root.Host
			}

			// TODO: Optimize this - dont convert back to string
			res, err := job.Scheduler.fetch(u.String())
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
	job.Work = "save-html"
	go job.Scheduler.Schedule(job)
}

func HandleSaveHtml(job Job) {
	fmt.Printf("Scheduler | Saving HTML -> %s\n", job.Data)

	err := job.Scheduler.archiveHtml(job.Archive.Body, job.Data)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// HandleDownloadMeta handles the download of meta data
func HandleDownloadMeta(job Job) {
	fmt.Printf("Scheduler | Download Metadata -> %s\n", job.Data)

	// Get entity from database
	e := job.Scheduler.DB.GetBookmarkByHash(job.Data)
	if e.ID == 0 {
		return
	}

	// Fetch content from url
	res, err := job.Scheduler.fetch(e.Bookmark.URL)
	if err != nil {
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
	ogpTitle, _ := doc.Find("meta[property=\"og:title\"]").First().Attr("content")
	ogpImage, _ := doc.Find("meta[property=\"og:image\"]").First().Attr("content")
	ogpDesc, _ := doc.Find("meta[property=\"og:description\"]").First().Attr("content")

	if ogpTitle != "" {
		title = ogpTitle
	}

	if ogpDesc != "" {
		desc = ogpDesc
	}

	job.Result = Result{
		Title:       title,
		Description: desc,
		Image:       ogpImage,
	}

	res.Body.Close()
	job.Work = "download-image"
	go job.Scheduler.Schedule(job)
}

// HandleDownloadImage handles the download of images
func HandleDownloadImage(job Job) {
	fmt.Printf("Scheduler | Download Image -> %s\n", job.Data)
	if job.Result.Image == "" {
		job.Work = "save"
		go job.Scheduler.Schedule(job)
		return
	}

	res, err := job.Scheduler.fetch(job.Result.Image)
	if err != nil {
		return
	}

	file, err := job.Scheduler.saveImage(res, job.Data)
	if err != nil {
		fmt.Println(err)
		return
	}

	job.Result.Image = file
	job.Work = "save"
	go job.Scheduler.Schedule(job)
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
