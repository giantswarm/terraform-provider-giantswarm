provider "giantswarm" {}

# Using AWS example, remove Azure resource
resource "giantswarm_cluster" "test" {
  owner = "giantswarm" 
  name = "Ferran est_cluster" 
  release_version= "8.6.0"
  availability_zones = 1
  
  workers_min = 2 
  workers_max = 3 

  worker_aws_ec2_instance_type = "m5.large" 
  count = 0
}
resource "giantswarm_app" "test_app"{
      cluster_id = "${giantswarm_cluster.test_azure.id}"
  		app_name = "kong-app"
			catalog = "giantswarm-incubator"
			name = "kong-app" 
			namespace = "kong"
			version = "0.2.0" 
      count = 0
}

# Using Azure example, remove AWS resource
resource "giantswarm_cluster" "test_azure" {
  owner = "giantswarm"
  name = "Ferran tf test_cluster"
  release_version= "8.5.0"
  workers_min = 3
  workers_max = 3
  worker_azure_vm_size = "Standard_D2s_v3"
  count = 1
}