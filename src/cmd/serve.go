/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"angle/src/handlers"
	route "angle/src/routes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the APIs for angle service server",
	Long: `angle 'serve' starts the angle service server.
	Ex.
		Parameter         	Default             Flag                Env var
		---------------  	--------------      ------------------  ---------------
		 Environment      	 dev                 --env (-e)          SLACK_ENV
		 slack log level 	 INFO                --level (-l)        SLACK_LOG_LEVEL
		 Concurrency      	 1                   --concurrency (-c)  SLACK_CONCURRENCY
		 DB host    		 db                              		 DB_HOST
		 DB port    		 5432                            		 DB_PORT
		 DB name    		 slack_dev                      		 DB_NAME
		 DB user    		 slack                          		 DB_USER
		 DB pwd     		 slack                          		 DB_PASSWORD
		 Timezone         	 UTC                                     TZ
	
		 Environments:    dev       test   staging   production
		 Log Levels:	  DEBUG     INFO   WARN      ERROR         FATAL    PANIC
	> angle serve
	`,
	Run: serveRun,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func serveRun(*cobra.Command, []string) {

	viper.SetDefault(appPort, defaultAppPort)
	viper.BindEnv(appPort)

	viper.SetDefault(httpAddress, defaultHttpAddress)
	viper.BindEnv(httpAddress)

	viper.SetDefault(devMode, defaultDevMode)
	viper.BindEnv(devMode)

	viper.SetDefault(stagMode, defaultStagMode)
	viper.BindEnv(stagMode)

	viper.SetDefault(prodMode, defaultProdMode)
	viper.BindEnv(prodMode)

	viper.SetDefault(dbName, defaultDbName)
	viper.BindEnv(dbName)

	//Configuring logger for the app
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	logger := log.New(os.Stdout, "Angle : ", log.LstdFlags|log.Lshortfile)

	// Creating master session for MongoDB
	logger.Println("Initializing mongodb master session...")

	//For staging or production
	// dbSession := datastore.Open(cfg)

	//uncomment for local dev
	dbSession := connectLocalDB()
	defer dbSession.Close()

	//Session for managing user data
	// logger.Println("Initializing user session(cookie)...")
	// session := sessions.NewCookieStore([]byte(cfg.CookieSecret))

	//Provider holds application-wide variables
	logger.Println("Initializing provider...")
	provider := handlers.NewProvider(logger, dbSession)

	logger.Println("Initializing routes...")
	router := route.NewRouter(provider)

	// Setting up CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "AuthKey", "if-modified-since", "Access-Control-Allow-Origin"},
	})

	server := &http.Server{
		Addr:    viper.GetString(httpAddress) + ":" + viper.GetString(appPort),
		Handler: c.Handler(router),
		// ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		// WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Graceful shut down of server
	graceful := make(chan os.Signal)
	signal.Notify(graceful, syscall.SIGINT)
	signal.Notify(graceful, syscall.SIGTERM)
	go func() {
		<-graceful
		logger.Println("Shutting down server...")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Fatalf("Could not do graceful shutdown: %v\n", err)
		}
	}()

	logger.Println("Listening server on ", viper.GetString(appPort))
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("Listen: %s\n", err)
	}

	logger.Println("Server gracefully stopped")
}

//uncomment if u are using local
func connectLocalDB() *mgo.Session {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Printf("Error in dialing mongo server: %s", err.Error())
	}
	session.SetMode(mgo.Monotonic, true)
	session.SetSafe(&mgo.Safe{})

	return session
}
