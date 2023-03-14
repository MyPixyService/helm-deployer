# Helm Deployer

Helm Deployer is a command-line tool written in Go that reads a YAML file containing Helm chart definitions and deploys the charts to a Kubernetes cluster. This tool simplifies the process of managing and deploying multiple Helm charts by automating the process and allowing users to define releases in a single YAML file. This project is heaviley inspired by the [KubeOne Helm integration][ko-helm].

[ko-helm]: https://docs.kubermatic.com/kubeone/v1.6/guides/helm-integration/

## Installation

To install `helm-deployer`, follow these steps:

1. Download and run the installation script using `wget`:
```bash
wget -qO- https://raw.githubusercontent.com/MyPixyService/helm-deployer/main/scripts/install.sh | bash
```

This script will download the `helm-deployer` binary and install it in your system. Note that running a script directly from a URL can be a security risk. Be sure to verify the source of the script before running it. You can inspect the contents of the script by opening it in a text editor.

2. Once the installation script has completed, you should be able to use the `helm-deployer` command from your terminal:
```bash
helm-deployer -h
```
This should display the help text for the `helm-deployer` command.

> **Note**:
> Please note that the installation script is provided "as is" and may not work for all environments. Be sure to review and understand the contents of the script before running it on your system. We cannot be responsible for any damage caused by the installation script.

## Usage

To use Helm Deployer, you need to create a YAML file containing the definitions for the Helm releases you want to deploy. Here's an example YAML file:

```yaml
# Format:
#- releaseName: <release name>
#  chart: <chart name>
#  repoURL: <repository URL>
#  namespace: <namespace>
#  version: <chart version>
#  values: {}

- releaseName: traefik
  chart: traefik
  repoURL: https://traefik.github.io/charts
  namespace: kube-system
  values:
    service:
      annotations:
        load-balancer.hetzner.cloud/name: cluster-lb
```


Each entry in the `releases` list represents a Helm release. The `releaseName` field is the name of the release, `repoURL` is the URL of the Helm repository, `chart` is the name of the Helm chart, `version` is the version of the chart you want to deploy, `namespace` is the Kubernetes namespace you want to deploy the chart to, and `values` is a map of chart values that you want to set.

To deploy the Helm releases defined in the YAML file, run the following command:

```bash
helm-deployer -f inputfile.yaml -k kubeconfig
```

Replace `inputfile.yaml` with the path to your YAML file and `kubeconfig` with the path to your Kubernetes configuration file. Helm Deployer will read the YAML file, generate temporary values files for each release, and deploy the Helm charts to your Kubernetes cluster.

You can also filter your `inputfile.yaml` to only deploy charts with a matching `releaseName` field: 

```bash
helm-deployer -f inputfile.yaml -k kubeconfig -r traefik
```

To uninstall one or all releases specified in your `inputfile.yaml` add the `-uninstall` flag:

```bash
helm-deployer -f inputfile.yaml -k kubeconfig -uninstall
```

This also works with the `-r <releaseName>` flag.

For more information about the command-line flags supported by Helm Deployer, run the following command:

```bash
helm-deployer -h
```

## License

Helm Deployer is licensed under the MIT license. See the LICENSE file for more information.
