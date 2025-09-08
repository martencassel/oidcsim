package delegation

type ConsentDecision int

const (
	ConsentDecisionNone ConsentDecision = iota
	ConsentDecisionApprove
	ConsentDecisionDeny
	ConsentStatusGranted
)

func (cd ConsentDecision) String() string {
	switch cd {
	case ConsentDecisionNone:
		return "none"
	case ConsentStatusGranted:
		return "granted"
	case ConsentDecisionApprove:
		return "approve"
	case ConsentDecisionDeny:
		return "deny"
	default:
		return "none"
	}
}
