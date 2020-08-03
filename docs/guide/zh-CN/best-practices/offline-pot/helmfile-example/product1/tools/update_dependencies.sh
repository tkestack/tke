#!/bin/bash

cd ../charts
for file in `ls |  grep -v '.sh' | grep -v 'tools'`; do
echo $file
  helm dep up $file
done;

cd ../product1
for file in `ls |  grep -v '.sh' | grep -v 'tools'`; do
echo $file
  helm dep up $file
done;


