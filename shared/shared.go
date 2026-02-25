package shared

type MessageAbortedError struct {
	Data MessageAbortedErrorData `json:"data,required"`
	Name MessageAbortedErrorName `json:"name,required"`
}


type MessageAbortedErrorData struct {
	Message string `json:"message,required"`
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
	Data ProviderAuthErrorData `json:"data,required"`
	Name ProviderAuthErrorName `json:"name,required"`
}


type ProviderAuthErrorData struct {
	Message    string `json:"message,required"`
	ProviderID string `json:"providerID,required"`
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
	Data UnknownErrorData `json:"data,required"`
	Name UnknownErrorName `json:"name,required"`
}


type UnknownErrorData struct {
	Message string `json:"message,required"`
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
