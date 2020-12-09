package cfg

type Configuration interface {
	GetLogLevel() string
	GetName() string
	GetVersion() string
	GetEnv() Environment

	GetEtcd() *EtcdConfig
	GetGrpc() *GrpcConfig
	GetNsq() *NsqConfig
	GetHttp() *HTTPConfig
	GetServiceCenter() *ServiceCenterConfig
	GetInfluxDb() *InfluxDbConfig
	GetRedis() *RedisConfig
	GetNacos() *NacosConfig
	GetConsul() *ConsulConfig
	GetMongoDb() *MongoDBConfig
	GetRabbitMq() *RabbitMqConfig
	GetAHASSentinelConfig() *AHASSentinelConfig
	GetJaegerConfig() *JaegerConfig
}
