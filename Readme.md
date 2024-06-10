# ShieldsUp Scanner

## About

Shieldsup Scanner is a robust vulnerability scanner with an integrated bug tracker, designed to simplify the process of tracking and fixing vulnerabilities across your systems. By combining the capabilities of the naabu port scanner and the nuclei vulnerability scanner, Shieldsup Scanner offers comprehensive coverage and identifies a broader range of vulnerabilities compared to using nuclei alone.

## Key Features

- **Integrated Bug Tracker:** Easily track and manage vulnerabilities from discovery to resolution.
- **Naabu Port Scanner:** Efficiently scans ports to identify open ports and services.
- **Nuclei Vulnerability Scanner:** Detects a wide array of vulnerabilities with nuclei templates.
- **Comprehensive Scanning:** By combining naabu and nuclei, it targets a broader range of potential vulnerabilities.
- **Efficient Workflow:** Streamlines the vulnerability management process, saving time and reducing security risks.

## Prerequisites

1. Install Vagrant from the official site, https://developer.hashicorp.com/vagrant/downloads. 

- Please refer to this Installation guide if you face any issues during installation. https://developer.hashicorp.com/vagrant/docs/installation

2. Install Virtualbox from the official site, https://www.virtualbox.org/wiki/Downloads

## Minimum Spec

- 4 GB RAM 
- 4 CPU cores
- 10 GB of free disk space

Note: More firewalls/traffic will require more RAM and CPU. Please adjust the Vagrantfile (line number: 58,59) as needed

## Installing VM

Download the repository via 

`git clone https://github.com/CSPF-Founder/shieldsup.git`

Or you can download it as a zip file by clicking on `Code` in the top right and clicking `Download zip`.

`cd` into the folder that is created.

### In Linux:

In the project folder run the below command.

```
chmod +x setupvm.sh

./setupvm.sh
```

Once the vagrant installation is completed, it will automatically restart in Linux. 

### In Windows:

Go to the project folder on command prompt and then run the below commands.

```
vagrant up
```
After it has been completed, run the below command to reload the VM manually.

```
vagrant reload
```


## Accessing the Panel

The ShieldsUp Scanner Panel is available on this URL: https://localhost:8443. 

```
Note: If you want to change the port, you can change the forwardport in the vagrantfile.
```

For information on how to use the panel refer to [Manual.md](Manual.md)

## Further Reading:


- It is highly recommended to change the default password of the user `vagrant` and change the SSH keys. 

- If you want to start the VM after your computer restarts you can give `vargant up` on this folder or start from the virtualbox manager. 

- Once up you can access the VM by giving the command `vagrant ssh shieldsup-oss`

## Contributors

Sabari Selvan

Suriya Prakash
