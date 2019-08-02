# terraform-provider-giantswarm
This repo contains a terraform provider which permits create, update, upgrade and delete Giant Swarm kubernetes clusters.

As Giant Swarm provides vainilla kubernetes, kubernetes/helm providers can be used along with this provider (see examples section).

## Installation
- clone the repo
- `make build`
- `cp terraform-provider-giantswarm_v0.1.0 $HOME/.terraform.d/plugins/YOUR_DISTRO/terraform-provider-giantswarm/`
- `terraform init`

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

### Prerequisites
- You need a user created in the Giant Swarm installation.
- Create a token associated to your user in order to login into the installation.
You can get one posting to the Giant Swarm API installation.

example:

```
#cat ferran.json
{
"email": "ferran@giantswarm.io",
"password_base64": "AfSeZ33emF5xQMQ=="
}
```
`curl -X POST -d @ferran.json https://GIANTSWARM_INSTALLATION_URL/v4/auth-tokens/`

- Once you have the token, export your credentials in `GIANTSWARM_INSTALLATION_ADDRESS` and `GIANTSWARM_TOKEN`

```
export GIANTSWARM_INSTALLATION_ADDRESS="https://GIANTSWARM_INSTALLATION_URL"
export GIANTSWARM_TOKEN=TOKEN_CREATED
```

### Import
`terraform import giantswarm_cluster.test YOUR_CLUSTER_ID`

### Examples
- Use terraform to create a giant swarm kubernetes cluster with a kubernetes deployment and helm chart using the terraform kubernetes and helm upstream providers.

[example/giantswarm-cluster-kubernetes-helm](https://github.com/ferrandinand/terraform-provider-giantswarm/tree/master/examples/giantswarm-cluster-kubernetes-helm)


- Use terraform to create a giant swarm kubernetes cluster in multiples environments.

[example/giantswarm-multiple-environments](https://github.com/ferrandinand/terraform-provider-giantswarm/tree/master/examples/giantswarm-multiple-environments)
