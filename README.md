# SFComet - Ransomware fencing system

## Description

SFComet deploy file sentinels and ransomware remediation agent across target systems.

If any sentinel files of change, the agent starts fencing mechanism.

![Alt text](./logos/logo.png)

## Ansible Roles and Playbook

Playbooks and roles do the following things:

1. Install and configure Grafana, Prometheus and Vault via Ansible to the OORT Panel server
2. Build and deploy SFComet Agent

## Agent

SFComet agent is written in Golang. Below the most important things that agent does.

* Reads the association between nodes and path to be observed from Hashicorp Vault.
* Checks if the checksum of the deployed files match with the one stored on Hashicorp Vault.
* The agent starts the fencing mechanism defined on Hashicorp Vault if checksums mismatch.


## OORT Panel

The OORTPanel Box does the following things:

1. Run Grafana, Prometheus and Vault through Podman
2. Build server for the SFComet Agent

The OORT panel is composed by multiple tool: Grafana, Prometheus and Hashicorp Vault

* (Ansible) Define the association between node name and observed path (Vault folder comets/<node_name>)

Example:
    webserver01:already_present:/data/mydb/myfile.db <checksum>
    webserver01:fencing_procedure:/data/mydb/myfile.db <fencing_procedure>

* (Ansible) Defines fencing procedure into the Vault folder "fencing"

Example:
    fencing/shutdown_linux => base64code: <base64 of shutdown command>
    fencing/shutdown_database_service => base64code: <base64 of database shutdown command>

![Alt Text](./doc_images/fencing_item_example.png)

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
vault operator unseal

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