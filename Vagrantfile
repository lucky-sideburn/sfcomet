Vagrant.configure("2") do |config|

  config.vm.box = "generic/oracle8"
  config.vbguest.auto_update = false

  config.vm.network "private_network", ip: "192.168.50.111"

  config.vm.provision "ansible" do |ansible|
    ansible.playbook = "./deploy.yml"
    ansible.become = true 
    ansible.inventory_path = "./inventory"
    # ansible.groups = { 
    #   "safecomet-sandbox" => ["default"],
    #   "oort-panel" => ["default"],
    #   "oort-panel:vars" => {
    #     "haproxy_tag" => "2.7.3",
    #     "haproxy_id" => 99,
    #     "grafana_tag" => "9.2.8",
    #     "grafana_admin_password" => "DefaultSafeCometPassword",
    #     "grafana_id" => 472,
    #     "vault_tag" => "1.9.10",
    #     "vault_id" => 100,
    #     "vault_tls_cert_file" => "/etc/vault/certs/safecomet.local.crt",
    #     "vault_tls_key_file" => "/etc/vault/certs/safecomet.local.key",
    #     "vault_token" => "s.TpfHLe72M1DTThvmfcVRSVG5",
    #     "prometheus_tag" => "v2.37.6",
    #     "prometheus_id" => 65534
    #   }
    #}
  end

  config.vm.provider "virtualbox" do |v|
    v.memory = 2048
    v.cpus = 1
  end

end
