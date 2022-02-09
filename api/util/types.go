package util

import (
	"time"
)

type PageInfo struct {
	HasNextPage bool
	EndCursor   string
}

type RepoInfoResponse struct {
	Data struct {
		Repository struct {
			StargazerCount   int
			DefaultBranchRef struct {
				Name string
			}
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

type CommitResponse struct {
	Data struct {
		Repository struct {
			Ref struct {
				Target struct {
					History struct {
						Edges []struct {
							Node struct {
								PushedDate time.Time
								Author     struct {
									Name string
								}
							}
						}
						PageInfo PageInfo
					}
				}
			}
		}
	}
}

type CommitResponseRest struct {
	Commit struct {
		Author struct {
			Name string
			Date time.Time
		}
	}
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
										Name  string
										Owner struct {
											Login string
										}
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
	Catalog string
	Owner   string
	Name    string
	Version string
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

type IssueResponseRest struct {
	State      string
	Assignees  []interface{}
	Comments   int
	Created_at time.Time
	Closed_at  time.Time
}

type ReleaseResponse struct {
	Data struct {
		Repository struct {
			Releases struct {
				Edges []struct {
					Node struct {
						CreatedAt time.Time
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

type Release struct {
	CreateDate time.Time
}

type Commit struct {
	PushedDate time.Time
	Author     string
}

type RepoRequestInfo struct {
	Name    string
	Owner   string
	Catalog string

	TimeFrame int // months
}

type RepoInfo struct {
	Name    string
	Owner   string
	Catalog string

	DefaultBranch string

	ScoreStatus       int //0 - not calcualted, 1 - queued, 2 ready
	RepoActivityScore Score
	RepoLicenseScore  Score

	DependencyActivityScore Score
	DependencyLicenseScore  Score

	UpdatedAt time.Time

	License       string
	CreateDate    time.Time
	LatestRelease time.Time
	Stars         int

	Releases     []Release
	Languages    []string
	Issues       Issues
	Dependencies []Dependency
	Commits      []Commit
}

type RepoInfoMessage struct {
	DataStatus int
	RepoInfo   RepoInfo
}

type Score struct {
	Score      float64
	Confidence float64
}
