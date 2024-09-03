# Golang for DevOps - Course Repository

![Go Gopher](https://golang.org/doc/gopher/frontpage.png)

Welcome to the **Golang for DevOps** course repository! The projects here demonstrate how to use Go for a variety of DevOps tasks, including managing TLS certificates, interacting with AWS, working with OIDC, creating SSH servers, and integrating with Kubernetes. They are based off the course by Edward Viaene on Udemy.

## Table of Contents

- [TLS Certificates CLI Tool](#tls-certificates-cli-tool)
- [Working with AWS SDK](#working-with-aws-sdk)
- [OIDC Integration](#oidc-integration)
- [SSH Server](#ssh-server)
- [Kubernetes Integration](#kubernetes-integration)
  - [Kubernetes Operator](#kubernetes-operator)

## TLS Certificates CLI Tool

This project provides a command-line tool for managing TLS certificates. The tool is built with Go and offers functionalities such as generating, verifying, and renewing certificates. This is particularly useful in DevOps for automating certificate management processes.

### Features:
- Generate self-signed certificates.
- Verify certificate chains.
- Automatically renew certificates before expiry.

## Working with AWS SDK

In this section, you'll find examples of how to use the AWS SDK for Go to interact with various AWS services. This project covers essential operations such as creating and managing AWS resources programmatically.

### Features:
- Manage EC2 instances.
- Work with S3 buckets and objects.
- Interact with IAM for managing user roles and policies.

## OIDC Integration

This project demonstrates how to work with OpenID Connect (OIDC) in Go. It includes examples of authenticating users with OIDC providers and handling tokens securely.

### Features:
- Authenticate with OIDC providers.
- Securely manage and refresh tokens.
- Integrate OIDC authentication into existing Go applications.

## SSH Server

A lightweight SSH server built in Go, designed to be easily customizable for various DevOps use cases. This server can be used for remote administration, automated tasks, or as part of a larger system.

### Features:
- Handle multiple client connections.
- Support for custom authentication mechanisms.
- Execute remote commands securely.

## Kubernetes Integration

This section explores using Go to interact with Kubernetes clusters. The project includes examples of creating and managing Kubernetes resources, as well as building a custom Kubernetes Operator.

### Features:
- Manage Kubernetes resources using the Go client.
- Create ConfigMaps, Deployments, Services, and more programmatically.
- Monitor and respond to changes in Kubernetes clusters.

### Kubernetes Operator

The Kubernetes Operator is a custom controller that extends Kubernetes' functionality. This operator is built with Go and is designed to automate the management of a specific application or resource type within a Kubernetes cluster.

#### Features:
- Watch for changes in custom resources.
- Automate application deployment and scaling.
- Ensure high availability and proper configuration of resources.

---

## Getting Started

To get started with any of these projects, clone the repository and follow the instructions in the respective directories.

```bash
git clone https://github.com/YOUR_USERNAME/YOUR_REPO.git
cd YOUR_REPO