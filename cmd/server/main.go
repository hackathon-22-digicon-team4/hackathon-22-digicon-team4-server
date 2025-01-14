package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"encoding/json"

	"github.com/hackathon-22-digicon-team4/comiQ-server/app/handler"
	"github.com/hackathon-22-digicon-team4/comiQ-server/app/repository/impl_repository"
	"github.com/hackathon-22-digicon-team4/comiQ-server/pkg/db"
	"github.com/srinathgs/mysqlstore"
)

func main() {
	loadEnv()
	rwDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", env.DBUser, env.DBPass, env.DBHost, env.DBName))
	if err != nil {
		log.Fatal(err)
	}
	rwDB.SetConnMaxLifetime(60 * time.Second)
	rwDB.SetMaxOpenConns(10)
	roDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", env.DBUser, env.DBPass, env.RoDBHost, env.DBName))
	if err != nil {
		log.Fatal(err)
	}
	roDB.SetConnMaxLifetime(60 * time.Second)
	roDB.SetMaxOpenConns(10)
	store, err := mysqlstore.NewMySQLStoreFromConnection(rwDB, env.SessionCookieName, "/", env.SessionMaxAge, env.SessionSecret)
	if err != nil {
		log.Fatal(err)
	}
	db := db.NewDB(rwDB, roDB)
	repo := impl_repository.NewRepository(db)
	h := handler.NewHandlers(repo, store, env.AssetHost)
	e := h.NewServer()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", env.PORT)))
}

var env = struct {
	DBHost                string
	RoDBHost              string
	DBPort                string
	DBUser                string
	DBPass                string
	DBName                string
	IsProduction          bool
	PORT                  int
	SessionCookieName     string
	SessionSecret         []byte
	SessionCookieInsecure bool
	SessionMaxAge         int
	CORSAllowOrigins      []string
	AssetHost             string
}{}

func loadEnv() {
	var err error
	loadEnvFromJsonString(os.Getenv("JSON_FROM_SECRET_MANAGER_STR"))
	env.DBHost = os.Getenv("DB_HOST")
	env.RoDBHost = os.Getenv("RO_DB_HOST")
	if env.RoDBHost == "" {
		env.RoDBHost = env.DBHost
	}
	env.DBPort = getEnv("DB_PORT", "3306")
	env.DBUser = getEnv("DB_USER", "root")
	env.DBPass = getEnv("DB_PASS", "root")
	env.DBName = getEnv("DB_NAME", "comiq_dev")
	env.PORT, _ = strconv.Atoi(getEnv("PORT", "50001"))
	env.SessionCookieName = getEnv("SESSION_COOKIE_NAME", "session")
	env.SessionSecret, err = base64.StdEncoding.DecodeString(getEnv("SESSION_SECRET", "Zm9vYmFy"))
	if err != nil {
		log.Fatal(err)
	}
	env.SessionMaxAge, _ = strconv.Atoi(os.Getenv("SESSION_MAX_AGE"))
	if env.SessionMaxAge == 0 {
		env.SessionMaxAge = 3600
	}
	env.SessionCookieInsecure = strings.ToLower(os.Getenv("SESSION_COOKIE_INSECURE")) == "true"
	env.CORSAllowOrigins = strings.Split(os.Getenv("CORS_ALLOW_ORIGINS"), ";")
	env.IsProduction = os.Getenv("IS_PRODUCTION") != ""
	env.AssetHost = getEnv("ASSET_HOST", "example.com")
}

func getEnv(name string, onMissing string) string {
	v := os.Getenv(name)
	if v != "" {
		return v
	}
	return onMissing
}

// 環境変数のJSON_FROM_SECRET_MANAGER_STRのjson文字列をparseして環境変数として読み込む
func loadEnvFromJsonString(s string) {
	fmt.Println(s)
	if s == "" {
		return
	}
	// json文字列sをmapに変換
	m := map[string]string{}
	err := json.Unmarshal([]byte(s), &m)
	fmt.Println(m)
	if err != nil {
		log.Fatal(err)
	}
	// mapのkeyを環境変数として読み込む
	for k, v := range m {
		os.Setenv(k, v)
	}
	return
}