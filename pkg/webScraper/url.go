package webcrawler

//URL ...
type URL struct {
	RootURL      string
	ParentURL    string
	CurrentURL   string
	CurrentDepth int
	MaxDepth     int
}
