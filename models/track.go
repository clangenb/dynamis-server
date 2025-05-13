package models

type Track struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	FilePath string `json:"file_path"`
	Tier     string `json:"tier"`
}
