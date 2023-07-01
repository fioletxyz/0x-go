package zerox

import (
	"encoding/json"
	"fmt"
)

type ApiErrorResponse struct {
	HTTPStatusCode int
	Message        map[string]json.RawMessage
}

func (e *ApiErrorResponse) Error() string {
	return fmt.Sprintf("Status code: %d, message: %s", e.HTTPStatusCode, e.Message)
}
