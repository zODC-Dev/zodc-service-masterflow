package types

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Metadata struct {
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	DeletedAt *string `json:"deletedAt"`
}

type Assignee struct {
	Id           int32  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	AvatarUrl    string `json:"avatarUrl"`
	IsSystemUser bool   `json:"isSystemUser"`
}

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
