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

type RepoInfo struct {
	Name    string
	Owner   string
	Catalog string

	DefaultBranch string

	ActivityScore          float32
	DependencyActivtyScore float32
	ActivityConfidence     float32

	LicenseScore           float32
	DependencyLicenseScore float32
	LicenseConfidence      float32

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

type RepoInfoDBResponse struct {
	Ready    bool
	RepoInfo RepoInfo
}
