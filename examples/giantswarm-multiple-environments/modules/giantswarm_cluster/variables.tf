
variable "owner" {
    description = "Name of the organization owning the cluster."
    default = "giantswarm"
}

variable "name" {
    description = "The name of the resource, also acts as it's unique ID"
}

variable "release_version" {
    description = "The release version to use in the new cluster. Modify to upgrade the cluster."
}

variable "availability_zones" {
    description = "Number of availability zones a cluster should be spread across. The default is provided via the info endpoint.Only in AWS."
    default = 1
}

variable "workers_min" {
    description = "Adjust the cluster node limits to make use of auto scaling or to have full control over the node count.In not in AWS then should be equal to workers_max"
    default = 3
}

variable "workers_max" {
    description = "Adjust the cluster node limits to make use of auto scaling or to have full control over the node count.In not in AWS then should be equal to workers_min."
    default = 3
}

variable "worker_aws_ec2_instance_type" {
    description = "EC2 instance type name. Must be the same for all worker nodes of a cluster.If not using AWS should not be set"
    default = ""
}

variable "worker_azure_vm_size" {
    description = "The domain of our web service."
    default = ""
}

variable "worker_num_cpus" {
    description = "Number of CPU cores."
    default = 0
}

variable "worker_storage_size" {
    description = "Node storage size in GB. Can be an integer or float."
    default = 0
}

variable "worker_memory_size" {
    description = ""
    default = 0
}