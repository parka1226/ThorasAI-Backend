package network

type NetworkTraffic struct {
	SourceIP        string        `bson:"source_ip"`
	SourcePort      int           `bson:"source_port"`
	DestinationIP   string        `bson:"destination_ip"`
	DestinationPort int           `bson:"destination_port"`
	Status          TrafficStatus `bson:"status"` // Custom type for status
}

type TrafficStatus string

const (
	StatusOK       TrafficStatus = "OK"
	StatusWarning  TrafficStatus = "Warning"
	StatusCritical TrafficStatus = "Critical"
)
