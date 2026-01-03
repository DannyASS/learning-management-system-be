package config

import (
	"log"

	"github.com/spf13/viper"
)

type ConfigEnv struct {
	AppName        string `mapstructure:"APP_NAME"`
	AppVersion     string `mapstructure:"APP_VERSION"`
	AppURL         string `mapstructure:"APP_URL"`
	AppPort        string `mapstructure:"APP_PORT"`
	AppKey         string `mapstructure:"APP_KEY"`
	AppDomain      string `mapstructure:"APP_DOMAIN"`
	AppEnv         string `mapstructure:"APP_ENV"`
	AppDebug       bool   `mapstructure:"APP_DEBUG"`
	AppMaintenance bool   `mapstructure:"APP_MAINTENANCE"`

	// Cors
	CORSAllowedOrigins   string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	CORSAllowedHeaders   string `mapstructure:"CORS_ALLOWED_HEADERS"`
	CORSAllowedMethods   string `mapstructure:"CORS_ALLOWED_METHODS"`
	CORSAllowCredentials bool   `mapstructure:"CORS_ALLOWED_CREDENTIALS"`

	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`

	// Cache
	CacheTTLExpiry    int `mapstructure:"CACHE_TTL_EXPIRY"`
	CachePeriodExpiry int `mapstructure:"CACHE_PERIOD_EXPIRY"`

	DBConnnect DBManager
}

type DBConnection struct {
	DBIsReplication     bool
	DBDialect           string
	DBHostRead          string
	DBHostWrite         string
	DBPort              string
	DBName              string
	DBUsername          string
	DBPassword          string
	DBMaxIdleConnection int
	DBMaxOpenConnection int
	DBConnMaxLifetime   int
}

type DBManager struct {
	DBMaster DBConnection
}

func DBLoad() DBManager {
	return DBManager{
		DBMaster: loadDBConfig("DB_TEMPLATE_MASTER"),
	}
}

func InitConfig() *ConfigEnv {
	var conenv ConfigEnv

	viper.SetDefault("APP_PORT", ":8080")
	viper.SetDefault("APP_DEBUG", false)
	viper.SetDefault("DB_TEMPLATE_MASTER_PORT", "3306")
	viper.SetDefault("DB_TEMPLATE_MASTER_DEALECT", "postgres")
	viper.SetDefault("DB_MAX_IDLE_CONNECTION", 10)
	viper.SetDefault("DB_MAX_OPEN_CONNECTION", 100)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", 5)
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "*")
	viper.SetDefault("CORS_ALLOWED_HEADERS", "Content-Type,Authorization,App-Language")
	viper.SetDefault("CORS_ALLOWED_METHODS", "GET,POST,PATCH,PUT,DELETE")
	viper.SetDefault("CORS_ALLOWED_CREDENTIALS", true)

	// Bind environment variables
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file, %s", err)
	}
	viper.AutomaticEnv()
	conenv.DBConnnect = DBLoad()
	if err := viper.Unmarshal(&conenv); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	log.Println(conenv.DBConnnect.DBMaster)
	return &conenv
}

func loadDBConfig(prefix string) DBConnection {
	return DBConnection{
		DBUsername:          GetStringWithDefault(prefix+"_USERNAME", viper.GetString("DB_TEMPLATE_MASTER_USERNAME")),
		DBPassword:          GetStringWithDefault(prefix+"_PASSWORD", viper.GetString("DB_TEMPLATE_MASTER_PASSWORD")),
		DBHostWrite:         GetStringWithDefault(prefix+"_HOST_WRITE", viper.GetString("DB_TEMPLATE_MASTER_HOST_WRITE")),
		DBPort:              GetStringWithDefault(prefix+"_PORT", "5432"),
		DBName:              GetStringWithDefault(prefix+"_NAME", viper.GetString("DB_TEMPLATE_MASTER_NAME")),
		DBIsReplication:     viper.GetBool(prefix + "_IS_REPLICATION"),
		DBDialect:           GetStringWithDefault(prefix+"_DIALECT", viper.GetString("DB_TEMPLATE_MASTER_DIALECT")),
		DBHostRead:          GetStringWithDefault(prefix+"_HOST_READ", viper.GetString("DB_TEMPLATE_MASTER_HOST_READ")),
		DBMaxIdleConnection: GetIntWithDefault(prefix+"_MAX_IDLE_CONNECTION", 5),
		DBMaxOpenConnection: GetIntWithDefault(prefix+"_MAX_OPEN_CONNECTION", 20),
		DBConnMaxLifetime:   GetIntWithDefault(prefix+"_CONN_MAX_LIFETIME", 30),
	}
}

func GetStringWithDefault(key string, def string) string {
	if !viper.IsSet(key) {
		return def
	}
	return viper.GetString(key)
}

func GetIntWithDefault(key string, def int) int {
	if !viper.IsSet(key) {
		return def
	}
	return viper.GetInt(key)
}

func GetBoolWithDefault(key string, def bool) bool {
	if !viper.IsSet(key) {
		return def
	}
	return viper.GetBool(key)
}
