package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/RackHD/inservice/uuid"
	"github.com/spf13/viper"
)

func keys(m map[string]interface{}) []string {
	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}

	return keys
}

var binaryName, buildDate, buildUser, commitHash, goVersion, osArch, releaseVersion string

func main() {
	log.Println(binaryName)
	log.Println("  Release version: " + releaseVersion)
	log.Println("  Built On: " + buildDate)
	log.Println("  Build By: " + buildUser)
	log.Println("  Commit Hash: " + commitHash)
	log.Println("  Go version: " + goVersion)
	log.Println("  OS/Arch: " + osArch)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalf("InService Agent Configuration Error: %s\n", err)
	}

	viper.SetConfigName("inservice")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/inservice.d")
	viper.AddConfigPath(dir)
	viper.AddConfigPath("$GOPATH/bin")
	viper.AddConfigPath("$HOME")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("InService Agent Configuration Error: %s\n", err)
	}

	log.Printf("InService Agent Configuration: %s\n", viper.ConfigFileUsed())

	// Host configuration end point for SSDP location advertisement.
	http := HTTPServer{
		Address:    viper.GetString("agent.http.address"),
		Port:       viper.GetInt("agent.http.port"),
		ConfigFile: viper.ConfigFileUsed(),
		URI:        viper.GetString("agent.http.uri"),
	}

	go http.Serve()

	// Host SSDP Server for advertising Agent/Plugin capabilities.
	ssdp := SSDPServer{
		ServiceType: viper.GetString("agent.ssdp.serviceType"),
		DeviceUUID:  uuid.GetUUID(viper.GetString("agent.ssdp.cacheFile")),
		Location: fmt.Sprintf(
			"http://%s:%d/%s",
			viper.GetString("agent.http.address"),
			viper.GetInt("agent.http.port"),
			viper.GetString("agent.http.uri"),
		),
		MaxAge: viper.GetInt("agent.ssdp.maxAge"),
	}

	go ssdp.Serve()

	// Host Plugins
	plugins, err := NewPluginServer(
		viper.GetString("agent.grpc.address"),
		viper.GetInt("agent.grpc.port"),
		keys(viper.GetStringMap("plugins")),
	)
	if err != nil {
		log.Fatalf("InService Agent Plugin Server Error: %s\n", err)
	}

	plugins.Serve()
}
