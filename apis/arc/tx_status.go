package arc

// TxStatus is the status of the transaction
type TxStatus string

// List of statuses available here: https://github.com/bitcoin-sv/arc
const (
	// Unknown contains value for unknown status
	Unknown TxStatus = "UNKNOWN" // 0
	// Queued contains value for queued status
	Queued TxStatus = "QUEUED" // 1
	// Received contains value for received status
	Received TxStatus = "RECEIVED" // 2
	// Stored contains value for stored status
	Stored TxStatus = "STORED" // 3
	// AnnouncedToNetwork contains value for announced to network status
	AnnouncedToNetwork TxStatus = "ANNOUNCED_TO_NETWORK" // 4
	// RequestedByNetwork contains value for requested by network status
	RequestedByNetwork TxStatus = "REQUESTED_BY_NETWORK" // 5
	// SentToNetwork contains value for sent to network status
	SentToNetwork TxStatus = "SENT_TO_NETWORK" // 6
	// AcceptedByNetwork contains value for accepted by network status
	AcceptedByNetwork TxStatus = "ACCEPTED_BY_NETWORK" // 7
	// SeenOnNetwork contains value for seen on network status
	SeenOnNetwork TxStatus = "SEEN_ON_NETWORK" // 8
	// Mined contains value for mined status
	Mined TxStatus = "MINED" // 9
	// Confirmed contains value for confirmed status
	Confirmed TxStatus = "CONFIRMED" // 108
	// Rejected contains value for rejected status
	Rejected TxStatus = "REJECTED" // 109
)

// String returns the string representation of the TxStatus
func (s TxStatus) String() string {
	statuses := map[TxStatus]string{
		Unknown:            "UNKNOWN",
		Queued:             "QUEUED",
		Received:           "RECEIVED",
		Stored:             "STORED",
		AnnouncedToNetwork: "ANNOUNCED_TO_NETWORK",
		RequestedByNetwork: "REQUESTED_BY_NETWORK",
		SentToNetwork:      "SENT_TO_NETWORK",
		AcceptedByNetwork:  "ACCEPTED_BY_NETWORK",
		SeenOnNetwork:      "SEEN_ON_NETWORK",
		Mined:              "MINED",
		Confirmed:          "CONFIRMED",
		Rejected:           "REJECTED",
	}

	if status, ok := statuses[s]; ok {
		return status
	}

	return "Can't parse status"
}

// MapTxStatusToInt maps the TxStatus to an int value
func MapTxStatusToInt(status TxStatus) (int, bool) {
	waitForStatusMap := map[TxStatus]int{
		Unknown:            0,
		Queued:             1,
		Received:           2,
		Stored:             3,
		AnnouncedToNetwork: 4,
		RequestedByNetwork: 5,
		SentToNetwork:      6,
		AcceptedByNetwork:  7,
		SeenOnNetwork:      8,
		Mined:              9,
		Confirmed:          108,
		Rejected:           109,
	}

	value, ok := waitForStatusMap[status]
	return value, ok
}
