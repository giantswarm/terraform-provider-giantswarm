provider "giantswarm" {}

# Using AWS example, remove Azure resource
resource "giantswarm_cluster" "test" {
  owner = "giantswarm" 
  name = "tf_test_cluster" 
  release_version= "8.4.0"
  availability_zones = 1
  
  workers_min = 4 
  workers_max = 7 

  worker_aws_ec2_instance_type = "m5.large" 
}
resource "giantswarm_app" "test_app"{
      cluster_id = "${giantswarm_cluster.test.id}"
  		app_name = "kong-app"
			catalog = "giantswarm-incubator"
			name = "kong-app" 
			namespace = "kong"
			version = "0.2.4" 
      count = 0
}

# Using Azure example, remove AWS resource
resource "giantswarm_cluster" "test_azure" {
  owner = "giantswarm"
  name = "Ferran tf test_cluster"
  release_version= "8.2.0"
  workers_min = 5
  workers_max = 5
  worker_azure_vm_size = "Standard_D2s_v3"
  count = 0
}