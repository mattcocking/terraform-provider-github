package github

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGithubOrganizationIpAllowList(t *testing.T) {
	// TODO: check these
	t.Run("creates and updates an IP allow list entry without error", func(t *testing.T) {
		config := fmt.Sprintf(`
			resource "github_organization_ip_allow_list" "test" {
				allow_list_value = "192.168.1.1/24"
				name            = "Test IP"
				is_active       = true
			}
		`)

		configUpdated := fmt.Sprintf(`
			resource "github_organization_ip_allow_list" "test" {
				allow_list_value = "192.168.1.1/24"
				name            = "Test IP Updated"
				is_active       = false
			}
		`)

		check := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"github_organization_ip_allow_list.test", "allow_list_value",
				"192.168.1.1/24",
			),
			resource.TestCheckResourceAttr(
				"github_organization_ip_allow_list.test", "name",
				"Test IP",
			),
			resource.TestCheckResourceAttr(
				"github_organization_ip_allow_list.test", "is_active",
				"true",
			),
		)

		checkUpdated := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				"github_organization_ip_allow_list.test", "allow_list_value",
				"192.168.1.1/24",
			),
			resource.TestCheckResourceAttr(
				"github_organization_ip_allow_list.test", "name",
				"Test IP Updated",
			),
			resource.TestCheckResourceAttr(
				"github_organization_ip_allow_list.test", "is_active",
				"false",
			),
		)

		testCase := resource.TestCase{
			PreCheck:  func() { testAccPreCheck(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config: config,
					Check:  check,
				},
				{
					Config: configUpdated,
					Check:  checkUpdated,
				},
			},
		}

		resource.Test(t, testCase)
	})
}
