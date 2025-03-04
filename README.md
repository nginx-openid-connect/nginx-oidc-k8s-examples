
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/nginxinc/kubernetes-ingress/badge)](https://api.securityscorecards.dev/projects/github.com/nginxinc/kubernetes-ingress) [![CI](https://github.com/nginxinc/kubernetes-ingress/actions/workflows/ci.yml/badge.svg)](https://github.com/nginxinc/kubernetes-ingress/actions/workflows/ci.yml)  [![FOSSA Status](https://app.fossa.com/api/projects/custom%2B5618%2Fgithub.com%2Fnginxinc%2Fkubernetes-ingress.svg?type=shield)](https://app.fossa.com/projects/custom%2B5618%2Fgithub.com%2Fnginxinc%2Fkubernetes-ingress?ref=badge_shield)  [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)  [![Go Report Card](https://goreportcard.com/badge/github.com/nginxinc/kubernetes-ingress)](https://goreportcard.com/report/github.com/nginxinc/kubernetes-ingress)  [![codecov](https://codecov.io/gh/nginxinc/kubernetes-ingress/branch/main/graph/badge.svg?token=snCn7Y0zC7)](https://codecov.io/gh/nginxinc/kubernetes-ingress)  [![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nginxinc/kubernetes-ingress?logo=github&sort=semver)](https://github.com/nginxinc/kubernetes-ingress/releases/latest)  ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nginxinc/kubernetes-ingress?logo=go)  [![Docker Pulls](https://img.shields.io/docker/pulls/nginx/nginx-ingress?logo=docker&logoColor=white)](https://hub.docker.com/r/nginx/nginx-ingress)  ![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/nginx/nginx-ingress?logo=docker&logoColor=white&sort=semver)  [![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/nginx-ingress)](https://artifacthub.io/packages/container/nginx-ingress/kubernetes-ingress)  [![Slack](https://img.shields.io/badge/slack-%23nginx--ingress--controller-green?logo=slack)](https://nginxcommunity.slack.com/channels/nginx-ingress-controller)

# 🚀 *Help make the NGINX Ingress Controller better by participating in our [survey](https://forms.office.com/Pages/ResponsePage.aspx?id=L_093Ttq0UCb4L-DJ9gcUKLQ7uTJaE1PitM_37KR881UMEs0Rk5PMkYzMTJTWVA0V1hUVTRLUUMyNS4u)!* 🚀

# NGINX Ingress Controller

This repo provides an implementation of an Ingress Controller for NGINX and NGINX Plus.

**Note**: this project is different from the NGINX Ingress Controller in [kubernetes/ingress-nginx](https://github.com/kubernetes/ingress-nginx) repo. See [this doc](https://docs.nginx.com/nginx-ingress-controller/intro/nginx-ingress-controllers) to find out about the key differences.

## What is the Ingress?

The Ingress is a Kubernetes resource that lets you configure an HTTP load balancer for applications running on Kubernetes, represented by one or more [Services](https://kubernetes.io/docs/concepts/services-networking/service/). Such a load balancer is necessary to deliver those applications to clients outside of the Kubernetes cluster.

The Ingress resource supports the following features:
* **Content-based routing**:
    * *Host-based routing*. For example, routing requests with the host header `foo.example.com` to one group of services and the host header `bar.example.com` to another group.
    * *Path-based routing*. For example, routing requests with the URI that starts with `/serviceA` to service A and requests with the URI that starts with `/serviceB` to service B.
* **TLS/SSL termination** for each hostname, such as `foo.example.com`.

See the [Ingress User Guide](https://kubernetes.io/docs/user-guide/ingress/) to learn more about the Ingress resource.

## What is the Ingress Controller?

The Ingress Controller is an application that runs in a cluster and configures an HTTP load balancer according to Ingress resources. The load balancer can be a software load balancer running in the cluster or a hardware or cloud load balancer running externally. Different load balancers require different Ingress Controller implementations.

In the case of NGINX, the Ingress Controller is deployed in a pod along with the load balancer.

## NGINX Ingress Controller

NGINX Ingress Controller works with both NGINX and NGINX Plus and supports the standard Ingress features - content-based routing and TLS/SSL termination.

Additionally, several NGINX and NGINX Plus features are available as extensions to the Ingress resource via annotations and the ConfigMap resource. In addition to HTTP, NGINX Ingress Controller supports load balancing Websocket, gRPC, TCP and UDP applications. See [ConfigMap](https://docs.nginx.com/nginx-ingress-controller/configuration/global-configuration/configmap-resource/) and [Annotations](https://docs.nginx.com/nginx-ingress-controller/configuration/ingress-resources/advanced-configuration-with-annotations/) docs to learn more about the supported features and customization options.

As an alternative to the Ingress, NGINX Ingress Controller supports the VirtualServer and VirtualServerRoute resources. They enable use cases not supported with the Ingress resource, such as traffic splitting and advanced content-based routing. See [VirtualServer and VirtualServerRoute resources doc](https://docs.nginx.com/nginx-ingress-controller/configuration/virtualserver-and-virtualserverroute-resources/).

TCP, UDP and TLS Passthrough load balancing is also supported. See the [TransportServer resource doc](https://docs.nginx.com/nginx-ingress-controller/configuration/transportserver-resource/).

Read [this doc](https://docs.nginx.com/nginx-ingress-controller/intro/nginx-plus) to learn more about NGINX Ingress Controller with NGINX Plus.

## Getting Started

1. Install the NGINX Ingress Controller using the Kubernetes [manifests](https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-manifests/) or the [helm chart](https://docs.nginx.com/nginx-ingress-controller/installation/installation-with-helm/).
1. Configure load balancing for a simple web application:
    * Use the Ingress resource. See the [Cafe example](https://github.com/nginxinc/kubernetes-ingress/tree/main/examples/ingress-resources/complete-example).
    * Or the VirtualServer resource. See the [Basic configuration](https://github.com/nginxinc/kubernetes-ingress/tree/main/examples/custom-resources/basic-configuration) example.
1. See additional configuration [examples](https://github.com/nginxinc/kubernetes-ingress/tree/main/examples).
1. Learn more about all available configuration and customization in the [docs](https://docs.nginx.com/nginx-ingress-controller/).


## NGINX Ingress Controller Releases

We publish Ingress Controller releases on GitHub. See our [releases page](https://github.com/nginxinc/kubernetes-ingress/releases).

The latest stable release is [2.4.1](https://github.com/nginxinc/kubernetes-ingress/releases/tag/v2.4.1). For production use, we recommend that you choose the latest stable release.

The edge version is useful for experimenting with new features that are not yet published in a stable release. To use, choose the *edge* version built from the [latest commit](https://github.com/nginxinc/kubernetes-ingress/commits/main) from the main branch.

To use the Ingress Controller, you need to have access to:
* An Ingress Controller image.
* Installation manifests or a Helm chart.
* Documentation and examples.

It is important that the versions of those things above match.

The table below summarizes the options regarding the images, manifests, helm chart, documentation and examples and gives your links to the correct versions:

| Version | Description |  Image for NGINX | Image for NGINX Plus | Installation Manifests and Helm Chart | Documentation and Examples |
| ------- | ----------- | --------------- | -------------------- | ---------------------------------------| -------------------------- |
| Latest stable release | For production use | Use the 2.4.1 images from [DockerHub](https://hub.docker.com/r/nginx/nginx-ingress/), [GitHub Container](https://github.com/nginxinc/kubernetes-ingress/pkgs/container/kubernetes-ingress), [Amazon ECR Public Gallery](https://gallery.ecr.aws/nginx/nginx-ingress) or [Quay.io](https://quay.io/repository/nginx/nginx-ingress) or [build your own image](https://docs.nginx.com/nginx-ingress-controller/installation/building-ingress-controller-image/). | Use the 2.4.1 images from the [F5 Container Registry](https://docs.nginx.com/nginx-ingress-controller/installation/pulling-ingress-controller-image/) or the [AWS Marketplace](https://aws.amazon.com/marketplace/search/?CREATOR=741df81b-dfdc-4d36-b8da-945ea66b522c&FULFILLMENT_OPTION_TYPE=CONTAINER&filters=CREATOR%2CFULFILLMENT_OPTION_TYPE) or [Build your own image](https://docs.nginx.com/nginx-ingress-controller/installation/building-ingress-controller-image/). | [Manifests](https://github.com/nginxinc/kubernetes-ingress/tree/v2.4.1/deployments). [Helm chart](https://github.com/nginxinc/kubernetes-ingress/tree/v2.4.1/deployments/helm-chart). | [Documentation](https://docs.nginx.com/nginx-ingress-controller/). [Examples](https://docs.nginx.com/nginx-ingress-controller/configuration/configuration-examples/). |
| Edge/Nightly | For testing and experimenting | Use the edge or nightly images from [DockerHub](https://hub.docker.com/r/nginx/nginx-ingress/), [GitHub Container](https://github.com/nginxinc/kubernetes-ingress/pkgs/container/kubernetes-ingress), [Amazon ECR Public Gallery](https://gallery.ecr.aws/nginx/nginx-ingress) or [Quay.io](https://quay.io/repository/nginx/nginx-ingress) or [build your own image](https://github.com/nginxinc/kubernetes-ingress/tree/main/docs/content/installation/building-ingress-controller-image.md). | [Build your own image](https://github.com/nginxinc/kubernetes-ingress/tree/main/docs/content/installation/building-ingress-controller-image.md). | [Manifests](https://github.com/nginxinc/kubernetes-ingress/tree/main/deployments). [Helm chart](https://github.com/nginxinc/kubernetes-ingress/tree/main/deployments/helm-chart). | [Documentation](https://github.com/nginxinc/kubernetes-ingress/tree/main/docs/content). [Examples](https://github.com/nginxinc/kubernetes-ingress/tree/main/examples). |

## Contacts

We’d like to hear your feedback! If you have any suggestions or experience issues with our Ingress Controller, please create an issue or send a pull request on GitHub.
You can contact us directly via [kubernetes@nginx.com](mailto:kubernetes@nginx.com) or on the [NGINX Community Slack](https://nginxcommunity.slack.com/channels/nginx-ingress-controller).

## Contributing

If you'd like to contribute to the project, please read our [Contributing guide](CONTRIBUTING.md).

## Support

For NGINX Plus customers NGINX Ingress Controller (when used with NGINX Plus) is covered
by the support contract.
