package sparkplug

import (
	"fmt"
	"time"
	"github.com/tjeumaster/go-sparkplug/sproto"
	"google.golang.org/protobuf/proto"
)

func (c *Client) buildNBIRTHPayload() ([]byte, error) {
	payload := &sproto.Payload{
		Timestamp: proto.Uint64(uint64(time.Now().UnixMilli())),
		Seq:       proto.Uint64(c.getSeq()),
		Metrics: []*sproto.Payload_Metric{
			ToMetric("bdSeq", c.BdSeq),
			ToMetric("Node Control/Rebirth", false),
			ToMetric("Node Control/Reboot", false),
		},
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal NBIRTH payload: %w", err)
	}

	return payloadBytes, nil
}

func (c *Client) buildNDEATHPayload() ([]byte, error) {
	bdSeq := c.BdSeq - 1
	payload := &sproto.Payload{
		Timestamp: proto.Uint64(uint64(time.Now().UnixMilli())),
		Seq:       proto.Uint64(c.getSeq()),
		Metrics: []*sproto.Payload_Metric{
			ToMetric("bdSeq", bdSeq),
			ToMetric("Node Control/Rebirth", false),
			ToMetric("Node Control/Reboot", false),
		},
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal NDEATH payload: %w", err)
	}

	return payloadBytes, nil
}

func (c *Client) buildDBIRTHPayload(d SparkplugDevice) ([]byte, error) {
	values := d.GetMetricValues()
	metrics := make([]*sproto.Payload_Metric, 0, len(values))
	for name, value := range values {
		metric := ToMetric(name, value)
		metrics = append(metrics, metric)
	}

	payload := &sproto.Payload{
		Timestamp: proto.Uint64(uint64(time.Now().UnixMilli())),
		Seq:       proto.Uint64(c.getSeq()),
		Metrics: metrics,
	}
	
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DBIRTH payload: %w", err)
	}

	return payloadBytes, nil
}

func (c *Client) buildDDEATHPayload(d SparkplugDevice) ([]byte, error) {
	payload := &sproto.Payload{
		Timestamp: proto.Uint64(uint64(time.Now().UnixMilli())),
		Seq:       proto.Uint64(c.getSeq()),
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DDEATH payload: %w", err)
	}

	return payloadBytes, nil
}