package scheduler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func HandleArchive(job Job) {

}

func HandleMetaDownload(job Job) {
	fmt.Printf("Scheduler | Meta Download -> %s\n", job.Data)
	e := job.Scheduler.DB.GetBookmarkByHash(job.Data)
	if e.ID == 0 {
		return
	}

	url, err := url.Parse(e.Bookmark.URL)
	if err != nil {
		return
	}

	req := &http.Request{
		Method:     "GET",
		URL:        url,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	res, err := job.Scheduler.HTTPClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode > 203 {
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
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

	job.Work = "download-image"
	go job.Scheduler.Schedule(job)
}

func HandleImageDownload(job Job) {
	fmt.Printf("Scheduler | Image Download -> %s\n", job.Data)
	url, err := url.Parse(job.Result.Image)
	if err != nil {
		return
	}

	req := &http.Request{
		Method:     "GET",
		URL:        url,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	res, err := job.Scheduler.HTTPClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.StatusCode > 203 {
		return
	}

	file, err := job.Scheduler.saveImageToDisk(res, job.Data)
	if err != nil {
		fmt.Println(err)
		return
	}

	job.Result.Image = file
	job.Work = "save"
	go job.Scheduler.Schedule(job)
}

func HandleSave(job Job) {
	fmt.Printf("Scheduler | Save -> %s\n", job.Data)
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
