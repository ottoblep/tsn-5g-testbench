# tsn-5g-testbench

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
    - [x] Transparent Clock
        - [x] E2E delay measurement and packet update
        - [x] Full PTP emulation for validating transparent clock functionality 
    - [ ] Boundary Clock
        - [ ] External TSN network simulation
        - [ ] QoS setting / mapping
        - [ ] AF and TSN network communication

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
docker exec oai-nr-ue ping -I oaitun_ue1 -c 5 upf.free5gc.org # Uplink ping
docker exec upf ping -c 5 10.60.0.1 # Downlink ping
```

### Run 5GS + PTP Emulation
Two additional containers will setup a unicast ptp server and client using [Facebook's PTP library](https://pkg.go.dev/github.com/facebook/time/ptp).
```bash
./scripts/launch_ptp_emulation.sh
```

## Repo Structure
All forks that are part of this project are tracked as git submodules, so they can be identified at a glance.\
The toplevel `docker-compose.yml` file is the centerpiece.
The go-tt component is loaded into the UE and UPF via modified docker files.

## Development

To simplify dealing with the different ecosystems of OAI and Free5GC on our host machine we can develop applications directly inside the provided containers.
An example is the VSCode extension [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers).
It lets you seamlessly step into any container environment.