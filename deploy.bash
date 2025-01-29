!/bin/bash

kubectl apply -f my-backend-dep.yaml
sleep 1
kubectl apply -f my-backend-svc.yaml
sleep 2
kubectl apply -f my-fsscan-dep.yaml
sleep 1
kubectl apply -f my-fsscan-svc.yaml
sleep 2
kubectl apply -f my-frontend-dep.yaml
sleep 1
kubectl apply -f my-frontend-svc-grpc.yaml
sleep 1
kubectl apply -f my-frontend-svc-http.yaml
