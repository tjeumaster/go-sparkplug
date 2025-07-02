package sparkplug

import (
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	MqttClient mqtt.Client
	Config	 *SparkplugConfig
	BdSeq uint64
	Seq uint64
	mu sync.Mutex
}

type SparkplugConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	ClientID string
	GroupID  string
	NodeID   string
}

type SparkplugDevice interface {
	GetDeviceID() string
	GetMetricValues() map[string]any
	GetMetricValueByName(name string) any
}

func NewClient(config *SparkplugConfig) *Client {
	return &Client{
		Config: config,
		BdSeq: 0,
		Seq:    0,
		MqttClient: nil,
	}
}

func (c *Client) getSeq() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	seq := c.Seq
	c.Seq = (c.Seq + 1) % 256
	return seq
}

func (c *Client) Start() {
	mqttBroker := fmt.Sprintf("tcp://%s:%d", c.Config.Host, c.Config.Port)
	opts := mqtt.NewClientOptions().
			AddBroker(mqttBroker).
			SetClientID(c.Config.ClientID).
			SetOnConnectHandler(func(client mqtt.Client) {
				c.PublishNBIRTH()
			})

	c.MqttClient = mqtt.NewClient(opts)

	for {
		token := c.MqttClient.Connect()
		token.Wait()
		if err := token.Error(); err == nil {
			break
		} else {
			fmt.Printf("Connection failed: %v. Retrying in 10 seconds...\n", err)
			time.Sleep(10 * time.Second)
		}
	}

	log.Printf("Connected to MQTT broker at %s", mqttBroker)

}

func (c *Client) publish(topic string, payload []byte, retained bool) error {
	token := c.MqttClient.Publish(topic, 0, retained, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("failed to publish to topic %s: %w", topic, err)
	}
	return nil
}

func (c *Client) PublishNBIRTH() error {
	payload, err := c.buildNBIRTHPayload()
	if err != nil {
		return fmt.Errorf("failed to build NBIRTH payload: %w", err)
	}

	topic := fmt.Sprintf("spBv1.0/%s/NBIRTH/%s", c.Config.GroupID, c.Config.NodeID)
	if err := c.publish(topic, payload, true); err != nil {
		return fmt.Errorf("failed to publish NBIRTH: %w", err)
	}

	log.Printf("Published NBIRTH to topic %s", topic)

	c.BdSeq++

	return nil
}

func (c *Client) PublishNDEATH() error {
	payload, err := c.buildNDEATHPayload()
	if err != nil {
		return fmt.Errorf("failed to build NDEATH payload: %w", err)
	}

	topic := fmt.Sprintf("spBv1.0/%s/NDEATH/%s", c.Config.GroupID, c.Config.NodeID)
	if err := c.publish(topic, payload, true); err != nil {
		return fmt.Errorf("failed to publish NDEATH: %w", err)
	}

	log.Printf("Published NDEATH to topic %s", topic)

	return nil
}

func (c *Client) PublishDBIRTH(device SparkplugDevice) error {
	payload, err := c.buildDBIRTHPayload(device)
	if err != nil {
		return fmt.Errorf("failed to build DBIRTH payload: %w", err)
	}

	topic := fmt.Sprintf("spBv1.0/%s/DBIRTH/%s/%s", c.Config.GroupID, c.Config.NodeID, device.GetDeviceID())
	if err := c.publish(topic, payload, false); err != nil {
		return fmt.Errorf("failed to publish DBIRTH: %w", err)
	}

	log.Printf("Published DBIRTH for device %s to topic %s", device.GetDeviceID(), topic)

	return nil
}

func (c *Client) PublishDDEATH(device SparkplugDevice) error {
	payload, err := c.buildDDEATHPayload(device)
	if err != nil {
		return fmt.Errorf("failed to build DDEATH payload: %w", err)
	}

	topic := fmt.Sprintf("spBv1.0/%s/DDEATH/%s/%s", c.Config.GroupID, c.Config.NodeID, device.GetDeviceID())
	if err := c.publish(topic, payload, false); err != nil {
		return fmt.Errorf("failed to publish DDEATH: %w", err)
	}

	log.Printf("Published DDEATH for device %s to topic %s", device.GetDeviceID(), topic)

	return nil
}







