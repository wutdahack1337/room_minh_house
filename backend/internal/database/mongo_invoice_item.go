package database

import (
	"context"
	"fmt"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) CreateInvoiceItem(ctx context.Context, item *types.InvoiceItem) error {
	// check if invoice exists
	invoiceColl := s.database.Collection(invoiceCollection)
	var invoice types.Invoice
	if err := invoiceColl.FindOne(ctx, bson.M{"contract_id": item.ContractId, "period": item.Period}).Decode(&invoice); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("database error: invoice with contract_id %s and period %s not found", item.ContractId, item.Period)
		}
		return err
	}

	collection := s.database.Collection(invoiceItemCollection)
	_, err := collection.InsertOne(ctx, item)
	return err
}

func (s *MongoService) UpdateInvoiceItem(ctx context.Context, item *types.InvoiceItem) error {
	collection := s.database.Collection(invoiceItemCollection)
	res, err := collection.UpdateOne(
		ctx,
		bson.M{"contract_id": item.ContractId, "period": item.Period, "item_type": item.ItemType},
		bson.M{"$set": item},
	)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *MongoService) ListInvoiceItems(ctx context.Context, contractID string, period string) ([]*types.InvoiceItem, error) {
	collection := s.database.Collection(invoiceItemCollection)
	cursor, err := collection.Find(ctx, bson.M{"contract_id": contractID, "period": period})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []*types.InvoiceItem
	for cursor.Next(ctx) {
		var item types.InvoiceItem
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		res = append(res, &item)
	}

	return res, cursor.Err()
}
