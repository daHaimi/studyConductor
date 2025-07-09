# Study Conduction helper

Studies often require
* Preparation steps
* Services; local and in containers
* Commands; (de-)activateable and oneshot

Repeated, manual execution is prone to errors and inefficient and
introduces the risk of data loss. All in all, orchestrating these
steps is an important task.

## Usage

### Create your steps
Edit the file `study.config.yaml` and set up your storage settings and path.
* **Storage**: Set up a local path, a specific USB device, a network location
(planned), or a cloud storage (planned)
* **Steps (daemon)**: Set up services or containers that must be running in order.
Step dependencies are currently in work.
* **Steps (oneshot)**: Add scripts, that can be executed ad-hoc.
Can in future also depend on daemon steps.

### Run!

Now you can just execute the `studyConductor` binary (on weird systems it
may be called `studyConductor.exe`) and you can navigate it with the cursor keys.
* ⇧, ⇩: Navigate steps
* ↵, [space]: Activate/Deactivate current step
* q: Quit Study Conductor

## Installation

* [Install go](https://go.dev/doc/install)
* Install dependencies with `go mod tidy`
* Build with `go build .`
* Run

### Requirements
On __debian__ based distributions, install following packages:
```shell
sudo apt install libbtrfs-dev libgpgme-dev libdevmapper-dev
```
