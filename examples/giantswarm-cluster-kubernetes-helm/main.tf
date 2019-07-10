provider "giantswarm" {}

resource "giantswarm_cluster" "test" {
  owner = "giantswarm"
  name = "tf_test_cluster"
  release_version= "8.0.0"
  workers_min = 4
  workers_max = 4
  //worker_aws_ec2_instance_type = "m5.large"
  worker_azure_vm_size = "Standard_D2s_v3"
  count = 1
}

provider "kubernetes" {
  host = "${giantswarm_cluster.test.api_endpoint}"

  client_certificate     = "${base64decode(giantswarm_cluster.test.client_cert)}"
  client_key             = "${base64decode(giantswarm_cluster.test.client_key)}"
  cluster_ca_certificate = "${base64decode(giantswarm_cluster.test.ca_cert)}"
}

provider "helm" {
    kubernetes {
        host = "${giantswarm_cluster.test.api_endpoint}"

        client_certificate     = "${base64decode(giantswarm_cluster.test.client_cert)}"
        client_key             = "${base64decode(giantswarm_cluster.test.client_key)}"
        cluster_ca_certificate = "${base64decode(giantswarm_cluster.test.ca_cert)}"
    }
    service_account = "tiller"
}
