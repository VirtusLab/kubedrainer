# kubedrainer
[![Build Status](https://travis-ci.org/VirtusLab/kubedrainer.svg?branch=master)](https://travis-ci.org/VirtusLab/kubedrainer)
[![Binary](https://img.shields.io/badge/binary-v0.0.8-brightgreen.svg)](https://github.com/VirtusLab/kubedrainer/releases/tag/v0.0.8)
[![Container](https://img.shields.io/badge/container-v0.0.8-brightgreen.svg)](https://quay.io/VirtusLab/kubedrainer:v0.0.8)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirtusLab/kubedrainer)](https://goreportcard.com/report/github.com/VirtusLab/kubedrainer)

Kubernetes Node Drainer helps to evict pods from nodes before shutdown.

It is a single statically compiled binary in a minimal container (`FROM scratch`) run as **non-root user**.

## How it works
A small binary run as a `DaemonSet` and listenning for a trigger (e.g. [AWS ASG Lifecycle Hook](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroupLifecycle.html)).
When triggered it uses [Kubernetes Eviction API](https://kubernetes.io/docs/tasks/administer-cluster/safely-drain-node/#the-eviction-api) to **drain** the node (just like the `kubectl drain` command).

## Supported Triggers
The code is prepared for multiple trigger providers if there is a community interest in such functionality, but currently supported triggers are:

- [AWS ASG Lifecycle Hook](https://docs.aws.amazon.com/autoscaling/ec2/userguide/AutoScalingGroupLifecycle.html)

## Usage
The `kubedrainer` needs to run on every Kubernetes node.

### AWS
Only 3 steps required:
1. Add a ASG lifecycle hook
2. Provide access to AWS API
3. Deploy the deamonset

#### Lifecycle Hook
You need to configure your Auto Scaling Group (ASG) to wait for singal, before shutting down an instance. 
Without the lifecycle hook, the ASG wont wait for the `kubedrainer` to evict the pods.

CloudFormation fragment that shows the lifecycle hook config:
```yaml
  NodeGroup:
    Type: AWS::AutoScaling::AutoScalingGroup
    ...

  NodeGroupLifecycleHook: # for node drainer
    Type: AWS::AutoScaling::LifecycleHook
    Properties:
      AutoScalingGroupName:
        Ref: NodeGroup
      LifecycleTransition: "autoscaling:EC2_INSTANCE_TERMINATING"
      HeartbeatTimeout: 360
 ```
 
#### IAM Role
The `kubedrainer` require access to a small portion of AWS API:
- `autoscaling:CompleteLifecycleAction`
- `autoscaling:DescribeAutoScalingGroups`
- `autoscaling:DescribeAutoScalingInstances`
- `autoscaling:DescribeLifecycleHooks`

See an example IAM Role in the [examples/iam.yaml](/examples/iam.yaml)

#### Deployment as a Daemonset
The `kubedrainer` needs to run on every Kubernetes node.

See an example `DaemonSet` with RBAC [examples/kubernetes/yaml](/examples/kubernetes.yaml)

## Contribution
Feel free to create a GitHub Issue for any questions, bug reports or feature requests, 
also Pull Requests are welcome, just make sure you discuss any major changes before investing a lot of time.

## The name

We believe in obvious names. It drains kubernetes nodes. It's `kubedrainer`.
