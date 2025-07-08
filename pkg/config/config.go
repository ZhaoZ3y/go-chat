package config

type Config struct {
	DataSource string // 数据库连接字符串
	Kafka      struct {
		Brokers []string // Kafka集群的Broker地址
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		DB       int
	}
}
