package provider

import (
	"fmt"

	"github.com/giantswarm/gsclientgen/models"

	"github.com/giantswarm/gsctl/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGiantswarmApp() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				Description: "Custom name of the app to be installed",
				Computed:    false,
				ForceNew:    true,
			},
			"catalog": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the organization owning the cluster",
				ForceNew:    true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster id where the app should be installed",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				Description: "Name of the chart to be installed",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				Description: "Namespace where the app will be installed",
				ForceNew:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				ForceNew:    false,
				Description: "Version of the app",
			},
		},
		Create: resourceCreateApp,
		Read:   resourceReadApp,
		Update: resourceUpdateApp,
		Delete: resourceDeleteApp,
		Exists: resourceExistsApp,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateApp(d *schema.ResourceData, m interface{}) error {

	apiClient := m.(*client.Wrapper)

	catalog := d.Get("catalog").(string)
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	version := d.Get("namespace").(string)
	clusterID := d.Get("cluster_id").(string)
	appName := d.Get("app_name").(string)

	AppDefinitionSpec := &models.V4CreateAppRequestSpec{
		Catalog:   &catalog,
		Name:      &name,
		Namespace: &namespace,
		Version:   &version,
	}
	AppDefinition := &models.V4CreateAppRequest{
		Spec: AppDefinitionSpec,
	}

	createClusterActivityName := "create-app"

	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = createClusterActivityName

	_, err := apiClient.CreateApp(clusterID, appName, AppDefinition, auxParams)
	if err != nil {
		return fmt.Errorf("Error creating app %s", err)
	}

	//d.Set("app_name", result.Payload.Metadata.Name)

	// Wait for the status to be available.
	//stateConf := &resource.StateChangeConf{
	//	Pending: []string{"Pending"},
	//	Target:  []string{"Created"},
	//	Refresh: func() (interface{}, string, error) {
	//
	//		resp, err := apiClient.GetAppStatus(clusterID, appName, auxParams)
	//		if err != nil {
	//			log.Printf("Error on Cluster status refresh: %s", err)
	//			return nil, "", err
	//		}
	//
	//		status := "Pending"
	//
	//		if resp == "DEPLOYED" {
	//			status = "Created"
	//		}
	//
	//		return resp, status, nil
	//	},
	//	Timeout:        60 * time.Minute,
	//	Delay:          20 * time.Second,
	//	MinTimeout:     5 * time.Second,
	//	PollInterval:   15 * time.Second,
	//	NotFoundChecks: 3,
	//}
	//
	//_, stateErr := stateConf.WaitForState()
	//if stateErr != nil {
	//	return fmt.Errorf(
	//		"Error waiting for app (%s) to become ready: %s", result, stateErr)
	//}

	return nil

}

func resourceReadApp(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceUpdateApp(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceDeleteApp(d *schema.ResourceData, m interface{}) error {

	clusterID := d.Get("cluster_id").(string)
	appName := d.Get("app_name").(string)

	apiClient := m.(*client.Wrapper)

	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = "delete-app"

	_, err := apiClient.DeleteApp(clusterID, appName, auxParams)
	if err != nil {
		return err
	}

	return nil
}

func resourceExistsApp(d *schema.ResourceData, m interface{}) (bool, error) {

	return true, nil
}
