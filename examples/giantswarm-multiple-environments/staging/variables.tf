variable "domain" {
    description = "Domain to be configured in the dns"
}

variable "g8s_nodes_type" {
    description = "Nodes instances to be configured"
    default = ""
}

variable "g8s_release_version" {
    description = "Giant Swarm release version"
    default = "8.0.0"
}

variable "environment" {
    description = "Environment name, will be used mainly to name resources."
}

variable "g8s_owner" {
    description = "Organization that will own the cluster."
}
variable "g8s_workers_min" {
    description = "Min nodes in kubernetes cluster"
}

variable "g8s_workers_max" {
    description = "Max nodes in kubernetes cluster."
}
variable "app_db_user" {}

variable "app_db_password" {}