package spb

import (
	"time"

	"github.com/tjeumaster/go-sparkplug/sproto"
	"google.golang.org/protobuf/proto"
)

func ToMetric(name string, value any) *sproto.Payload_Metric {
	var metricType sproto.DataType
	var metricValue sproto.Payload_Metric_Value

	switch v := value.(type) {
	case int:
		metricType = sproto.DataType_Int64
		metricValue = &sproto.Payload_Metric_LongValue{
			LongValue: uint64(v),
		}

	case int32:
		metricType = sproto.DataType_Int32
		metricValue = &sproto.Payload_Metric_IntValue{
			IntValue: uint32(v),
		}
	case int64:
		metricType = sproto.DataType_Int64
		metricValue = &sproto.Payload_Metric_LongValue{
			LongValue: uint64(v),
		}

	case uint32:
		metricType = sproto.DataType_UInt32
		metricValue = &sproto.Payload_Metric_IntValue{
			IntValue: uint32(v),
		}

	case uint64:
		metricType = sproto.DataType_UInt64
		metricValue = &sproto.Payload_Metric_LongValue{
			LongValue: uint64(v),
		}

	case float32:
		metricType = sproto.DataType_Float
		metricValue = &sproto.Payload_Metric_FloatValue{
			FloatValue: float32(v),
		}

	case float64:
		metricType = sproto.DataType_Double
		metricValue = &sproto.Payload_Metric_DoubleValue{
			DoubleValue: float64(v),
		}

	case string:
		metricType = sproto.DataType_String
		metricValue = &sproto.Payload_Metric_StringValue{
			StringValue: string(v),
		}

	case bool:
		metricType = sproto.DataType_Boolean
		metricValue = &sproto.Payload_Metric_BooleanValue{
			BooleanValue: bool(v),
		}

	case []byte:
		metricType = sproto.DataType_Bytes
		metricValue = &sproto.Payload_Metric_BytesValue{
			BytesValue: v,
		}

	default:
		return nil
	}

	metric := &sproto.Payload_Metric{
		Name:      proto.String(name),
		Timestamp: proto.Uint64(uint64(time.Now().UnixMilli())),
		Datatype:  proto.Uint32(uint32(metricType)),
		Value:     metricValue,
	}

	return metric
}
