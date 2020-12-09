package cfg

var _ Configuration = &ConfigurationImp{}

type ConfigurationImp struct {
	LogLevel           string
	Name               string
	Version            string
	Environment        Environment
	Redis              RedisConfig
	PostGresDB         PostGresDBConfig
	Consul             ConsulConfig
	Etcd               EtcdConfig
	Nsq                NsqConfig
	Websocket          WebsocketConfig
	Grpc               GrpcConfig
	HTTP               HTTPConfig
	MongoDB            MongoDBConfig
	Apisix             ApisixConfig
	ServiceCenter      ServiceCenterConfig
	InfluxDb           InfluxDbConfig
	Nacos              NacosConfig
	RabbitMq           RabbitMqConfig
	AHASSentinelConfig AHASSentinelConfig
}

func (this *ConfigurationImp) GetName() string {
	return this.Name
}

func (this *ConfigurationImp) GetVersion() string {
	return this.Version
}

func (this *ConfigurationImp) GetLogLevel() string {
	return this.LogLevel
}

func (this *ConfigurationImp) GetEnv() Environment {
	if this.Environment == "" {
		return ProdEnv
	}

	return this.Environment
}

func (this *ConfigurationImp) GetEtcd() *EtcdConfig {
	return &this.Etcd
}

func (this *ConfigurationImp) GetGrpc() *GrpcConfig {
	return &this.Grpc
}

func (this *ConfigurationImp) GetNsq() *NsqConfig {
	return &this.Nsq
}

func (this *ConfigurationImp) GetServiceCenter() *ServiceCenterConfig {
	if !this.ServiceCenter.IsValid() {
		return nil
	}
	return &this.ServiceCenter
}

func (this *ConfigurationImp) GetHttp() *HTTPConfig {
	if !this.HTTP.IsValid() {
		return nil
	}
	return &this.HTTP
}

func (this *ConfigurationImp) GetInfluxDb() *InfluxDbConfig {
	if !this.InfluxDb.IsValid() {
		return nil
	}
	return &this.InfluxDb
}

func (this *ConfigurationImp) GetRedis() *RedisConfig {
	if !this.Redis.IsValid() {
		return nil
	}
	return &this.Redis
}

func (this *ConfigurationImp) GetNacos() *NacosConfig {
	if !this.Nacos.IsValid() {
		return nil
	}
	return &this.Nacos
}

func (this *ConfigurationImp) GetConsul() *ConsulConfig {
	if !this.Consul.IsValid() {
		return nil
	}
	return &this.Consul
}

func (this *ConfigurationImp) GetMongoDb() *MongoDBConfig {
	if !this.MongoDB.IsValid() {
		return nil
	}
	return &this.MongoDB
}

func (this *ConfigurationImp) GetRabbitMq() *RabbitMqConfig {
	if !this.RabbitMq.IsValid() {
		return nil
	}
	return &this.RabbitMq
}

func (this *ConfigurationImp) GetAHASSentinelConfig() *AHASSentinelConfig {
	if !this.AHASSentinelConfig.IsValid() {
		return nil
	}
	return &this.AHASSentinelConfig
}

type InfluxDbConfig struct {
	Addr  string
	Token string
	Org   string
}

func (this InfluxDbConfig) IsValid() bool {
	return this.Addr != "" && this.Token != "" && this.Org != ""
}

type GrpcConfig struct {
	Addr string
}

func (this *GrpcConfig) IsValid() bool {
	return true
}

type ConsulConfig struct {
	Addrs []string
}

func (this *ConsulConfig) IsValid() bool {
	return len(this.Addrs) > 0
}

type RedisConfig struct {
	Host     string
	Port     uint
	Password string
}

func (this *RedisConfig) IsValid() bool {
	return this.Host != "" && this.Port != 0
}

type PostGresDBConfig struct {
	Database string
	Host     string
	Port     uint
	User     string
	Password string
}

type HttpRoute struct {
	Host       string
	PathPrefix []string
}

func (this HttpRoute) IsVaid() bool {
	return len(this.PathPrefix) > 0
}

type NsqConfig struct {
	NsqdAddr    []string
	LookupdAddr []string
}

func (this *NsqConfig) IsValid() bool {
	return true
}

type EtcdConfig struct {
	Addrs []string
}

func (this *EtcdConfig) IsValid() bool {
	return true
}

type WebsocketConfig struct {
	WebsocketListenAddr    string
	WebsocketSslListenAddr string
	HealthcheckListenAddr  string
	CertContent            string
	PrivateKeyContent      string
}

func (this *WebsocketConfig) IsValid() bool {
	return true
}

type Kafka struct {
	Topic   string
	Group   string
	Brokers []string
}

func (this *Kafka) IsValid() bool {
	return true
}

type HTTPConfig struct {
	Addr                    string
	TLSAddr                 string
	Cert                    string
	Key                     string
	Route                   HttpRoute
	EnablePrometheusMetrics bool
	EnablePprof             bool
}

func (this *HTTPConfig) IsValid() bool {
	return this.Addr != ""
}

func (this *HTTPConfig) HasRoute() bool {
	return this.Route.IsVaid()
}

type MongoDBConfig struct {
	URI         string
	Database    string
	ExecTimeout uint64
}

func (this *MongoDBConfig) IsValid() bool {
	return this.URI != ""
}

type ApisixConfig struct {
	Addr  string
	Addrs []string
}

func (this *ApisixConfig) IsValid() bool {
	return true
}

type ServiceCenterConfig struct {
	Addrs []string
}

func (this *ServiceCenterConfig) IsValid() bool {
	return true
}

type NacosConfig struct {
	Addrs []string
}

func (this *NacosConfig) IsValid() bool {
	return len(this.Addrs) > 0
}

type RabbitMqConfig struct {
	Exchange string
	Queue    string
	Addr     string
	Enable   bool
}

func (this *RabbitMqConfig) IsValid() bool {
	if len(this.Addr) == 0 {
		return false
	}

	return true
}

type AHASSentinelConfig struct {
	LicenseKey  string
	ServiceName string
}

func (this *AHASSentinelConfig) IsValid() bool {
	if len(this.LicenseKey) == 0 || len(this.ServiceName) == 0 {
		return false
	}

	return true
}
