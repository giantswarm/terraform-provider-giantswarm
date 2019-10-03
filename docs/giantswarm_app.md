# Giantswarm App
## Usage
```hcl
resource "giantswarm_app" "test_app"{
    cluster_id = "${giantswarm_cluster.test.id}"
  	app_name = "kong-app3"
	catalog = "giantswarm-incubator"
	name = "kong-app" 
	namespace = "kong"
	version = "0.2.0" 
}
```

## Attributes Reference

* `app_name` - (Forces new resource) Custom name for the app. Used also as a resource ID.
* `cluster_id` - (Forces new resource) Existent Cluster ID in an installation.
* `catalog` - (Forces new resource) The catalog where the chart for this app can be found.
* `name` - (Forces new resource) Name of the chart that should be used to install this app.
* `namespace` - (Forces new resource) Namespace that this app will be installed to.
* `version` - Version of the chart that should be used to install this app.
