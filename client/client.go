package client

import (
	"context"
	"fmt"
	"github.com/chargehive/example/chargehive"
	"github.com/chargehive/example/config"
	"github.com/chargehive/proto/golang/chargehive/chtype"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

var client Client

func Init(config *config.Config) {
	client.conf = config
}

func Get() *Client {
	return &client
}

type Client struct {
	ctx  context.Context
	chc  *chargehive.ChargeHiveClient
	conn *grpc.ClientConn
	conf *config.Config
}

// GetContext returns a context object for given chargehive details
func (a *Client) getCtx() *context.Context {
	if a.ctx == nil {
		ctx := context.Background()
		a.ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
			"chive-project-id":   []string{a.conf.ProjectId},
			"chive-access-token": []string{a.conf.ApiAccessToken},
		})
	}
	return &a.ctx
}

func (a *Client) getConn() *grpc.ClientConn {
	if a.conn == nil {
		var err error
		if a.conn, err = grpc.Dial(a.conf.ApiHost, grpc.WithInsecure()); err != nil {
			log.Fatalf("failed to connect to charghive api: %s", err)
		}
	}
	return a.conn
}

func (a *Client) getChc() *chargehive.ChargeHiveClient {
	if a.chc == nil {
		c := chargehive.NewChargeHiveClient(a.getConn())
		a.chc = &c
	}
	return a.chc
}

func (a *Client) Ping(input string) string {
	response, err := (*a.getChc()).Ping(*a.getCtx(), &chargehive.StringTransport{Value: input})
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%x", response.GetValue())
}

func (a *Client) ChargeCancel(chargeId string, reason chtype.Reason) string {
	response, err := (*a.getChc()).ChargeCancel(*a.getCtx(), &chargehive.ChargeCancelRequest{ChargeId: chargeId, Reason: &reason})
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("Success:%v  Result:%s", response.GetCancelSuccess(), response.GetCancelResult())
}

func (a *Client) ChargeCapture(chargeId, currency string, units int64) string {
	amount := chtype.Amount{Units: units, Currency: currency}
	response, err := (*a.getChc()).ChargeCapture(*a.getCtx(), &chargehive.ChargeCaptureRequest{ChargeId: chargeId, Amount: &amount})
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("Acknowledged:%v  ProcessId:%s", response.GetAcknowledged(), response.GetProcessId())
}

func (a *Client) ChargeRefund(chargeId, currency string, units int64, reason chtype.Reason, txns []*chargehive.ChargeRefundTransaction) string {
	amount := chtype.Amount{Units: units, Currency: currency}
	response, err := (*a.getChc()).ChargeRefund(*a.getCtx(), &chargehive.ChargeRefundRequest{
		ChargeId:     chargeId,
		Amount:       &amount,
		Reason:       &reason,
		Transactions: txns,
	})
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("Acknowledged:%v  ProcessId:%s", response.GetAcknowledged(), response.GetProcessId())
}
