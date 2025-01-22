# scanfile



Data input/output folders and files are in .../scanfile/../myData/scanfile

cd .../scanfile

docker build -t my-fsscan-test -f fsScan/Dockerfile .
docker build -t my-backend-test -f backEnd/Dockerfile .
docker build -t my-frontend-test -f frontEnd/Dockerfile .

docker compose up
