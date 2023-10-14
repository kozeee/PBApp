package router

import (
	"PBAPP/common"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

func AddPaddleGroup(app *fiber.App) {
	ctmGroup := app.Group("/paddle")
	ctmGroup.Post("/", padCreateCtm)
	ctmGroup.Get("/prices", padGetPrices)
	ctmGroup.Get("/subscriptions/:id", padGetSubs)
	ctmGroup.Post("/subscriptions/cancel/:id", padCancelSub)
	ctmGroup.Post("/subscriptions/update/:id", padUpdateSub)
	ctmGroup.Post("/subscriptions/update/pause/:id", padPauseSub)
	ctmGroup.Post("/subscriptions/update/unpause/:id", padUnpauseSub)
	ctmGroup.Post("/subscriptions/transactions/:id", padFetchTransactions)
	ctmGroup.Post("/subscriptions/change/:id", padChangePlans)
	ctmGroup.Post("/subscriptions/change/preview/:id", padPreviewPlans)

}

// Used to make the request to paddle
type padCtm struct {
	Email  string `json:"email" bson:"email"`
	Name   string `json:"name" bson:"name"`
	Locale string `json:"locale" bson:"locale"`
}

// Create a Customer Record in Paddle + Initialize in internal DB
func padCreateCtm(c *fiber.Ctx) error {

	//Get our endpoint and create an http client
	url, bearer := common.HttpHelper()
	endpoint := url + "/customers"
	client := &http.Client{}

	//Convert the data to padCTM and marshall to json
	b := new(padCtm)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}
	sendReq, _ := json.Marshal(b)

	//Check if the email exists in the db
	exists := (common.DoesCtmExist(b.Email))
	if exists != "Not Found" {
		return c.Status(400).JSON(fiber.Map{
			"Customer Exists": exists,
		})
	}

	//Make our request to paddle, create the customer in the db then respond with the ctm id
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(sendReq))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create CTM",
			"message": err.Error(),
		})
	}
	req.Header.Set("Authorization", bearer)
	res, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create CTM",
			"message": err.Error(),
		})
	}
	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)
	var results map[string]interface{}
	json.Unmarshal([]byte(response), &results)
	customer := results["data"].(map[string]interface{})
	customerID := customer["id"].(string)
	internalCreateCtm(b, customerID)

	return c.Status(200).JSON(fiber.Map{"data": customer["id"]})
}

// Get all subscriptions based on a CTM id
func padGetSubs(c *fiber.Ctx) error {

	//Get our endpoint and create an http client
	url, bearer := common.HttpHelper()
	client := &http.Client{}

	//Convert the data to padCTM and marshall to json
	id := c.Params("id")
	endpoint := url + "/subscriptions?customer_id=" + id

	// Make our request to paddle, create the customer in the db then respond with the ctm id
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", bearer)
	res, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create CTM",
			"message": err.Error(),
		})
	}
	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)
	var results map[string]interface{}
	json.Unmarshal([]byte(response), &results)

	return c.Status(200).JSON(fiber.Map{"data": results["data"]})
}

// Takes in a req from the front-end and cancels on the paddle side
type cancelMin struct {
	EffectiveFrom string `json:"effective_from,omitempty"`
}

func padCancelSub(c *fiber.Ctx) error {
	url, bearer := common.HttpHelper()
	client := &http.Client{}
	id := c.Params("id")
	endpoint := url + "/subscriptions/" + id + "/cancel"

	requestBody := cancelMin{EffectiveFrom: "next_billing_period"}
	sendReq, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(sendReq))
	req.Header.Set("Authorization", bearer)
	res, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to cancel sub",
			"message": err.Error(),
		})
	}
	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)
	var results map[string]interface{}
	json.Unmarshal([]byte(response), &results)

	return c.Status(200).JSON(fiber.Map{"data": results["data"]})
}

// Return price data - product set by an env variable currently. Could be modified to fetch more dynamic data
func padGetPrices(c *fiber.Ctx) error {
	url, bearer := common.HttpHelper()
	endpoint := url + "/prices?product_id=" + os.Getenv("testProduct")
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatal("Failed to create request:", err)
	}
	req.Header.Set("Authorization", bearer)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("GET request failed:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var response map[string]json.RawMessage
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return c.Status(200).JSON(fiber.Map{"data": response["data"]})
}

// Creates a $0 transaction to let buyers update their sub
func padUpdateSub(c *fiber.Ctx) error {

	//Get our endpoint and create an http client
	url, bearer := common.HttpHelper()
	client := &http.Client{}

	//fe will pass the subscription id back to us to make this request
	id := c.Params("id")
	endpoint := url + "/subscriptions/" + id + "/transactions"

	// Make our request to paddle
	req, _ := http.NewRequest("POST", endpoint, nil)
	req.Header.Set("Authorization", bearer)
	res, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create Transaction",
			"message": err.Error(),
		})
	}
	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)
	fmt.Println(string(response))
	var results map[string]interface{}
	json.Unmarshal([]byte(response), &results)

	return c.Status(200).JSON(fiber.Map{"data": results["data"]})
}

func padFetchTransactions(c *fiber.Ctx) error {
	//Get our endpoint and create an http client
	url, bearer := common.HttpHelper()
	client := &http.Client{}

	//Convert the data to padCTM and marshall to json
	id := c.Params("id")
	endpoint := url + "/transactions?status=completed&subscription_id=" + id

	// Make our request to paddle, create the customer in the db then respond with the ctm id
	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", bearer)
	res, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Could not retrieve transaction data",
			"message": err.Error(),
		})
	}
	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)
	var results map[string]interface{}
	json.Unmarshal([]byte(response), &results)

	dataSlice, dataOk := results["data"].([]interface{})
	if dataOk {
		// Iterate through the "data" slice
		for _, item := range dataSlice {
			// Each item in the slice is a map[string]interface{}
			itemMap, mapOk := item.(map[string]interface{})
			if mapOk {
				// Access the "ID" field from the item map
				id, idOk := itemMap["id"].(string)
				if idOk {
					// You have successfully accessed the "ID" field
					endpoint = url + "/transactions/" + id + "/invoice"
					fmt.Println(endpoint)
					req, _ := http.NewRequest("GET", endpoint, nil)
					req.Header.Set("Authorization", bearer)
					res, err := client.Do(req)
					if err != nil {
						return c.Status(500).JSON(fiber.Map{
							"error":   "Could not retrieve transaction data",
							"message": err.Error(),
						})
					}
					defer res.Body.Close()
					response, _ := io.ReadAll(res.Body)
					var results map[string]interface{}
					json.Unmarshal([]byte(response), &results)
					return c.Status(200).JSON(fiber.Map{"data": results})

				} else {
					fmt.Println(endpoint)
				}
			} else {
				// Handle the case where an item is not a map
			}
		}
	} else {
		// Handle the case where "data" is not a slice of interface{}
	}
	return c.Status(400).JSON(fiber.Map{"data": "Failed to return an ID"})

}

type padSubPlan struct {
	Items []struct {
		PriceID  string `json:"price_id"`
		Quantity int    `json:"quantity"`
	} `json:"items"`
	ProrationBillingMode string `json:"proration_billing_mode"`
}

func padPreviewPlans(c *fiber.Ctx) error {
	b := new(padSubPlan)
	id := c.Params("id")
	url, bearer := common.HttpHelper()
	endpoint := url + "/subscriptions/" + id + "/preview"
	client := &http.Client{}

	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}
	sendReq, _ := json.Marshal(b)

	req, err := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(sendReq))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create CTM",
			"message": err.Error(),
		})
	}
	req.Header.Set("Authorization", bearer)
	res, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create CTM",
			"message": err.Error(),
		})
	}
	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)
	var results map[string]interface{}
	json.Unmarshal([]byte(response), &results)
	return c.Status(200).JSON(fiber.Map{"data": results["data"]})
}

func padChangePlans(c *fiber.Ctx) error {
	b := new(padSubPlan)
	id := c.Params("id")
	url, bearer := common.HttpHelper()
	endpoint := url + "/subscriptions/" + id
	client := &http.Client{}

	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}
	sendReq, _ := json.Marshal(b)

	req, err := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(sendReq))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create CTM",
			"message": err.Error(),
		})
	}
	req.Header.Set("Authorization", bearer)
	res, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create CTM",
			"message": err.Error(),
		})
	}
	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)
	var results map[string]interface{}
	json.Unmarshal([]byte(response), &results)
	return c.Status(200).JSON(fiber.Map{"data": results["data"]})
}

type padPauseController struct {
	EffectiveFrom string `json:"effective_from"`
}

func padPauseSub(c *fiber.Ctx) error {
	b := new(padPauseController)
	id := c.Params("id")
	url, bearer := common.HttpHelper()
	endpoint := url + "/subscriptions/" + id + "/pause"
	client := &http.Client{}

	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}
	sendReq, _ := json.Marshal(b)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(sendReq))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create CTM",
			"message": err.Error(),
		})
	}
	req.Header.Set("Authorization", bearer)
	res, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create CTM",
			"message": err.Error(),
		})
	}
	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)
	var results map[string]interface{}
	json.Unmarshal([]byte(response), &results)
	return c.Status(200).JSON(fiber.Map{"data": results["data"]})
}

func padUnpauseSub(c *fiber.Ctx) error {
	b := new(padPauseController)
	id := c.Params("id")
	url, bearer := common.HttpHelper()
	endpoint := url + "/subscriptions/" + id + "/resume"
	client := &http.Client{}

	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err,
		})
	}
	sendReq, _ := json.Marshal(b)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(sendReq))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create CTM",
			"message": err.Error(),
		})
	}
	req.Header.Set("Authorization", bearer)
	res, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create CTM",
			"message": err.Error(),
		})
	}
	defer res.Body.Close()
	response, _ := io.ReadAll(res.Body)
	var results map[string]interface{}
	json.Unmarshal([]byte(response), &results)
	return c.Status(200).JSON(fiber.Map{"data": results["data"]})
}
