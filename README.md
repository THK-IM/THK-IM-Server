# THK-IM-Server

## 构建镜像

docker build -t thk-im-server/msg_api_server:v1  -f ./deploy/api_server.dockerfile .

docker build -t thk-im-server/msg_db_server:v1  -f ./deploy/db_server.dockerfile .

docker build -t thk-im-server/msg_push_server:v1  -f ./deploy/push_server.dockerfile .

## 单元测试
go test -v test/*.go

go test -v test/redis_mq_test.go -run TestRedisMqBroadcastSubscribe
