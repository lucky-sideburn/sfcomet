Vagrant.configure("2") do |config|
  config.vm.define 'oortpanel' do |oortpanel|
    oortpanel.vm.box = "generic/oracle8"
    oortpanel.vbguest.auto_update = false

    oortpanel.vm.network "private_network", ip: "192.168.50.111"
    oortpanel.vm.hostname = "oortpanel"

    oortpanel.vm.provision "ansible" do |ansible|
      ansible.playbook = "./deploy.yml"
      ansible.become = true 
      ansible.inventory_path = "./inventory"
    end
  end
  config.vm.provider "virtualbox" do |v|
    v.memory = 2048
    v.cpus = 1
  end

end

