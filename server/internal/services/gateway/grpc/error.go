package grpc

func GetMessage[T any](details []any) (*T, bool) {
	for _, detail := range details {
		if info, ok := detail.(*T); ok {
			return info, true
		}
	}
	return nil, false
}
