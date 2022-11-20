# Matching Customer & Partner
The goal is to match the right partner (for ex:a craftsman) to a customer based on their project 
requirements. In this particular challenge we concentrate on flooring. The aim is to propose partners to the customers 
based on the customer's project requirements.

## Information regarding Solution
I have created the APIs that offers the following functionality:
1. Based on a customer-request, return a list of partners that offer the service. The list
should be sorted by “best match”. The quality of the match is determined first on
average rating and second by distance to the customer. (If the customer is not in the operating radius(kms) of the 
partner, the matched list won't have that partner in it. )
`http:://localhost:8080/request` <br />
2. List all the partners <br />
`http:://localhost:8080/partner`<br />
3. For a specific partner, return the detailed partner data <br />
 
`http:://localhost:8080/partner/{partnerid}`

For the challenge I have created a docker-compose.yml file.
This creates a PostgreSQL database and inserts dummy data. The database contains multiple tables, 
the schema is in the image.

![schema](https://github.com/chandanacharya1/customer-matching/blob/main/db/db.png?raw=true)

The APIs are authenticated using JWT authentication. 
So without a valid token the APIs cannot be accessed.

The posgres server uses port 55432 and the configs are stored in .env file
The web server runs on port 8080
### Having Server up and running
1. `docker-compose up`
2. `go run main.go`

### Calling APIs
#### To get a valid token 

`curl --location --request POST 'http://localhost:8080/login' --header 'Access: test'`

copy the **TOKEN** returned as a response.

#### To get list best matched partners
We have to use the token received in the above call,
replace **TOKEN** in the below curl with the token copied.

`curl --location --request POST 'http://localhost:8081/api/request' \
--header 'Content-Type: application/json' \
--header 'Token: TOKEN' \
--data-raw '{
"name": "tiles",
"latitude": 52.59,
"longitude": 13.37,
"squaremeters": 800,
"phonenumber": 0123456789
}'`

Returns a list of best matched partners.

(If there is a case when 2 matched partners have same rating, the partner with shortest distance will be 
displayed on top)

#### To get a list of all partners
We have to use the same token used in the above call,
replace **TOKEN** in the below curl with the token copied.

`curl --location --request GET 'http://localhost:8080/partner' \
--header 'Token: TOKEN' \
--header 'Content-Type: application/json' \
--data-raw ''`

Returns a list of all partners.

#### To get a details of particular partner
We have to use the same token used in the above call,
replace **TOKEN** in the below curl with the token copied.

`curl --location --request GET 'http://localhost:8081/api/partner/6' \
--header 'Token: TOKEN' \
--header 'Content-Type: application/json' \
--data-raw ''`

Returns a partner whose id is 6 
(available partnerids in the Postgres db are 1 to 6)

``
