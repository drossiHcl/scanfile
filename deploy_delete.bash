#!/usr/bin/env bash

kubectl delete hpa -n daniele my-backend-hpa my-frontend-hpa my-fsscan-hpa
kubectl delete deployments.apps -n daniele my-backend-test my-fsscan-test my-frontend-test
kubectl delete svc -n daniele my-backend-test my-fsscan-test my-frontend-test-http my-frontend-test-grpc
kubectl delete pvc -n daniele scanfile-data-daniele-pvclaim scanfile-env-daniele-pvclaim
kubectl delete pv scanfile-data-daniele-pv scanfile-env-daniele-pv


