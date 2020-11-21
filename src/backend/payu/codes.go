package payugo

// https://developers.payu.com/pl/restapi.html#references_statuses
type Status string

const (
	// 201 Created
	StatusSuccess   Status = "SUCCESS"
	StatusCompleted Status = "COMPLETED"
	// 302 Found
	StatusWarningContinueRedirect Status = "WARNING_CONTINUE_REDIRECT"
	StatusWarningContinue3DS      Status = "WARNING_CONTINUE_3DS"
	StatusWarningContinueCVV      Status = "WARNING_CONTINUE_CVV"
	// 400 Bad request
	StatusErrorSyntax       Status = "ERROR_SYNTAX"
	StatusErrorValueInvalid Status = "ERROR_VALUE_INVALID"
	// 401 Unauthorized
	StatusUnauthorized Status = "UNAUTHORIZED"
	// 404 Not found
	StatusDataNotFound Status = "DATA_NOT_FOUND"
)

func (s Status) String() string { return string(s) }
