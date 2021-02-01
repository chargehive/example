package client

import (
	"context"
	"fmt"
	"github.com/chargehive/example/chargehive"
	"github.com/chargehive/example/config"
	"github.com/chargehive/proto/golang/chargehive/chtype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

var (
	chClient *chargehive.ChargeHiveClient
	conn     *grpc.ClientConn
	conf     *config.Config
)

func Init(config *config.Config) {
	conf = config
}

// GetContext returns a context object for given chargehive details
func getCtx(ctx context.Context, c *gin.Context) context.Context {
	requestID := uuid.New().String()
	ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
		"chive-project-id":      []string{conf.ProjectId},
		"chive-access-token":    []string{conf.ApiAccessToken},
		"chive-remote-address":  []string{c.Request.RemoteAddr},
		"chive-user-agent":      []string{c.Request.UserAgent()},
		"chive-transport-agent": []string{"example site"},
		"chive-request-id":      []string{requestID},
	})
	return ctx
}

func getConn() *grpc.ClientConn {
	if conn == nil {
		var err error
		if conn, err = grpc.Dial(conf.ApiHost, grpc.WithInsecure()); err != nil {
			log.Fatalf("failed to connect to charghive api: %s", err)
		}
	}
	return conn
}

func getChClient() chargehive.ChargeHiveClient {
	if chClient == nil {
		c := chargehive.NewChargeHiveClient(getConn())
		chClient = &c
	}
	return *chClient
}

func Ping(c *gin.Context) string {
	ctx := getCtx(context.Background(), c)
	response, err := getChClient().Ping(ctx, &chargehive.StringTransport{Value: "ping"})
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%x", response.GetValue())
}

func ChargeCancel(c *gin.Context, chargeId string, reason chtype.Reason) string {
	ctx := getCtx(context.Background(), c)
	response, err := getChClient().ChargeCancel(ctx, &chargehive.ChargeCancelRequest{ChargeId: chargeId, Reason: &reason})
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("Success:%v  Result:%s", response.GetCancelSuccess(), response.GetCancelResult())
}

func ChargeCapture(c *gin.Context, chargeId, currency string, units int64) string {
	amount := chtype.Amount{Units: units, Currency: currency}
	ctx := getCtx(context.Background(), c)
	response, err := getChClient().ChargeCapture(ctx, &chargehive.ChargeCaptureRequest{ChargeId: chargeId, Amount: &amount})
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("Acknowledged:%v  ProcessId:%s", response.GetAcknowledged(), response.GetProcessId())
}

func ChargeRefund(c *gin.Context, chargeId, currency string, units int64, reason chtype.Reason, txns []*chargehive.ChargeRefundTransaction) string {
	amount := chtype.Amount{Units: units, Currency: currency}
	ctx := getCtx(context.Background(), c)
	response, err := getChClient().ChargeRefund(ctx, &chargehive.ChargeRefundRequest{
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
