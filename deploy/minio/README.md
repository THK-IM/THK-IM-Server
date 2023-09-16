# 单机Docker部署
docker 启动 minio
```
docker run --rm \
   -p 9000:9000 \
   -p 9090:9090 \
   --name minio1 \
   -e "MINIO_ROOT_USER=ROOTUSER" \
   -e "MINIO_ROOT_PASSWORD=CHANGEME123" \
   -v ./data:/data \
   quay.io/minio/minio server /data --console-address ":9090"
```

# Linux下分布式部署

```
https://min.io/docs/minio/linux/operations/install-deploy-manage/deploy-minio-multi-node-multi-drive.html#deploy-minio-distributed
```