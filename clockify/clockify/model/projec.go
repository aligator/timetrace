package model

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WorkspaceID string `json:"workspaceId"`
	Workspace   *Workspace
}

type CreateProject struct {
	Name     string `json:"name"`
	IsPublic bool   `json:"isPublic"`
}
