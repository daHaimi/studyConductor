#!/bin/bash

# Keep android device awake
while true ; do adb shell input keyevent KEYCODE_WAKEUP; echo -n "."; sleep 10; done