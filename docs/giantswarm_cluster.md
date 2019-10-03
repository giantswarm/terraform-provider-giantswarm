# Giantswarm Cluster
## Usage
```hcl
resource "giantswarm_cluster" "test" {
  owner = "giantswarm" 
  name = "tf_test_cluster" 
  release_version= "8.1.0"
  availability_zones = 1
  
  workers_min = 3 
  workers_max = 3 

  worker_aws_ec2_instance_type = "m5.large" 
}
```

## Attributes Reference

* `owner` - (Required) Organization present in the cluster.
* `name` - Name of the cluster.
* `release_version` - The release version to use in the new cluster.Change to upgrade a cluster.
* `availability_zones` - (Forces new resource) Number of availavility zones. Available only in AWS.
* `workers_min` - Adjust the cluster node limits to make use of auto scaling or to have full control over the node count. If not in AWS then workers_min value should be equal to workers_max.
* `workers_max` - Adjust the cluster node limits to make use of auto scaling or to have full control over the node count. If not in AWS then workers_max value should be equal to workers_min.
* `worker_num_cpus` - Number of CPU cores.
* `worker_storage_size` - Node storage size in GB. Can be an integer or float.
* `worker_memory_size` - Memory in GB.

* `worker_aws_ec2_instance_type` - (Forces new resource) Only in AWS. EC2 instance type name. Must be the same for all worker nodes of a cluster.
If you are using AWS, worker_azure_vm_size should not be defined

* `worker_azure_vm_size` - (Forces new resource) Only in Azure.Azure Virtual Machine size. Must be the same for all worker nodes of a cluster.
If you are using Azure worker_aws_ec2_instance_type should not be defined.

If you are using on-prem KVM worker_aws_ec2_instance_type and worker_azure_vm_size should not be defined.

* `api_endpoint` - (Computed) URL endpoint of the k8s api.
* `ca_cert` - (Computed) Certificate authority data created.
* `client_cert` - (Computed) Client certificate created.
* `client_key` - (Computed) Client key created.
