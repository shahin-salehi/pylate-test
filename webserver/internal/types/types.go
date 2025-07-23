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
	FileUrl string
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

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID int64 `db:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `db:"password_hash" json:"password_hash"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Key string

const (
	ContextUser  Key = "user_id"
	ContextGroup Key = "group_id"
)

type NewPDF struct{
	Url string `json:"url"`
	Filename string `json:"filename"`
	Category string `json:"category"`
	Owner int64 `json:"owner"`
}
