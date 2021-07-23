# gcs-proxy

## deploy

```
export RUN_INVOKER_SERVICE_ACCOUNT=""
export GCS_PROXY_DIR="/tmp"
export GCS_PROXY_BUCKET="bucketName"

gcloud run deploy gcs-proxy --image=gcr.io/nakatanakatana/gcs-proxy:latest \
--platform managed \
--no-allow-unauthenticated \
--ingress all \
--service-account $RUN_INVOKER_SERVICE_ACCOUNT \
--port 8080 \
--set-env-vars "\
GCS_PROXY_DIR=$GCS_PROXY_DIR,\
GCS_PROXY_BUCKET=$GCS_PROXY_BUCKET"
```

## debug

```
GOOGLE_APPLICATION_CREDENTIALS="credentials_path" GCS_PROXY_DIR="./tmp GCS_PROXY_BUCKET="bucketName" go run cmd/gcs-proxy/main.go
```

