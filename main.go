package main

import (
	"log"
	"os"
	"path/filepath"
)

var config = conf{}
var api = Api{}

func main() {
	webserver(&config)
}

func init() {
	if wd, err := os.Getwd(); err != nil {
		log.Fatal(err)
	} else if configFileName, err = filepath.Abs(filepath.Join(wd, configFileName)); err != nil {
		log.Fatal(err)
	} else if err = config.loadConfig(configFileName); err != nil {
		log.Fatal(err)
	}

	// populate missing fields with .chive.yaml if available
	if config.ProjectId == "" || config.ApiAccessToken == "" || config.ApiHost == "" {
		var chc chiveConf
		if err := chc.loadChiveConfig(); err != nil {
			log.Fatal("missing required fields in config and cannot load .chive.yaml")
		} else if err = config.populateConfigFromChive(&chc); err != nil {
			log.Fatal("missing required fields in config and cannot populate from .chive.yaml: " + err.Error())
		}
	}

	log.Printf("Using config %+v\n", config)
}

//
// func testtt() {
// 	var conn *grpc.ClientConn
// 	conn, err := grpc.Dial(config.ApiHost, grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatalf("did not connect: %s", err)
// 	}
// 	defer conn.Close()
// 	c := ch.NewChargeHiveClient(conn)
//
// 	// ping
// 	response, err := c.Ping(context.Background(), &ch.StringTransport{Value: "dick"})
// 	if err != nil {
// 		log.Fatalf("Error when calling SayHello: %s", err)
// 	}
// 	log.Printf("Response from server: %s", response.Value)
//
// 	amt := chtype.Amount{
// 		Units:                5,
// 		Currency:             "USD",
// 		XXX_NoUnkeyedLiteral: struct{}{},
// 		XXX_unrecognized:     nil,
// 		XXX_sizecache:        0,
// 	}
// 	// connection validation
// 	cvReq := ch.ChargeCaptureRequest{ChargeId: "XX:06ff00ec-60ba-11eb-80c0-acbc32c65a39:1", Amount: &amt}
// 	cvResp, err := c.ChargeCapture(context.Background(), &cvReq)
// 	if err != nil {
// 		log.Fatalf("Error when calling SayHello: %s", err)
// 	}
// 	log.Printf("Response from server: %b", &cvResp.Acknowledged)
// }
