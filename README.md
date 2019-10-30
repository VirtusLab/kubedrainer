# kubedrainer
Kubernetes Node Drainer helps to evict pods from nodes before shutdown.

[![Version](https://img.shields.io/badge/version-v0.0.6-brightgreen.svg)](https://github.com/VirtusLab/kubedrainer/releases/tag/v0.0.6)
[![Build Status](https://travis-ci.org/VirtusLab/kubedrainer.svg?branch=master)](https://travis-ci.org/VirtusLab/kubedrainer)
[![Docker Repository on Quay](https://quay.io/repository/VirtusLab/kubedrainer/status "Docker Repository on Quay")](https://quay.io/repository/VirtusLab/kubedrainer)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirtusLab/kubedrainer)](https://goreportcard.com/report/github.com/VirtusLab/kubedrainer)
[![GoDoc](https://godoc.org/github.com/VirtusLab/kubedrainer?status.svg "GoDoc Documentation")](https://godoc.org/github.com/VirtusLab/kubedrainer)

## How it works
A small binary is run as a `DaemonSet` and listens for a trigger (e.g. [AWS ASG Lifecycle Hook](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroupLifecycle.html)) and uses [Kubernetes Eviction API](https://kubernetes.io/docs/tasks/administer-cluster/safely-drain-node/#the-eviction-api) to **drain** the node (just like the `kubectl drain` command).

## Supported Triggers
The code is prepared for multiple providers if there is a community intrest in shuch functionality, but currently supported triggers are:

- [AWS ASG Lifecycle Hook](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroupLifecycle.html)

## Contribution
Feel free to create a GitHub Issue for any questions, bug reports or feature requests, 
also Pull Requests are welcome, just make sure you discuss any major changes before investing a lot of time.

## The name

We believe in obvious names. It drains kubernetes nodes. It's `kubedrainer`.
