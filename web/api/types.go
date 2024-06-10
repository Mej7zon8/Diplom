package api

type ResponseAny = Response[any]

type Response[T any] struct {
	Data   T      `json:"data"`
	Status int    `json:"status"`
	Error  string `json:"error"`

	// RealError is the error that caused the response to fail.
	// The error is used for logging purposes.
	RealError error `json:"-"`
}

func ResponseBad(status int, realError error, text string) ResponseAny {
	return ResponseAny{Status: status, Error: text, RealError: realError}
}

func ResponseBadRequest(realError error) ResponseAny {
	return ResponseAny{Status: 400, Error: "bad request", RealError: realError}
}

func ResponseInternalServerError(realError error) ResponseAny {
	return ResponseAny{Status: 500, Error: "internal server error", RealError: realError}
}
