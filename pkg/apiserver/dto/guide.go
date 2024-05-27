package dto

type Guide struct {
	Name     string `json:"name"`
	Label    string `json:"label"`
	Click    bool   `json:"click"`
	Subtitle bool   `json:"subtitle"`
}

type GetGuideReq struct {
	Lang      string `json:"lang"`
	GuideName string `json:"guideName"`
}

type GetGuideRes struct {
	Guide []*Guide `json:"guide"`
}

type GetArticleReq struct {
	Lang        string `json:"lang"`
	GuideName   string `json:"guideName"`
	ArticleName string `json:"articleName"`
}
type GetArticleRes struct {
	Path string `json:"path"`
}
type GetFirstArticleNameReq struct {
	Lang      string `json:"lang"`
	GuideName string `json:"guideName"`
}
type GetFirstArticleNameRes struct {
	ArticleName string `json:"articleName"`
}
