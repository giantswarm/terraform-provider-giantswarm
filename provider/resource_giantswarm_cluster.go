package provider

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/giantswarm/gsclientgen/models"
	"github.com/giantswarm/gsctl/client"
	"github.com/giantswarm/gsctl/commands/types"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGiantswarmCluster() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the organization owning the cluster",
				ForceNew:    false,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The name of the resource, also acts as it's unique ID",
				ForceNew:    false,
			},
			"release_version": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "The release version to use in the new cluster",
			},
			"availability_zones": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "Number of availability zones a cluster should be spread across. The default is provided via the info endpoint.",
			},
			"workers_min": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "Adjust the cluster node limits to make use of auto scaling or to have full control over the node count. ",
			},
			"workers_max": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "Adjust the cluster node limits to make use of auto scaling or to have full control over the node count.",
			},
			"worker_num_cpus": {
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "Number of CPU cores",
			},
			"worker_storage_size": {
				Type:        schema.TypeFloat,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "Node storage size in GB. Can be an integer or float.",
			},
			"worker_memory_size": {
				Type:        schema.TypeFloat,
				Required:    false,
				Optional:    true,
				Description: "",
			},
			"worker_aws_ec2_instance_type": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "EC2 instance type name. Must be the same for all worker nodes of a cluster.",
			},
			"worker_azure_vm_size": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "Azure Virtual Machine size. Must be the same for all worker nodes of a cluster.",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL endpoint of the k8s api",
			},
			"ca_cert": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "certificate_authority_data",
			},
			"client_cert": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "client_certificate_data",
			},
			"client_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "client_key_data",
			},
		},
		Create: resourceCreateCluster,
		Read:   resourceReadCluster,
		Update: resourceUpdateCluster,
		Delete: resourceDeleteCluster,
		Exists: resourceExistsCluster,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// creates a models.V4AddClusterRequest from clusterDefinition
func createAddClusterBody(d types.ClusterDefinition) *models.V4AddClusterRequest {
	a := &models.V4AddClusterRequest{}
	a.AvailabilityZones = int64(d.AvailabilityZones)
	a.Name = d.Name
	a.Owner = &d.Owner
	a.ReleaseVersion = d.ReleaseVersion
	a.Scaling = &models.V4AddClusterRequestScaling{
		Min: d.Scaling.Min,
		Max: d.Scaling.Max,
	}

	if len(d.Workers) == 1 {
		ndmWorker := &models.V4AddClusterRequestWorkersItems{}
		ndmWorker.Memory = &models.V4AddClusterRequestWorkersItemsMemory{SizeGb: float64(d.Workers[0].Memory.SizeGB)}
		ndmWorker.CPU = &models.V4AddClusterRequestWorkersItemsCPU{Cores: int64(d.Workers[0].CPU.Cores)}
		ndmWorker.Storage = &models.V4AddClusterRequestWorkersItemsStorage{SizeGb: float64(d.Workers[0].Storage.SizeGB)}
		ndmWorker.Labels = d.Workers[0].Labels
		ndmWorker.Aws = &models.V4AddClusterRequestWorkersItemsAws{InstanceType: d.Workers[0].AWS.InstanceType}
		ndmWorker.Azure = &models.V4AddClusterRequestWorkersItemsAzure{VMSize: d.Workers[0].Azure.VMSize}
		a.Workers = append(a.Workers, ndmWorker)
	}

	return a
}

func resourceCreateCluster(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.WrapperV2)

	var creationResult struct {
		// cluster ID
		id string
		// location to fetch details on new cluster from
		location string
		// cluster definition assembled
		definition types.ClusterDefinition
	}
	result := creationResult

	createClusterActivityName := "create-cluster"

	workers := []types.NodeDefinition{}

	worker := types.NodeDefinition{
		CPU:     types.CPUDefinition{Cores: d.Get("worker_num_cpus").(int)},
		Storage: types.StorageDefinition{SizeGB: float32(d.Get("worker_storage_size").(float64))},
		Memory:  types.MemoryDefinition{SizeGB: float32(d.Get("worker_memory_size").(float64))},
		AWS:     types.AWSSpecificDefinition{InstanceType: d.Get("worker_aws_ec2_instance_type").(string)},
		Azure:   types.AzureSpecificDefinition{VMSize: d.Get("worker_azure_vm_size").(string)},
	}

	scaling := types.ScalingDefinition{
		Min: int64(d.Get("workers_min").(int)),
		Max: int64(d.Get("workers_max").(int)),
	}

	workers = append(workers, worker)

	clusterDefinition := types.ClusterDefinition{
		Name:              d.Get("name").(string),
		Owner:             d.Get("owner").(string),
		ReleaseVersion:    d.Get("release_version").(string),
		AvailabilityZones: d.Get("availability_zones").(int),
		Scaling:           scaling,
		Workers:           workers,
	}

	result.definition = clusterDefinition

	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = createClusterActivityName

	addClusterBody := createAddClusterBody(result.definition)
	response, err := apiClient.CreateCluster(addClusterBody, auxParams)

	if err != nil {
		return fmt.Errorf("Error creating cluster %s", err)
	}

	result.location = response.Location
	result.id = strings.Split(result.location, "/")[3]
	d.Set("id", result.id)
	d.SetId(result.id)

	// Wait for the status to be available.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Created"},
		Refresh: func() (interface{}, string, error) {

			// Sometimes api just has consistency issues and doesn't see
			// our cluster yet. Return an empty state.
			resp, err := apiClient.GetClusterStatus(result.id, auxParams)
			if err != nil {
				log.Printf("Error on Cluster status refresh: %s", err)
				if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
					clusterExists := false
					response, err := apiClient.GetClusters(auxParams)
					if err != nil {
						return nil, "", fmt.Errorf("error getting api endpoint for %s cluster", result.id)
					}
					for _, cluster := range response.Payload {
						if cluster.ID == result.id {
							clusterExists = true
						}
					}
					if !clusterExists {
						// Get some time to consolidate resource creation
						time.Sleep(15 * time.Second)
						return nil, "", nil
					}
				} else {
					return nil, "", err
				}
			}

			status := "Pending"
			if resp.Cluster.HasCreatedCondition() && !resp.Cluster.HasDeletedCondition() {
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

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for cluster (%s) to become ready: %s", result.id, stateErr)
	}

	//Get k8s api endpoint
	cluster, err := apiClient.GetCluster(result.id, auxParams)
	if err != nil {
		return fmt.Errorf("error getting api endpoint for %s cluster", result.id)
	}
	d.Set("api_endpoint", cluster.Payload.APIEndpoint)

	//Set keypair data
	err = createKeypair(m, d)
	if err != nil {
		return fmt.Errorf("error getting api endpoint for %s cluster", result.id)
	}

	return nil
}

func resourceReadCluster(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.WrapperV2)

	ClusterID := d.Id()

	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = "read-cluster"
	cluster, err := apiClient.GetCluster(ClusterID, auxParams)
	if err != nil {
		if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
			d.SetId("")
		} else {
			return fmt.Errorf("error finding Cluster with ID %s", ClusterID)
		}
	}

	numberZones := len(cluster.Payload.AvailabilityZones)

	d.SetId(ClusterID)
	d.Set("owner", cluster.Payload.Owner)
	d.Set("name", cluster.Payload.Name)
	d.Set("release_version", cluster.Payload.ReleaseVersion)
	d.Set("availability_zones", numberZones)
	d.Set("workers_min", cluster.Payload.Scaling.Min)
	d.Set("workers_max", cluster.Payload.Scaling.Max)
	d.Set("worker_num_cpus", cluster.Payload.Workers[0].CPU)
	d.Set("worker_storage_size", cluster.Payload.Workers[0].Storage)
	d.Set("worker_memory_size", cluster.Payload.Workers[0].Memory)

	if cluster.Payload.Workers[0].Aws != nil {
		d.Set("worker_aws_ec2_instance_type", cluster.Payload.Workers[0].Aws.InstanceType)
	}
	if cluster.Payload.Workers[0].Azure != nil {
		d.Set("worker_azure_vm_size", cluster.Payload.Workers[0].Azure.VMSize)
	}
	d.Set("api_endpoint", cluster.Payload.APIEndpoint)

	//Set keypair data
	err = createKeypair(m, d)
	if err != nil {
		return fmt.Errorf("error getting api endpoint for %s cluster", ClusterID)
	}

	return nil
}

func resourceUpdateCluster(d *schema.ResourceData, m interface{}) error {

	apiClient := m.(*client.WrapperV2)

	ClusterID := d.Id()
	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = "update-cluster"

	//Update scale cluster settings
	scaling := models.V4ModifyClusterRequestScaling{
		Min: int64(d.Get("workers_min").(int)),
		Max: int64(d.Get("workers_max").(int)),
	}

	clusterDefinition := &models.V4ModifyClusterRequest{
		Name:           d.Get("name").(string),
		Owner:          d.Get("owner").(string),
		ReleaseVersion: d.Get("release_version").(string),
		Scaling:        &scaling,
	}
	_, err := apiClient.ModifyCluster(ClusterID, clusterDefinition, auxParams)
	if err != nil {
		return fmt.Errorf("error upgrading %s cluster", ClusterID)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Updated"},
		Refresh: func() (interface{}, string, error) {

			// Sometimes api just has consistency issues and doesn't see
			// our cluster yet. Return an empty state.
			resp, err := apiClient.GetClusterStatus(ClusterID, auxParams)
			if err != nil {
				log.Printf("Error on Cluster status refresh: %s", err)
				if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
					response, err := apiClient.GetClusters(auxParams)
					if err != nil {
						return nil, "", fmt.Errorf("error getting api endpoint for %s cluster", ClusterID)
					}

					clusterExists := false
					for _, cluster := range response.Payload {
						if cluster.ID == ClusterID {
							clusterExists = true
						}
					}
					if !clusterExists {
						// Get some time to consolidate resource creation
						time.Sleep(15 * time.Second)
						return nil, "", nil
					}
				} else {
					return nil, "", err
				}
			}

			status := "Pending"
			if resp.Cluster.HasUpdatingCondition() {
				status = "Pending"
			}
			if resp.Cluster.HasUpdatedCondition() {
				status = "Updated"
			}

			return resp, status, nil
		},
		Timeout:        120 * time.Minute,
		Delay:          20 * time.Second,
		MinTimeout:     5 * time.Second,
		PollInterval:   15 * time.Second,
		NotFoundChecks: 3,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for cluster (%s) to become ready: %s", ClusterID, stateErr)
	}

	return nil
}

func resourceDeleteCluster(d *schema.ResourceData, m interface{}) error {

	clusterID := d.Id()

	apiClient := m.(*client.WrapperV2)

	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = "delete-cluster"

	_, err := apiClient.DeleteCluster(clusterID, auxParams)
	if err != nil {
		return err
	}

	// Wait for the status to be available.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Deleted"},
		Refresh: func() (interface{}, string, error) {
			status := "Pending"
			resp, err := apiClient.GetClusterStatus(clusterID, auxParams)
			if err != nil {
				log.Printf("Error on Cluster status refresh: %s", err)
				if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
					response, err := apiClient.GetClusters(auxParams)
					if err != nil {
						return nil, "", fmt.Errorf("error getting api endpoint for %s cluster", clusterID)
					}

					//Check if the cluster has been already deleted
					clusterExists := false
					for _, cluster := range response.Payload {
						if cluster.ID == clusterID {
							clusterExists = true
						}
					}
					if !clusterExists {
						status = "Deleted"
					}
				}
			}
			if resp.Cluster.HasDeletedCondition() || resp.Cluster.HasDeletingCondition() {
				status = "Deleted"
			}

			return resp, status, nil
		},
		Timeout:        60 * time.Minute,
		Delay:          20 * time.Second,
		MinTimeout:     5 * time.Second,
		PollInterval:   15 * time.Second,
		NotFoundChecks: 3,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf("Error waiting for cluster (%s) to be deleted: %s", clusterID, err)
	}

	return nil
}

func resourceExistsCluster(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.WrapperV2)

	clusterID := d.Id()

	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = "read-cluster"
	_, err := apiClient.GetCluster(clusterID, auxParams)
	if err != nil {
		if strings.Contains(err.Error(), "RESOURCE_NOT_FOUND") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func createKeypair(m interface{}, d *schema.ResourceData) error {

	apiClient := m.(*client.WrapperV2)
	clusterID := d.Id()

	description := "Generated by Terraform for cluster " + clusterID
	addKeyPairBody := &models.V4AddKeyPairRequest{
		Description:              &description,
		TTLHours:                 0,
		CnPrefix:                 "",
		CertificateOrganizations: "system:masters," + d.Get("owner").(string),
	}

	auxParams := apiClient.DefaultAuxiliaryParams()
	auxParams.ActivityName = "create-keypair"

	//Check if keypair already exists
	keyPairExists := false
	responseGet, err := apiClient.GetKeyPairs(clusterID, auxParams)
	if err != nil {
		return fmt.Errorf("Error getting keypair for: %s %s", clusterID, err)
	}

	for _, keyPair := range responseGet.Payload {
		if keyPair.Description == description {
			keyPairExists = true
		}
	}

	if !keyPairExists {
		responseCreate, err := apiClient.CreateKeyPair(clusterID, addKeyPairBody, auxParams)
		if err != nil {
			return fmt.Errorf("Error creating keypair for: %s %s", clusterID, err)
		}

		// store credentials into vars
		d.Set("ca_cert", base64.StdEncoding.EncodeToString([]byte(responseCreate.Payload.CertificateAuthorityData)))
		d.Set("client_cert", base64.StdEncoding.EncodeToString([]byte(responseCreate.Payload.ClientCertificateData)))
		d.Set("client_key", base64.StdEncoding.EncodeToString([]byte(responseCreate.Payload.ClientKeyData)))
	}
	return nil
}
