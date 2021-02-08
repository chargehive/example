package client

import (
	"context"
	"encoding/json"
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

func Ping(c *gin.Context, value string) string {
	ctx := getCtx(context.Background(), c)
	response, err := getChClient().Ping(ctx, &chargehive.StringTransport{Value: value})
	if err != nil {
		log.Println(err)
	}
	jResponse, _ := json.Marshal(response)
	return string(jResponse)
}

func ChargeCreate(c *gin.Context, currency string, units int64, paymentMethodId string) string {
	ctx := getCtx(context.Background(), c)
	defaultAddress := chtype.Address{
		LineOne:     "123 House",
		LineTwo:     "Cartel Street",
		Town:        "Chiswick",
		County:      "London",
		Country:     "UK",
		PostalCode:  "SW1 1RP",
		Fao:         "Bob Hogan",
		CompanyName: "Hogan Ltd",
	}
	person := chtype.Person{
		Title:       "Mr",
		FirstName:   "Hulk",
		LastName:    "Hogan",
		FullName:    "Hulk Hogan",
		Email:       "hulk@hogan.com",
		PhoneNumber: "+1800-hulk-hogan",
		Language:    "EN",
	}
	item := chtype.ChargeItem{
		SubscriptionId: "",
		RenewalNumber:  0,
		Duration:       0,
		StartDate:      nil,
		EndDate:        nil,
		ProductType:    0,
		SkuType:        0,
		Delivery:       nil,
		Quantity:       0,
		UnitPrice:      nil,
		TaxAmount:      nil,
		DiscountAmount: nil,
		Name:           "",
		Description:    "",
		ProductCode:    "",
		SkuCode:        "",
		TermUnits:      1,
		TermType:       chtype.TERM_TYPE_MONTH,
	}
	amt := chtype.Amount{Units: units, Currency: currency}
	meta := chtype.ChargeMeta{
		BillingAddress:  &defaultAddress,
		DeliveryAddress: &defaultAddress,
		Items:           []*chtype.ChargeItem{&item},
		Terms:           "",
		Note:            "",
		MerchantMemo:    "",
		InvoiceDate:     nil,
		DueDate:         nil,
		DiscountAmount:  nil,
		DeliveryAmount:  nil,
		TaxAmount:       nil,
		TotalAmount:     &amt,
		Person:          &person,
		Company:         nil,
		IpAddress:       c.Request.RemoteAddr,
		Delivery:        nil,
		Device:          nil,
		CustomerId:      "",
		PlacementId:     conf.PlacementToken,
	}
	req := chargehive.ChargeCreateRequest{
		MerchantReference: "RANDOM_REF!",
		Amount:            &amt,
		PaymentMethodIds:  []string{paymentMethodId},
		ExpiryTime:        nil,
		ContractType:      chtype.CONTRACT_TYPE_SUBSCRIPTION_INITIAL,
		Environment:       chtype.CHARGE_ENVIRONMENT_ECOMMERCE,
		ChargeMeta:        &meta,
	}
	response, err := getChClient().ChargeCreate(ctx, &req)
	if err != nil {
		log.Println(err)
	}
	jResponse, _ := json.Marshal(response)
	return string(jResponse)
}

func MethodUpdate(c *gin.Context, token string, schema chtype.PaymentMethodSchema, jsonDetails []byte, methodType chtype.PaymentMethodType, methodProvider chtype.PaymentMethodProvider) string {
	ctx := getCtx(context.Background(), c)
	method := chtype.PaymentMethod{
		Schema:    schema,
		Json:      jsonDetails,
		Type:      methodType,
		Provider:  methodProvider,
		InputType: chtype.INPUT_TYPE_VIRTUAL,
	}
	response, err := getChClient().MethodUpdate(ctx, &chargehive.MethodUpdateRequest{
		Token:                token,
		PaymentMethodUpdates: &method,
	})
	if err != nil {
		log.Println(err)
	}
	jResponse, _ := json.Marshal(response)
	return string(jResponse)
}

func ChargeCancel(c *gin.Context, chargeId string, reason chtype.Reason) string {
	ctx := getCtx(context.Background(), c)
	response, err := getChClient().ChargeCancel(ctx, &chargehive.ChargeCancelRequest{ChargeId: chargeId, Reason: &reason})
	if err != nil {
		log.Println(err)
	}
	jResponse, _ := json.Marshal(response)
	return string(jResponse)
}

func ChargeCapture(c *gin.Context, chargeId, currency string, units int64) string {
	amount := chtype.Amount{Units: units, Currency: currency}
	ctx := getCtx(context.Background(), c)
	response, err := getChClient().ChargeCapture(ctx, &chargehive.ChargeCaptureRequest{ChargeId: chargeId, Amount: &amount})
	if err != nil {
		log.Println(err)
	}
	jResponse, _ := json.Marshal(response)
	return string(jResponse)
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
	jResponse, _ := json.Marshal(response)
	return string(jResponse)
}
