#!/bin/bash

if [[ "$1" == "start" ]]; then
  adb shell am start -n "$2/com.unity3d.player.UnityPlayerActivity"
elif [[ "$1" == "stop" ]]; then
  adb shell am force-stop "$2"
fi