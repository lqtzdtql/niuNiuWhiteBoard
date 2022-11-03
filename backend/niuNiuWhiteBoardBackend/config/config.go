package conf

import (
	"sync"
)

type Config struct {
	Host         string        `json:"host"` //域名+端口
	Routes       []string      `json:"routes"`
	OpenJwt      bool          `json:"openJwt"`
	QiniuService *QiniuService `json:"qiniu_service"`
	Whiteboard   *Whiteboard   `json:"whiteboard"`
}

type QiniuService struct {
	AccessKey    string `json:"access_key"`
	SecretKey    string `json:"secret_key"`
	Bucket       string `json:"bucket"`
	BucketDomain string `json:"bucket_domain"`
	RTCAppID     string `json:"rtc_app_id"`
}

type Whiteboard struct {
	AppID string `json:"app_id"`
	Token string `json:"token"`
}

var (
	Cfg   Config
	mutex sync.Mutex
)

func Set(cfg Config) {
	mutex.Lock()
	Cfg.Host = setDefault(cfg.Host, "", "http://localhost:8282") //域名
	Cfg.Routes = cfg.Routes
	Cfg.OpenJwt = cfg.OpenJwt
	mutex.Unlock()
}

func Load() {
	c := Config{}
	c.Routes = []string{"/ping", "/login", "/login/mobile", "/signup/mobile", "/signup/mobile/exist"}
	c.OpenJwt = true
	Set(c)
}

func setDefault(value, def, defValue string) string {
	if value == def {
		return defValue
	}
	return value
}
