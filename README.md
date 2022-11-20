# Matching Customer & Partner (with a simple UI)
The goal is to match the right partner (for ex:a craftsman) to a customer based on their project
requirements. In this particular challenge we concentrate on flooring. The aim is to propose partners to the customers
based on the customer's project requirements.

## Information regarding Solution

I have created a simple UI that offers the following functionality:
1. List all the partners.
2. Finds partners and display partner details, based on the id provided by the user.
3. Based on the user's input data, displays a list of partners that offer the service. The list
   should be sorted by “best match”. The quality of the match is determined first on
   average rating and second by distance to the customer. (If the customer is not in the operating radius(kms) of the
   partner, the matched list won't have that partner in it. )

For the challenge I have created a docker-compose.yml file.
This creates a PostgreSQL database and inserts dummy data. The database contains multiple tables,
the schema is in the image.

![schema](https://github.com/chandanacharya1/customer-matching/blob/feature/version-2-with-UI/db/db.png?raw=true)

The APIs are authenticated using JWT authentication. and the TOKEN is stored in the cookie, the expiration is set to 15 minutes.
So without a valid token the APIs cannot be accessed. For the purpose of the challenge, I store the credentials in a map
instead in a database. In real scenario to be stored in the database.

The posgres server uses port 55432 and the configs are stored in .env file
The web server runs on port 8080
### Having Server up and running
1. `docker-compose up`
2. `go run main.go`

### Testing
#### Login
Open any browser (Chrome, Firefox) and navigate to `http://localhost:8080/`
A login page appears. Enter username and password and click Login.
`Username`: `user1` and `Password`: `password1`

Entering wrong credentials or accesing any urls like `http://localhost:8080/request`, `http://localhost:8080/partner`,
`http://localhost:8080/getpartner` returns 401 status code.
#### To get a list of all partners
Under the heading **List Partners** click on *List all Partners*
This displays list of all the partners in the table

#### To get a details of particular partner
Under the heading **Find  Partner by ID**, enter a Partner ID (For ex: 2, available partnerids in the Postgres db are 1 
to 6). 
Click on *Find Partner* button
This displays the details of partner with partner id 2.
(If a partner is located at multiple addresses, the list will have 2 entries)

If a partner with entered PartnerID is not present,status code 404 is returned.

#### To get list best matched partners
Under the heading **Find Partner Matches**, enter the details.

`Material`:` "tiles"`<br />
`Customer Latitude`: `52.59` <br />
`Customer Longitude`: `13.37` <br />
`Squaremeters`: `800` <br />
`Phone Number`: `0123456789` <br />

Click on *Find Matches* button.

This displays a list of best matched partners.
(If there is a case when 2 matched partners have same rating, the partner with shortest distance will be
displayed on top)

If there are no best matches for entered details, status code 404 is returned.

``
