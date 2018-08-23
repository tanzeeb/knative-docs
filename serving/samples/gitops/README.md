# GitOps with Knative

Demonstrate GitOps with Knative Serving, Build, and Eventing.

## Prerequisites

* Knative installation
* Github account
* Dockerhub account

## Getting Started

1. Copy the [code/](./code/) folder into a new repository. This will be our `CODE_REPO`.
1. Copy the [config/](./config/) folder into a new repository. This will
   be our `CONFIG_REPO`.
1. Update `secrets.yaml` with your account details. Apply the secrets:
```
$ kubectl apply -f secrets.yaml
```
1. Apply a tag to the `CODE_REPO`:
```
$ git tag v1
$ git push --tags
```
1. Update the configuration in [deployment.yaml](./config/deployment.yaml) with your
   specific configuration. Commit and push the changes to `CONFIG_REPO`.
1. From the `CONFIG_REPO`, perform the initial deployment:
```
$ kubectl apply -f deployment.yaml
```

You should now be able to see your app running.

## Continuous Deployment

1. Make a change to the `CODE_REPO`, and tag the new commit `v2`. Push
   your commit.
1. Update the git tag in `deployment.yaml` to `v2`. Submit the change as
   a pull request.

Your deployed app should now be updated to `v2`, and the PR closed.
