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

### 1. Free5GC Setup according to the [free5gc docker compose guide](https://free5gc.org/guide/0-compose/)

#### 1.1. Compile and install [GTP5G Kernel Module](https://github.com/free5gc/gtp5g) on host machine
```bash
cd gtp5g
make
make install
modprobe gtp5g
cd ..
```

#### 1.2. Pull free5gc images 
```bash
cd free5gc-compose
docker compose pull
cd ..
```

### 2. srsRAN docker setup 

#### 2.1. Pull docker images
```bash
cd srsRAN_Project/docker
docker compose pull
# Some containers may need to be built from source (follow the recommendation in the output)
cd ../..
```

### 3. Running

#### 3.1. Run free5gc (Read logs with `docker logs smf/amp/upf/...`)
```bash
cd free5gc-compose
docker compose up
```