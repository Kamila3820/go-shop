package files

import "mime/multipart"

type FileReq struct {
	File        *multipart.FileHeader `form:"file"`        // รองรับ File ตอน upload
	Destination string                `form:"destination"` // Path File ที่อยู่
	Extension   string                // นามสกุล File
	FileName    string
}

type FileRes struct {
	FileName string `json:"filename"`
	Url      string `json:"url"`
}

type DeleteFileReq struct {
	Destination string `json:"destination"`
}
