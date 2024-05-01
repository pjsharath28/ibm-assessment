# Nginx Deployment CLI Tool

A command-line tool written in Go to deploy Nginx server on a Kubernetes cluster.

## Features

- Deploy Nginx server to a Kubernetes cluster
- Scale the number of replicas
- Set the version of Nginx deployed
- Authentication using a kubeconfig file

## Prerequisites

- Go installed on your system
- Access to a Kubernetes cluster
- Knowledge of your cluster's kubeconfig file path

## Installation

Clone the repository:

```bash
git clone https://github.com/pjsharath28/ibm-assessment.git
```

Build the binary:

```bash
cd ibm-assessment
go build -o dist/deploynginx main.go
```

## Usage

```bash
./deploynginx --version <nginx_version> --scale <replica_count> --kubeconfig </path/to/kubeconfig> --namespace <namespace>
```

- `--version`: Version of Nginx to deploy.
- `--scale`: Number of replicas to scale to.
- `--kubeconfig`: Path to kubeconfig file.
- `--namespace`: Kubernetes namespace to deploy into (optional, defaults to "default").

## Example

```bash
./dist/deploynginx --version 1.13.10 --scale 3 --kubeconfig /Users/sharath/.kube/config
```

This command will deploy Nginx server with version 1.13.10, scale the number of replicas to 3, using the provided kubeconfig file, and deploy it into the default namespace.

## Testing

To run tests, use the following command:

```bash
go test -v ./...
```

