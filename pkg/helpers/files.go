package helpers

import (
	"encoding/base64"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/config"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/servererrors"
)

func GetExtAndBase64FromJsonString(body string) (string, string, error) {
	split := strings.Split(body, ";")
	header := strings.Split(split[0], "/")
	extension := "." + header[len(header)-1]
	if len(split) < 2 {
		return "", "", servererrors.ErrorBase64
	}

	strBase64 := strings.Split(split[1], ",")[1]

	if len(strBase64) == 0 || len(extension) == 0 {
		return "", "", servererrors.ErrorBase64
	}

	return extension, strBase64, nil
}

func GetImageFileNameAndUrl(suffix string) (string, string) {
	cfg := config.GetInstance()

	id := uuid.New().String()
	prefix := strings.Split(id, "-")[0]
	if len(prefix) > 4 {
		prefix = prefix[:4]
	}
	name := "images/" + prefix + "/" + id + suffix
	return cfg.FilesPath + name, cfg.FilesUrlPrefix + name
}

func EnsureBaseDir(fpath string) error {
	baseDir := path.Dir(fpath)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, 0755)
}

func SaveBase64ToFile(fileName string, fileBase64Body string) error {
	decodedBuffer, err := base64.StdEncoding.DecodeString(fileBase64Body)
	if err != nil {
		return err
	}

	if err := EnsureBaseDir(fileName); err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(decodedBuffer); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}

	return nil
}

func DeleteFileFromUrl(url string) error {
	cfg := config.GetInstance()
	filePath := strings.Replace(url, cfg.FilesUrlPrefix, cfg.FilesPath, 1)
	return os.Remove(filePath)
}
