package common

import (
	"PBAPP/models"
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

// Sets up our base URL and Bearer token for easy access throughout all the files
func HttpHelper() (string, string) {
	baseUrl := "https://sandbox-api.paddle.com"
	bearer := "Bearer " + os.Getenv("bearer")
	return baseUrl, bearer
}

// Search for a customer record based on email and return the CTM record as an interface
func DoesCtmExist(email string) interface{} {
	coll := GetDBCollection("CTMs")
	CTM := models.CTM{}
	err := coll.FindOne(context.TODO(), bson.M{"email": email}).Decode(&CTM)
	if err != nil {
		return "Not Found"
	}
	return CTM
}

// Updates the ctm record in the DB (works with ctmRegistration to initalize the ctm with address/business objects after creation)
func UpdateCtm(ctmid string, updateFields *models.CTM) error {
	coll := GetDBCollection("CTMs")
	filter := bson.M{"customer": ctmid}
	update := bson.M{"$set": updateFields}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// Called in main.go to initialize the .env file
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}
