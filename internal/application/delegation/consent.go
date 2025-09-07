package delegationapp

type ConsentStatus string

const (
	ConsentGranted  ConsentStatus = "granted"
	ConsentDenied   ConsentStatus = "denied"
	ConsentRequired ConsentStatus = "required"
)


