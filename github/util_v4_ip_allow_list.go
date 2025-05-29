package github

import (
	"github.com/shurcooL/githubv4"
)

type IpAllowListEntry struct {
	ID             githubv4.String
	Name           githubv4.String
	AllowListValue githubv4.String
	IsActive       githubv4.Boolean
	CreatedAt      githubv4.String
	UpdatedAt      githubv4.String
}

type IpAllowListEntries struct {
	Nodes      []IpAllowListEntry
	PageInfo   PageInfo
	TotalCount githubv4.Int
}
