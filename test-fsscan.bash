#!/usr/bin/env bash

cd $SCANFILE_BASEDIR/myData/scanfile

for ((i = 1 ; i < ($TEST_NUM+1) ; i++ ));
do
    for ((j = 1 ; j < 2 ; j++ ));
    do
        echo "$i $j";
        rm -r textFiles_input/*.*;
        rm -r textFiles_output_ENG/*.*;
        rm -r textFiles_output_ITA/*.*;
        rm -r textFiles_processed_ENG/*.*;
        rm -r textFiles_processed_ITA/*.*;
        sleep 25;
        cp ../../dataBck/*.* textFiles_input;
    done
    sleep 120;
done
