package scheduler

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	cnst "github.com/Techassi/gomark/internal/constants"
	"github.com/Techassi/gomark/internal/util"
)

func (s *Scheduler) saveImage(res *http.Response, hash string) (string, error) {
	// TODO: Check if the file is already saved and quit if yes
	save, hashed := s.savepath(res, cnst.FS_IMAGE_DIR, hash)
	os.MkdirAll(save, os.ModePerm)
	err := s.save(filepath.Join(save, hashed), res.Body)
	return hashed, err
}

func (s *Scheduler) archiveSource(res *http.Response, hash string) (string, error) {
	// TODO: Check if the file is already saved and quit if yes
	save, hashed := s.savepath(res, cnst.FS_ARCHIVE_DIR, hash)
	os.MkdirAll(save, os.ModePerm)
	err := s.save(filepath.Join(save, hashed), res.Body)
	return s.Config.BaseURL + "archive/" + hash + "/" + hashed, err
}

func (s *Scheduler) archiveHtml(in []byte, hash string) error {
	// TODO: Check if the file is already saved and quit if yes
	save := filepath.Join(s.Config.WebRoot, cnst.FS_ARCHIVE_DIR, hash)
	os.MkdirAll(save, os.ModePerm)
	err := s.save(filepath.Join(save, "index.html"), bytes.NewReader(in))
	return err
}

func (s *Scheduler) savepath(res *http.Response, dir, hash string) (string, string) {
	var (
		fileWithExt string = filepath.Base(res.Request.URL.String())
		fileExt     string = filepath.Ext(fileWithExt)
		fileName    string = strings.TrimSuffix(fileWithExt, fileExt)
		savePath    string = filepath.Join(s.Config.WebRoot, dir, hash)
		hashedName  string = util.Adler32(fileName)
	)
	return savePath, hashedName + fileExt
}

func (s *Scheduler) save(path string, in io.Reader) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	out.Close()

	return err
}
