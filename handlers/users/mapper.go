package users

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
)

var ErrUserNotFound = errors.New("user not found")

type Mapper struct {
	db         *mongo.Client
	collection *mongo.Collection
}

func NewMapper(db *mongo.Client) data.Mapper {
	collection := db.Database(config.AppName).Collection("users")
	return &Mapper{
		db,
		collection,
	}
}

func (m *Mapper) Insert(ctx context.Context, document any, result any, opts ...*options.InsertOneOptions) (any, error) {
	_, err := m.collection.InsertOne(ctx, document, opts...)
	return nil, err
}

func (m *Mapper) FindOne(ctx context.Context, filter any, result any, opts ...*options.FindOneOptions) (any, error) {
	collOpts := options.FindOne().SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2,
	})
	opts = append(opts, collOpts)

	err := m.collection.FindOne(ctx, filter, opts...).Decode(result)
	if err == mongo.ErrNoDocuments {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Mapper) FindOneById(ctx context.Context, id string, result any, opts ...*options.FindOneOptions) (any, error) {
	filter := bson.D{{"$or", bson.A{
		bson.D{{"id", id}},
		bson.D{{"username", id}},
	}}}
	return m.FindOne(ctx, filter, result, opts...)
}

func (m *Mapper) Find(ctx context.Context, filter any, result any, opts ...*options.FindOptions) (any, error) {
	if filter == nil {
		filter = bson.D{}
	}

	collOpts := options.Find().SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2,
	})
	opts = append(opts, collOpts)

	cur, err := m.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Mapper) Aggregate(ctx context.Context, filter any, limit int, skip int, result any, opts ...*options.AggregateOptions) (any, error) {
	// TODO implement me
	panic("implement me")
}

func (m *Mapper) Count(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
	if filter == nil {
		filter = bson.D{}
	}

	count, err := m.collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *Mapper) Update(ctx context.Context, filter any, update any, result any, opts ...*options.UpdateOptions) (any, error) {
	res := m.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if res.Err() != nil {
		return nil, res.Err()
	}

	if result != nil {
		err := res.Decode(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (m *Mapper) UpdateById(ctx context.Context, id string, document any, result any, opts ...*options.UpdateOptions) (any, error) {
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", document}}

	return m.Update(ctx, filter, update, result, opts...)
}