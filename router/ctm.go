package router

import (
	"PBAPP/common"
	"PBAPP/models"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Set the routes to be added in main
func AddCTMGroup(app *fiber.App) {
	ctmGroup := app.Group("/ctm")
	ctmGroup.Post("/", ctmRegister)
	ctmGroup.Get("/email/:email", getCTMEmail)
}

// Takes in an email and passes to utils. Returns a CTM object (defined by customer model)
func getCTMEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "email is required",
		})
	}

	CTM := common.DoesCtmExist(email)
	if CTM == "Not Found" {
		return c.Status(200).JSON(fiber.Map{"data": nil, "add": nil, "biz": nil})
	}
	ctmID := CTM.(models.CTM).Customer
	ADD := common.DoesAddExist(ctmID)
	if ADD == "Not Found" {
		return c.Status(200).JSON(fiber.Map{"data": CTM, "add": nil, "biz": nil})
	}
	BIZ := common.DoesBizExist(ctmID)
	if BIZ == "Not Found" {
		return c.Status(200).JSON(fiber.Map{"data": CTM, "add": ADD, "biz": nil})
	}

	return c.Status(200).JSON(fiber.Map{"data": CTM, "add": ADD, "biz": BIZ})
}

// Internal customer struct - doesn't initiate with address/business
type internalCtm struct {
	Email    string `json:"email" bson:"email"`
	Name     string `json:"name" bson:"name"`
	Locale   string `json:"locale" bson:"locale"`
	Customer string `json:"customer" bson:"customer"`
}

// Deprecated used for testing ctm creation - see ctmRegistration/ctmCreate
func internalCreateCtm(b *padCtm, customer string) string {
	// Use the padCTM to fill out the internalCTM struct
	coll := common.GetDBCollection("CTMs")
	customerDetails := internalCtm{Customer: customer, Name: b.Name, Email: b.Email, Locale: b.Locale}
	_, err := coll.InsertOne(context.TODO(), customerDetails)
	if err != nil {
		return "Failed to create ctm"
	}
	//err = common.CreateAddress(customer)
	return "Success"
}

// Handles the registration request data, split off and used to create the ctm add biz objects
type ctmRegistration struct {
	Email         string `json:"email,omitempty"`
	FirstLine     string `json:"street_address,omitempty"`
	City          string `json:"city,omitempty"`
	Region        string `json:"region,omitempty"`
	PostalCode    string `json:"postal_code,omitempty"`
	CountryCode   string `json:"country_code,omitempty"`
	Name          string `json:"name,omitempty"`
	CompanyNumber string `json:"company_number,omitempty"`
	TaxIdentifier string `json:"tax_identifier,omitempty"`
	Customer      string `json:"customer,omitempty"`
}

// Attempts to create the ctm, add, and biz object in that order and return the ID of each to the front-end (allows us to checkout immediately)
func ctmRegister(c *fiber.Ctx) error {
	b := new(ctmRegistration)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	ctmid := createPaddleCTM(b.Email)
	if ctmid == "x" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Could not create customer in paddle",
		})
	}

	ctmAdd := models.ADD{
		FirstLine:   b.FirstLine,
		City:        b.City,
		PostalCode:  b.PostalCode,
		Region:      b.Region,
		CountryCode: b.CountryCode,
		Description: "Created by a PB test app",
	}

	addID := createPaddleAddress(ctmid, &ctmAdd)
	if addID == "x" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Could not create address in paddle",
		})
	}

	ctmBiz := models.BIZ{
		Name:          b.Name,
		CompanyNumber: b.CompanyNumber,
		TaxIdentifier: b.TaxIdentifier,
	}

	bizID := createPaddleBusiness(ctmid, &ctmBiz)
	if bizID == "x" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Could not create business in paddle",
		})
	}

	updateCTM := models.CTM{
		Address:  addID,
		Business: bizID,
	}

	err := common.UpdateCtm(ctmid, &updateCTM)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Customer Update Failed",
		})
	}
	return c.Status(200).JSON(fiber.Map{"customerID": ctmid, "addressID": addID, "businessID": bizID})
}

// creates a CTM on paddle side and if successful pushes to the local db
func createPaddleCTM(email string) string {
	url, bearer := common.HttpHelper()
	endpoint := url + "/customers"
	client := &http.Client{}
	requestBody := ctmRegistration{
		Email: email,
	}
	sendReq, _ := json.Marshal(requestBody)
	//Make our request to paddle, create the customer in the db then respond with the ctm id
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(sendReq))
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
		customerID := dataMap["id"].(string)
		requestBody.Customer = customerID
		createCtm(&requestBody)
		return customerID
	}

	return "x"
}

// creates an address on paddle side and if successful pushes to the local db
func createPaddleAddress(ctmid string, initFields *models.ADD) string {
	url, bearer := common.HttpHelper()
	endpoint := url + "/customers/" + ctmid + "/addresses"
	client := &http.Client{}
	requestBody := initFields
	sendReq, _ := json.Marshal(requestBody)
	//Make our request to paddle, create the customer in the db then respond with the ctm id
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(sendReq))
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
		initFields.PadID = addressID
		initFields.Customer = ctmid
		common.CreateAddress(initFields)
		return addressID
	}

	return "x"
}

// creates a business on paddle side and if successful pushes to the local db
func createPaddleBusiness(ctmid string, initFields *models.BIZ) string {
	url, bearer := common.HttpHelper()
	endpoint := url + "/customers/" + ctmid + "/businesses"
	client := &http.Client{}
	requestBody := initFields
	sendReq, _ := json.Marshal(requestBody)
	//Make our request to paddle, create the customer in the db then respond with the ctm id
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(sendReq))
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
		initFields.PadID = businessID
		initFields.Customer = ctmid
		common.CreateBiz(initFields)
		return businessID
	}

	return "x"
}

func createCtm(ctm *ctmRegistration) string {
	coll := common.GetDBCollection("CTMs")
	internalCTM := models.CTM{
		Email:    ctm.Email,
		Customer: ctm.Customer,
	}
	_, err := coll.InsertOne(context.TODO(), internalCTM)
	if err != nil {
		return "Not Found"
	}
	return ctm.Customer
}
