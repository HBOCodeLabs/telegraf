package appdynamics

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"

	"github.com/parnurzeal/gorequest"
)

var debugEnabled bool

type AppDynamics struct {
	ControllerHostname string
	ControllerPort int
	ControllerSslEnabled bool
	AccountName string
	AccountAccessKey string
	ApplicationName string
	TierName string
	TierId int
	NodeName string
	Debug bool
	MaxProcessSnapshotsPerPeriod int
	ProcessSnapshotCountResetPeriodSeconds int
	AutoSnapshotDurationSeconds int
	RoUser string
	RoPassword string
	ExcludeAgentFromCallGraph bool
	MachineAgentHostname string
	ReplaceAllDots bool
	Timeout int

	conn net.Conn
	machineAgentEndpoint string
}

type Tier struct {
	Id int `json:"id"`
	Name string `json:"name"`
	NumberOfNodes int `json:"numberOfNodes"`
	AgentType string `json:"agentType"`
}

func (a *AppDynamics) Description() string {
	return "Configuration for sending metrics to AppDynamics Machine Agent"
}

func (a *AppDynamics) SampleConfig() string {
	var sampleConfig = `
		## Hostname of AppDynamics Machine Agent
		host = "http://localhost:8293"
		## Connection timeout, default to 5 seconds
		timeout = 5
	`
	return sampleConfig
}

func (a *AppDynamics) Connect() error {
	// Make a connection to the AppDynamics agent
	if a.Timeout <=0 {
		a.Timeout = 5
	}
	conn, err := net.DialTimeout("tcp", a.MachineAgentHostname, time.Duration(a.Timeout) * time.Second)
	if err == nil {
		a.conn = conn
		return nil
	} else {
		return err
	}
}

func (a *AppDynamics) Close() error {
	// Close connection
	err := a.conn.Close()
	return err
}

func (a *AppDynamics) Write(metrics []telegraf.Metric) error {
	debug(fmt.Sprintf("writing metrics at %s", time.Now().UTC()))
	for _, metric := range metrics {
		// write `pt` to the output sink here
	}
	return nil
}

func init() {
	appd := AppDynamics{}
	debugEnabled = appd.Debug
	debug(fmt.Sprintf("%+v", &appd))
	outputs.Add("appdynamics", func() telegraf.Output { return &appd })
}

func debug(msg string) {
	if debugEnabled {
		fmt.Println(msg)
	}
}

func (a *AppDynamics) getTierId() {
	if a.TierId <= 0 {
		path := fmt.Sprintf("/controller/rest/applications/%s/tiers/%s?output=JSON", a.AccountName, a.TierName)
		request := gorequest.New().SetBasicAuth(fmt.Sprintf("%s@%s", a.RoUser, a.AccountName), a.RoPassword)
		_, body, err := request.Get(fmt.Sprintf("%s:%d%s", a.ControllerHostname, a.ControllerPort, path)).End()
		if err == nil {
			var data []Tier
			err := json.Unmarshal([]byte(body), &data)
			if err == nil {
				a.TierId = data[0].Id
				debug(fmt.Sprintf("TierId=%d", a.TierId))
			}
		}
	} else {
		debug(fmt.Sprintf("TierId=%d", a.TierId))
	}
	a.machineAgentEndpoint = fmt.Sprintf(
		"%s/machineagent/metrics?name=Server|Component:%s|Custom Metrics|",
		a.MachineAgentHostname,
		a.TierId,
	)
	debug(fmt.Sprintf("machineAgentEndpoint=%s", a.machineAgentEndpoint))
}

var nodeinfo = `
{
    "controllerHostName": "hbo-go.saas.appdynamics.com",
    "controllerPort": 443,
    "controllerSslEnabled": true,
    "accountName": "hbo-go",
    "accountAccessKey": "d27dd42da851",
    "applicationName": "hurley.production",
    "tierName": "accounts",
    "nodeName": "us-west-2.accounts.i-308478ed.1.1.120",
    "debug": false,
    "maxProcessSnapshotsPerPeriod": 1,
    "processSnapshotCountResetPeriodSeconds": 60,
    "autoSnapshotDurationSeconds": 5,
    "rouser": "apirouser",
    "ropassword": "HBORocks2!",
    "serviceName": "accounts",
    "excludeAgentFromCallGraph": true
}
appdynamics: {
	nodeinfoPath: '/tmp/nodeinfo.js',
	debug: true,
	retry: 20,
	replaceAll: {{ replaceAllDots|lower }},
	endpoint: 'http://localhost:8293'
}
`