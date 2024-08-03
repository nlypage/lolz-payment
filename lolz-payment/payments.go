package lolzpayment

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

type PaymentsHistoryRequest struct {
	//UserID, just the user id, in official docks it says that the field is required, but it's not XD.
	UserID int
	// Type field is optional. Type of operation.
	Type string
	// Pmin field is optional. Minimal price of account (Inclusive).
	Pmin int
	// Pmax field is optional. Maximum price of account (Inclusive).
	Pmax int
	// Page field is optional. The number of the page to display results from.
	Page int
	// OperationID field is optional. ID of the operation from which the result begins.
	OperationID int
	// Receiver field is optional. Username of user, which receive money from you.
	Receiver string
	// Sender field is optional. Username of user, which sent money to you.
	Sender string
	// StartDate field is optional. Start date of operation.
	StartDate time.Time
	// EndDate field is optional. End date of operation.
	EndDate time.Time
	// Wallet field is optional. Wallet, which used for money payouts.
	Wallet string
	// Comment field is optional. Comment for money transfers.
	Comment string
	/* IsHold field is optional.
	With IsHold = true - api will return only payments with hold,
	With IsHold = false - api will return only payments without hold
	*/
	IsHold *bool
	// ShowPaymentsStats field is optional. Display payment stats for selected period (outgoing value, incoming value).
	ShowPaymentsStats *bool
}

type PaymentData struct {
	UserID              int    `json:"user_id,omitempty"`
	Username            string `json:"username,omitempty"`
	Comment             string `json:"comment,omitempty"`
	IsBanned            int    `json:"is_banned,omitempty"`
	DisplayStyleGroupId int    `json:"display_style_group_id,omitempty"`
	UniqUsernameCss     string `json:"uniq_username_css,omitempty"`
	AvatarDate          int    `json:"avatar_date,omitempty"`
	UserGroupId         int    `json:"user_group_id,omitempty"`
}

type UserBalance struct {
	UserID              int `json:"user_id,omitempty"`
	UserBalance         int `json:"user_balance,omitempty"`
	UserHold            int `json:"user_hold,omitempty"`
	UserBalanceWithHold int `json:"user_balance_with_hold,omitempty"`
}

type Payment struct {
	OperationID   int         `json:"operation_id,omitempty"`
	OperationDate int64       `json:"operation_date,omitempty"`
	OperationType string      `json:"operation_type,omitempty"`
	OutgoingSum   int         `json:"outgoing_sum,omitempty"`
	IncomingSum   int         `json:"incoming_sum,omitempty"`
	ItemID        int         `json:"item_id,omitempty"`
	Wallet        string      `json:"wallet,omitempty"`
	IsFinished    int         `json:"is_finished,omitempty"`
	IsHold        int         `json:"is_hold,omitempty"`
	PaymentSystem string      `json:"payment_system,omitempty"`
	Data          PaymentData `json:"data,omitempty"` // The list can be returned in this key, I've given up on it. Read the note on the PaymentsHistory function
	HoldEndDate   int         `json:"hold_end_date,omitempty"`
	Api           int         `json:"api,omitempty"`
	PaymentStatus string      `json:"payment_status,omitempty"`
	User          UserBalance `json:"user,omitempty"`
}

// Payments map key - is payment ID.
type Payments map[string]Payment

type paymentsHistoryResponse struct {
	Payments Payments `json:"payments,omitempty"` // The list can be returned in this key, I've given up on it. Read the note on the PaymentsHistory function
}

// PaymentsHistory is a function to get user payments history using user/:userID/payments endpoint.
/*
	Note to the function, due to the fact that php, if the structure is empty, returns an empty list, errors may occur when json unmarshall in payments and data
	I chose to just ignore them :D
	If you have a desire, you can do something about it <3.
*/
func (c *Client) PaymentsHistory(ctx context.Context, historyRequest PaymentsHistoryRequest) (Payments, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: fmt.Sprintf("user/%d/payments", historyRequest.UserID),
	}
	if historyRequest.Type != "" {
		r.setParam("type", historyRequest.Type)
	}
	if historyRequest.Pmin != 0 {
		r.setParam("pmin", historyRequest.Pmin)
	}
	if historyRequest.Pmax != 0 {
		r.setParam("pmax", historyRequest.Pmax)
	}
	if historyRequest.Page != 0 {
		r.setParam("page", historyRequest.Page)
	}
	if historyRequest.OperationID != 0 {
		r.setParam("operation_id_lt", historyRequest.OperationID)
	}
	if historyRequest.Receiver != "" {
		r.setParam("receiver", historyRequest.Receiver)
	}
	if historyRequest.Sender != "" {
		r.setParam("sender", historyRequest.Sender)
	}
	if !historyRequest.StartDate.IsZero() {
		r.setParam("startDate", historyRequest.StartDate.Format(time.RFC3339))
	}
	if !historyRequest.EndDate.IsZero() {
		r.setParam("endDate", historyRequest.EndDate.Format(time.RFC3339))
	}
	if historyRequest.Wallet != "" {
		r.setParam("wallet", historyRequest.Wallet)
	}
	if historyRequest.Comment != "" {
		r.setParam("comment", historyRequest.Comment)
	}
	if historyRequest.IsHold != nil {
		r.setParam("is_hold", *historyRequest.IsHold)
	}
	if historyRequest.ShowPaymentsStats != nil {
		r.setParam("show_payment_stats", historyRequest.ShowPaymentsStats)
	}

	data, err := c.do(ctx, r)
	if err != nil {
		return nil, err
	}

	var res paymentsHistoryResponse
	if err = json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Payments, nil
}

// CreatePaymentLink returns a link to transfer funds to the account whose token you specified.
func (c *Client) CreatePaymentLink(amount float64, comment string, redirectURL string) string {
	return fmt.Sprintf("https://lzt.market/balance/transfer?username=%s&hold=0&amount=%f&comment=%s&redirect=%s", c.username, math.Ceil(amount), comment, redirectURL)
}

type PaymentsHandlerOptions struct {
	// Type field is optional. Default is "receiving_money"
	Type string
	// Pmin field is optional. Minimal price of account (Inclusive).
	Pmin int
	// Pmax field is optional. Maximum price of account (Inclusive).
	Pmax int
	// Sender field is optional. Username of user, which sent money to you.
	Sender string
	// Comment field is optional. Comment for money transfers.
	Comment string
	/* IsHold field is optional.
	With IsHold = true - api will return only payments with hold,
	With IsHold = false - api will return only payments without hold
	*/
	IsHold *bool
	// Period is optional. Default is 1s. The period for checking the availability of new payments.
	Period time.Duration
}

var DefaultPaymentsHandlerOptions *PaymentsHandlerOptions = &PaymentsHandlerOptions{
	Type:   "receiving_money",
	Period: time.Second,
}

type HandlerFunc func(payment Payment) error

// HandlePayments will send new payments to your handlerFunction.
/*
	it starts periodic verification of new payments in goroutine,
	if your code ends after calling this one, you need to set WaitGroup.
*/
func (c *Client) HandlePayments(handlerFunc HandlerFunc, options *PaymentsHandlerOptions) {
	go func() {
		var (
			lastPaymentDate int64 = time.Now().Unix()
		)

		if options == nil {
			options = DefaultPaymentsHandlerOptions
		}

		if options.Period.Seconds() == 0 {
			options.Period = time.Second
		}

		if options.Type == "" {
			options.Type = "receiving_money"
		}

		for {
			newPayments, _ := c.PaymentsHistory(context.Background(), PaymentsHistoryRequest{
				UserID:    c.userID,
				Type:      options.Type,
				StartDate: time.Unix(lastPaymentDate, 0).Add(time.Second),
				EndDate:   time.Now().Add(time.Hour * 24),
				Pmin:      options.Pmin,
				Pmax:      options.Pmax,
				Sender:    options.Sender,
				Comment:   options.Comment,
				IsHold:    options.IsHold,
			})

			for _, payment := range newPayments {
				go func() {
					errHandle := handlerFunc(payment)
					if errHandle != nil {
						log.Println(fmt.Errorf("got error while handling payment %d: %w", payment.OperationID, errHandle))
					}
				}()
				lastPaymentDate = max(lastPaymentDate, payment.OperationDate)
			}

			time.Sleep(options.Period)
		}
	}()
}
