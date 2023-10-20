# Connecting Flights

Go serverless aplication running on AWS Lambda.

##Testing the application:

THe aplication is deployed on AWS Lambda and can be tested by sending a POST request to the following URL:
https://k7rt6og6p7.execute-api.us-east-1.amazonaws.com/dev/

Send a post to URL with a body like this:
```json
{
  "departure": "DUB",
  "departureDateTime": "2023-10-19T15:00:00Z",
  "arrival": "WRO",
  "arrivalDateTime": "2024-12-19T18:00:00Z"
}
```
