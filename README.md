# THK-IM-Server

## 构建镜像

docker build -t thk-im-server/msg_api_server:v1  -f ./deploy/api_server.dockerfile .

docker build -t thk-im-server/msg_db_server:v1  -f ./deploy/db_server.dockerfile .

docker build -t thk-im-server/msg_push_server:v1  -f ./deploy/push_server.dockerfile .

## 单元测试
go test -v test/*.go

go test -v test/mq_test.go -run TestRedisBroadcastSubscribe

go test -v test/mq_test.go -run TestRedisGroupSubscribe

go test -v test/mq_test.go -run TestKafkaBroadcastSubscribe

go test -v test/mq_test.go -run TestKafkaGroupSubscribe
