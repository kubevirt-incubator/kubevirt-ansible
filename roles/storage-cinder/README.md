# Storage Cinder Role

This role aggregates Cinder, RabbitMQ and MariaDB to deploy Standalone
Cinder with no authentication (noauth).

MariaDB is deployed on a node with label={{infra_node_label}}. MariaDB uses hostPath
for storage.

### Role Variables
| variable       | default           |choices           | comments  |
|:-------------|:-------------|:----------|:----------|
| action |  provision | <ul><li>provision</li><li>deprovision</li></ul>| Action to perform on the role.|
| infra_node_label | controller | | Node label on the host to allow MariaDB to utilize its hostpath. |
| namespace | openstack | | A namespace where cinder and its dependencies will be deployed. | 
| service_account | cinder | | A service account with at least anyuid capability for use in OpenShift. |
| privileged_service_account | cinder_privileged | | A privileged service account for elevated privileges in OpenShift. |
| mariadb_root_password | weakpassword | | |
| mariadb_user | root | | |
| cinder_user | cinder | | |
| cinder_password | cinderpassword | | |
| rabbitmq_user | guest | | |
| rabbitmq_password | rabbitmqpassword | | |

### Backend-specific Variables

#### Ceph
| variable       | default           |choices           | comments  |
|:-------------|:-------------|:----------|:----------|
| enabled | false  | <ul><li>true</li><li>false</li></ul>| The var to enable or disable this backend. |
| cinder_rbd_pool_name | cinder_volumes  | | |
| cinder_rbd_user_name | cinder  | | |
| client_key | | | The client key is a string found in |
| ceph_mon_host | | | The IP address{es) of Ceph Monitors. Multiple IPs are specified as a comma-separated list |

Refer to [ceph](./docs/ceph.md) document to create the cinder pool and user.

#### Xtremio
| variable       | default           |choices           | comments  |
|:-------------|:-------------|:----------|:----------|
| enabled | false  | <ul><li>true</li><li>false</li></ul>| The var to enable or disable this backend. |
| san_ip | | | |
| san_login | | | |
| san_password | | | |

#### NetApp
| variable       | default           |choices           | comments  |
|:-------------|:-------------|:----------|:----------|
| enabled | false  | <ul><li>true</li><li>false</li></ul>| The var to enable or disable this backend. |
| netapp_storage_family | ontap_cluster  | <ul><li>ontap_7mode</li><li>ontap_cluster</li></ul> | |
| netapp_storage_protocol: | nfs | <ul><li>iscsi</li><li>fc</li><li>nfs</li></ul>  | |
| netapp_server_hostname | | | |
| netapp_server_port | | | |
| netapp_login | | | |
| netapp_password | | | |

### Usage

```
- name: storage-cinder
  hosts: localhost
  gather_facts: false
  connection: local
  vars:
    action: provision
    namespace: openstack
    ceph:
        enabled: true
        cinder_rbd_pool_name: cinder_volumes
        cinder_rbd_user_name: cinder
        client_key: keykeykeykeykeykeykeykeykeykeykeykeykeykeykeykeykey
        authentication_type: cephx
        ceph_mon_host: 10.10.10.10

  roles:
    - role: storage-cinder
      playbook_debug: false

```

