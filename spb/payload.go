package spb

import (
	"fmt"
	"time"

	"github.com/tjeumaster/go-sparkplug/sproto"
	"google.golang.org/protobuf/proto"
)

func (c *Client) buildNBIRTHPayload() ([]byte, error) {
	payload := &sproto.Payload{
		Timestamp: proto.Uint64(uint64(time.Now().UnixMilli())),
		Seq:       proto.Uint64(c.Seq),
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
	payload := &sproto.Payload{
		Timestamp: proto.Uint64(uint64(time.Now().UnixMilli())),
		Seq:       proto.Uint64(c.Seq),
		Metrics: []*sproto.Payload_Metric{
			ToMetric("bdSeq", c.BdSeq),
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

func (c *Client) buildDBIRTHPayload(d Device) ([]byte, error) {
	values := d.GetMetricValues()
	metrics := make([]*sproto.Payload_Metric, 0, len(values))
	for name, value := range values {
		metric := ToMetric(name, value)
		metrics = append(metrics, metric)
	}

	payload := &sproto.Payload{
		Timestamp: proto.Uint64(uint64(time.Now().UnixMilli())),
		Seq:       proto.Uint64(c.Seq),
		Metrics: metrics,
	}
	
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DBIRTH payload: %w", err)
	}

	return payloadBytes, nil
}

func (c *Client) buildDDEATHPayload() ([]byte, error) {
	payload := &sproto.Payload{
		Timestamp: proto.Uint64(uint64(time.Now().UnixMilli())),
		Seq:       proto.Uint64(c.Seq),
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DDEATH payload: %w", err)
	}

	return payloadBytes, nil
}

func (c *Client) buildNDATAPayload(metricValues map[string]any) ([]byte, error) {
	if len(metricValues) == 0 {
		return nil, fmt.Errorf("no metrics provided for NDATA payload")
	}

	metrics := make([]*sproto.Payload_Metric, 0, len(metricValues))
	for name, value := range metricValues {
		metric := ToMetric(name, value)
		if metric != nil {
			metrics = append(metrics, metric)
		}
	}

	payload := &sproto.Payload{
		Timestamp: proto.Uint64(uint64(time.Now().UnixMilli())),
		Seq:       proto.Uint64(c.Seq),
		Metrics:   metrics,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal NDATA payload: %w", err)
	}

	return payloadBytes, nil
}

func (c *Client) buildDDATAPayload(metricValues map[string]any) ([]byte, error) {
	if len(metricValues) == 0 {
		return nil, fmt.Errorf("no metrics provided for DDATA payload")
	}

	metrics := make([]*sproto.Payload_Metric, 0, len(metricValues))
	for name, value := range metricValues {
		metric := ToMetric(name, value)
		if metric != nil {
			metrics = append(metrics, metric)
		}
	}

	payload := &sproto.Payload{
		Timestamp: proto.Uint64(uint64(time.Now().UnixMilli())),
		Seq:       proto.Uint64(c.Seq),
		Metrics: metrics,
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal DDATA payload: %w", err)
	}

	return payloadBytes, nil
}