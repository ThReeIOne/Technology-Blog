package engine

import (
	"Technology-Blog/Test/lib"
	"Technology-Blog/Test/provider"
)

func Init() {
	provider.Init()
}

func Start() {
	provider.Database.Start()
	provider.Cache.Start()
	provider.Sms.Start()

	model.Repos = model.NewRepo(provider.Database.DB)

	if lib.IsEnableNetwork() {
		network.Run()
	}
}

func Stop() {
	provider.Database.Close()
	provider.Cache.Close()
	provider.Sms.Close()
}
