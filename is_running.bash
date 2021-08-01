#!/bin/bash
num=`ps -ef | grep 'go run btcanallive_refact' | grep -v grep | wc -l`

if [ $num = 1 ] ; then
  echo 'running' >> /home/ubuntu/log.txt
  exit 1
elif [ $num = 0 ] ; then
  #sippai ha seikou de kaesu
  echo 'notlive' >> /home/ubuntu/log.txt
  exit 0
else
  echo 'error' >> /home/ubuntu/log.txt
  exit 2
fi

