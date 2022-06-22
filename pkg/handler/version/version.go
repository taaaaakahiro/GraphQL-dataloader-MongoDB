package version

import "go.uber.org/zap"

type Handler struct {
}

func NewHandler(logger *zap.Logger, ver string) *Handler {
	return &Handler{}
}
