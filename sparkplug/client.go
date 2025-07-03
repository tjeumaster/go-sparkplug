package sparkplug

import (
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tjeumaster/go-sparkplug/sproto"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	MqttClient mqtt.Client
	Config	 *SparkplugConfig
	BdSeq uint64
	currentBdSeq uint64
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

func (c *Client) getBdSeq() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.currentBdSeq = c.BdSeq
	c.BdSeq++
	return c.currentBdSeq
}


const (
	NBIRTH = "NBIRTH"
	NDEATH = "NDEATH"
	DBIRTH = "DBIRTH"
	DDEATH = "DDEATH"
	NDATA  = "NDATA"
	DDATA  = "DDATA"
	NCMD   = "NCMD"
	DCMD   = "DCMD"
)

func (c *Client) getTopic(d SparkplugDevice, topicType string) string {
	switch topicType {
	case NBIRTH:
		return fmt.Sprintf("spBv1.0/%s/NBIRTH/%s", c.Config.GroupID, c.Config.NodeID)

	case NDEATH:
		return fmt.Sprintf("spBv1.0/%s/NDEATH/%s", c.Config.GroupID, c.Config.NodeID)

	case DBIRTH:
		if d == nil {
			return ""
		}
		return fmt.Sprintf("spBv1.0/%s/DBIRTH/%s/%s", c.Config.GroupID, c.Config.NodeID, d.GetDeviceID())

	case DDEATH:
		if d == nil {
			return ""
		}
		return fmt.Sprintf("spBv1.0/%s/DDEATH/%s/%s", c.Config.GroupID, c.Config.NodeID, d.GetDeviceID())

	case NDATA:
		return fmt.Sprintf("spBv1.0/%s/NDATA/%s", c.Config.GroupID, c.Config.NodeID)

	case DDATA:
		if d == nil {
			return ""
		}
		
		return fmt.Sprintf("spBv1.0/%s/DDATA/%s/%s", c.Config.GroupID, c.Config.NodeID, d.GetDeviceID())

	case NCMD:
		return fmt.Sprintf("spBv1.0/%s/NCMD/%s", c.Config.GroupID, c.Config.NodeID)

	case DCMD:
		return fmt.Sprintf("spBv1.0/%s/DCMD/%s/+", c.Config.GroupID, c.Config.NodeID)

	default:
		return ""
	}
}

func (c *Client) Start() {
	mqttBroker := fmt.Sprintf("tcp://%s:%d", c.Config.Host, c.Config.Port)

	ndeathPayload, err := c.buildNDEATHPayload()
	if err != nil {
		log.Fatalf("Failed to build NDEATH payload: %v", err)
	}
	
	ndeathTopic := c.getTopic(nil, NDEATH)
	if ndeathTopic == "" {
		log.Fatalf("Failed to get NDEATH topic")
	}

	opts := mqtt.NewClientOptions().
			AddBroker(mqttBroker).
			SetClientID(c.Config.ClientID).
			SetOnConnectHandler(func(client mqtt.Client) {
				c.PublishNBIRTH()
			}).
			SetWill(ndeathTopic, string(ndeathPayload), 0, true)

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

	c.MqttClient.Subscribe(c.getTopic(nil, NCMD), 0, c.onCommandReceived)
	c.MqttClient.Subscribe(c.getTopic(nil, DCMD), 0, c.onCommandReceived)

	log.Printf("Connected to MQTT broker at %s", mqttBroker)
}

func (c *Client) Stop() {
	if c.MqttClient != nil && c.MqttClient.IsConnected() {
		ddeathPayload, err := c.buildNDEATHPayload()
		if err != nil {
			log.Printf("Failed to build NDEATH payload: %v", err)
		} else {
			ddeathTopic := c.getTopic(nil, NDEATH)
			if ddeathTopic != "" {
				token := c.MqttClient.Publish(ddeathTopic, 0, true, string(ddeathPayload))
				token.Wait()
				if err := token.Error(); err != nil {
					log.Printf("Failed to publish NDEATH: %v", err)
				} else {
					log.Printf("Published NDEATH to topic %s", ddeathTopic)
				}
			} else {
				log.Printf("Failed to get NDEATH topic")
			}
		}
		c.MqttClient.Disconnect(250)
		log.Printf("Disconnected from MQTT broker")
	} else {
		log.Printf("MQTT client is not connected, nothing to stop")
	}
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

	topic := c.getTopic(nil, NBIRTH)
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
	payload, err := c.buildDDEATHPayload()
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

func (c *Client) PublishNDATA(metricValues map[string]any) error {
	payload, err := c.buildNDATAPayload(metricValues)
	if err != nil {
		return fmt.Errorf("failed to build NDATA payload: %w", err)
	}

	topic := fmt.Sprintf("spBv1.0/%s/NDATA/%s", c.Config.GroupID, c.Config.NodeID)
	if err := c.publish(topic, payload, false); err != nil {
		return fmt.Errorf("failed to publish NDATA: %w", err)
	}

	log.Printf("Published NDATA to topic %s", topic)

	return nil
}

func (c *Client) PublishDDATA(device SparkplugDevice, metricValues map[string]any) error {
	payload, err := c.buildDDATAPayload(metricValues)
	if err != nil {
		return fmt.Errorf("failed to build DDATA payload: %w", err)
	}

	topic := fmt.Sprintf("spBv1.0/%s/DDATA/%s/%s", c.Config.GroupID, c.Config.NodeID, device.GetDeviceID())
	if err := c.publish(topic, payload, false); err != nil {
		return fmt.Errorf("failed to publish DDATA: %w", err)
	}

	log.Printf("Published DDATA for device %s to topic %s", device.GetDeviceID(), topic)

	return nil
}

func (c *Client) onCommandReceived(client mqtt.Client, msg mqtt.Message) {
	payloadBytes := msg.Payload()

    var payload sproto.Payload
    err := proto.Unmarshal(payloadBytes, &payload)
    if err != nil {
        log.Printf("Failed to decode command payload: %v", err)
        return
    }

    for _, metric := range payload.Metrics {
        c.handleCommandMetric(metric, msg.Topic())
    }
}

func (c *Client) handleCommandMetric(metric *sproto.Payload_Metric, topic string) error{
	name := metric.GetName()
	
	switch name {
	case "Node Control/Rebirth":
		log.Printf("Received Rebirth command on topic %s", topic)
		err := c.PublishNBIRTH()
		if err != nil {
			return err
		}

		log.Printf("Published NBIRTH in response to Rebirth command on topic %s", topic)
		return nil
		
	case "Node Control/Reboot":
		log.Printf("Received Reboot command on topic %s", topic)
		return nil

	default:
		log.Printf("Received unknown command '%s' on topic %s", name, topic)
		return nil
	}
}



