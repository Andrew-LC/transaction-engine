package domain

type Response struct {
	Status       string `json:"status,omitempty"`
	RespCode     string `json:"resp_code,omitempty"`
	Message      string `json:"message,omitempty"`
	Balance      int64  `json:"balance"`
	Transactions any    `json:"transactions,omitempty"`
}

func NewResponse(status, respCode, message string, balance int64) Response {
	return Response{
		Status:   status,
		RespCode: respCode,
		Message:  message,
		Balance:  balance,
	}
}

func NewTransactionsResponse(transactions any) Response {
	return Response{
		Status:       "SUCCESS",
		RespCode:     "00",
		Transactions: transactions,
	}
}
