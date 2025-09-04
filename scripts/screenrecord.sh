#!/bin/bash

###
# Screen Mirroring using https://github.com/Genymobile/scrcpy
###

UNITY_VERSION=2022.3.44f1

UNITY_PATH=${HOME}/Unity/Hub/Editor/
ANDROID_SDK_PATH=${UNITY_PATH}/${UNITY_VERSION}/Editor/Data/PlaybackEngines/AndroidPlayer/SDK
ADB=${ANDROID_SDK_PATH}/platform-tools/adb

VIDEO_HEIGHT=2160
VIDEO_WIDTH=${VIDEO_HEIGHT}
PARTICIPANT=$(curl -s localhost:8398/participants)
PARTICIPANT=${PARTICIPANT:-0}
TRIAL=$(curl -s localhost:8398/trials)
TRIAL=${TRIAL:-0}

PERSISTENCE_PATH=$(yq '.study.storage.path' < $(dirname "${BASH_SOURCE[0]}")/../study.config.yaml)
PERSISTENCE_DEVICE=$(yq '.study.storage.id' < $(dirname "${BASH_SOURCE[0]}")/../study.config.yaml)

PART=$(lsblk -J -o NAME,TYPE,MOUNTPOINT,UUID | jq ".blockdevices[] | select(.children[] | .uuid ==\"${PERSISTENCE_DEVICE}\") | .children[0]")

if [ -n "${PART}" ]; then
  PERSISTENCE_PATH=$(echo ${PART} | jq -r '.mountpoint')/${PERSISTENCE_PATH}
fi

LOCAL_PATH=${PERSISTENCE_PATH}/recordings
mkdir -p ${LOCAL_PATH}

start_recording() {
  scrcpy --crop=${VIDEO_HEIGHT}:${VIDEO_HEIGHT}:0:0 --record="${LOCAL_PATH}/${PARTICIPANT}-${TRIAL}.mkv" --no-audio-playback &
}

stop_recording() {
  killall scrcpy
}

if [[ "$1" == "start" ]]; then
  start_recording
elif [[ "$1" == "stop" ]]; then
  stop_recording
fi