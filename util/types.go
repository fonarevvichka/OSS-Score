package util

import "time"

type RepoInfoResponse struct {
	Data struct {
		Repository struct {
			LicenseInfo struct {
				Key string
			}
			CreatedAt time.Time
			Languages struct {
				Edges []struct {
					Node struct {
						Name string
					}
				}
			}
		}
	}
}

type RepoInfo struct {
	License      string
	CreateDate   time.Time
	Languages    []string
	Dependencies []struct {
		Dependency struct {
			Owner   string
			Name    string
			Url     string
			Version string
		}
	}
}
