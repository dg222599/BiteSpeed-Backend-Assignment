## BiteSpeed Backend Assignment

* This repo contains the codebase for Backend assignment task for BiteSpeed - [BiteSpeed Assignment](https://bitespeed.notion.site/Bitespeed-Backend-Task-Identity-Reconciliation-53392ab01fe149fab989422300423199)<br>
* Backend Service is written in Golang and Gin framework is used.
* Database used is PostgreSQL and web service connects and interacts with DB via [GORM](https://gorm.io/)
* Service is deployed on [Render](https://render.com/).
* `/identify` endpoint is exposed for POST requests
* **Please find the complete URL for hosted endpoint**  **[Service](https://bitespeed-backend-assignment-api.onrender.com/identify)**
*   *Please note that first few hits to the remote URL can be a bit slow to begin with since render spins down the instance after inactivity, after first few requests the service will work fine.*

#### Notes
* Two other endpoints `/getcontacts` `/remove/id` are also deployed for internal purposes and can be used as well.
* `/getcontacts` passes a `GET` request and the JSON output basically shows the state of DB at that time.
* `/remove/id` passes a `DELETE` request to delete a contact with specific ID.
* Separate URLs for `GET` and `DELETE` requests are [Get All Contacts](https://bitespeed-backend-assignment-api.onrender.com/getcontacts) [Delete Contact with ID](https://bitespeed-backend-assignment-api.onrender.com/remove/4) , user can change the ID directly in the link to delete specific user.
* ***Please note that docker-compose file is also made since the intention was to deploy via docker-compose  , but since render does not support docker-compose installation both the services are deployed separately , DB service is deployed directly as a standalone PgSQL instance which is provided by render itself***

* More changes and updates will be done in future to all the endpoints and for various testcases.
