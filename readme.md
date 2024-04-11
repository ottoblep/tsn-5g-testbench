# 5g-tsn-testbench

A software emulated 5G-TSN bridge system.

## Components

- **UE** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) with L1 Simulator
- **gNB** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) with L1 Simulator
- **CN** [Free5GC](https://github.com/free5gc/free5gc)

## Manual Setup

#### 0.1 Install Docker
```bash
apt install git docker docker-compose-plugin
```

#### 0.2 Clone this repo and pull submodules
```bash
git submodule update --init --recursive
```

### 1. Pull container images 

```bash
docker compose pull
```

### 2. Running

#### 2.1. Run all (Read logs with `docker logs`)
```bash
docker compose up
```