package responses

type NodeFormDetailResponse struct {
	Template    FormTemplateFindAll         `json:"template"`
	Fields      []FormTemplateFieldsFindAll `json:"fields"`
	Data        []NodeFormDataResponse      `json:"data"`
	DataId      string                      `json:"dataId"`
	IsSubmitted bool                        `json:"isSubmitted"`
	IsApproved  bool                        `json:"isApproved"`
}

type JiraFormDetailResponse struct {
	Template FormTemplateFindAll         `json:"template"`
	Fields   []FormTemplateFieldsFindAll `json:"fields"`
	Data     []NodeFormDataResponse      `json:"data"`
}
