package constants

// Level of each message types
type FeeLevel int

const (
	NONE FeeLevel = 0
	HIGH FeeLevel = 1
	MED  FeeLevel = 2
	LOW  FeeLevel = 3
)

var LEVELS = map[string]FeeLevel{
	"MsgSend":     LOW,
	"MsgCreate":   HIGH,
	"MsgUpdate":   MED,
	"MsgDelete":   LOW,
	"MsgBook":     HIGH,
	"MsgComplete": MED,
}

var FEE_LEVELS = map[FeeLevel]int{
	HIGH: 3,
	MED:  2,
	LOW:  1,
	NONE: 0,
}

const FEE_DENOM = "SHR"
