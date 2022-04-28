# **fetchRewardsAssessment**

# Running the Web Server
1. Clone the repo
2. Locate the repo in your terminal
3. Run the following command: `go run server.go`

# API Documentation

## 1. Get Transactions (Designed For Testing Purposes)
```
curl http://localhost:8080/transaction
```
Example Response:
```
[
    {
        "payer": "DANNON",
        "points": 300,
        "timestamp": "2020-10-31T10:00:00Z"
    },
    {
        "payer": "UNILEVER",
        "points": 200,
        "timestamp": "2020-10-31T11:00:00Z"
    },
    {
        "payer": "DANNON",
        "points": -200,
        "timestamp": "2020-10-31T15:00:00Z"
    },
    {
        "payer": "MILLER COORS",
        "points": 10000,
        "timestamp": "2020-11-01T14:00:00Z"
    },
    {
        "payer": "DANNON",
        "points": 1000,
        "timestamp": "2020-11-02T14:00:00Z"
    }
]
```
HTTP Request
```
GET http://localhost:8080/transaction
```
## 2. Add Transaction
```
curl http://localhost:8080/transaction \
    -H "Content-Type: application/json" \
    -X POST \
	  -d '
    {
      "payer" : "DANNON",
      "points" : 1000,
      "timestamp" : "2020-11-02T14:00:00Z"
    }'
```
Example Response:
```
{
    "payer" : "DANNON",
    "points" : 1000,
    "timestamp" : "2020-11-02T14:00:00Z"
}
```
HTTP Request
```
POST http://localhost:8080/transaction
```
Request Arguments
| Field | Type | Description |
| ----------- | ----------- | ----------- |
| `payer` | string (**Required**) | The partner/payer involved in the transaction
| `points` | integer (**Required**) | Number of points lost or gained in a transaction
| `timestamp` | date (**Required**) | The date and time the transaction happened

## 3. Spend Points
```
curl http://localhost:8080/points \
    -H "Content-Type: application/json" \
    -X PUT \
    -d '
	  {
	    "points": 5000
	  }'
```
Example Response:
```
[
    {
        "payer": "DANNON",
        "points": -100
    },
    {
        "payer": "UNILEVER",
        "points": -200
    },
    {
        "payer": "MILLER COORS",
        "points": -4700
    }
]
```
HTTP Request
```
PUT http://localhost:8080/points
```
Request Arguments
| Field | Type | Description |
| ----------- | ----------- | ----------- |
| `points` | integer (**Required**) | Number of points to spend (Can't be negative and can't be more than the total number of points the user has)

## 4. Get Payer Point Balances
```
curl http://localhost:8080/points
```
Example Response:
```
{
    "DANNON": 1000,
    "MILLER COORS": 5300,
    "UNILEVER": 0
}
```
HTTP Request
```
GET http://localhost:8080/points
```
