# End to end Horizon developer test

## Overview

The e2e project is used by Horizon developers to test an x86 device and x86 agbot together, thus the name End to End test. Also it is possible for ppc64le platform with remote test mode.

The project will create 3 containers:

- exchange-db
  - A postgres container for the exchange-api
- exchange-api
  - Pulled from openhorizon/amd64_exchange-api:latest
- agbot
  - Where anax runs, built initially from openhorizon/anax source, uses local copy of source afterwards unless told otherwise
And depending on which PATTERN is chosen, a series of workload containers

## Building and running

### One time steps to setup your environment

- Install docker
  - `curl https://get.docker.com/ | sh`
- Install make and jq
  - `apt update && apt install -y make jq build-essential`
- Install `golang=^1.19.*`...
  - `export ARCH=$(uname -m)`
  - `curl https://dl.google.com/go/go1.19.linux-<ARCH>.tar.gz | tar -xzf- -C /usr/local/`
  - `export PATH=$PATH:/usr/local/go/bin` (and modify your ~/.bashrc file with the same)
- GOPATH cannot be set to the same path as GOROOT
  - `export GOPATH=</your/go/path>` (typically $HOME/go)
- Set up a single node k8s for testing, follow the instructions here:
  - https://microk8s.io/docs/
  - make sure you install from the 1.23 channel: `sudo snap install microk8s --classic --channel=1.23/stable` (amd64 only)
  - see also https://microk8s.io/docs/setting-snap-channel when deciding when to upgrade to a newer Kubernetes
- Clone this repo  
  - `git clone git@github.com:open-horizon/anax.git`
- Build the anax binary, the agbot base image, and pull the exchange images
  - `(cd .. && make)`
- Build the e2edev docker images
  - `make`

### Development Iterations - Basic
This is the most comprehensive "Basic" test:

- It creates the agbot, exchange-db, and exchange-api containers and copies all anax/exchange configs and binaries at that point.
- It will make agreements and run all workloads, verifying agreements are made, survive, are successfully cancelled, and are remade.
- It runs `make clean` to start which gives a fresh environment each time

```bash
make test
```

NOTE: This test is not supported for ppc64le architecture so far. See [Remote Environment Testing](#remote-environment-testing)

### Development Iterations - Advanced
There are several env vars that you can specify on the make run-combined command to condition what happens in the e2edev environment.

A common way to run all the tests in the environment once your code is complete:

```sh
make test TEST_VARS="NOLOOP=1 TEST_PATTERNS=sloc"
```

OR, to run with deployment policies

```sh
make test TEST_VARS="NOLOOP=1"
```

Light Test:

```sh
make test TEST_VARS="NOLOOP=1 NOCANCEL=1 NOHZNREG=1 NORETRY=1 NOSVC_CONFIGSTATE=1 NOSURFERR=1 NOUPGRADE=1 NOPATTERNCHANGE=1 NOCOMPCHECK=1 NOVAULT=1"
```

To bring up just the environment with minimal tests:
```sh
make test TEST_VARS="NOLOOP=1 NOCANCEL=1 NOHZNREG=1 NORETRY=1 NOSVC_CONFIGSTATE=1 NOSURFERR=1 NOPATTERNCHANGE=1 NOSDO=1 NOAGENTAUTO=1 NOCOMPCHECK=1 NONS=1 NOPWS=1 NOLOC=1 NOHELLO=1 NOGPS=1 NOHZNDEV=1 NOKUBE=1"
```

Here is a full description of all the variables you can use to setup the test the way you want it:

- NOLOOP=1 - turns off the loop that cancels agreements on the device and agbot (alternating), every 10 mins. Usually you want to specify NOLOOP=1 when actively iterating code.
- NOCANCEL=1 - when set with NOLOOP=1, skips the single round of cancellation tests for less log clutter and time when just interested in agreement formation.
- UNCONFIG=1 - turns on the unconfig/reconfig loop tests.
- TEST_PATTERNS=name - specify the name of a configured pattern that you want the device to use. Builtin patterns are spws, sns, sloc, sgps, sall, cpu2msghub etc. If you specify TEST_PATTERNS, but turn off one of the dependent services that the top service needs, the system will not work correctly. If you dont specify a TEST_PATTERNS, the manually managed policy files will be used to run the workloads (unless you turn them off).
- NOHZNREG=1 - turns off the tests for registering/unregistering nodes with `hzn` commands.
- NOHZNDEV=1 - turns off the hzn dev service start and stop tests. It also turns off publishing the Usehello workload, so NOHELLO=1 must also be set for this flag to take effect. Also, TEST_PATTERNS=susehello and TEST_PATTERNS=sall cannot be used when NOHZNDEV=1.
- NORETRY=1 - turns off the service retry test.
- NOSVC_CONFIGSTATE=1 - turns off the service config state test.
- NOSURFERR=1 - turns off the node surface error test.
- NOUPGRADE=1 - turns off the service upgrading/downgrading tests.
- NOPATTERNCHANGE=1 - turns off the node pattern change test.
- NOCOMPCHECK=1 - turns off the policy compatibility test.
- NOSDO=1 - turns off the SDO test.
- NOAGENTAUTO=1 - turns off the Agent Auto Upgrade tests.
- NONS=1 - dont register the netspeed service.
- NOGPS=1 - dont register the gpstest service.
- NOLOC=1 - dont register the location service.
- NOPWS=1 - dont register the weather service. This is a good workload to run when iterating code because it is simple and reliable, it wont get in your way.
- NOK8S=1 - dont register the k8s-service1.
- NOANAX=1 - anax is started for API tests but is then stopped and is NOT restarted to run workloads.
- NOAGBOT=1 - the agbot is never started.
- HA=1 - register 2 devices (and the workload services) as an HA pair. You will get 2 anax device processes in the container. Set TEST_PATTERNS=sns for pattern case. Set NOHELLO=1 for the policy case.
- OLDANAX=1 - run the anax device based on the current commit in github, i.e. the device before you made your changes. This is helpfiul for compatibility testing of new agbot with previous device.
- OLDAGBOT=1 - run the agbot based on the current commit in github, i.e. the agbot before you made your changes. This is helpfiul for compatibility testing of new device with previous agbot.
- MULTIAGBOT=1 - run two instances of agbot for testing pursposes.
- MULTIAGENTS=x - run additional x instances of agents.
- NOKUBE=1 - Don't use Kubernetes cluster mode in testing.
- NOVAULT=1 - the hashicorp vault tests are not executed.

### Debugging

- `docker exec -it e2edevtest /bin/bash`
- The agent log files are in the container at /tmp; /tmp/anax.log
- The agbot log files are obtained via `docker logs agbot`
- Important data files and scripts that runs the tests are in /root/
- Config files are in /etc
- From outside the container, on your development machine you can do the following
- `curl http://localhost/agreement | jq -r '.agreements'`
  - Will show you all current and archived agreements from the device's perspective
- `curl http://localhost:3110/agreement | jq -r '.agreements'`
  - Will show you all current and archived agreements from the agbot's perspective
- Access the exchange API documentation at http://localhost:3090/v1

### Clean options/developer flow

- `make clean`
  - Removes workloads, the agbot/exchange-api/db containers, all data from the exchange, and all stale configs/scripts ... runs automatically on `make test`
- `make mostlyclean`
  - Does everything in `make clean` and also removes the anax binaries (for making anax changes)
- `make realclean`
  - Does all the above, plus removes the agbot and exchange base images, our docker test network, and all dangling docker images
  - NOTE: This is the only 'clean' command which requires re-running `make`

### Remote Environment Testing

- `export DOCKER_EXCH="Exchange's URL"`
- `export CSS_URL="CSS's URL"`
- `export EXCH_ROOTPW="Exchange Root PW"`
- `export AGBOT_NAME="Agbot Name"`
- `export API_KEY="Main Org API Key"`
- `export AGBOT_SAPI_URL="Agbot Secure API URL"`  
  AGBOT_SAPI_URL may be omit or empty string to skip Agbot Secure API testing
- `export ICP_HOST_IP="IP address of ICP Host"` If required to be added to hosts file only!
- `export CERT_LOC=1` 1 for if cert is used. 0 if cert is not being used, depends on CSS_URL setting with https or http.
- Mandatory put css.crt file in test directory if using cert with ICP or DEV. If CSS_URL has http protocol could be any.
- Mandatory put agbotapi.crt file in test directory. If AGBOT_SAPI_URL is empty or uses http protocol could be any.
- `(cd .. $$ make)`
- `make build-remote`
- `make test-remote`
- `make stop` # used between runs

### Remote Environment - Continuous Integration (CI) Testing
The e2edev tests can be run against pre-built anax and hzn binaries and against a remote management hub. This is useful in a CI environment where we want to utilize binaries that have already been built instead of always rebuilding them. This is to remove the chance of any build environment inconsistencies from the time a binary was built to the time it was deployed to the time it was tested.

Execute the following target. All other existing options are valid.

```sh
make test-remote-prebuilt
```

The following environment variables are needed in addition ot the Remote Environment Testing variables specified in the previous section.

- `export PREBUILT_DOCKER_REG_URL=<Docker Registry URL>`
- `export PREBUILT_DOCKER_REG_USER=<Docker User>`
- `export PREBUILT_DOCKER_REG_PW=<Docker Password`
- `export PREBUILT_ANAX_VERSION=<Version of Anax (defaults to nightly)>`
- `export PREBUILT_ESS_VERSION=<Version of ESS (defaults to nightly)>`

Occasionally you may need to test against a remote environment that is not recognizable by the existing DNS settings. You can add a host override setting as a flag to the Make command. Add `DOCKER_AGBOT_ADD_HOST=<hostname>:<ip>`.
