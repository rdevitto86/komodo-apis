package httpapi

func errorObj(msg string) map[string]string {
	return map[string]string{ "error": msg }
}
