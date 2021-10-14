package util

type Repository struct {
	Repository RepoInfo
}

type RepoInfo struct {
	Name string
	Url  string
}

type Data struct {
	Data Repository
}
