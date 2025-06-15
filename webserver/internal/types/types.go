package types


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
