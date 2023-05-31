package arc

type TxStatus string

const (
	// List of statuses available here: https://github.com/bitcoin-sv/arc
	UNKNOWN              TxStatus = "UNKNOWN"              // 0
	QUEUED               TxStatus = "QUEUED"               // 1
	RECEIVED             TxStatus = "RECEIVED"             // 2
	STORED               TxStatus = "STORED"               // 3
	ANNOUNCED_TO_NETWORK TxStatus = "ANNOUNCED_TO_NETWORK" // 4
	REQUESTED_BY_NETWORK TxStatus = "REQUESTED_BY_NETWORK" // 5
	SENT_TO_NETWORK      TxStatus = "SENT_TO_NETWORK"      // 6
	ACCEPTED_BY_NETWORK  TxStatus = "ACCEPTED_BY_NETWORK"  // 7
	SEEN_ON_NETWORK      TxStatus = "SEEN_ON_NETWORK"      // 8
	MINED                TxStatus = "MINED"                // 9
	CONFIRMED            TxStatus = "CONFIRMED"            // 108
	REJECTED             TxStatus = "REJECTED"             // 109
)

func (s TxStatus) String() string {
	statuses := map[TxStatus]string{
		UNKNOWN:              "UNKNOWN",
		QUEUED:               "QUEUED",
		RECEIVED:             "RECEIVED",
		STORED:               "STORED",
		ANNOUNCED_TO_NETWORK: "ANNOUNCED_TO_NETWORK",
		REQUESTED_BY_NETWORK: "REQUESTED_BY_NETWORK",
		SENT_TO_NETWORK:      "SENT_TO_NETWORK",
		ACCEPTED_BY_NETWORK:  "ACCEPTED_BY_NETWORK",
		SEEN_ON_NETWORK:      "SEEN_ON_NETWORK",
		MINED:                "MINED",
		CONFIRMED:            "CONFIRMED",
		REJECTED:             "REJECTED",
	}

	if status, ok := statuses[s]; ok {
		return status
	}

	return "Can't parse status"
}

func MapTxStatusToInt(status TxStatus) (int, bool) {
	waitForStatusMap := map[TxStatus]int{
		UNKNOWN:              0,
		QUEUED:               1,
		RECEIVED:             2,
		STORED:               3,
		ANNOUNCED_TO_NETWORK: 4,
		REQUESTED_BY_NETWORK: 5,
		SENT_TO_NETWORK:      6,
		ACCEPTED_BY_NETWORK:  7,
		SEEN_ON_NETWORK:      8,
		MINED:                9,
		CONFIRMED:            108,
		REJECTED:             109,
	}

	value, ok := waitForStatusMap[status]
	return value, ok
}
