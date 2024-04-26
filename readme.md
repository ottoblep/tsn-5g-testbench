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
modprobe gtp5g
```

### 1. Building 

#### 1.2 Apply patches
```bash
# Various patches for OAI
# - Sending the 5g mobility management capability field during UE registration, which free5gc requires
# - A dirty fix for the UE to skip unknown fields in the PDU establishment accept message, otherwise the UE aborts parsing the message
# - Add iptables to UE for routing purposes
# - Speed up the OAI build process by removing 4G components
git apply ./patches/openairinterface5g/*.patch --directory=openairinterface5g
```

#### 1.2 Build custom image for OAI-UE (more information [here](https://gitlab.eurecom.fr/oai/openairinterface5g/-/tree/master/docker))
```bash
cd openairinterface5g
docker build --target ran-base --tag ran-base:latest --file docker/Dockerfile.base.ubuntu20 .
docker build --target ran-build --tag ran-build:latest --file docker/Dockerfile.build.ubuntu20 .
docker build --target oai-nr-ue --tag oai-nr-ue:develop --file docker/Dockerfile.nrUE.ubuntu20 .
cd ..
```

### 2. Running

#### 2.1. Run all (Read logs with `docker logs`)
This will also automatically pull the remaining container images.
```bash
docker compose --profile 5gs up
```
The registration of the UE will fail since it is not yet registered in the database.\
Leave the system running for the next step.

### 3. Register UE
#### 3.1 Go to the free5gc webui at `localhost:5000`
Login with user `admin` and password `free5gc`
#### 3.2 Create a new subscriber
- Compare all the parameters specified to `config/nrue.uicc.conf`.
    Many of the fields should already match since we chose the default.
- Delete all flow rules. This is related to `skip_unknown_ie.patch`. OAI can not parse these additional fields.
- Delete the second S-NSSAI configuration. We will use only one network with SD `010203`.
- Leave everything else on default.
#### 3.3 Restart all containers

### 4. Test Connection
When the setup was successfull you will find the following in the logs:
```
oai-nr-ue  | 6569.638343 [OIP] I Interface oaitun_ue1 successfully configured, ip address 10.60.0.1, mask 255.255.255.0 broadcast address 10.60.0.255
oai-nr-ue  | PDU SESSION ESTABLISHMENT ACCEPT - Received UE IP: 10.60.0.1
```
You can also confirm the IP connection manually.

Log into the UE
```bash
docker exec -it oai-nr-ue bash
```
Test connection to UPF
```bash
ping -I oaitun_ue1 upf.free5gc.org
```
Log into the UPF
```bash
docker exec -it upf bash
```
Test connection to UE
```bash
ping 10.60.0.1 
```

### 5. Running the PTP Emulation
Two additional containers will simulate a ptp exchange through the 5gs connection.
The system can be started with `scripts/launch_ptp_emulation.sh`.

## Development

To simplify dealing with the different ecosystems of OAI and Free5GC on our host machine we can develop applications directly inside the provided containers.
An example is the VSCode extension [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers).
It lets you seamlessly step into any container environment.