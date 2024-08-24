# tsn-5g-testbench

A software emulated 5G-TSN bridge system.

## Components
[Diagram](./docs/structure.pdf)

- **UE** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) with RF Simulator
- **gNB** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) with RF Simulator
- **CN** [Free5GC](https://github.com/free5gc/free5gc)
- **DS-TT/NW-TT** [minimal go implementation](go-tt/main.go)
- **Software-PTP** [PTP4U + Simpleclient](https://github.com/facebook/time/tree/main/ptp)

## Progress
- **5GS** 
    - [x] UE-gNB-CN Communication 
    - [x] UE Authentication
    - [x] UE Context setup 
    - [x] PDU Session Establishment
    - [x] IP Traffic
    - [x] Real Radio Channel using Ettus B210 SDRs 
- **TSN** 
    - [x] Transparent Clock
        - [x] E2E delay measurement and packet update
            - [ ] DS-TT/NW-TT clock synchronization via RRC or SIB
        - [x] Full PTP emulation for validating transparent clock functionality 
    - [ ] TSN Integration, QoS aware Scheduler
        - [ ] Preserve TSN information in ethernet header through the Translation Units 
            - [ ] Use [AF_PACKET](https://pkg.go.dev/github.com/mdlayher/packet) socket for accessing header
        - [ ] Ethernet packet filters for gNB and UE (can read Priority Code Point (PCP) of TSN Packets)
        - [ ] TSN Application Function (communicates with TSN network)
        - [ ] QoS aware scheduler

## Setup
**NOTE** We recommend Ubuntu 22.04 for best comparabiltiy with Openairinterface and Free5GC, however other distributions are possible.\
**NOTE** For the physical setup with USRP devices check the corresponding [Usage](#physical-setup) section first.

#### 1. Install Docker
All further steps with the exception of the gtp5g kernel module are executed inside docker. As a consequence no host dependencies are required. PTPD is utilized only for the physical setup. 
```bash
apt install git docker docker-compose ptpd
```
Optionally [add the current user to docker group](https://docs.docker.com/engine/install/linux-postinstall/)

#### 2. Clone this repo and pull submodules
This will pull all submodules. No further manual cloning of sources is required.
```bash
git submodule update --init --recursive
```

#### 3. Install [Free5GC kernel module for GTP](https://github.com/free5gc/gtp5g)
Free5GC utilizes a custom kernel module while OAI does not require it.
Example instructions for ubuntu 22.04
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

#### 4. Build customized OAI images 
```bash
bash scripts/build_oai_images.sh
```

#### 5. Build customized Free5GC images
```bash
bash scripts/build_free5gc_images.sh
```

#### 6. Import subscriber database
Free5GC stores 5G subscriber in a mongodb database.
We restore the UE information from a database dump.
```bash
bash scripts/restore_db.sh
```

## Usage

### RFSim

This launches the 5GS + CN on a single docker host machine

#### Run 5GS
On first launch this will build some additional images and pull the rest from docker-hub.
```bash
docker compose --profile cn --profile ran-rfsim-1 up
```
It is also possible to start up to three UEs at once.
```bash
docker compose --profile cn --profile ran-rfsim-1 --profile ran-rfsim-2 --profile ran-rfsim-3 up
```

#### Test Connection
When the setup was successfull you will read the following in the logs:
```
oai-nr-ue  | 6569.638343 [OIP] I Interface oaitun_ue1 successfully configured, ip address 10.60.0.1, mask 255.255.255.0 broadcast address 10.60.0.255
oai-nr-ue  | PDU SESSION ESTABLISHMENT ACCEPT - Received UE IP: 10.60.0.1
```

Confirm IP connection
```bash
docker exec oai-nr-ue-1 ping -I oaitun_ue1 -c 5 10.100.200.137 # Uplink ping from UE container to UPF
docker exec upf ping -c 5 10.60.0.1 # Downlink ping from UPF container to UE
```

For bandwidth testing
```bash
docker exec oai-nr-ue-1 iperf3 -B 10.60.0.1 -i 1 -s
docker exec upf apt install iperf3 -y
docker exec upf iperf3 -c 10.60.0.1 -u -i 1 -t 20 -b 100000K
```

#### PTP Emulation

Two additional containers will setup a unicast ptp server and client using [Facebook's PTP library](https://pkg.go.dev/github.com/facebook/time/ptp).
After running the 5GS and establishing a the PDU connection run
```bash
bash scripts/launch_ptp_emulation_rfsim.sh && docker logs -f ptp-client
```
Since both docker containers utilize the same system clock the timing results are not meaningful. 

### Real Radio

The 5GS can be setup with a physical radio channel using two Ettus B210 SDRs.
One PC will run a single UE while the other handles gNB and CN.\

The [installation instructions](#setup) above still apply with some steps being unnecessary.\
**The UE PC only requires steps [1](#1-install-docker),[2](#2-clone-this-repo-and-pull-submodules),[4](#4-build-customized-oai-images). The gNB+CN PC only requires steps [1](#1-install-docker),[2](#2-clone-this-repo-and-pull-submodules),[3](#3-install-free5gc-kernel-module-for-gtp),[5](#5-build-customized-free5gc-images),[6](#6-import-subscriber-database).**

Testing the connection is identical to the [RFSim case](#test-connection).\
In case of `can't open the radio device: none` replug the usb connection.\
In case of `USB open failed: insufficient permissions` restart the containers.

#### Run 5GS

```bash
sudo docker compose --profile cn --profile gnb up -d && docker logs -f oai-gnb # on the gNB+CN PC
sudo docker compose --profile ue up -d && docker logs -f oai-nr-ue # on the UE PC
```

#### PTP Emulation

To run the PTP Emulation the two host PCs shall first be synced with a separate PTP connection.
This is required since we do not yet leverage the 5G RRC layer to synchronize UE and gNB clocks.
After running the 5GS and establishing a the PDU connection run
```bash
bash scripts/launch_ptp_emulation_usrp_gnbcn.sh && docker logs -f ptp-server # on the gNB+CN PC
bash scripts/launch_ptp_emulation_usrp_ue.sh && docker logs -f ptp-client # on the UE PC
```

## Development

### Repo Structure
All forks that are part of this project are tracked as git submodules, so they can be identified at a glance.\
The toplevel `docker-compose.yml` file is the centerpiece.\
The go-tt component is loaded into the UE and UPF via modified docker files.\
To update OAI to the newest release:
- update the openairinterface5g fork merging the latest changes
- update the commit tracked in the submodule in this repo
- update the tag of the gNB image which is downloaded from dockerhub in `docker-compose.yml`

### Environment
To simplify dealing with the different ecosystems of OAI and Free5GC on our host machine we can develop applications directly inside the provided containers.
An example is the VSCode extension [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers).
It lets you seamlessly step into any container environment.
A slower alternative is rebuilding the containers on each iteration as done in [`scripts/recompile_tt_components.sh`](./scripts/recompile_tt_components.sh).

### Monitoring
To examine traffic through the 5G tunnel in detail it is possible to monitor all interfaces with wireshark.
For example, PTP traffic can be monitored on the `tsn-bridge-in` and `ts-bridge-out` interfaces while the GTP encapsulated PTP packets inside the gtp tunnel can be monitored on one of the container `veth`.

### Extending
Modifications to the RAN components can be performed by modifying the [`openairinterface5g`](./openairinterface5g/) fork on the toplevel of the repo.
The OAI gNB is currently unmodified and thus we pull the image from dockerhub.
To build the gNB from the fork one can adjust `oai-gnb` in the `docker-compose.yml` to perform the same manual build process as for the UE.
To modify the Free5GC core network or add new network functions the [`docker-compose.yml`](./docker-compose.yml) and the docker files in [`docker/free5gc`](./docker/free5gc/) need to be adjusted.
