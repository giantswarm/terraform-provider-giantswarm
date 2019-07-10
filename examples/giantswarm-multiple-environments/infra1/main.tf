provider cloudflare {}
provider "giantswarm" {}
provider "kubernetes" {
  host = "${module.giantswarm_cluster.api_endpoint}"

  client_certificate     = "${base64decode(module.giantswarm_cluster.client_cert)}"
  client_key             = "${base64decode(module.giantswarm_cluster.client_key)}"
  cluster_ca_certificate = "${base64decode(module.giantswarm_cluster.ca_cert)}"
}

provider "helm" {
    kubernetes {
        host = "${module.giantswarm_cluster.api_endpoint}"

        client_certificate     = "${base64decode(module.giantswarm_cluster.client_cert)}"
        client_key             = "${base64decode(module.giantswarm_cluster.client_key)}"
        cluster_ca_certificate = "${base64decode(module.giantswarm_cluster.ca_cert)}"
    }
    service_account = "tiller"
}


resource "kubernetes_service_account" "tiller" {
  metadata {
    name = "tiller"
    namespace = "kube-system"
  }
}

resource "kubernetes_cluster_role_binding" "example" {
    metadata {
        name = "tiller-clusterrolebinding"
    }
    role_ref {
        api_group = "rbac.authorization.k8s.io"
        kind = "ClusterRole"
        name = "cluster-admin"
    }
    subject {
        kind = "ServiceAccount"
        name = "tiller"
        namespace = "kube-system"
    }
}


module "cloudflare_dns" {
  source      = "../modules/cloudflare_dns"
  domain      = "${var.domain}"
  dns_record  = "www"
  dns_value = "${replace(module.giantswarm_cluster.api_endpoint, "api", "service_name")}"
  record_type  = "CNAME"
  record_ttl  = 3600
}

module "giantswarm_cluster" {
  source      = "../modules/giantswarm_cluster"

  owner = "${var.g8s_owner}"
  name = "${var.environment}"
  release_version = "${var.g8s_release_version}"
  workers_min = "${var.g8s_workers_min}"
  workers_max = "${var.g8s_workers_max}"
  worker_aws_ec2_instance_type =  "${var.g8s_nodes_type}"
}


module "mariadb" {
  source      = "../modules/helm_maria_db"

  name = "${var.environment}"
  username = "${var.app_db_user}"
  password = "${var.app_db_password}"
}

