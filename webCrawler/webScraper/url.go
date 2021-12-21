package webcrawler

//URL ...
type URL struct {
	ParentURL    string
	URL          string
	CurrentDepth int
	MaxDepth     int
}
