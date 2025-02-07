package requests

type NodeRequest struct {
	Id       string
	NodeType string
	Title    string
	X        int32
	Y        int32
	DueIn    int32
	EndType  string
}

type GroupRequest struct {
	Id       string
	Title    string
	X        int32
	Y        int32
	W        int32
	H        int32
	ParentId string
	Type     string
}

type ConnectionRequest struct {
	Id   string
	From string
	To   string
	Type string
}

type WorkflowRequest struct {
	Title       string
	Type        string
	CategoryId  int32
	Version     int32
	Description string
	Decoration  string
	Nodes       []NodeRequest
	Groups      []GroupRequest
	Connections []ConnectionRequest
}
