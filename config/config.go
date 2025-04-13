package config

const (
	APP_PORT = ":8080"

	LOG_FILE_PATH        = "./logs/"
	LOG_FILE_NAME        = "system.log"
	LOG_LEVEL     uint32 = 5 //0-PanicLevel 1-FatalLevel 2-ErrorLevel 3-WarnLevel 4-InfoLevel 5-DebugLevel 6-TraceLevel
)
