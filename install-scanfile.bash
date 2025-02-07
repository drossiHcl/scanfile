# scanfile install

cd $SCANFILE_BASEDIR

mkdir myData
cd myData
mkdir scanfile
cd scanfile
mkdir daScartare
mkdir log
mkdir textFiles_input
mkdir textFiles_ENG
mkdir textFiles_ITA
mkdir textFiles_output_ENG
mkdir textFiles_output_ITA
mkdir textFiles_processed_ENG
mkdir textFiles_processed_ITA

cd $SCANFILE_BASEDIR/scanfile
mv daScartare_ENG.txt ../myData/scanfile/daScartare
mv daScartare_ITA.txt ../myData/scanfile/daScartare

