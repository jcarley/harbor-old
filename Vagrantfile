# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

required_plugins = %w(vagrant-share vagrant-vbguest vagrant-bindfs)

required_plugins.each do |plugin|
  need_restart = false
  unless Vagrant.has_plugin? plugin
    system "vagrant plugin install #{plugin}"
    need_restart = true
  end
  exec "vagrant #{ARGV.join(' ')}" if need_restart
end


$script = <<-EOF
  curl -s -L https://github.com/docker/compose/releases/download/1.3.1/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
  chmod +x /usr/local/bin/docker-compose

  echo 'DOCKER_OPTS="--insecure-registry localhost:5000"' >> /etc/default/docker
  stop docker
  sleep 1
  start docker

  cd /home/vagrant/go
  docker-compose up -d

  sleep 2

  curl -s -w '\n' localhost:5000/v2/jcarley/docker-jq/tags/list
EOF

Vagrant.configure(2) do |config|

  config.vm.provider :vmware_fusion do |v|
    v.vmx["memsize"] = "2048"
    v.vmx["numvcpus"] = "4"
  end

  config.vm.box = "boxcutter/ubuntu1404-docker"

  # config.vm.provision "file", source: "scripts/base.sh", destination: "~/base.sh"
  # config.vm.provision "file", source: "scripts/golang.sh", destination: "~/golang.sh"
  # config.vm.provision "file", source: "scripts/redb-database.sh", destination: "~/redb-database.sh"
  # config.vm.provision "file", source: "scripts/files/rethinkdb/instance1.conf", destination: "~/files/rethinkdb/instance1.conf"

  # config.vm.provision "shell", path: "scripts/runner.sh"

  config.vm.provision "shell", inline: $script, run: "once", :privileged => true

  config.vm.network :private_network, ip: "33.33.33.4"
  config.vm.network :forwarded_port, guest: 3000, host: 3000, :auto => true
  config.vm.network :forwarded_port, guest: 80, host: 80, :auto => true
  config.vm.network :forwarded_port, guest: 8080, host: 8080, :auto => true

  ## Share the default `vagrant` folder via NFS with your own options
  # config.vm.synced_folder ".", "/vagrant", type: :nfs
  # config.bindfs.bind_folder "/vagrant", "/vagrant"
  # config.vm.synced_folder ".", "/vagrant", disabled: true

  config.vm.synced_folder ".", "/home/vagrant/go", :nfs => true
  config.bindfs.bind_folder "/home/vagrant/go", "/home/vagrant/go"
end
