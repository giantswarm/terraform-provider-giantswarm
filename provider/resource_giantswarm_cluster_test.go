package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/giantswarm/gsctl/client"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCluster_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckClusterBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleClusterExists("giantswarm_cluster.test"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test", "name", "Ferran test_cluster"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test", "owner", "giantswarm"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test", "release_version", "8.5.0"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test", "workers_min", "2"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test", "workers_max", "2"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test", "worker_azure_vm_size", "Standard_D2s_v3"),
				),
			},
		},
	})
}

func TestAccCluster_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckClusterUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleClusterExists("giantswarm_cluster.test2"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test2", "name", "Ferran test_cluster update"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test2", "owner", "giantswarm"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test2", "release_version", "8.5.0"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test2", "workers_min", "2"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test2", "workers_max", "2"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test2", "worker_azure_vm_size", "Standard_D2s_v3"),
				),
			},
			{
				Config: testAccCheckClusterUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleClusterExists("giantswarm_cluster.test2"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test_update", "name", "Ferran test_cluster update"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test_update", "owner", "giantswarm"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test_update", "release_version", "8.5.0"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test_update", "workers_min", "3"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test_update", "workers_max", "3"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test_update", "worker_azure_vm_size", "Standard_D2s_v3"),
				),
			},
		},
	})
}

var whiteSpaceRegex = regexp.MustCompile("name cannot contain whitespace")

func testAccCheckClusterDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Wrapper)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "giantswarm_cluster" {
			continue
		}

		auxParams := apiClient.DefaultAuxiliaryParams()
		auxParams.ActivityName = "read-cluster"

		_, err := apiClient.GetClusterV4(rs.Primary.ID, auxParams)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}

func testAccCheckExampleClusterExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		name := rs.Primary.ID
		apiClient := testAccProvider.Meta().(*client.Wrapper)
		auxParams := apiClient.DefaultAuxiliaryParams()
		auxParams.ActivityName = "read-cluster"
		_, err := apiClient.GetClusterV4(name, auxParams)
		if err != nil {
			return fmt.Errorf("error fetching Cluster with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckClusterBasic() string {
	return fmt.Sprintf(`
resource "giantswarm_cluster" "test" {
  owner = "giantswarm"
  name = "Ferran test_cluster"
  release_version= "8.5.0"
  workers_min = 2
  workers_max = 2
  worker_azure_vm_size = "Standard_D2s_v3"
}
`)
}

func testAccCheckClusterUpdatePre() string {
	return fmt.Sprintf(`
resource "giantswarm_cluster" "test2" {
	owner = "giantswarm"
	name = "Ferran test_cluster update"
	release_version= "8.5.0"
	workers_min = 2
	workers_max = 2
	worker_azure_vm_size = "Standard_D2s_v3"
}
`)
}

func testAccCheckClusterUpdatePost() string {
	return fmt.Sprintf(`
resource "giantswarm_cluster" "test2" {
	owner = "giantswarm"
	name = "Ferran test_cluster update"
	release_version= "8.5.0"
	workers_min = 3
	workers_max = 3
	worker_azure_vm_size = "Standard_D2s_v3"
}
`)
}
