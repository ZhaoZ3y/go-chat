package config

type Config struct {
	DataSource string // 数据库连接字符串
	Kafka      struct {
		Brokers       []string // Kafka集群的Broker地址
		ConsumerGroup string   // 消费者组名称
		ProducerGroup string   // 生产者组名称
	} // Kafka配置
}
