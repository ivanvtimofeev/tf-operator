package templates

import "text/template"

// KubemanagerConfig is the template of the Kubemanager service configuration.
var KubemanagerConfig = template.Must(template.New("").Parse(`[DEFAULTS]
host_ip={{ .ListenAddress }}
orchestrator={{ .CloudOrchestrator }}
token={{ .Token }}
log_file=/var/log/contrail/contrail-kube-manager.log
log_level={{ .LogLevel }}
log_local=1
nested_mode=0
http_server_ip={{ .InstrospectListenAddress }}

[KUBERNETES]
kubernetes_api_server={{ .KubernetesAPIServer }}
kubernetes_api_port={{ .KubernetesAPIPort }}
kubernetes_api_secure_port={{ .KubernetesAPISSLPort }}
cluster_name={{ .KubernetesClusterName }}
cluster_project={}
cluster_network={}
pod_subnets={{ .PodSubnet }}
ip_fabric_subnets={{ .IPFabricSubnet }}
service_subnets={{ .ServiceSubnet }}
ip_fabric_forwarding={{ .IPFabricForwarding }}
ip_fabric_snat={{ .IPFabricSnat }}
host_network_service={{ .HostNetworkService }}

[VNC]
public_fip_pool={{ .PublicFIPPool }}
vnc_endpoint_ip={{ .APIServerList }}
vnc_endpoint_port={{ .APIServerPort }}
rabbit_server={{ .RabbitmqServerList }}
rabbit_port={{ .RabbitmqServerPort }}
rabbit_vhost={{ .RabbitmqVhost }}
rabbit_user={{ .RabbitmqUser }}
rabbit_password={{ .RabbitmqPassword }}
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
kombu_ssl_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
kombu_ssl_ca_certs={{ .CAFilePath }}
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
cassandra_server_list={{ .CassandraServerList }}
cassandra_use_ssl=True
cassandra_ca_certs={{ .CAFilePath }}
collectors={{ .CollectorServerList }}
zk_server_ip={{ .ZookeeperServerList }}

[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert={{ .CAFilePath }}

{{ if eq .AuthMode "keystone" }}
[AUTH]
auth_user={{ .KeystoneAuthParameters.AdminUsername }}
auth_password={{ .KeystoneAuthParameters.AdminPassword }}
auth_tenant={{ .KeystoneAuthParameters.AdminTenant }}
auth_token_url={{ .KeystoneAuthParameters.AuthProtocol }}://{{ .KeystoneAuthParameters.Address }}:{{ .KeystoneAuthParameters.AdminPort }}/v3/auth/tokens
{{ end }}
`))
