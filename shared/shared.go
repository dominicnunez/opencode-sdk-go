package shared

type MessageAbortedError struct {
	Data MessageAbortedErrorData `json:"data"`
	Name MessageAbortedErrorName `json:"name"`
}


type MessageAbortedErrorData struct {
	Message string `json:"message"`
}

type MessageAbortedErrorName string

const (
	MessageAbortedErrorNameMessageAbortedError MessageAbortedErrorName = "MessageAbortedError"
)

func (r MessageAbortedErrorName) IsKnown() bool {
	switch r {
	case MessageAbortedErrorNameMessageAbortedError:
		return true
	}
	return false
}

type ProviderAuthError struct {
	Data ProviderAuthErrorData `json:"data"`
	Name ProviderAuthErrorName `json:"name"`
}


type ProviderAuthErrorData struct {
	Message    string `json:"message"`
	ProviderID string `json:"providerID"`
}

type ProviderAuthErrorName string

const (
	ProviderAuthErrorNameProviderAuthError ProviderAuthErrorName = "ProviderAuthError"
)

func (r ProviderAuthErrorName) IsKnown() bool {
	switch r {
	case ProviderAuthErrorNameProviderAuthError:
		return true
	}
	return false
}

type UnknownError struct {
	Data UnknownErrorData `json:"data"`
	Name UnknownErrorName `json:"name"`
}


type UnknownErrorData struct {
	Message string `json:"message"`
}

type UnknownErrorName string

const (
	UnknownErrorNameUnknownError UnknownErrorName = "UnknownError"
)

func (r UnknownErrorName) IsKnown() bool {
	switch r {
	case UnknownErrorNameUnknownError:
		return true
	}
	return false
}
