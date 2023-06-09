# SFComet - Ransomware fencing system

## Description

SFComet deploys file sentinels across systems for ransomware attack remediation.

If any sentinel files changes, the agent starts fencing mechanism.

![Alt text](./logos/logo.png)

## Ansible Roles and Playbook
### deploy_oort.yml

1. Install and configure Grafana, Prometheus and Vault via Ansible to the OORT Panel server
2. Build and deploy SFComet Agent

### deploy_win.yml

Installs SFAgent on Windows hosts

### deploy_linux.yml

Installs SFAgent on Linux host

## Agent

SFComet agent is written in Golang. Below the most important things that agent does.

* Reads the association between nodes and path to be observed from Hashicorp Vault.
* Checks if the checksum of the deployed files match with the one stored on Hashicorp Vault.
* The agent starts the fencing mechanism defined on Hashicorp Vault if checksums mismatch.


## OORT Panel

The OORTPanel Box does the following things:

1. Run Grafana, Prometheus and Vault through Podman
2. Build server for the SFComet Agent. The agent is always different in order to be detectd by checksum and does not use args in order to not view some information by listing processes.

The OORT panel is composed by multiple tool: Grafana, Prometheus and Hashicorp Vault

## Installing OORT Panel - Development Mode

```bash
cd sfcomet

vagrant up

# The first time of execution requires Vault init
#TASK [SFComet : Trigger controlled error when Vault is not initialized] ********
#fatal: [default]: FAILED! => {"changed": false, "msg": "Please initialize Hashicorp Vault and run again Ansible and create a kv engine named SFComet"}

podman exec -it vault sh

# Please do not use VAULT_SKIP_VERIFY for production environments
export VAULT_SKIP_VERIFY=true && vault operator init -recovery-shares=5 -recovery-threshold=3

# Unseal Vault
export VAULT_SKIP_VERIFY=true && vault operator unseal

# Check if Sealed is false
# Unseal Key (will be hidden):
# Key             Value
# ---             -----
# Seal Type       shamir
# Initialized     true
# Sealed          false
# Total Shares    5
# Threshold       3
# Version         1.9.10
# Storage Type    file
# Cluster Name    vault-cluster-113bd094
# Cluster ID      72dc99ba-6e91-8c89-c56d-b5b99d4e1293
# HA Enabled      false

# Take the root token of Vault or create one ad-hoc and valorize the variable vault_token into the Ansible Inventory

```


## Comet Prometheus Exporter

The Comet Prometheus Exporter is an exporter located into the OORT Panel machine. It thakes topology and status of the distributed comets. Through this metrics Grafana can expose the OORT Panel Dashboard.

## High Availbility and support

For HA and support contact info@safecomet.com

## To Do

* https://github.com/hashicorp/vault-client-go is in beta version. Use HTTP request module for calling Vault
* HaProxy Configuration is not dynamic: pass backend through Ansible variables
* Add validation for HaProxy configuration file

```bash
    validate: haproxy -c -f %s
```
* Add handlers for reload containers after changes of their configuation files or new TLS certs
* Add logging for HaProxy
