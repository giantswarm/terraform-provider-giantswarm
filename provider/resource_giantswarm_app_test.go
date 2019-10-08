package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/giantswarm/gsctl/client"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// Valid cluster where create app tests.
const clusterID = "pgut4"

func TestAccApp_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleAppExists("giantswarm_app.test_app"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_app", "app_name", "kong-app-basic"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_app", "cluster_id", clusterID),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_app", "version", "0.2.0"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_app", "namespace", "kong-basic"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_app", "catalog", "giantswarm-incubator"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_app", "name", "kong-app"),
				),
			},
		},
	})
}

func TestAccApp_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleAppExists("giantswarm_app.test_update"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_update", "app_name", "kong-app-update"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_update", "cluster_id", clusterID),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_update", "version", "0.1.0"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_update", "namespace", "kong-update"),
					resource.TestCheckResourceAttr(
						"giantswarm_cluster.test_update", "name", "kong-app"),
				),
			},
			{
				Config: testAccCheckAppUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleAppExists("giantswarm_app.test_update"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_update", "app_name", "kong-app-update"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_update", "cluster_id", clusterID),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_update", "version", "0.2.0"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_update", "namespace", "kong-update"),
					resource.TestCheckResourceAttr(
						"giantswarm_app.test_update", "name", "kong-app"),
				),
			},
		},
	})
}

func TestAccApp_Multiple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAppMultiple(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleAppExists("giantswarm_app.test_app_mult"),
					testAccCheckExampleAppExists("giantswarm_app.test_app_mult2"),
				),
			},
		},
	})
}

func testAccCheckAppDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*client.Wrapper)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "giantswarm_app" {
			continue
		}

		auxParams := apiClient.DefaultAuxiliaryParams()
		auxParams.ActivityName = "detroy-app"

		_, err := apiClient.GetApp(clusterID, rs.Primary.ID, auxParams)
		if err == nil {
			return fmt.Errorf("App still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}

func testAccCheckExampleAppExists(resource string) resource.TestCheckFunc {
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
		auxParams.ActivityName = "read-app"

		_, err := apiClient.GetApp(clusterID, name, auxParams)
		if err != nil {
			return fmt.Errorf("error fetching App with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckAppBasic() string {
	return fmt.Sprintf(`
resource "giantswarm_app" "test_app"{
	cluster_id = "0x39t"
	app_name = "kong-app-basic"
	catalog = "giantswarm-incubator"
	name = "kong-app" 
	namespace = "kong-basic"
	version = "0.2.0" 
}
`)
}

func testAccCheckAppUpdatePre() string {
	return fmt.Sprintf(`
resource "giantswarm_app" "test_update"{
	cluster_id = "2v36m"
	app_name = "kong-app-update"
	catalog = "giantswarm-incubator"
	name = "kong-app" 
	namespace = "kong-update"
	version = "0.1.0" 
}
`)
}

func testAccCheckAppUpdatePost() string {
	return fmt.Sprintf(`
resource "giantswarm_app" "test_update"{
	cluster_id = "2v36m"
	app_name = "kong-app-update"
	catalog = "giantswarm-incubator"
	name = "kong-app" 
	namespace = "kong-update"
	version = "0.2.0" 
}
`)
}

func testAccCheckAppMultiple() string {
	return fmt.Sprintf(`
resource "giantswarm_app" "test_app_mult"{
    cluster_id = "2v36m"
  	app_name = "kong-app-mult"
	catalog = "giantswarm-incubator"
	name = "kong-app" 
	namespace = "kong-mult"
	version = "0.2.0" 
}
resource "giantswarm_app" "test_app_mult2"{
    cluster_id = "2v36m"
  	app_name = "kong-app-mult2"
	catalog = "giantswarm-incubator"
	name = "kong-app" 
	namespace = "kong-mult2"
	version = "0.2.0" 
}
`)
}
