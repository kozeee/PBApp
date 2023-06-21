package common

import (
	"PBAPP/models"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// Checks if the Address exists in our DB, if it does we sync with paddle and return the most recent record
func DoesAddExist(ctmid string) interface{} {
	coll := GetDBCollection("ADDs")
	ADD := models.ADD{}
	err := coll.FindOne(context.TODO(), bson.M{"customer": ctmid}).Decode(&ADD)
	if err != nil {
		return "Not Found"
	}

	if ADD.PadID != "" {
		url, bearer := HttpHelper()
		endpoint := url + "/customers/" + ADD.Customer + "/addresses/" + ADD.PadID
		client := &http.Client{}
		req, err := http.NewRequest("GET", endpoint, nil)

		req.Header.Set("Authorization", bearer)
		res, err := client.Do(req)
		if err != nil {
			return "Pad Get Failed"
		}
		defer res.Body.Close()
		response, _ := io.ReadAll(res.Body)
		var results map[string]interface{}
		json.Unmarshal([]byte(response), &results)

		dataMap, ok := results["data"].(map[string]interface{})
		if ok {
			padAdd := models.ADD{
				Status:      dataMap["status"].(string),
				Description: dataMap["description"].(string),
				FirstLine:   dataMap["first_line"].(string),
				City:        dataMap["city"].(string),
				PostalCode:  dataMap["postal_code"].(string),
				Region:      dataMap["region"].(string),
				CountryCode: dataMap["country_code"].(string),
				PadID:       dataMap["id"].(string),
			}
			UpdateAddress(ctmid, &padAdd)
		}
		err = coll.FindOne(context.TODO(), bson.M{"customer": ctmid}).Decode(&ADD)
		return ADD
	}
	return ADD
}

// Updates our internal DB, mostly used to sync from Paddle, but could be used to manually update a record
func UpdateAddress(ctmid string, updateFields *models.ADD) error {
	coll := GetDBCollection("ADDs")
	filter := bson.M{"customer": ctmid}
	update := bson.M{"$set": updateFields}
	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func PadUpdateAddress(ctmid string, initFields *models.ADD) string {
	url, bearer := HttpHelper()
	endpoint := url + "/customers/" + ctmid + "/addresses/" + initFields.PadID
	client := &http.Client{}
	requestBody := models.ADD{CountryCode: initFields.CountryCode, City: initFields.City, FirstLine: initFields.FirstLine, PostalCode: initFields.PostalCode, Region: initFields.Region}
	sendReq, _ := json.Marshal(requestBody)
	//Make our request to paddle, create the customer in the db then respond with the ctm id
	req, err := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(sendReq))
	if err != nil {
		return "x"
	}
	req.Header.Set("Authorization", bearer)
	res, err := client.Do(req)
	if err != nil {
		return "x"
	}
	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)
	var results map[string]interface{}
	json.Unmarshal([]byte(response), &results)
	dataMap, ok := results["data"].(map[string]interface{})
	if ok {
		addressID := dataMap["id"].(string)
		UpdateAddress(ctmid, initFields)
		return addressID
	}

	return "x"
}

// Creates a record in paddle then syncs our DB using UpdateAddress
func CreateAddress(initFields *models.ADD) error {
	//Get our endpoint and create an http client
	coll := GetDBCollection("ADDs")
	_, err := coll.InsertOne(context.TODO(), initFields)
	if err != nil {
		return err
	}
	return nil
}
