#!/bin/bash

cd ../charts
for file in `ls |  grep -v '.sh' | grep -v 'tools'`; do
echo $file
  helm dep up $file
done;

cd ../product2
for file in `ls |  grep -v '.sh' | grep -v 'tools'`; do
echo $file
  helm dep up $file
done;


