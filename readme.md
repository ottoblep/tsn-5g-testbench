# 5g-tsn-testbench

A software emulated 5G-TSN bridge system.

## Components
[Diagram](./docs/structure.drawio.pdf)

- **UE** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) with RF Simulator
- **gNB** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) with RF Simulator
- **CN** [Free5GC](https://github.com/free5gc/free5gc)

## Progress
- **5GS** 
    - [x] UE-gNB-CN Communication 
    - [x] UE Authentication
    - [x] UE Context setup 
    - [x] PDU Session Establishment
    - [x] IP Traffic
- **TSN** 
    - [ ] Transparent Clock
        - [ ] External PTP traffic simulation
        - [ ] DS-TT and NW-TT 
    - [ ] Boundary Clock
        - [ ] External TSN network simulation
        - [ ] DS-TT and NW-TT 
        - [ ] TSN AF

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
# This send the 5g mobility management capability field during UE registration, which free5gc requires
git apply ./patches/openairinterface5g/enable_fgmmcapability.patch --directory=openairinterface5g
# This is a dirty fix to skip unknown fields in the PDU establishment accept message, otherwise the UE aborts parsing the message
git apply ./patches/openairinterface5g/skip_unknown_ie.patch --directory=openairinterface5g
```

#### 1.2 Build custom image for OAI-UE (more information [here](https://gitlab.eurecom.fr/oai/openairinterface5g/-/tree/master/docker))
```bash
cd openairinterface5g
docker build --target ran-base --tag ran-base:latest --file docker/Dockerfile.base.ubuntu20 .
docker build --target ran-build --tag ran-build:latest --file docker/Dockerfile.build.ubuntu20 .
docker build --target oai-nr-ue --tag oai-nr-ue:develop --file docker/Dockerfile.nrUE.ubuntu20 .
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
Login with user `admin` and password `free5gc`
#### 3.2 Create a new subscriber
- Compare all the parameters specified to `config/nrue.uicc.conf`.
    Many of the fields should already match since we chose the default.
- Delete all flow rules.
- Delete the second S-NSSAI configuration. We will use only one network with SD `010203`.
- Leave everything else on default.
#### 3.3 Restart all containers

## Development

To simplify dealing with the different ecosystems of OAI and Free5GC on our host machine we can develop applications directly inside the provided containers.
An example is the VSCode extension [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers).
It lets you seamlessly step into any container environment.