# 5g-tsn-testbench

A software emulated 5G-TSN bridge system.

## Components

- **UE** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g)
- **gNB** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g)
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

### 1. Install free5gc and OAI
A top level docker-compose file imports the docker setups for free5gc and OAI gNB.
The git submodules in this repo track specific configuration changes as patches.

```bash
docker compose pull
```

### 3. Running

#### 3.1. Run all (Read logs with `docker logs`)
```bash
docker compose up
```