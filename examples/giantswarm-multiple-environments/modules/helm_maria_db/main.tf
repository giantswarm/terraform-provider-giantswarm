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

resource "helm_release" "mariadb" {
    name      = "${var.name}"
    chart     = "stable/mariadb"
    version   = "${var.version}"

    set {
        name  = "mariadbUser"
        value = "${var.username}"
    }

    set {
        name = "mariadbPassword"
        value = "${var.password}"
    }
}

