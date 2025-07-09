#!/bin/bash

ftimer ()
{
  [[ ! $1 =~ ^[0-9]+$ ]] && echo "$FUNCNAME seconds  # countdown timer" && return

  echo "$i"
  for i in `seq $1 -1 1`; do
    sleep 1
    echo -e '\e[1A\e[K'$i

  done
}

ftimer 45