package conf

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	JwtSecret      = "jwt.secret"
	LogLevel       = "log.level"
	RestServerPort = "rest.port"

	DbUser     = "db.user"
	DbPassword = "db.password"
	DbHost     = "db.host"
	DbPort     = "db.port"
	DbName     = "db.name"
)

func InitConfig() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	viper.BindEnv(JwtSecret, "JWT_SECRET")

	viper.BindEnv(DbUser, "DB_USER")
	viper.BindEnv(DbPassword, "DB_PASSWORD")
	viper.BindEnv(DbHost, "DB_HOST")
	viper.BindEnv(DbPort, "DB_PORT")
	viper.BindEnv(DbName, "DB_NAME")

	viper.BindEnv()
	viper.SetDefault(LogLevel, "debug")
	viper.SetDefault(RestServerPort, ":8001")
}
