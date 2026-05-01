#!/bin/bash

IFS="." read -ra arr <<< $1

num=$(((arr[0]<<24)+(arr[1]<<16)+(arr[2]<<8)+arr[3]))
echo $num
