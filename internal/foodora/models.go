package foodora

import "time"

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (t AuthToken) ExpiresAt(now time.Time) time.Time {
	if t.ExpiresIn <= 0 {
		return time.Time{}
	}
	return now.Add(time.Duration(t.ExpiresIn) * time.Second)
}

type OAuthPasswordRequest struct {
	Username     string
	Password     string
	ClientSecret string
	ClientID     string
	OTPMethod    string
	OTPCode      string
	MfaToken     string
}

type OAuthRefreshRequest struct {
	RefreshToken string
	ClientSecret string
	ClientID     string
}

type ActiveOrdersResponse struct {
	Status int              `json:"status"`
	Data   ActiveOrdersData `json:"data"`
}

type ActiveOrdersData struct {
	Count         int           `json:"count"`
	ActiveOrders  []ActiveOrder `json:"active_orders"`
	PollInSeconds *int          `json:"poll_in_sec"`
}

type ActiveOrder struct {
	Code        string            `json:"code"`
	IsDelivered bool              `json:"is_delivered"`
	Vendor      ActiveOrderVendor `json:"vendor"`
	Status      StatusMessages    `json:"status_messages"`
}

type ActiveOrderVendor struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type StatusMessages struct {
	Subtitle string        `json:"subtitle"`
	Titles   []StatusTitle `json:"titles"`
}

type StatusTitle struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
	Filled bool   `json:"is_filled"`
}

type OrderStatusResponse struct {
	Status int            `json:"status"`
	Data   map[string]any `json:"data"`
}

type OrderHistoryRequest struct {
	Include        string
	Offset         int
	Limit          int
	PandaGoEnabled bool
}

type OrderHistoryResponse struct {
	Status int              `json:"status"`
	Data   OrderHistoryData `json:"data"`
}

type OrderHistoryData struct {
	TotalCount FlexibleInt        `json:"total_count"`
	Items      []OrderHistoryItem `json:"items"`
}

type OrderHistoryItem struct {
	OrderCode             string              `json:"order_code"`
	CurrentStatus         *OrderHistoryStatus `json:"current_status"`
	ConfirmedDeliveryTime *OrderHistoryTime   `json:"confirmed_delivery_time"`
	Vendor                *OrderHistoryVendor `json:"vendor"`
	TotalValue            float64             `json:"total_value"`
}

type OrderHistoryVendor struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type OrderHistoryStatus struct {
	Code               FlexibleString `json:"code"`
	Message            string         `json:"message"`
	InternalStatusCode FlexibleString `json:"internal_status_code"`
}

type OrderHistoryTime struct {
	Date     FlexibleTime `json:"date"`
	Timezone string       `json:"timezone"`
}
