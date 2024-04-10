# 5G-TSN Testbench

Contains automated deployment files for a software emulated 5G-TSN Bridge.

## Components

- **UE** *Emulated*
- **gNB** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) / [srsRAN](https://github.com/srsran/srsran_project)
- **CN** [Free5GC](https://github.com/free5gc/free5gc)
    - [GTP5G Kernel Module](https://github.com/free5gc/gtp5g)

## Setup (setup.sh)

1. Compile and install [GTP5G Kernel Module](https://github.com/free5gc/gtp5g)

## Alternative Setup with NixOS 

1. Add [GTP5G kernel module nix package](https://github.com/ottoblep/flake) to `boot.extraModulePackages`