kubectl apply -f my-scanfile-configmap.yaml
kubectl apply -f my-scanfile-secret.yaml
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
