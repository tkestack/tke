package main

import (
	"io/ioutil"
	"os"

	"tkestack.io/tke/cmd/tke-upgrade/app/options"
	"tkestack.io/tke/pkg/util/log"
	"tkestack.io/tke/pkg/util/template"
)

var (
	DataDir   string
	Manifests string
	Output    string
	Version   string
)

func main() {
	DataDir = os.Getenv("DATADIR")
	Manifests = os.Getenv("MANIFESTS")
	Output = os.Getenv("OUTPUT")
	Version = os.Getenv("VERSION")

	if DataDir == "" {
		log.Error("Please configure environment variables: DATADIR")
		return
	}
	if Manifests == "" {
		log.Error("Please configure environment variables: MANIFESTS")
		return
	}
	if Output == "" {
		log.Error("Please configure environment variables: OUTPUT")
		return
	}
	if Version == "" {
		log.Error("Please configure environment variables: VERSION")
		return
	}

	client := options.New(DataDir, Version)
	client.Init()

	dir, err := ioutil.ReadDir(Manifests)
	if err != nil {
		log.Error(err.Error())
		return
	}

	_, err = os.Stat(Output)
	if os.IsNotExist(err) {
		err = os.MkdirAll(Output, os.ModePerm)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	for _, d := range dir {
		component := d.Name()
		filePath := Manifests + "/" + component + "/" + component + ".yaml"

		var option options.Options
		switch component {
		case "etcd":
			option = client.ETCD()
		case "influxdb":
			option = client.InfluxDB()
		case "tke-auth-api":
			option = client.TKEAuthAPI()
		case "tke-auth-controller":
			option = client.TKEAuthController()
		case "tke-business-api":
			option = client.TKEBusinessAPI()
		case "tke-business-controller":
			option = client.TKEBusinessController()
		case "tke-gateway":
			option = client.TKEGateway()
		case "tke-monitor-api":
			option = client.TKEMonitorAPI()
		case "tke-monitor-controller":
			option = client.TKEMonitorController()
		case "tke-notify-api":
			option = client.TKENotifyAPI()
		case "tke-notify-controller":
			option = client.TKENotifyController()
		case "tke-platform-api":
			option = client.TKEPlatformAPI()
		case "tke-platform-controller":
			option = client.TKEPlatformController()
		}
		if option != nil {
			yaml, _ := template.ParseFile(filePath, option)
			fileName := Output + "/" + component + ".yaml"
			err := ioutil.WriteFile(fileName, yaml, 0644)
			if err != nil {
				log.Error(err.Error())
				return
			}
			log.Info(fileName + " is Created")
		}
	}
}
