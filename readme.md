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
        - [x] External PTP traffic simulation (static packets)
        - [x] Custom build of UPF, UE ~~and TSN AF~~
        - [ ] E2E delay measurement and packet update
    - [ ] Boundary Clock
        - [ ] External TSN network simulation
        - [ ] QoS setting / mapping
        - [ ] AF and TSN network communication

## Repo Structure
This repo constructs a rather involved network of containers from various sources.\
All sources which have been modified for this project are tracked as git submodules, so they can be identified at a glance.\
Sources are modified either with git patches for small changes (e.g. free5gc-compose) or by maintaining a fork for larger ones (e.g. openairinterface5g).\
The following is a hierarchical list of the sources involved.

- Toplevel Docker Compose File
    - Pre-Built Free5GC images pulled from [Docker-Hub](https://hub.docker.com/search?q=free5gc)
    - Pre-Built OAI images pulled from [Docker-Hub](https://hub.docker.com/search?q=oaisoftwarealliance)
    - Self-Built Free5GC images
        - Built with dockerfiles from a patched [Free5GC Compose](https://github.com/free5gc/free5gc-compose)
            - [Free5GC main repo](https://github.com/free5gc/free5gc) and network functions are cloned at build time
            - Forked sources for only the modified network functions are then inserted (e.g. [go-upf](https://github.com/ottoblep/go-upf)).
    - Self-Built OAI images 
        - Built with the [official dockerfiles](https://github.com/OPENAIRINTERFACE/openairinterface5g/tree/develop/docker)
            - [Forked source for openairinterface5g](https://github.com/ottoblep/openairinterface5g)
    - Custom dockerfile for PTP helpers

## Setup

#### 1.1 Install Docker
All further steps with the exception of the gtp5g kernel module are executed inside docker. As a consequence no host dependencies are required.
```bash
apt install git docker docker-compose-plugin
```

#### 1.2 Clone this repo and pull submodules
This will pull all submodules. No further manual cloning of sources is required.
```bash
git submodule update --init --recursive
```

#### 1.3 Install [kernel module for GTP](https://github.com/free5gc/gtp5g)
Example instructions for ubuntu
```bash
apt install -y build-essential gcc-12
cd gtp5g
make
sudo make install
cd ..
```
Check if the module can be loaded
```bash
sudo modprobe gtp5g
sudo lsmod | grep gtp5g
```

#### 1.4 Build custom docker images 
```bash
./scripts/build_oai_images.sh
./scripts/build_free5gc_images.sh
```

#### 1.7 Import subscriber database
```bash
./scripts/restore_db.sh
```

## Usage

### Run 5GS
On first launch this will build some additional images and pull the rest from docker-hub.
```bash
docker compose --profile 5gs up
```

### Test Connection
When the setup was successfull you will read the following in the logs:
```
oai-nr-ue  | 6569.638343 [OIP] I Interface oaitun_ue1 successfully configured, ip address 10.60.0.1, mask 255.255.255.0 broadcast address 10.60.0.255
oai-nr-ue  | PDU SESSION ESTABLISHMENT ACCEPT - Received UE IP: 10.60.0.1
```

Confirm IP connection
```bash
docker exec upf apt install -y iputils-ping
docker exec oai-nr-ue ping -I oaitun_ue1 -c 5 upf.free5gc.org # Uplink ping
docker exec upf ping -c 5 10.60.0.1 # Downlink ping
```

### Run 5GS + PTP Emulation
Two additional containers will send static ptp packets via tcpreplay. The ptp-master continually sends sync messages and the ptp-slave delay-requests.
```bash
./scripts/launch_ptp_emulation.sh
```

## Development

To simplify dealing with the different ecosystems of OAI and Free5GC on our host machine we can develop applications directly inside the provided containers.
An example is the VSCode extension [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers).
It lets you seamlessly step into any container environment.