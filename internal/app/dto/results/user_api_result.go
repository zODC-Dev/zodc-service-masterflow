package results

import "time"

type UserApiResult struct {
	Message string              `json:"message"`
	Data    []UserApiDataResult `json:"data"`
}

type UserApiDataResult struct {
	ID              int32                      `json:"id"`
	Email           string                     `json:"email"`
	Name            string                     `json:"name"`
	SystemRole      string                     `json:"systemRole"`
	IsActive        bool                       `json:"isActive"`
	CreatedAt       time.Time                  `json:"createdAt"`
	IsJiraLinked    bool                       `json:"isJiraLinked"`
	IsSystemUser    bool                       `json:"isSystemUser"`
	PermissionNames []string                   `json:"permissionNames"`
	ProjectRoles    []UserApiProjectRoleResult `json:"projectRoles"`
	AvatarUrl       string                     `json:"avatarUrl"`
}

type UserApiProjectRoleResult struct {
	ProjectKey      string   `json:"projectKey"`
	Role            string   `json:"role"`
	PermissionNames []string `json:"permissionNames"`
}
