# Oct 2023
This project will no longer be updated - if you have any questions about the project feel free to reach out.

# Paddle Billing Demo
Demo includes the ability to create customers, subscribers, and businesses
Subscription creation through checkout, and cancellation through management API
Allows users to view subscription data and update plans (including line-item changes)
Syncs data with Paddle before changing db records - Strictly enforces parity with Paddle (assumes paddle data is always most up-to-date)


## Env
.env file is expected in the root PBApp Directory

MONGO_URI= Mongo connection URI
PORT= Desired port for routing 
LocalURL= Mostly used to pass to the front-end (most are hard-coded in the .tmpl pages at the moment)
testProduct = pro_ product ID from Paddle - just to serve some example product/prices. Future improvements should allow the prices being served to be handled by the db
bearer = Paddle bearer auth token for Paddle API requests
