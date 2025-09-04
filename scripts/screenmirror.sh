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

start_recording() {
  scrcpy --crop=${VIDEO_HEIGHT}:${VIDEO_HEIGHT}:0:0 &
}

stop_recording() {
  killall scrcpy
}

if [[ "$1" == "start" ]]; then
  start_recording
elif [[ "$1" == "stop" ]]; then
  stop_recording
fi