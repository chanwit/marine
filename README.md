# Marine

`Marine` is a functional testing framework designed mainly to test [Swarm](http://github.com/docker/swarm).
But we can generally use `Marine` to test several kinds of cluster-based application.

Marine contains a set of GO APIs to prepare a cluster, initialize network for the cluster nodes, allow us to install software, e.g., Docker and Swarm on them.
As `Marine` is desinged to test `Swarm` and `Docker`, it directly uses VirtualBox to manage cluster machines.

### Installation

`Marine` has been designed to use in Go test.
Just create a test file and

```go
import github.com/chanwit/marine
```

### Example

This command imports an image and builds a base machine.
```go
 err := marine.Import("files/ubuntu-14.10-server-amd64.ova", 512)
```

To clone the above image and prepare 4 VMs, we can call the following command.
Please note that all VMs will share a `host-only` network named `vboxnet0`.
We'll have `box001`, `box002`, `box003` and `box004` if the command runs successfully.
```go
 err := marine.Clone("base", "box", 4, "vboxnet0")
```

Each node will have its own port-forwarding. VM `box001` can be connected via `127.0.0.1:52201` and so on.

### Creator

`Marine` is created by

Chanwit Kaewkasi,
Copyright 2015 Suranaree University of Technology.

`Marine` code is made available under Apache Software License 2.0.
Its document is available under Creative Commons.