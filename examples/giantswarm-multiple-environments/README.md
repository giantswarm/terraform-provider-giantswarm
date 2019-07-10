# Using Terraform to manage multiple environments
In this scenario we are going to setup an example which will manage different environments with different providers.
Take this as an example, not as a production ready setup.

The setup used will be:

- A terraform module that will include:
    - Giant Swarm managed kubernetes cluster.
    - Helm Chart that will deploy a mariadb helm chart for dev environments.
    - AWS RDS database instance for staging and production.
    - DNS entry in Cloudflare.

- Configurations for each environment with a main.tf with the providers configured and variables tfvars definition.


# Scenarios
We can have multiples scenarios, take this just as a one possibility with a learning scope.

- `dev1` is a test environment where we might want to have a fully dev environment available.

- `infra1` is a test environment where a team is testing a new Giant Swarm release.

- `staging` environment. Same as production environment but less nodes.

- `production` environment. 