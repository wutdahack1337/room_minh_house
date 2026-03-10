package database

import (
	"context"
	"fmt"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) CreateInvoice(ctx context.Context, invoice *types.Invoice) error {
	// check if contract exists
	contractColl := s.database.Collection(contractCollection)
	var contract types.Contract
	if err := contractColl.FindOne(ctx, bson.M{"contract_id": invoice.ContractId}).Decode(&contract); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("database error: contract with id %s not found", invoice.ContractId)
		}
		return err
	}

	collection := s.database.Collection(invoiceCollection)
	_, err := collection.InsertOne(ctx, invoice)
	return err
}

func (s *MongoService) UpdateInvoice(ctx context.Context, invoice *types.Invoice) error {
	collection := s.database.Collection(invoiceCollection)
	res, err := collection.UpdateOne(ctx, bson.M{"contract_id": invoice.ContractId, "period": invoice.Period}, bson.M{"$set": invoice})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *MongoService) DeleteInvoice(ctx context.Context, contractID string, period string) error {
	collection := s.database.Collection(invoiceCollection)
	result, err := collection.DeleteOne(ctx, bson.M{"contract_id": contractID, "period": period})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	// also delete all associated invoice items
	itemColl := s.database.Collection(invoiceItemCollection)
	_, err = itemColl.DeleteMany(ctx, bson.M{"contract_id": contractID, "period": period})
	return err
}

func (s *MongoService) GetInvoice(ctx context.Context, contractID string, period string) (*types.Invoice, error) {
	var invoice types.Invoice
	collection := s.database.Collection(invoiceCollection)
	if err := collection.FindOne(ctx, bson.M{"contract_id": contractID, "period": period}).Decode(&invoice); err != nil {
		return nil, err
	}

	return &invoice, nil
}

func (s *MongoService) ListInvoices(ctx context.Context, contractID string) ([]*types.Invoice, error) {
	collection := s.database.Collection(invoiceCollection)
	cursor, err := collection.Find(ctx, bson.M{"contract_id": contractID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []*types.Invoice
	for cursor.Next(ctx) {
		var invoice types.Invoice
		if err := cursor.Decode(&invoice); err != nil {
			return nil, err
		}
		res = append(res, &invoice)
	}

	return res, cursor.Err()
}

func (s *MongoService) DeleteInvoicesByContractId(ctx context.Context, contractID string) error {
	invoices, err := s.ListInvoices(ctx, contractID)
	if err != nil {
		return err
	}
	for _, inv := range invoices {
		if err := s.DeleteInvoice(ctx, contractID, inv.Period); err != nil {
			return err
		}
	}
	return nil
}

func (s *MongoService) UpdateInvoiceTotalAmount(ctx context.Context, contractID string, period string) error {
	items, err := s.ListInvoiceItems(ctx, contractID, period)
	if err != nil {
		return err
	}

	var total float64
	for _, item := range items {
		total += item.Amount
	}

	collection := s.database.Collection(invoiceCollection)
	res, err := collection.UpdateOne(
		ctx,
		bson.M{"contract_id": contractID, "period": period},
		bson.M{"$set": bson.M{"total_amount": total}},
	)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
