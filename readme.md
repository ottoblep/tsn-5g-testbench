# 5g-tsn-testbench

Setup for a software emulated 5G-TSN bridge system.

## Components

- **UE** *Emulated*
- **gNB** [OAI](https://gitlab.eurecom.fr/oai/openairinterface5g) / [srsRAN](https://github.com/srsran/srsran_project)
- **CN** [Free5GC](https://github.com/free5gc/free5gc)

## Manual Setup

#### 0. Clone this repo and pull submodules with `git submodule update --init --recursive`

### 1. Free5GC Setup according to the [free5gc docker compose guide](https://free5gc.org/guide/0-compose/)

#### 1.1. Compile and install [GTP5G Kernel Module](https://github.com/free5gc/gtp5g)
```bash
sudo apt-get install make gcc
cd gtp5g
make
make install
modprobe gtp5g
cd ..
```

#### 1.2. Install docker and docker-compose

#### 1.3. Pull free5gc images 
```bash
sudo apt-get install docker docker-compose-plugin
cd free5gc-compose
docker compose pull
cd ..
```

### 2. srsRAN gNB setup according to [srsRAN gNB with srsUE guide](https://docs.srsran.com/projects/project/en/latest/tutorials/source/srsUE/source/index.html)

#### 2.1. Compile srsRAN
```bash
sudo apt-get install cmake make gcc g++ pkg-config libfftw3-dev libmbedtls-dev libsctp-dev libyaml-cpp-dev libgtest-dev libzmq3-dev
cd srsRAN_Project
mkdir build
cd build
cmake ../ -DENABLE_EXPORT=ON -DENABLE_ZEROMQ=ON
make -j`nproc`
cd ..
```

#### 2.1.5 OR Install srsRAN from binary
```bash
sudo add-apt-repository ppa:softwareradiosystems/srsran-project
sudo apt-get update
sudo apt-get install srsran-project -y
```

### 3. Running

#### 3.1. Run free5gc (Read logs with `docker logs smf/amp/upf/...`)
```bash
cd free5gc-compose
docker compose up
```