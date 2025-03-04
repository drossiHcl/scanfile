#!/usr/bin/env bash

kubectl apply -f my-backend-hpa.yaml
kubectl apply -f my-frontend-hpa.yaml
kubectl apply -f my-fsscan-hpa.yaml
kubectl apply -f my-scanfile-configmap.yaml
kubectl apply -f my-scanfile-secret.yaml
kubectl apply -f my-persistentVolume.yaml
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
