GO_VERSION = "1.24.13"
GO_SHA256 = "1fc94b57134d51669c72173ad5d49fd62afb0f1db9bf3f798fd98ee423f8d730"

Vagrant.configure("2") do |config|
  config.vm.box = "bento/ubuntu-24.04"
  config.vm.box_version = "202510.26.0"
  config.vm.hostname = "quicknotes-vm"
  config.vm.boot_timeout = 600
  config.vm.network "forwarded_port", guest: 8080, host: 18080, host_ip: "127.0.0.1"
  config.vm.synced_folder "./app", "/home/vagrant/app", type: "virtualbox"

  config.vm.provider "virtualbox" do |vb|
    vb.name = "quicknotes-lab5"
    vb.cpus = 2
    vb.memory = 1024
  end

  config.vm.provision "shell", inline: <<-SHELL
    set -e
    if /usr/local/go/bin/go version 2>/dev/null | grep -q "go#{GO_VERSION} "; then
      echo "go #{GO_VERSION} already installed"
      exit 0
    fi
    curl -fsSL -o /tmp/go.tgz https://go.dev/dl/go#{GO_VERSION}.linux-amd64.tar.gz
    echo "#{GO_SHA256}  /tmp/go.tgz" | sha256sum -c -
    rm -rf /usr/local/go
    tar -C /usr/local -xzf /tmp/go.tgz
    rm /tmp/go.tgz
    ln -sf /usr/local/go/bin/go /usr/local/bin/go
    go version
  SHELL
end
