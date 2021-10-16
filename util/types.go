package util

import "time"

type RepoInfoResponse struct {
	Data struct {
		Repository struct {
			LicenseInfo struct {
				Key string
			}
			CreatedAt time.Time
		}
	}
}

type RepoInfo struct {
	License    string
	CreateDate time.Time
}
