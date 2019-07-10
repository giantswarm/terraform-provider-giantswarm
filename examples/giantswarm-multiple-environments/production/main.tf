provider cloudflare {}
provider "giantswarm" {}
provider "aws" {}
provider "kubernetes" {
  host = "${module.giantswarm_cluster.api_endpoint}"

  client_certificate     = "${base64decode(module.giantswarm_cluster.client_cert)}"
  client_key             = "${base64decode(module.giantswarm_cluster.client_key)}"
  cluster_ca_certificate = "${base64decode(module.giantswarm_cluster.ca_cert)}"
}

module "cloudflare_dns" {
  source      = "../modules/cloudflare_dns"
  domain      = "${var.domain}"
  dns_record  = "www"
  dns_value = ["${replace(module.giantswarm_cluster.api_endpoint, "api", "service_name")}"]
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


# Database in staging and production using RDS upstream terraform module
data "aws_vpc" "default" {
  default = true
}

data "aws_subnet_ids" "all" {
  vpc_id = "${data.aws_vpc.default.id}"
}

data "aws_security_group" "default" {
  vpc_id = "${data.aws_vpc.default.id}"
  name   = "default"
}

module "db" {
  source  = "terraform-aws-modules/rds/aws"
  version = "1.0.0"

  identifier = "${var.environment}"
  engine            = "mysql"
  engine_version    = "5.7.19"
  instance_class    = "db.t2.large"
  allocated_storage = 5

  name     = "${var.environment}"
  username = "${var.app_db_user}"
  password = "${var.app_db_password}"
  port     = "3306"

  vpc_security_group_ids = ["${data.aws_security_group.default.id}"]

  maintenance_window = "Mon:00:00-Mon:03:00"
  backup_window      = "03:00-06:00"

  multi_az = true

  # disable backups to create DB faster
  backup_retention_period = 0

  # DB subnet group
  subnet_ids = "${data.aws_subnet_ids.all.ids}"

  # DB parameter group
  family = "mysql5.7"

  # DB option group
  major_engine_version = "5.7"

  # Database Deletion Protection
  deletion_protection = true

}