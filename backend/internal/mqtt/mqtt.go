// Package mqtt publishes each host to Home Assistant as a connectivity binary_sensor
// using MQTT discovery (https://www.home-assistant.io/integrations/mqtt/#mqtt-discovery).
package mqtt

import (
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/models"
)

var (
	mu      sync.Mutex
	client  paho.Client
	usedURL string
)

// Add - publish discovery config, state and attributes for one host
func Add(host models.Host) {

	c := connect()
	if c == nil {
		return
	}

	id := "lanlookout_" + strings.ReplaceAll(host.Mac, ":", "")
	base := "lanlookout/" + id
	name := host.Name
	if name == "" {
		name = host.Mac
	}

	cfg, _ := json.Marshal(map[string]any{
		"name":                  nil, // entity takes the device name
		"unique_id":             id,
		"object_id":             id,
		"device_class":          "connectivity",
		"state_topic":           base + "/state",
		"json_attributes_topic": base + "/attr",
		"device": map[string]any{
			"identifiers":  []string{id},
			"name":         name,
			"connections":  [][]string{{"mac", host.Mac}},
			"manufacturer": host.Hw,
		},
	})
	attr, _ := json.Marshal(map[string]any{
		"ip": host.IP, "mac": host.Mac, "iface": host.Iface, "vendor": host.Hw,
		"known": host.Known == 1, "type": host.Type, "tags": host.Tags, "last_seen": host.Date,
	})
	state := "OFF"
	if host.Now == 1 {
		state = "ON"
	}

	pub(c, "homeassistant/binary_sensor/"+id+"/config", cfg)
	pub(c, base+"/attr", attr)
	pub(c, base+"/state", []byte(state))
}

func pub(c paho.Client, topic string, payload []byte) {
	t := c.Publish(topic, 0, true, payload)
	if t.WaitTimeout(5*time.Second) && t.Error() != nil {
		slog.Error("MQTT publish failed", "topic", topic, "err", t.Error())
	}
}

// connect - reuse the client, reconnect when the URL or login changed
func connect() paho.Client {

	mu.Lock()
	defer mu.Unlock()

	key := conf.AppConfig.MqttURL + "|" + conf.AppConfig.MqttUser + "|" + conf.AppConfig.MqttPass
	if client != nil && usedURL == key && client.IsConnected() {
		return client
	}
	if client != nil {
		client.Disconnect(250)
		client = nil
	}

	opts := paho.NewClientOptions().
		AddBroker(conf.AppConfig.MqttURL).
		SetClientID("lanlookout").
		SetUsername(conf.AppConfig.MqttUser).
		SetPassword(conf.AppConfig.MqttPass).
		SetConnectTimeout(5 * time.Second).
		SetAutoReconnect(true)

	c := paho.NewClient(opts)
	t := c.Connect()
	if !t.WaitTimeout(5*time.Second) || t.Error() != nil {
		slog.Error("MQTT connect failed", "url", conf.AppConfig.MqttURL, "err", t.Error())
		return nil
	}
	slog.Info("Connected to MQTT", "url", conf.AppConfig.MqttURL)

	client, usedURL = c, key
	return client
}
