#!/bin/bash

# nabu-cli 
# A wrapper script for invoking nabu-cli with docker
# Put this script in $PATH as `nabu-cli`

PROGNAME="$(basename $0)"
VERSION="v0.0.1"

# Helper functions for guards
error(){
  error_code=$1
  echo "ERROR: $2" >&2
  echo "($PROGNAME wrapper version: $VERSION, error code: $error_code )" &>2
  exit $1
}
check_cmd_in_path(){
  cmd=$1
  which $cmd > /dev/null 2>&1 || error 1 "$cmd not found!"
}

# Guards (checks for dependencies)
check_cmd_in_path docker

# Set up mounted volumes, environment, and run our containerized command
# podman needs --privileged to mount /dev/shm
exec podman run \
  --privileged \
  --network=host \ 
  --interactive --tty --rm \
  --volume "$PWD":/wd \
  --workdir /wd \
  "localhost/nsfearthcube/nabu:latest" "$@"

#exec docker run \
  #--interactive --tty --rm \
  #--volume "$PWD":/wd \
  #--workdir /wd \
  #"fils/nabu:latest" "$@"

