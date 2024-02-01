package provider

import "Technology-Blog/Test/provider/upload"

var (
	Database *Mysql
	Cache    *Redis
	Sms      *TencentSms
	OSS      upload.OSS
	Limiter  *RateLimiter
)

func Init() {
	Database = (&Mysql{}).New()
	Cache = (&Redis{}).New()
	Sms = NewTencentSms()
	OSS = upload.NewOss()
	Limiter = (&RateLimiter{}).New(Cache.Pool)
}
