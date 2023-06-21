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

func DoesBizExist(ctmid string) interface{} {
	coll := GetDBCollection("BIZs")
	BIZ := models.BIZ{}
	err := coll.FindOne(context.TODO(), bson.M{"customer": ctmid}).Decode(&BIZ)
	if err != nil {
		return "Not Found"
	}

	if BIZ.PadID != "" {
		url, bearer := HttpHelper()
		endpoint := url + "/customers/" + ctmid + "/businesses/" + BIZ.PadID
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
			padBiz := models.BIZ{
				Name:          dataMap["name"].(string),
				CompanyNumber: dataMap["company_number"].(string),
				TaxIdentifier: dataMap["tax_identifier"].(string),
			}
			UpdateBiz(ctmid, &padBiz)
		}
		err = coll.FindOne(context.TODO(), bson.M{"customer": ctmid}).Decode(&BIZ)
		return BIZ
	}

	return BIZ
}

func UpdateBiz(ctmid string, updateFields *models.BIZ) error {
	coll := GetDBCollection("BIZs")
	filter := bson.M{"customer": ctmid}

	update := bson.M{"$set": updateFields}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func PadUpdateBiz(ctmid string, initFields *models.BIZ) string {
	url, bearer := HttpHelper()
	endpoint := url + "/customers/" + ctmid + "/businesses/" + initFields.PadID
	client := &http.Client{}
	requestBody := models.BIZ{Name: initFields.Name, CompanyNumber: initFields.CompanyNumber, TaxIdentifier: initFields.TaxIdentifier}
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
		businessID := dataMap["id"].(string)
		UpdateBiz(ctmid, initFields)
		return businessID
	}

	return "x"
}

func CreateBiz(initFields *models.BIZ) error {

	coll := GetDBCollection("BIZs")
	_, err := coll.InsertOne(context.TODO(), initFields)
	if err != nil {
		return err
	}
	return nil
}
