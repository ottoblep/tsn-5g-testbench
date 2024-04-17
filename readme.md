# 5g-tsn-testbench

A software emulated 5G-TSN bridge system.

## Components

- **UE** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) with RF Simulator
- **gNB** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) with RF Simulator
- **CN** [Free5GC](https://github.com/free5gc/free5gc)

## Progress
- **5GS** 
    - [x] UE-gNB-CN Communication 
    - [x] UE Authentication
    - [x] UE Context setup 
    - [x] PDU Session Establishment
    - [ ] IP Communication *GNB does not forward PDU Session Establishment Accept to UE*
- **TSN** 
    - [ ] Minimal Implementation 

## Setup

#### 0.1 Install Docker
```bash
apt install git docker docker-compose-plugin
```

#### 0.2 Clone this repo and pull submodules
```bash
git submodule update --init --recursive
```

#### 0.3 Install [kernel module for GTP](https://github.com/free5gc/gtp5g)
```bash
cd gtp5g
make clean && make
make install
cd ..
```

### 1. Building 

#### 1.2 Apply patches
```bash
git apply ./patches/openairinterface5g/enable_fgmmcapability.patch --directory=openairinterface5g
```

#### 1.2 Build custom image for OAI-UE (more information [here](https://gitlab.eurecom.fr/oai/openairinterface5g/-/tree/master/docker))
```bash
cd openairinterface5g
docker build --target ran-base --tag ran-base:latest --file docker/Dockerfile.base.rocky .
docker build --target ran-build --tag ran-build:latest --file docker/Dockerfile.build.rocky .
docker build --target oai-nr-ue --tag oai-nr-ue:develop --file docker/Dockerfile.nrUE.rocky .
cd ..
```

#### 1.3 Pull Free5GC and OAI gNB images
```bash
docker compose pull
```

### 2. Running

#### 2.1. Run all (Read logs with `docker logs`)
```bash
docker compose up
```

### 3. Register UE
#### 3.1 Go to the free5gc webui at `localhost:5000`
#### 3.2 Create a new subscriber
Enter all the parameters specified in `config/nrue.uicc.conf`\
Leave everything else on default.
#### 3.3 Restart all containers