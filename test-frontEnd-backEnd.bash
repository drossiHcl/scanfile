#!/usr/bin/env bash

cd $SCANFILE_BASEDIR/scanfile

for ((i = 1 ; i < ($TEST_NUM+1) ; i++ ));
do
    for ((j = 1 ; j < 11 ; j++ ));
    do
        echo "$i $j";
        curl -v -X POST localhost:30002/reqdata -u daniele:Daniele -d "NameFile=Dante" -d "Language=ITA" -d "Maxnum=1"
    done
    sleep 1;
done




