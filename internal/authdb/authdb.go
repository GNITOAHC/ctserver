package authdb

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RefreshDB struct {
	c *mongo.Collection
}

func New(uri, dbName, collectionName, jwtsecret string) *RefreshDB {
	client, err := mongo.Connect(context.TODO(), options.Client(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	collection := client.Database(dbName).Collection(collectionName)
	// err = newExpireIndex(collection) // Create index for expireAt field
	return &RefreshDB{c: collection}
}

type Session struct {
	Mail      string    `bson:"mail"`
	Ip        string    `bson:"ip"`
	Location  string    `bson:"location"`
	Device    string    `bson:"device"`
	used      bool      `bson:"used"`
	sessionId string    `bson:"_id"`
	expireAt  time.Time `bson:"expireAt"`
}

// Create a session, return signed jwt refresh token with the given prefix for the session
func (r *RefreshDB) NewSession(s Session, expireAfter time.Duration, prefix, secret string) (string, error) {
	tokenId := uuid.New().String()

	_, err := r.c.InsertOne(context.TODO(), bson.M{
		"mail":     s.Mail,
		"ip":       s.Ip,
		"location": s.Location,
		"device":   s.Device,
		"used":     false,
		"_id":      tokenId,
		"expireAt": time.Now().Add(expireAfter),
	})
	if err != nil {
		return "", err
	}

	signed, err := Sign(&RefreshToken{Mail: s.Mail, SessionToken: tokenId}, secret, prefix, expireAfter)
	if err != nil {
		return "", err
	}

	return signed, nil
}

// DeleteSession deletes a given the session token
func (r *RefreshDB) DeleteSession(st string) error {
	_, err := r.c.DeleteOne(context.TODO(), bson.M{"_id": st})
	if err != nil {
		return err
	}
	return nil
}

// DeleteAllSessions deletes all sessions given the user's mail
func (r *RefreshDB) DeleteAllSessions(mail string) error {
	_, err := r.c.DeleteMany(context.TODO(), bson.M{"mail": mail})
	if err != nil {
		return err
	}
	return nil
}

// ListSessions returns all sessions given the user's mail
func (r *RefreshDB) ListSessions(mail string) ([]Session, error) {
	res, err := r.c.Find(context.Background(), bson.M{"mail": mail})
	if err != nil {
		return nil, err
	}
	var b []bson.M
	if err = res.All(context.Background(), &b); err != nil {
		return nil, err
	}
	var sessions []Session
	for _, v := range b {
		if v["used"].(bool) {
			continue
		}
		sessions = append(sessions, Session{
			Mail:      v["mail"].(string),
			Ip:        v["ip"].(string),
			Location:  v["location"].(string),
			Device:    v["device"].(string),
			sessionId: v["_id"].(string),
		})
	}
	return sessions, nil
}

// ValidSession checks if a session is valid and marks it as used
func (r *RefreshDB) ValidSession(refreshToken, secret, prefix string) (bool, error) {
	decoded, err := Decode(refreshToken, secret, prefix)
	if err != nil {
		return false, err
	}
	sessiontoken := decoded.SessionToken
	res, err := r.c.Find(context.Background(), bson.M{"_id": sessiontoken})
	if err != nil {
		return false, err
	}
	// If session exists and isn't used, mark as used and return true
	if res.Next(context.Background()) {
		var session bson.M
		if err := res.Decode(&session); err != nil {
			return false, err
		}
		if session["used"].(bool) {
			log.Print("session already used")
			return false, nil
		}
		// Mark the session as used
		_, err := r.c.UpdateOne(
			context.TODO(),
			bson.M{"_id": sessiontoken},
			bson.M{"$set": bson.M{"used": true}},
		)
		if err != nil {
			return false, err
		}
		return true, nil

	}
	return false, nil
}
