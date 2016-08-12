package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bndr/gotabulate"
	"github.com/coreos/etcd/client"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

type ConfigBackend struct {
	Type string
}

type Config struct {
	Network   string
	SubnetLen int
	Backend   ConfigBackend
}

type BackendData struct {
	VtepMac string
}

type Subnet struct {
	PublicIP    string
	BackendType string
	BackendData BackendData
}

func main() {
	app := cli.NewApp()

	app.Name = "flannelctl"
	app.Usage = "flannel command utility"

	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(c)

	app.Commands = []cli.Command{
		{
			Name:  "config",
			Usage: "view flannel configuration",
			Action: func(c *cli.Context) error {
				resp, err := kapi.Get(context.Background(), "/coreos.com/network/config", nil)
				if err != nil {
					log.Fatal(err)
				}

				config := Config{}
				json.Unmarshal([]byte(resp.Node.Value), &config)

				output := [][]string{}
				output = append(output, []string{"Network", config.Network})
				output = append(output, []string{"SubnetLen", strconv.Itoa(config.SubnetLen)})
				backend, _ := json.Marshal(config.Backend)
				output = append(output, []string{"Backend", string(backend)})

				tabulate := gotabulate.Create(output)
				tabulate.SetHeaders([]string{"", "Value"})
				fmt.Println(tabulate.Render("grid"))

				return nil
			},
		},
		{
			Name:  "subnet",
			Usage: "list fannel subnets",
			Action: func(c *cli.Context) error {
				subnetkey := "/coreos.com/network/subnets"
				resp, err := kapi.Get(context.Background(), subnetkey, nil)
				if err != nil {
					log.Fatal(err)
				}
				output := [][]string{}
				for _, v := range resp.Node.Nodes {
					subnet := strings.Replace(v.Key[len(subnetkey)+1:], "-", "/", -1)
					subnetInfo := Subnet{}
					json.Unmarshal([]byte(v.Value), &subnetInfo)

					output = append(output, []string{subnet, subnetInfo.PublicIP, subnetInfo.BackendType})
				}

				tabulate := gotabulate.Create(output)
				tabulate.SetHeaders([]string{"Subnet", "PublicIP", "BackendType"})
				fmt.Println(tabulate.Render("grid"))
				return nil
			},
		},
	}

	app.Run(os.Args)

}
