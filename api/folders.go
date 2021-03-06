package api

type CreateFolderRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentId    string `json:"parentId"`
}

type FolderResponse struct {
	CreatedAt   string           `json:"createdAt"`
	CreatedBy   string           `json:"createdBy"`
	ModifiedAt  string           `json:"modifiedAt"`
	ModifiedBy  string           `json:"modifiedBy"`
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	ItemType    string           `json:"itemType"`
	ParentId    string           `json:"parentId"`
	Permissions []string         `json:"permissions"`
	Description string           `json:"description"`
	Children    []folderChildren `json:"children"`
}

type GlobalFolderResponse struct {
	Id string `json:"id"`
}

type GlobalFolderResultResponse struct {
	Data []folderChildren `json:"data"`
}

type GlobalFolderStatusRequest struct {
	JobId string `json:"jobId"`
}

type GlobalFolderStatusResponse struct {
	Status        string `json:"status"`
	StatusMessage string `json:"statusMessage"`
}

type UpdateFolderRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type folderChildren struct {
	CreatedAt   string   `json:"createdAt"`
	CreatedBy   string   `json:"createdBy"`
	ModifiedAt  string   `json:"modifiedAt"`
	ModifiedBy  string   `json:"modifiedBy"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	ItemType    string   `json:"itemType"`
	ParentId    string   `json:"parentId"`
	Permissions []string `json:"permissions"`
}
