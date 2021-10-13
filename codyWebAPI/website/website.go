package website

//Website ..
type Website interface {
	InitWebsite()
	SearchWebsite(item string) ([]string, error)
}
