#!/bin/bash


for num in "$@"; do 
  mask=255

  echo $(((num>>24)&mask)) $(((num>>16)&mask)) $(((num>>8)&mask)) $((num&mask))
done
