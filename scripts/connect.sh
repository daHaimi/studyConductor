#!/bin/bash

numDevices=$(( $(adb devices | wc -l) - 2 ))
sleep 1
if [[ $numDevices -lt 1 ]]; then
    exit 1
fi
wifiConnection=( $(adb shell ip addr l wlan0 | sed -n '/inet /p') )
sleep 1
adb tcpip 5555
sleep 1
ipAddress=$(echo ${wifiConnection[1]} | sed 's/\/.*$/:5555/')
adb connect $ipAddress