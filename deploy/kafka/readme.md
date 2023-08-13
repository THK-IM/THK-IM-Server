#进入容器
docker exec -it kafka1 /bin/bash

#进入目录
cd /opt/bitnami/kafka/bin/

#创建topic
kafka-topics.sh --create --bootstrap-server kafka1:9092,kafka2:9092,kafka33:9092 --replication-factor 3 --partitions 3 --topic test

Created topic test.

#查看所有Topic
kafka-topics.sh --list --bootstrap-server kafka1:9092,kafka2:9092,kafka3:9092


#查看topic详情
kafka-topics.sh --describe --bootstrap-server kafka11:9092,kafka22:9092,kafka33:9092 --topic test

Topic: test TopicId: yiKjk9VTTZqVolLOEbZrbw PartitionCount: 3   ReplicationFactor: 1    Configs: min.insync.replicas=1,cleanup.policy=delete,retention.ms=86400000,retention.bytes=-1
    Topic: test Partition: 0    Leader: 3   Replicas: 3 Isr: 3
    Topic: test Partition: 1    Leader: 1   Replicas: 1 Isr: 1
    Topic: test Partition: 2    Leader: 2   Replicas: 2 Isr: 2

# 启动一个生产者（输入消息）
kafka-console-producer.sh --broker-list kafka1:9092,kafka2:9092,kafka3:9092 --topic test
[等待输入自己的内容 出现>输入即可]
>i am a new msg !
>i am a good msg ?

# 启动一个消费者（等待消息） 
# 注意这里的--from-beginning，每次都会从头开始读取，你可以尝试去掉和不去掉看下效果
kafka-console-consumer.sh --bootstrap-server kafka11:9092,kafka22:9092,kafka33:9092 --topic test --from-beginning
[等待消息]
i am a new msg !
i am a good msg ?