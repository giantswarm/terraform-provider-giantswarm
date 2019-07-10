output "api_endpoint" {
    value = "${giantswarm_cluster.cluster.api_endpoint}"
}
output "client_cert" {
    value = "${giantswarm_cluster.cluster.client_cert}"
}
output "ca_cert" {
    value = "${giantswarm_cluster.cluster.ca_cert}"
}
output "client_key" {
    value = "${giantswarm_cluster.cluster.client_key}"
}