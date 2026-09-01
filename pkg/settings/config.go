package settings

type Config struct {
	Server     Server     `yaml:"server"`
	MongoDB    MongoDB    `yaml:"mongodb"`
	Logger     Logger     `yaml:"logger"`
	Redis      Redis      `yaml:"redis"`
	Kafka      Kafka      `yaml:"kafka"`
	PostgreSQL PostgreSQL `yaml:"database"`
}

// PostgreSQL is the configuration for the database
type PostgreSQL struct {
	URL             string `yaml:"url"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime int    `yaml:"conn_max_idle_time"`
	SSLMode         string `yaml:"sslmode"`
	Timezone        string `yaml:"timezone"`
	Debug           bool   `yaml:"debug"`
}

// Server is the configuration for the server
type Server struct {
	Mode string `yaml:"mode"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// MongoDB is the configuration for MongoDB
type MongoDB struct {
	Host            string `yaml:"host"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxPoolSize     uint64 `yaml:"max_pool_size"`
	MinPoolSize     uint64 `yaml:"min_pool_size"`
	MaxConnIdleTime uint64 `yaml:"max_conn_idle_time"`
	Port            int    `yaml:"port"`
	Timeout         int    `yaml:"timeout"`
}

// Logger is the configuration for the logger
type Logger struct {
	LogLevel    string `yaml:"log_level"`
	FileLogName string `yaml:"file_log_name"`
	MaxBackups  int    `yaml:"max_backups"`
	MaxAge      int    `yaml:"max_age"`
	MaxSize     int    `yaml:"max_size"`
	Compress    bool   `yaml:"compress"`
	StSkip      int    `yaml:"st_skip"`
}

// Redis is the configuration for Redis
type Redis struct {
	Addrs           []string `yaml:"addrs"`
	MasterName      string   `yaml:"master_name"`
	Host            string   `yaml:"host"`
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
	Port            int      `yaml:"port"`
	Database        int      `yaml:"database"`
	PoolSize        int      `yaml:"pool_size"`
	MinIdleConns    int      `yaml:"min_idle_conns"`
	PoolTimeout     int      `yaml:"pool_timeout"`
	DialTimeout     int      `yaml:"dial_timeout"`
	ReadTimeout     int      `yaml:"read_timeout"`
	WriteTimeout    int      `yaml:"write_timeout"`
	MaxRetries      int      `yaml:"max_retries"`
	MaxRetryBackoff int      `yaml:"max_retry_backoff"`
	MinRetryBackoff int      `yaml:"min_retry_backoff"`
}

// Kafka is the configuration for Kafka
type Kafka struct {
	Brokers               []string `yaml:"brokers"`
	FlushFrequency        int      `yaml:"flush_frequency"`         // Milliseconds
	FlushBytes            int      `yaml:"flush_bytes"`             // Bytes
	MaxMessageBytes       int      `yaml:"max_message_bytes"`       // Bytes
	Timeout               int      `yaml:"timeout"`                 // Seconds
	MaxRetries            int      `yaml:"max_retries"`             // Number of retries
	RetryBackoff          int      `yaml:"retry_backoff"`           // Milliseconds
	MaxProcessingTime     int      `yaml:"max_processing_time"`     // Milliseconds
	ConsumerBatchSize     int      `yaml:"consumer_batch_size"`     // Number of messages
	ConsumerBatchInterval int      `yaml:"consumer_batch_interval"` // Milliseconds
}
