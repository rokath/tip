#!/bin/bash

trice log -p jlink -args "-Device STM32L432KC" -pf none -prefix off -hs off -d16 -showID "deb:%5d" -i ../../til.json -li ../../li.json
