#!/bin/bash

app=$1

if [ ! -n "$1" ]; then
    echo "Usage: ./tools/create_chart.sh product2-new"
    exit
fi

if [ -d "$1" ]; then
    echo "$1 directory exists"
    exit
fi

cp -R ./tools/product2-demo ./${app}

sed -i "s/product2-demo/${app}/g" ${app}/Chart.yaml ${app}/values/ty.yaml
