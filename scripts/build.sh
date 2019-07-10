#!/bin/bash

VERSION=$(cat version)
echo "building terraform-provider-giantswarm_${VERSION}"
go build -o terraform-provider-giantswarm_${VERSION}