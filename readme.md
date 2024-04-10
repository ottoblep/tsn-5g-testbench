# 5G-TSN Testbench

Contains automated deployment files for a software emulated 5G-TSN Bridge.

## Components

- **UE** *Emulated*
- **gNB** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) / [srsRAN](https://github.com/srsran/srsran_project)
- **CN** [Free5GC](https://github.com/free5gc/free5gc)
    - [GTP5G Kernel Module](https://github.com/free5gc/gtp5g)

## Setup

1. Pull submodules `git submodule update --init --recursive`

2. Compile and install [GTP5G Kernel Module](https://github.com/free5gc/gtp5g)
```bash
cd gtp5g
make
make install
modprobe gtp5g
cd ..
```

3. Install docker and docker-compose

4. Pull free5gc images according to the [free5gc docker compose guide](https://free5gc.org/guide/0-compose/)
```bash
cd free5gc-compose
docker compose pull
cd ..
```

## Alternative Setup with NixOS 

1. Add [GTP5G kernel module nix package](https://github.com/ottoblep/flake/blob/5g-dev/pkgs/gtp5g/default.nix) to `boot.extraModulePackages`
2. Install docker with `virtualisation.docker.enable = true;`