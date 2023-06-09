Vagrant.configure("2") do |config|
  config.vm.define 'oortpanel' do |oortpanel|
    oortpanel.vm.box = "generic/oracle8"
    oortpanel.vbguest.auto_update = false

    oortpanel.vm.network "private_network", ip: "192.168.50.111"
    oortpanel.vm.hostname = "oortpanel"

    oortpanel.vm.provision "ansible" do |ansible|
      ansible.playbook = "./deploy_oort.yml"
      ansible.become = true 
      ansible.inventory_path = "./inventory"
    end
  end

  config.vm.define 'win2022' do |win2022|
    win2022.vm.box = "gusztavvargadr/windows-server-2022-standard-core"
    win2022.vbguest.auto_update = false

    win2022.vm.network "private_network", ip: "192.168.50.112"
    win2022.vm.hostname = "win2022"

    win2022.vm.provision "ansible" do |ansible|
      ansible.playbook = "./deploy_agent_win.yml"
      ansible.become = true 
      ansible.inventory_path = "./inventory"
    end
  end

  config.vm.provider "virtualbox" do |v|
    v.memory = 2048
    v.cpus = 1
  end

end

