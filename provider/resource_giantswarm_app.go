package provider

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/giantswarm/gsclientgen/models"

	"github.com/giantswarm/gsctl/client"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceGiantswarmApp() *schema.Resource {
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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Cluster id where the app should be installed",
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
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
	version := d.Get("version").(string)
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

	createAppActivityName := "create-app"

	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = createAppActivityName

	_, err := apiClient.CreateApp(clusterID, appName, AppDefinition, auxParams)
	if err != nil {
		return fmt.Errorf("Error creating app %s", err)
	}

	// Wait for the status to be available.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Created"},
		Refresh: func() (interface{}, string, error) {

			resp, err := apiClient.GetAppStatus(clusterID, appName, auxParams)
			if err != nil {
				log.Printf("Error on App status refresh: %s", err)
				return nil, "", err
			}

			status := "Pending"

			if resp == "DEPLOYED" {
				status = "Created"
			}

			return resp, status, nil
		},
		Timeout:        60 * time.Minute,
		Delay:          20 * time.Second,
		MinTimeout:     5 * time.Second,
		PollInterval:   15 * time.Second,
		NotFoundChecks: 3,
	}

	resp, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for app (%s) to become ready: %s", resp, stateErr)
	}
	d.SetId(appName)

	return nil

}

func resourceReadApp(d *schema.ResourceData, m interface{}) error {

	apiClient := m.(*client.Wrapper)

	ClusterID := d.Get("cluster_id").(string)
	appName := d.Get("app_name").(string)

	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = "read-app"
	app, err := apiClient.GetApp(ClusterID, appName, auxParams)
	if err != nil {
		if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding applicaiton with name %s", appName)
		}
	}

	d.SetId(appName)
	d.Set("app_name", app.Metadata.Name)
	d.Set("catalog", app.Spec.Catalog)
	d.Set("name", app.Spec.Name)
	d.Set("namespace", app.Spec.Namespace)
	d.Set("version", app.Spec.Version)

	return nil
}

func resourceUpdateApp(d *schema.ResourceData, m interface{}) error {

	apiClient := m.(*client.Wrapper)
	version := d.Get("version").(string)
	clusterID := d.Get("cluster_id").(string)
	appName := d.Get("app_name").(string)

	AppDefinitionSpec := &models.V4ModifyAppRequestSpec{
		Version: version,
	}
	AppDefinition := &models.V4ModifyAppRequest{
		Spec: AppDefinitionSpec,
	}

	updateAppActivityName := "update-app"

	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = updateAppActivityName

	_, err := apiClient.ModifyApp(clusterID, appName, AppDefinition, auxParams)

	if err != nil {
		return fmt.Errorf("Error modifying app %s", err)
	}

	// Wait for the status to be available.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Created"},
		Refresh: func() (interface{}, string, error) {

			resp, err := apiClient.GetAppStatus(clusterID, appName, auxParams)
			if err != nil {
				log.Printf("Error on App status refresh: %s", err)
				return nil, "", err
			}

			status := "Pending"

			log.Printf("Status of the app: %s", resp)
			if resp == "DEPLOYED" {
				status = "Created"
			}

			return resp, status, nil
		},
		Timeout:        60 * time.Minute,
		Delay:          20 * time.Second,
		MinTimeout:     5 * time.Second,
		PollInterval:   15 * time.Second,
		NotFoundChecks: 3,
	}

	resp, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for app (%s) to become ready: %s", resp, stateErr)
	}

	d.SetId(appName)

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

	d.SetId("")
	return nil
}

func resourceExistsApp(d *schema.ResourceData, m interface{}) (bool, error) {

	apiClient := m.(*client.Wrapper)

	clusterID := d.Get("cluster_id").(string)
	appName := d.Id()

	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = "read-app"
	_, err := apiClient.GetApp(clusterID, appName, auxParams)
	if err != nil {
		if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
