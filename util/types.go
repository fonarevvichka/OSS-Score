package util

import (
	"time"
)

type RepoInfoResponse struct {
	Data struct {
		Repository struct {
			LatestRelease struct {
				CreatedAt time.Time
			}
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
							PageInfo PageInfo
						}
					}
				}
				PageInfo PageInfo
			}
		}
	}
}

type Dependency struct {
	PacakgeName   string
	NameWithOwner string
	Version       string
}

type IssueResponse struct {
	Data struct {
		Repository struct {
			Issues struct {
				Edges []struct {
					Node struct {
						Closed    bool
						CreatedAt time.Time
						ClosedAt  time.Time
						Assignees struct {
							TotalCount int
						}
						Participants struct {
							TotalCount int
						}
						Comments struct {
							TotalCount int
						}
					}
				}
				PageInfo PageInfo
			}
		}
	}
}

type OpenIssue struct {
	CreateDate time.Time

	Participants int
	Comments     int
	Assignees    int
}

type ClosedIssue struct {
	CreateDate time.Time
	CloseDate  time.Time

	Participants int
	Comments     int
}

type Issues struct {
	OpenIssues   []OpenIssue
	ClosedIssues []ClosedIssue
}

type RepoInfo struct {
	License        string
	CreateDate     time.Time
	LatestRealease time.Time
	Languages      []string
	Issues         Issues
	Dependencies   []Dependency
}

type QueryError struct {
	message string
}
