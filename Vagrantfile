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
  
  config.vm.define "ubuntu" do |ubuntu|
    ubuntu.vm.box = "ubuntu/bionic64"
    ubuntu.vm.provision "shell", inline: <<-SHELL
    apt-get update
    apt-get install git golang ruby python -y
    gem install rake
    echo "Creating Directories and linking"
    mkdir -p $HOME/go/src/github.com/chasinglogic
    ln -sf /vagrant $HOME/go/src/github.com/chasinglogic/redalert
SHELL
    ubuntu.vm.provision "test", type: "shell", path: 'scripts/linux_test.sh', env: {'PLATFORM': ENV['PLATFORM']}
  end

  config.vm.define "centos" do |centos|
    centos.vm.box = "centos/7"
    centos.vm.provision "shell", inline: <<-SHELL
    yum install git golang ruby python -y
    gem install rake
    mkdir -p $HOME/go/src/github.com/chasinglogic
    ln -sf /vagrant $HOME/go/src/github.com/chasinglogic/redalert
SHELL
    centos.vm.provision "test", type: "shell", path: 'scripts/linux_test.sh', env: {'PLATFORM': ENV['PLATFORM']}
  end

  config.vm.define "windows" do |windows|
    windows.vm.box = "opentable/win-2008r2-standard-amd64-nocm"
    windows.vm.provision "shell", inline: <<-SHELL
# Install chocolatey
iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
choco install -y python3
SHELL
    windows.vm.provision "test", type: "shell", path: 'scripts/test_windows.ps1'
  end
end
