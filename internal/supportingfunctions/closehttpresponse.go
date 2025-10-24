package supportingfunctions

import "net/http"

// CloseHTTPResponse закрытие http соединения с предварительной проверкой
func CloseHTTPResponse(res *http.Response) {
	if res == nil || res.Body == nil {
		return
	}

	res.Body.Close()
}
