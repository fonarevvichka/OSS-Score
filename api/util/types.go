package util

import (
	"time"
)

type PageInfo struct {
	HasNextPage bool
	EndCursor   string
}

// GraphQL Responses
//deprecated
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

// deprecated
type PullResponse struct {
	Data struct {
		Repository struct {
			PullRequests struct {
				Edges []struct {
					Node struct {
						Closed    bool
						Merged    bool
						CreatedAt time.Time
						ClosedAt  time.Time
					}
				}
				PageInfo PageInfo
			}
		}
	}
}

// deprecated
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
	CreatedAt time.Time

	Assignees int
}

type ClosedIssue struct {
	CreatedAt time.Time
	ClosedAt  time.Time
}

type Issues struct {
	OpenIssues   []OpenIssue
	ClosedIssues []ClosedIssue
}

type OpenPR struct {
	CreatedAt time.Time
}

type ClosedPR struct {
	CreatedAt time.Time
	ClosedAt  time.Time
}

type PullRequests struct {
	OpenPR   []OpenPR
	ClosedPR []ClosedPR
}

type Release struct {
	CreatedAt time.Time
}

type Commit struct {
	PushedDate time.Time
	Author     string
}

type Dependency struct {
	Catalog string
	Owner   string
	Name    string
	Version string
}

type RepoRequestInfo struct {
	Name    string
	Owner   string
	Catalog string

	TimeFrame int // months
}

type RepoInfo struct {
	Name    string `dynamodbav:"name"`
	Owner   string `dynamodbav:"owner"`
	Catalog string

	DefaultBranch string

	Status int //0 - not calculated, 1 - queued, 2 - pulled from queue, 3 - ready, 4 -error

	DataStartPoint time.Time
	UpdatedAt      time.Time

	License       string
	CreatedAt     time.Time
	LatestRelease time.Time
	Stars         int

	Releases     []Release
	Languages    []string
	Issues       Issues
	PullRequests PullRequests
	Dependencies []Dependency
	Commits      []Commit
}

type Score struct {
	Score      float64 `json:"score"`
	Confidence float64 `json:"confidence"`
}

type NameOwner struct {
	Owner string
	Name  string
}

type ScoreRequestBody struct {
	TimeFrame int `json:"timeFrame"`
}

type ScoreCategory struct {
	Min    float64
	Max    float64
	Weight float64
}

// Deprecated GraphQL Responses
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

// Deprecated REST Responses
type CommitResponseRest struct {
	Commit struct {
		Author struct {
			Name string
			Date time.Time
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

type PullResponseRest struct {
	State      string
	Created_at time.Time
	Closed_at  time.Time
}
