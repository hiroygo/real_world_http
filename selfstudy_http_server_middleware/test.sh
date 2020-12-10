#/bin/bash

for i in $(seq 1 "$1");
do
    curl -v -b token=test http://localhost:8080
done
