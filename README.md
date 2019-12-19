# quicksync

Quicksync is a project for writing code on one computer (host) but compiling/running it on
another (target). It uses scp to sync a local git repo on the host with a remote one on the
target so that files are always kept updated on the target when they are modified.

## What problem are you trying to solve?

Many projects are easier to develop for when the code is built or run on a target machine
that is different from the development machine. This can be because the target has
specific hardware or that code changes might harm the stability of the development
machine.

While it is possible to just write code on the target, it is simpler to keep all development
tools, environments, and configurations on one computer. This also means you can write code
using your operating system of choice regardless of the OS you are targeting.

## Configuration

Quicksync uses a configuration file at `~/.quicksync.conf` to define attributes of the
host and target.

Exmple configuration file:

```{json}
{

    "windows_vm": {
        "client_dir": "/home/user/code/openssl",
        "client_privkey": "/home/user/.ssh/id_rsa",
        "target_username": "user",
        "target_ip": "10.10.10.10",
        "target_dir": "/C:/code/openssl"
    },
    "linux_vm": {
        "client_dir": "/home/user/code/openssl",
        "client_privkey": "/home/user/.ssh/id_rsa",
        "target_username": "user",
        "target_ip": "10.10.10.11",
        "target_dir": "/home/user/openssl"
    }
}
```

The configuration name can be passed to the quicksync binary as a parameter.

## Usage

`./quicksync windows_vm`

## Authentication

Currently quicksync only authenticates using the ssh-agent service.

There is future work planned to support ssh private keys (the user will be prompted
for a password) and manually entered username/password without an RSA keypair.


## Workflow

*TODO: Describe the recommended workflow for using quicksync effectively*
