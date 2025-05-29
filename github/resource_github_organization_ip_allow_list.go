package github

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shurcooL/githubv4"
)

func resourceGithubOrganizationIpAllowList() *schema.Resource {
	return &schema.Resource{
		Create: resourceGithubOrganizationIpAllowListCreate,
		Read:   resourceGithubOrganizationIpAllowListRead,
		Update: resourceGithubOrganizationIpAllowListUpdate,
		Delete: resourceGithubOrganizationIpAllowListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the IP allow list entry.",
			},
			"allow_list_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A single IP address or range of IP addresses in CIDR notation.",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the entry is currently active.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifies the date and time when the object was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Identifies the date and time when the object was updated.",
			},
		},
	}
}

func resourceGithubOrganizationIpAllowListCreate(d *schema.ResourceData, meta interface{}) error {
	err := checkOrganization(meta)
	if err != nil {
		return err
	}

	ctx := context.Background()
	owner := meta.(*Owner)

	client := owner.v4client
	orgName := owner.name

	// Get organization ID to pass to the mutation
	var orgQuery struct {
		Organization struct {
			ID githubv4.ID
		} `graphql:"organization(login: $login)"`
	}
	variables := map[string]interface{}{
		"login": githubv4.String(orgName),
	}
	err = client.Query(ctx, &orgQuery, variables)
	if err != nil {
		return err
	}

	// Create IP allow list entry
	var mutate struct {
		CreateIpAllowListEntry struct {
			IpAllowListEntry struct {
				ID githubv4.ID
			}
		} `graphql:"createIpAllowListEntry(input: $input)"`
	}
	input := githubv4.CreateIpAllowListEntryInput{
		OwnerID:        orgQuery.Organization.ID,
		AllowListValue: githubv4.String(d.Get("allow_list_value").(string)),
		IsActive:       githubv4.Boolean(d.Get("is_active").(bool)),
	}
	name := d.Get("name").(string)
	if name != "" {
		namePtr := githubv4.String(name)
		input.Name = &namePtr
	}
	err = client.Mutate(ctx, &mutate, input, nil)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s", mutate.CreateIpAllowListEntry.IpAllowListEntry.ID))
	return resourceGithubOrganizationIpAllowListRead(d, meta)
}

func resourceGithubOrganizationIpAllowListRead(d *schema.ResourceData, meta interface{}) error {
	err := checkOrganization(meta)
	if err != nil {
		return err
	}

	var query struct {
		Organization struct {
			IpAllowListEntries IpAllowListEntries `graphql:"ipAllowListEntries(first: 100, after:$cursor)"`
		} `graphql:"organization(login: $login)"`
	}

	orgName := meta.(*Owner).name
	variables := map[string]interface{}{
		"login":  githubv4.String(orgName),
		"cursor": (*githubv4.String)(nil),
	}

	client := meta.(*Owner).v4client
	for {
		err = client.Query(context.Background(), &query, variables)
		if err != nil {
			return err
		}

		for _, node := range query.Organization.IpAllowListEntries.Nodes {
			if node.ID == d.Id() {
				err := d.Set("name", node.Name)
				if err != nil {
					return err
				}
				err = d.Set("allow_list_value", node.AllowListValue)
				if err != nil {
					return err
				}
				err = d.Set("is_active", node.IsActive)
				if err != nil {
					return err
				}
				err = d.Set("created_at", node.CreatedAt)
				if err != nil {
					return err
				}
				err = d.Set("updated_at", node.UpdatedAt)
				return err
			}
		}

		if !query.Organization.IpAllowListEntries.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.NewString(query.Organization.IpAllowListEntries.PageInfo.EndCursor)
	}

	d.SetId("")
	return nil
}

func resourceGithubOrganizationIpAllowListUpdate(d *schema.ResourceData, meta interface{}) error {
	err := checkOrganization(meta)
	if err != nil {
		return err
	}

	var mutate struct {
		UpdateIpAllowListEntry struct {
			IpAllowListEntry struct {
				ID githubv4.ID
			}
		} `graphql:"updateIpAllowListEntry(input: $input)"`
	}

	input := githubv4.UpdateIpAllowListEntryInput{
		IPAllowListEntryID: d.Id(),
	}
	// TODO: validate this
	if d.HasChange("name") {
		namePtr := githubv4.String(d.Get("name").(string))
		input.Name = &namePtr
	}
	if d.HasChange("allow_list_value") {
		allowListValue := d.Get("allow_list_value").(string)
		input.AllowListValue = githubv4.String(allowListValue)
	}
	if d.HasChange("is_active") {
		isActive := d.Get("is_active").(bool)
		input.IsActive = githubv4.Boolean(isActive)
	}

	client := meta.(*Owner).v4client
	err = client.Mutate(context.Background(), &mutate, input, nil)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s", mutate.UpdateIpAllowListEntry.IpAllowListEntry.ID))
	return resourceGithubOrganizationIpAllowListRead(d, meta)
}

func resourceGithubOrganizationIpAllowListDelete(d *schema.ResourceData, meta interface{}) error {
	err := checkOrganization(meta)
	if err != nil {
		return err
	}

	var mutate struct {
		DeleteIpAllowListEntry struct {
			ClientMutationId githubv4.ID
		} `graphql:"deleteIpAllowListEntry(input: $input)"`
	}
	input := githubv4.DeleteIpAllowListEntryInput{
		IPAllowListEntryID: d.Id(),
	}

	client := meta.(*Owner).v4client
	err = client.Mutate(context.Background(), &mutate, input, nil)
	return err
}
