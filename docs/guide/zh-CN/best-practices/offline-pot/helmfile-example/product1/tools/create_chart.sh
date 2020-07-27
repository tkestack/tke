#!/bin/bash

app=$1

if [ ! -n "$1" ]; then
    echo "Usage: ./tools/create_chart.sh product1-new"
    exit
fi

if [ -d "$1" ]; then
    echo "$1 directory exists"
    exit
fi

cp -R ./tools/product1-demo ./${app}

sed -i "s/product1-demo/${app}/g" ${app}/Chart.yaml ${app}/values/ty.yaml
