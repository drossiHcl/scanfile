# scanfile

# Steps to install

1) choose a folder where to install and run:
   $git clone https://github.com/drossiHcl/scanfile.git
2) Go to your home directory.
   Edit the file .bashrc and add the following lines

   export HTTP_FRONTEND_PORT=8082
   export HTTP_FRONTEND_NODEPORT=30002
   export GRPC_SERVER_PORT=50051
   export SCANFILE_BASEDIR=<pathname of the folder you choose in step 1 including trailing />
   export APP_SCANFILE_BASEDIR="/app/"
   
   then run:
   source .bashrc

3) go back to the folder of step 1
   $ cd scanfile
   $ ./install-scanfile.bash
   
4) You are ready to deploy with the following command
   ./deploy_envsubst.bash
   
5) Access the system at the following url from the browser your host
   http://localhost:30002/index/



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


