package model

import (
	"encoding/json"
	"fmt"
	"os"

	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/apiserver/dto"
	"gitlab.autocarat.de/gitlab-instance-9a4b3570/go-components/server/pkg/config"
)

func ReadGuideFromFile(lang, guideName string) (*dto.GetGuideRes, error) {
	cfg := config.GetInstance()
	data, err := os.ReadFile(fmt.Sprintf("%s%s/%s/guide.json", cfg.GuidePath, lang, guideName))
	if err != nil {
		return nil, err
	}
	var guide *dto.GetGuideRes
	err = json.Unmarshal(data, &guide)
	if err != nil {
		return nil, err
	}
	return guide, nil
}
func GetArticlePath(lang, guideName, articleName string) (*dto.GetArticleRes, error) {
	cfg := config.GetInstance()

	path := fmt.Sprintf("%sguide/%s/%s/%s.html", cfg.FilesUrlPrefix, lang, guideName, articleName)
	return &(dto.GetArticleRes{Path: path}), nil

}
func GetFirstArticleName(lang, guideName string) (*dto.GetFirstArticleNameRes, error) {
	guide, err := ReadGuideFromFile(lang, guideName)

	if err != nil {
		return nil, err
	}

	for _, value := range guide.Guide {
		if value.Click {
			return &dto.GetFirstArticleNameRes{ArticleName: value.Name}, nil
		}
	}
	return nil, fmt.Errorf("There ara no avalible values")
}
