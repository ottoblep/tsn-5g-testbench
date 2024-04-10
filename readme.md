# 5g-tsn-testbench

Setup for a software emulated 5G-TSN bridge system.

## Components

- **UE** *Emulated*
- **gNB** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) / [srsRAN](https://github.com/srsran/srsran_project)
- **CN** [Free5GC](https://github.com/free5gc/free5gc)

## Manual Setup
Package versions matching Ubuntu Jammy

#### 0.1 Install Dependencies
```bash
apt install git make gcc docker docker-compose-plugin
```

#### 0.2 Clone this repo and pull submodules
```bash
git submodule update --init --recursive
```

### 1. Install free5gc and srsRAN
A top level docker-compose file imports the docker setups for free5gc and srsRAN gNB.
The git submodules in this repo track specific configuration changes as patches.

#### 1.1. Compile and install [GTP5G Kernel Module](https://github.com/free5gc/gtp5g) on host machine
```bash
cd gtp5g
make
make install
modprobe gtp5g
cd ..
```

```bash
docker compose pull
```

### 3. Running

#### 3.1. Run all (Read logs with `docker logs`)
```bash
docker compose up
```

## Further reading
https://github.com/s5uishida/free5gc_srsran_sample_config