package file

//TODO mark as event somehow

type TmpFileUploaded struct {
	FileID string `json:"file_id"`
	Path   string `json:"path"`
}

type CommitFile struct {
	Id       string            `json:"id"`
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata"`
}

type FileCreated struct {
	TempID   string            `json:"temp_id"`
	PermID   string            `json:"perm_id"`
	Existed  bool              `json:"existed"`
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata"`
}
