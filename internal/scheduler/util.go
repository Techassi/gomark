package scheduler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	cnst "github.com/Techassi/gomark/internal/constants"
	"github.com/Techassi/gomark/internal/util"
)

func (s *Scheduler) saveImageToDisk(res *http.Response, hash string) (string, error) {
	// TODO: Check if the file is already saved and quit if yes
	save, hashed := s.imageSavepath(res, hash)
	os.MkdirAll(save, os.ModePerm)
	out, err := os.Create(save)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return "", err
	}
	out.Close()
	res.Body.Close()
	return hashed, nil

}

func (s *Scheduler) imageSavepath(res *http.Response, hash string) (string, string) {
	fileWithExt := filepath.Base(res.Request.URL.String())
	fileExt := filepath.Ext(fileWithExt)
	fileName := strings.TrimSuffix(fileWithExt, fileExt)
	savePath := filepath.Join(s.Config.WebRoot, cnst.FS_IMAGE_DIR, hash)
	hashedName := util.Adler32(fileName)
	return filepath.Join(savePath, hashedName+fileExt), hashedName + fileExt
}
