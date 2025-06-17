package constants

type FieldStatusString string

const (
	AvailableStatus FieldStatusString = "pending"
	BrokerStatus    FieldStatusString = "settlement"
)

func (p FieldStatusString) String() string {
	return string(p)
}
