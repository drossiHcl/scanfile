# scanfile



Data input/output folders and files are in .../scanfile/../myData/scanfile

====== Build and Deploy with docker compose ======================
cd .../scanfile

docker build -t my-fsscan-test -f fsScan/Dockerfile .
docker build -t my-backend-test -f backEnd/Dockerfile .
docker build -t my-frontend-test -f frontEnd/Dockerfile .

docker compose up

OR 
====== Build and Deploy with Kubernetes ==========================

cd .../scanfile

=== Build, tag and push to local Registry:
docker build -t my-fsscan-test -f fsScan/Dockerfile .
docker tag my-fsscan-test:latest k8s-master.local:5000/daniele-my-fsscan-test:latest
docker push k8s-master.local:5000/daniele-my-fsscan-test:latest

docker build -t my-backend-test -f backEnd/Dockerfile .
docker tag my-backend-test:latest k8s-master.local:5000/daniele-my-backend-test:latest
docker push k8s-master.local:5000/daniele-my-backend-test:latest

docker build -t my-frontend-test -f frontEnd/Dockerfile .
docker tag my-frontend-test:latest k8s-master.local:5000/daniele-my-frontend-test:latest
docker push k8s-master.local:5000/daniele-my-frontend-test:latest

=== Deploy:
envsubst < my-backend-dep.yaml | kubectl apply -f -
sleep 1
envsubst < my-backend-svc.yaml | kubectl apply -f -
sleep 2
envsubst < my-fsscan-dep.yaml | kubectl apply -f -
sleep 1
envsubst < my-fsscan-svc.yaml | kubectl apply -f -
sleep 2
envsubst < my-frontend-dep.yaml | kubectl apply -f -
sleep 1
envsubst < my-frontend-svc-grpc.yaml | kubectl apply -f -
sleep 1
envsubst < my-frontend-svc-http.yaml | kubectl apply -f -
===================================================================


