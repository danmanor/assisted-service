# Running Assisted-service Subsystem Tests on Kubernetes

## Overview

This document details the steps required to run subsystem tests for the Assisted-service deployed in two modes:
- REST-API mode
- Kube-API mode

## Deployment for Subsystem Tests

Assisted-service components will be deployed in a Kind cluster using the Podman provider. Necessary tools like `kind`, `strato-skipper`, and `waiting` will be installed automatically if not present.

## Running the tests

To test REST-API mode, run:

```bash
make subsystem-tests
```

To test kube-api mode, run:

```bash
ENABLE_KUBE_API=true make subsystem-tests
```

Optionally the following environment variables can be exported:

* `FOCUS="install_cluster"` - An optional flag used for [focused specs](https://onsi.github.io/ginkgo/#focused-specs) with regular expression.
* `SKIP="install_cluster"` - An optional flag to skip scopes with regular expressions.
* `VERBOSE=true` - An optional flag to print verbosed data.

## Quick Update and Test
To quickly update the service and run tests after making code changes, use:

```bash
skipper make patch-service
```

This command builds a new service image, pushes it to the container registry, and triggers a rollout of the updated service Deployment.
