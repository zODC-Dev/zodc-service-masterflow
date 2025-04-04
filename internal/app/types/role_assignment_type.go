// Define a new type for role assignment request
package types

// RoleAssignmentRequest represents a request to assign a role to a user
type RoleAssignmentRequest struct {
	UserID     int32  `json:"user_id"`     // The ID of the user to assign the role to
	ProjectKey string `json:"project_key"` // The project key
	RoleName   string `json:"role_name"`   // The role name to assign
}

// RoleAssignmentResponse represents a response from role assignment
type RoleAssignmentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
