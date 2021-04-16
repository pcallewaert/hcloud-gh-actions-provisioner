# Kubernetes CronJob

This example will add 2 CronJobs to your kubernetes cluster. They are both very simular, but one is for scaling up the servers, and the other one for scaling down.
A secret is needed that contains the senstive configuration (HCLOUD_TOKEN, GITHUB_OWNER and GITHUB_PAT)

## Apply

1. Create a secret: `kubectl create secret generic hcloud-gh-actions-provisioner --from-literal=HCLOUD_TOKEN=<hcloud token> --from-literal=GITHUB_OWNER=<organisation owner> --from-literal=GITHUB_PAT=<github personal access token>`
2. Change the `cronjob.yaml` file with your configuration.
3. `kubectl apply -f cronjob.yaml`