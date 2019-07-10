resource "giantswarm_cluster" "cluster" {
  owner = "${var.owner}"
  name = "${var.name}"
  release_version = "${var.release_version}"
  availability_zones = "${var.availability_zones}"
  workers_min = "${var.workers_min}"
  workers_max = "${var.workers_max}"
  worker_aws_ec2_instance_type = "${var.worker_aws_ec2_instance_type}"
  worker_azure_vm_size = "${var.worker_azure_vm_size}"
  worker_num_cpus = "${var.worker_num_cpus}"
  worker_storage_size = "${var.worker_storage_size}"
  worker_memory_size = "${var.worker_memory_size}"
}
