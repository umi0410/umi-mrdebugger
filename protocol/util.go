package protocol

// MakeCEKResponse creates CEKResponse instance with given params
func MakeCEKResponse(sessionAttributes map[string]string, responsePayload CEKResponsePayload) CEKResponse {
	response := CEKResponse{
		SessionAttributes: sessionAttributes,
		Response:          responsePayload,
	}

	return response
}

// MakeOutputSpeech creates OutputSpeech instance with given params
func MakeOutputSpeech(msg string) OutputSpeech {
	return OutputSpeech{
		Type: "SimpleSpeech",
		Values: Value{
			Lang:  "ko",
			Value: msg,
			Type:  "PlainText",
		},
	}
}
