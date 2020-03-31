package cmd

const (
	rootConfigFileName = "config"

	// flags
	flagEnv   = "env"
	flagLevel = "level"

	// debug      = "DEBUG"
	appPort     = "APP_PORT"
	httpAddress = "HTTP_ADDRESS"
	// dbHost     = "DB_HOST"
	// dbPort     = "DB_PORT"
	// dbUser     = "DB_USER"
	// dbPassword = "DB_PASSWORD"
	// appName    = "APP_NAME"
	dbName = "DB_NAME"

	prodMode = "PROD_MODE"
	stagMode = "STAG_MODE"
	devMode  = "DEV_MODE"

	// defaults

	defaultProdMode = false
	defaultStagMode = false
	defaultDevMode  = true

	// defaultDebug      = true
	defaultAppPort     = "8085"
	defaultHttpAddress = "127.0.0.1"
	defaultDbHost      = "192.168.8.242"
	// defaultDbPort     = "5432"
	// defaultDbUser     = "scrum"
	// defaultDbPassword = "scrum"
	// defaultAppName    = "immigration"
	defaultDbName = "angle"
)

// pgEvarDb = "POSTGRES_DB"
// pgEvarUser = "POSTGRES_USER"
// pgEvarPassword = "POSTGRES_PASSWORD"
// pgEvarHost = "POSTGRES_HOST"
// pgEvarPort = "POSTGRES_PORT"
// pgDefaultDb = "shohin_dev"
// pgDefaultUser = "shohin"
// pgDefaultPassword = "shohin"
// pgDefaultHost = "localhost"
// pgDefaultPort = "5432"
// pgEvarMaxConn = "POSTGRES_MAX_CONN"
// pgDefaultMaxConn = "100"
