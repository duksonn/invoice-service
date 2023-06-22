<h1 align="center">Invoice Financing Application</h1>

<p>This repository contains the source code and documentation for an Invoice Financing application. The application is designed to manage invoices, facilitate bidding and trading processes between investors and invoice issuers. It enforces transactional consistency and provides concurrency handling for invoice management.</p>

## Layers
This project was built with hexagonal architecture and have the next layers:
* Handler
* Dependencies
* Config
* Service
* Domain
* Repository

Also, DDD was used to design data structs.

### External dependencies
* _PostgreSQL_

##1. API Definition
   The application provides an API with the following endpoints:

###1.1. Create an Invoice
Endpoint: `POST /v1/invoice`

Create a new invoice. The request body should contain the invoice details.
#### Request
- All fields must be in payload
```shell
{
    "items": [
        {
            "id": "item_1",
            "description": "Truck",
            "price": 25000,
            "quantity": 1
        },
        {
            "id": "item_2",
            "description": "Fire item",
            "price": 5000,
            "quantity": 2
        }
    ],
    "issuer_id": 1
}
```
#### Response
```shell
{
    "id": "bcd1fd61-b25f-4eb0-b5fd-65fac0774f96",
    "due_date": "2023-07-13T00:03:12-03:00",
    "asking_price": 35000,
    "status": "ACTIVE",
    "items": [
        {
            "id": "item_1",
            "description": "Truck",
            "price": 25000,
            "quantity": 1
        },
        {
            "id": "item_2",
            "description": "Fire item",
            "price": 5000,
            "quantity": 2
        }
    ],
    "created_at": "2023-06-22T00:03:12-03:00",
    "issuer_id": 1
}
```

###1.2. Get Invoice
Endpoint: `GET /v1/invoice/{invoice_id}`

Retrieve an invoice by its ID, including its current status and the investor who traded it (if applicable).
#### Response
```shell
{
    "id": "bcd1fd61-b25f-4eb0-b5fd-65fac0774f96",
    "due_date": "2023-07-13T00:03:12-03:00",
    "asking_price": 35000,
    "status": "ACTIVE",
    "items": [
        {
            "id": "item_1",
            "description": "Truck",
            "price": 25000,
            "quantity": 1
        },
        {
            "id": "item_2",
            "description": "Fire item",
            "price": 5000,
            "quantity": 2
        }
    ],
    "created_at": "2023-06-22T00:03:12-03:00",
    "issuer_id": 1,
    "investors_ids": null
}
```
_If the invoice has been purchased the investors_ids will appear in response._

###1.3. Get Issuer
Endpoint: `GET /v1/issuer/{issuer_id}`

Retrieve information about an issuer, including their balance.
#### Response
```shell
{
    "id": 2,
    "company_name": "Amazon INC",
    "available_funds": 3000
}
```

###1.4. Get Investors
Endpoint: `GET /v1/investor`

Retrieve a list of all investors and their balances.
#### Response
```shell
{
    "investors": [
        {
            "id": 1,
            "name": "John Doe",
            "available_funds": 0
        },
        {
            "id": 2,
            "name": "Robert McKenzie",
            "available_funds": 1000000
        },
        {
            "id": 3,
            "name": "Margaret Hill",
            "available_funds": 5000
        }
    ]
}
```

###1.5. Place a Bid
Endpoint: `POST /v1/bid/place`

Place a bid to purchase an invoice. The request body should contain the details of the bid.
#### Request
- All fields must be in payload
```shell
{
    "invoice_id": "a654a192-7aa7-48ae-b888-e876e99785e2",
    "investor_id": 3,
    "bid_amount": 500
}
```
#### Response 204 NoContent

###1.6. Approve a Trade
Endpoint: `PUT /v1/trade/{trade_id}?approved=true`

Approve a trade by specifying the trade ID and passing approved true for approvals or false for rejections.
#### Response 204 NoContent

###1.7. Get Trades
Endpoint: `GET /v1/trade or GET /v1/trade?status=WAITING_APPROVAL`

Retrieve information about trades. Use the first endpoint to get a list of all trades and the second endpoint to get details about a specific status.
Only allowed status can be passed by query arg.
#### Response
```shell
{
    "trades": [
        {
            "id": "0aefe073-5786-4fdb-adac-ef2537e3a105",
            "invoice_id": "a654a192-7aa7-48ae-b888-e876e99785e2",
            "investors_ids": [
                2,
                3
            ],
            "trade_status": "WAITING_APPROVAL",
            "created_at": "2023-06-21T01:08:50-03:00",
            "updated_at": null
        }
    ]
}
```

_Please refer to the source code for more detailed information about request and response formats._

##2. Getting Started

To run the application locally, follow these steps:

1. Clone the repository: git clone <repo-url>
2. Install the required dependencies.
3. Configure the necessary environment variables.
4. Download any postgreSQL server, run it in port 5432 and create schema called bankable.
5. Run the migration file `schema.sql` provided in the source code inside migrations folder.
6. Start the application.
7. Access the API using the provided endpoints or importing the collection included in the source code named `invoice-service.postman_collection.json`.

_If any problem with configurations, please check `env_dev.json` file._

##3. Running the tests

In order to run the project tests you need to execute the following command:

```shell
go test ./...
```