package main

import (
	"context"
	"fmt"
	"github.com/chargehive/example/chargehive"
	"github.com/chargehive/proto/golang/chargehive/chtype"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

type Api struct {
	ctx    context.Context
	client *chargehive.ChargeHiveClient
	conn   *grpc.ClientConn
}

// GetContext returns a context object for given chargehive details
func (a *Api) getCtx() *context.Context {
	if a.ctx == nil {
		ctx := context.Background()
		a.ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
			"chive-project-id":   []string{config.ProjectId},
			"chive-access-token": []string{config.ApiAccessToken},
		})
	}
	return &a.ctx
}

func (a *Api) getConn() *grpc.ClientConn {
	if a.conn == nil {
		var err error
		if a.conn, err = grpc.Dial(config.ApiHost, grpc.WithInsecure()); err != nil {
			log.Fatalf("failed to connect to charghive api: %s", err)
		}
	}
	return a.conn
}

func (a *Api) getClient() *chargehive.ChargeHiveClient {
	if a.client == nil {
		c := chargehive.NewChargeHiveClient(a.getConn())
		a.client = &c
	}
	return a.client
}

func (a *Api) Ping(input string) string {
	response, err := (*a.getClient()).Ping(*a.getCtx(), &chargehive.StringTransport{Value: input})
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%x", response.GetValue())
}

func (a *Api) ChargeCancel(chargeId string, reason chtype.Reason) string {
	response, err := (*a.getClient()).ChargeCancel(*a.getCtx(), &chargehive.ChargeCancelRequest{ChargeId: chargeId, Reason: &reason})
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("Success:%v  Result:%s", response.GetCancelSuccess(), response.GetCancelResult())
}

func (a *Api) ChargeCapture(chargeId, currency string, units int64) string {
	amount := chtype.Amount{Units: units, Currency: currency}
	response, err := (*a.getClient()).ChargeCapture(*a.getCtx(), &chargehive.ChargeCaptureRequest{ChargeId: chargeId, Amount: &amount})
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("Acknowledged:%v  ProcessId:%s", response.GetAcknowledged(), response.GetProcessId())
}

func (a *Api) ChargeRefund(chargeId, currency string, units int64, reason chtype.Reason, txns []*chargehive.ChargeRefundTransaction) string {
	amount := chtype.Amount{Units: units, Currency: currency}
	response, err := (*a.getClient()).ChargeRefund(*a.getCtx(), &chargehive.ChargeRefundRequest{
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
