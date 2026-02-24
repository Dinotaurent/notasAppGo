package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Nota struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Titulo    string        `json:"titulo" validate:"required,min=3,max=20"`
	Contenido string        `json:"contenido" validate:"required,min=6,max=60"`
}

type NotaModel struct {
	Client *mongo.Client
}

func (m *NotaModel) Insert(n Nota) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.Client.Database("notasdb").Collection("notas").InsertOne(ctx, n)
	return err
}

//func (m *NotaModel) GetAll() ([]Nota, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	cursor, err := m.Client.Database("notasdb").Collection("notas").Find(ctx, bson.M{})
//	if err != nil {
//		return nil, err
//	}
//	defer cursor.Close(ctx)
//
//	var notas []Nota
//	err = cursor.All(ctx, &notas)
//	if err != nil {
//		return nil, err
//	}
//	return notas, nil
//}

func (m *NotaModel) GetAll(page, size int64) ([]Nota, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	collection := m.Client.Database("notasdb").Collection("notas")

	// Contar el total de documentos (Indispensable para el frontend)
	total, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	// Aplicar Skip y Limit
	opts := options.Find().
		SetSkip((page - 1) * size).
		SetLimit(size)

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var notas []Nota
	if err := cursor.All(ctx, &notas); err != nil {
		return nil, 0, err
	}

	return notas, total, nil
}

func (m *NotaModel) GetByID(id string) (Nota, error) {
	objID, _ := bson.ObjectIDFromHex(id)
	var n Nota
	err := m.Client.Database("notasdb").Collection("notas").FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&n)
	return n, err
}

func (m *NotaModel) Update(id string, n Nota) error {
	objID, _ := bson.ObjectIDFromHex(id)
	_, err := m.Client.Database("notasdb").Collection("notas").UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"titulo": n.Titulo, "contenido": n.Contenido}},
	)
	return err
}

func (m *NotaModel) Delete(id string) error {
	objID, _ := bson.ObjectIDFromHex(id)
	_, err := m.Client.Database("notasdb").Collection("notas").DeleteOne(context.TODO(), bson.M{"_id": objID})
	return err
}
