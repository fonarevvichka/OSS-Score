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
type PageInfo struct {
	HasNextPage bool
	EndCursor   string
}
type DependencyResponse struct {
	Data struct {
		Repository struct {
			DependencyGraphManifests struct {
				TotalCount int
				Edges      []struct {
					Node struct {
						Dependencies struct {
							TotalCount int
							Edges      []struct {
								Node struct {
									PacakgeName  string
									Requirements string
									Repository   struct {
										NameWithOwner string
									}
								}
							}
							PageInfo
						}
					}
				}
				PageInfo
			}
		}
	}
}
type Dependency struct {
	PacakgeName   string
	NameWithOwner string
	Version       string
}

type RepoInfo struct {
	License      string
	CreateDate   time.Time
	Languages    []string
	Dependencies []Dependency
}

type QueryError struct {
	message string
}
