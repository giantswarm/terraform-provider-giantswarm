provider "giantswarm" {}

# Using AWS example, remove Azure resource
resource "giantswarm_cluster" "test" {
  owner = "giantswarm" 
  name = "tf_test_cluster" 
  release_version= "8.1.0"
  availability_zones = 1
  
  workers_min = 4 
  workers_max = 7 

  worker_aws_ec2_instance_type = "m5.large" 
}

# Using Azure example, remove AWS resource
resource "giantswarm_cluster" "test" {
  owner = "giantswarm"
  name = "Ferran tf test_cluster"
  release_version= "8.0.0"
  workers_min = 5
  workers_max = 5
  worker_azure_vm_size = "Standard_D2s_v3"
}