package helper

import (
	"errors"
	"reflect"
	"test_backend_esb/model"
)

var TypeRegistry = map[string]reflect.Type{
	"user":               reflect.TypeOf(model.User{}),
	"item":               reflect.TypeOf(model.Item{}),
	"category":           reflect.TypeOf(model.Category{}),
	"store":              reflect.TypeOf(model.Store{}),
	"transaction":        reflect.TypeOf(model.Transaction{}),
	"payment_method":     reflect.TypeOf(model.PaymentMethod{}),
	"transaction_detail": reflect.TypeOf(model.TransactionDetail{}),
	"order_method":       reflect.TypeOf(model.OrderMethod{}),
	"batch_store":        reflect.TypeOf([]model.Store{}),
	"batch_item":         reflect.TypeOf([]model.Item{}),
}

var RelationshipRegistry = map[reflect.Type][]string{
	reflect.TypeOf(model.Transaction{}): {"Order.OrderDetail.Item.Category", "Terminal.Store.Merchant"},
	reflect.TypeOf(model.Category{}):    {"Store"},
	// add other models and their relationships here
}

func CreateGORMModel(modelName string) (model.GORMModel, error) {
	modelType, ok := TypeRegistry[modelName]
	if !ok {
		return nil, errors.New("gorm_model: model not found")
	}
	element := reflect.New(modelType).Interface().(model.GORMModel)
	return element, nil
}
