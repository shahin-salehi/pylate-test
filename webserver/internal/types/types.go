package types

import "time"


type Match struct {
	Filename string
	PageNumber int32
	Title string
	Category string
	Content string
	HTML string
	Score float32
	Meta string
}

type ViewerInstructions struct {
	Path string `json:"path"`
	Highlight string `json:"highlight"`
}

type File struct{
	PdfID int64 `json:"pdfID" db:"pdf_id"`
	Filename string `json:"filename"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}

type DeleteRequest struct {
	ID int64 `json:"id"`
}
