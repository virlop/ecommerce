package consume

type ConsumeMessage struct {
	RoutingKey string `json:"routing_key" example:"Remote RoutingKey to Reply"`
	Exchange   string `json:"exchange"`
}
