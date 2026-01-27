package device

type Handler struct {
	deviceService *Service
}

func NewHandler(ds *Service) *Handler {
	return &Handler{
		deviceService: ds,
	}
}
