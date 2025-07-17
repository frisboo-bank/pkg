package contracts

type StatusType = string

type CheckStatus struct {
	statusUp   StatusType
	statusDown StatusType
	Status     StatusType `json:"status"`
}

type CheckAllStatus map[string]CheckStatus

func (ch CheckAllStatus) IsAllUP() bool {
	for _, status := range ch {
		if !status.IsUP() {
			return false
		}
	}

	return true
}

func NewCheckStatus(err error, statusUp StatusType, statusDown StatusType) CheckStatus {
	status := CheckStatus{
		statusUp:   statusUp,
		statusDown: statusDown,
	}

	if err != nil {
		status.Status = statusDown
		return status
	}

	status.Status = statusUp

	return status
}

func (s CheckStatus) IsUP() bool {
	return s.Status == s.statusUp
}
