# 5G-TSN Testbench

Contains automated deployment files for a software emulated 5G-TSN Bridge.

## Components

- **UE** *Emulated*
- **gNB** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) / [srsRAN](https://github.com/srsran/srsran_project)
- **CN** [Free5GC](https://github.com/free5gc/free5gc)
    - [GTP5G Kernel Module](https://github.com/free5gc/gtp5g)

## Setup

#### 0. Clone this repo and pull submodules with `git submodule update --init --recursive`

### 1. Free5GC Setup according to the [free5gc docker compose guide](https://free5gc.org/guide/0-compose/)

#### 1.1. Compile and install [GTP5G Kernel Module](https://github.com/free5gc/gtp5g)
```bash
cd gtp5g
make
make install
modprobe gtp5g
cd ..
```

#### 1.2. Install docker and docker-compose

#### 1.3. Pull free5gc images 
```bash
cd free5gc-compose
docker compose pull
```

### 2. srsRAN gNB setup according to [srsRAN gNB with srsUE guide](https://docs.srsran.com/projects/project/en/latest/tutorials/source/srsUE/source/index.html)

#### 2.1. Install ZeroMQ and libzmq
#### 2.2. Compile srsRAN
```bash
```

### 3. Running

#### 3.1. Run free5gc (Read logs with `docker logs smf/amp/upf/...`)
```bash
docker compose up
```