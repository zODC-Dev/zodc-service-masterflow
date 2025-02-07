package responses

type NodeResponse struct {
	Id       string `json:"id"`
	NodeType string `json:"nodeType"`
	Title    string `json:"title"`
	X        int32  `json:"x"`
	Y        int32  `json:"y"`
	DueIn    int32  `json:"dueIn"`
	EndType  string `json:"endType"`
}

type GroupResponse struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	X        int32  `json:"x"`
	Y        int32  `json:"y"`
	W        int32  `json:"w"`
	H        int32  `json:"h"`
	ParentId string `json:"parentId"`
	Type     string `json:"type"`
}

type ConnectionResponse struct {
	Id   string `json:"id"`
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

type WorkflowResponse struct {
	Title       string               `json:"title"`
	Type        string               `json:"type"`
	CategoryId  int32                `json:"categoryId"`
	Version     int32                `json:"version"`
	Description string               `json:"description"`
	Decoration  string               `json:"decoration"`
	Nodes       []NodeResponse       `json:"nodes"`
	Groups      []GroupResponse      `json:"groups"`
	Connections []ConnectionResponse `json:"connections"`
}
