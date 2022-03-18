package awsClient

// import (
// 	"context"
// 	"log"
// 	"reflect"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/service/ec2"
// )

// type ClientFactory struct {
// 	Ctx     context.Context
// 	cfg     *aws.Config
// 	clients map[reflect.Type]*interface{}
// }

// func (cf ClientFactory) getClient(typ reflect.Type) (*interface{}, error) {
// 	var err error
// 	if cf.cfg == nil {

// 		*cf.cfg, err = config.LoadDefaultConfig(cf.Ctx)
// 		log.Printf("debug: Creating new aws.Config")
// 	}

// 	val, ok := cf.clients[typ]
// 	if ok == true {
// 		return val, nil
// 	}

// 	client := reflect.ValueOf(typ)
// 	switch typ {
// 	case ec2.Client:

// 	}
// 	client.NewFromConfig(cf.cfg)

// }
