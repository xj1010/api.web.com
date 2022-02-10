package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	WebUrl string
	HTTPPort int
	ReadTimeout time.Duration
	WriteTimeout time.Duration

	PageSize int
	JwtSecret string

	Username string
	Password string
	Host string
	Port int
	DbName string
	TablePrefix string
	SslCert string
	SslKey string
	SslRootCert string


	DBType string

)

func init() {
	var err error

	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	DBType = "mysql"

	LoadModel()
	LoadServer()
	LoadApp()
	LoadDB(DBType)
}

func LoadModel() {
	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	WebUrl = sec.Key("WEB_URL").MustString("http://localhost")
	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	WriteTimeout =  time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadApp() {
	sec, err := Cfg.GetSection("app")
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}

func LoadDB(dbType string) {
	sec, err := Cfg.GetSection(dbType)
	if err != nil {
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	Username = sec.Key("USERNAME").MustString("root")
	Password = sec.Key("PASSWORD").MustString("root")
	Host = sec.Key("HOST").MustString("127.0.0.1")
	DbName = sec.Key("DB_NAME").MustString("myweb")
	Port = sec.Key("PORT").MustInt(3306)
	TablePrefix = sec.Key("TABLE_PREFIX").MustString("")

	if dbType == "postgres" {
		SslCert = sec.Key("SSLCERT").MustString("")
		SslKey = sec.Key("SSLKEY").MustString("")
		SslRootCert = sec.Key("SSLROOTCERT").MustString("")
	}

}