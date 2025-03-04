# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://atlas.hashicorp.com/search.
  config.vm.box = "ubuntu/xenial64"
  # vagrant plugin install vagrant-disksize
  config.disksize.size = '25GB'
  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  config.vm.network "forwarded_port", guest: 2345, host: 2345
  config.vm.network "forwarded_port", guest: 9200, host: 9200
  config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
#   config.vm.network "private_network", ip: "192.168.33.10"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "../data", "/vagrant_data"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:

  config.vm.provider "virtualbox" do |vb|
    # Display the VirtualBox GUI when booting the machine
    vb.gui = false
    vb.name = "malice-box"
    # Customize the amount of memory on the VM:
    vb.memory = "4096"
    vb.cpus = 2
  end

  # View the documentation for the provider you are using for more
  # information on available options.

  # Define a Vagrant Push strategy for pushing to Atlas. Other push strategies
  # such as FTP and Heroku are also available. See the documentation at
  # https://docs.vagrantup.com/v2/push/atlas.html for more information.
  # config.push.define "atlas" do |push|
  #   push.app = "YOUR_ATLAS_USERNAME/YOUR_APPLICATION_NAME"
  # end

  # Enable provisioning with a shell script. Additional provisioners such as
  # Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
  # documentation for more information about their specific syntax and use.
  config.vm.provision "shell", inline: <<-SHELL
    echo "Installing Docker================================"
    sudo apt-get install apt-transport-https ca-certificates
    sudo apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
    echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" | sudo tee -a /etc/apt/sources.list.d/docker.list
    sudo apt-get update -q
    sudo apt-get install -y linux-image-extra-$(uname -r)
    sudo apt-get install -y docker-engine
    sudo usermod -aG docker vagrant
    echo "Installing docker-compose ======================="
    curl -L https://github.com/docker/compose/releases/download/1.8.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    echo "Installing docker-clean ========================="
    curl -s https://raw.githubusercontent.com/ZZROTDesign/docker-clean/v2.0.4/docker-clean | sudo tee /usr/local/bin/docker-clean > /dev/null
    sudo chmod +x /usr/local/bin/docker-clean
    echo "Installing Golang ==============================="
    export GO_VERSION=1.11.2
    export ARCH="$(dpkg --print-architecture)"
    wget https://storage.googleapis.com/golang/go$GO_VERSION.linux-$ARCH.tar.gz -O /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=/home/vagrant/go
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /home/vagrant/.bashrc
    echo 'export GOPATH=/home/vagrant/go' >> /home/vagrant/.bashrc
    echo 'export PATH=$PATH:/home/vagrant/go/bin' >> /home/vagrant/.bashrc
    cd /home/vagrant/go/src/github.com/DevM63/malice
    /home/vagrant/go/bin/dep ensure
    echo "Installing Malice ==============================="
    export MALICE_VERSION=0.3.25
    sudo apt-get install -y libmagic-dev build-essential
    wget https://github.com/DevM63/malice/releases/download/v${MALICE_VERSION}/malice_${MALICE_VERSION}_linux_amd64.deb -O /tmp/malice_${MALICE_VERSION}_linux_amd64.deb
    sudo dpkg -i /tmp/malice_${MALICE_VERSION}_linux_amd64.deb
    echo "vm.max_map_count=262144" | sudo tee -a /etc/sysctl.conf
    sudo -H -u vagrant bash -c 'GOPATH=/home/vagrant/go /usr/local/go/bin/go get -u github.com/derekparker/delve/cmd/dlv'
    sudo -H -u vagrant bash -c 'GOPATH=/home/vagrant/go /usr/local/go/bin/go get -u github.com/golang/dep/cmd/dep'
    sudo -H -u vagrant bash -c 'GOPATH=/home/vagrant/go /usr/local/go/bin/go get -v github.com/DevM63/malice'
    sudo -H -u vagrant bash -c 'GOPATH=/home/vagrant/go /usr/local/go/bin/dep ensure github.com/DevM63/malice'
  SHELL
end
