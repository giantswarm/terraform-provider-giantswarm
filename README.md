# terraform-provider-giantswarm
This repo contains a terraform provider which permits create, update, upgrade and delete Giant Swarm kubernetes clusters and apps.

As Giant Swarm provides vainilla kubernetes, kubernetes/helm providers can be used along with this provider (see examples section).

## Installation
- clone the repo
- `make build`
- `cp terraform-provider-giantswarm_v0.1.0 $HOME/.terraform.d/plugins/YOUR_DISTRO/terraform-provider-giantswarm/`
- `terraform init`

## Usage
Go through the docs section in order to see attributes reference and examples section to see it in action.

### Prerequisites
- You need a user created in the Giant Swarm installation.
- Create a token associated to your user in order to login into the installation.
You can get one posting to the Giant Swarm API installation.

example:

```
#cat ferran.json
{
"email": "ferran@giantswarm.io",
"password_base64": "AfSeZ3Xds3weF5xQMQ=="
}
```
`curl -X POST -d @ferran.json https://GIANTSWARM_INSTALLATION_URL/v4/auth-tokens/`

- Once you have the token, export your credentials in `GIANTSWARM_INSTALLATION_ADDRESS` and `GIANTSWARM_TOKEN`

```
export GIANTSWARM_INSTALLATION_ADDRESS="https://GIANTSWARM_INSTALLATION_URL"
export GIANTSWARM_TOKEN=TOKEN_CREATED
```

### Import
`terraform import giantswarm_cluster.test YOUR_CLUSTER_ID`

### Examples

- Check the simplest example (creating an Azure/AWS giantswarm cluster) 
[example/simple-example](https://github.com/giantswarm/terraform-provider-giantswarm/tree/master/examples/simple-example)

- Use terraform to create a giant swarm kubernetes cluster with a kubernetes deployment and helm chart using the terraform kubernetes and helm upstream providers.

[example/giantswarm-cluster-kubernetes-helm](https://github.com/giantswarm/terraform-provider-giantswarm/tree/master/examples/giantswarm-cluster-kubernetes-helm)


- Use terraform to create a giant swarm kubernetes cluster in multiples environments.

[example/giantswarm-multiple-environments](https://github.com/giantswarm/terraform-provider-giantswarm/tree/master/examples/giantswarm-multiple-environments)
