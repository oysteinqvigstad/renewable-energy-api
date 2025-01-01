# Countries Renewable Historic Information API

## Project Overview

A RESTful API service providing historical renewable energy share percentages by country, built with Golang. The service combines data from REST Countries API and a Renewable Energy Dataset to deliver historical trends with customizable date ranges, country comparisons, mean calculations, and comprehensive sorting capabilities for renewable energy statistics.

This project was developed as part of the Cloud Technologies (PROG2005) course at NTNU during the 2023 academic year. The team, consisting of Idar Løkset Mykløy, Aksel Matashev, and Øystein Qvigstad, successfully implemented the service and deployed it on NTNU's OpenStack instance. The project was focused on implementing proper RESTful semantics and HTTP response codes according to RFC 7231 specifications, ensuring standardized API behavior and communication patterns.

### Key Features
#### Data Access
- Current renewable energy percentages by country
- Historical renewable energy data with customizable date ranges
- Neighboring countries' energy statistics
- Mean value calculations across countries and time periods
#### Notifications
- Webhook registration system
- Persistent webhook storage via Firebase
- Configurable notification triggers based on API call frequency



### API Endpoints Overview



#### Renewables Current
```
GET /energy/v1/renewables/current/{country?}
Optional: ?neighbours=bool
```
#### Renewables History
```
GET /energy/v1/renewables/history/{country?}
Optional: ?begin=year&end=year
```
#### Notifications
```
POST /energy/v1/notifications/
GET /energy/v1/notifications/
GET /energy/v1/notifications/{webhook_id}
DELETE /energy/v1/notifications/{webhook_id}
```
#### Status
```
GET /energy/v1/status/
```

Detailed examples of requests and responses can be found below






## Project Structure

#### Core Components
- `/src/api`: API client implementations
- `/src/cmd`: Main application and stub commands
- `/src/internal`: Core business logic
- `/src/res`: Resource files and data

```
    .
    ├── compose.yaml                                // Docker compose file
    ├── progress.md                                 // Heuristic testing report
    ├── README.md                                   // Provides project documentation and guidelines.
    ├── deploy                                      // Related to service deployment.
    │   ├── Dockerfile                              // Defines how to build a Docker image for the project.
    │   ├── README.md                               // Contains instructions for deploying the project.
    └── src                                         // Contains the source code of the application.
        ├── api                                     // Directory for API-related code.
        │   └── restcountries.go                    // Implements the RESTful countries API client.
        ├── cmd                                     // Applications
        │   ├── app                                 // Main application command.
        │   │   └── main.go                         // Main entry point for the application.
        │   └── stub                                // Stub command for testing purposes.
        │       └── stub_countries_api.go           // Stub implementation for the countries API.
        ├── internal                                // Internal code for the project
        │   ├── firebase_client                     // Directory for Firebase client-related code.
        │   │   ├── bundle_update.go                // Firebase-related code for bundle updates.
        │   │   ├── client.go                       // Firebase client implementation.
        │   │   └── constants.go                    // Constants related to Firebase.
        │   ├── stub                                // Stub command for testing purposes.
        │   │   └── stub_countries_api              // Directory for stub implementation of countries API.
        │   │       ├── constants.go                // Constants for stub countries API.
        │   │       ├── countries_db.go             // In-memory countries database for stub API.
        │   │       ├── handlers.go                 // Handlers for stub countries API.
        │   │       └── handlers_test.go            // Tests for stub countries API handlers.
        │   ├── types                               // Directory for type definitions for various data structures
        │   │   ├── registrations.go                // Data structures for webhook registrations and updates.
        │   │   └── renewable_db.go                 // Data structures and operations related to renewabe energy
        │   ├── utils                               // Utility functions
        │   │   ├── time.go                         // Utility functions for time manipulation.
        │   │   ├── time_test.go                    // Tests for time utility functions.
        │   │   ├── url_parser.go                   // Utility functions for URL parsing.
        │   │   └── url_parser_test.go              // Tests for URL parsing utility functions.
        │   ├── web                                 // Http requests and responses for the web service
        │   │   ├── constants.go                    // Constants related to web handling.
        │   │   ├── cover_test.out                  // Test coverage output for web package.
        │   │   ├── handlers.go                     // Handlers for web-related functions.
        │   │   ├── handlers_test.go                // Tests for web handlers.
        │   │   ├── middleware.go                   // Middleware functions for web handling.
        │   │   ├── routes.go                       // Route definitions for web handling.
        │   │   ├── state.go                        // State management for web handling.
        │   │   ├── structs.go                      // Structs related to web handling.
        │   │   └── webhook.go                      // Webhook-related code for web handling.
        │   └── web_client                          // Internal web client package
        │       ├── client.go                       // Wrapper for http client
        │       └── client_test.go                  // Tests for web client.
        └── res                                     // Resource files containing data used in the project.
        	├── renewable-share-energy.csv          // CSV file with renewable share energy data.
        	└── rest_countries.json                 // JSON file for RESTful countries API.
  ```



## Deployment
#### Prerequisites
- Docker
- Firebase account and credentials
- Firebase secret key file (secret_key.json) placed in the project root directory
#### Local Setup
```
# Start the service with Docker
sudo docker compose up -d
```



# Endpoints (Detailed)

The service exposes four main endpoints through its REST API:

- `/energy/v1/renewables/current`
- `/energy/v1/renewables/history`
- `/energy/v1/notifications/`
- `/energy/v1/status/`



## 1. Endpoint: Current percentage of renewables

This endpoint focuses on returning the latest percentages of renewables in the energy mix.

    Method: GET
    Path: /energy/v1/renewables/current/{country?}

`{country?}`refers to an optional country 3-letter code.

`{?neighbours=bool?}`refers to an optional parameter indicating whether neighbouring countries' values should be shown.

### Request and response examples

**Request:**

`/energy/v1/renewables/current/nor`

**Response**

```
{
    "name": "Norway",
    "isoCode": "NOR",
    "year": "2021",
    "percentage": 71.558365
}
```
**Request:**

`/energy/v1/renewables/current/nor?neighbours=true`

**Response**
```
[
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "2021",
        "percentage": 71.558365
    },
    {
        "name": "Finland",
        "isoCode": "FIN",
        "year": "2021",
        "percentage": 34.61129
    },
    {
        "name": "Russia",
        "isoCode": "RUS",
        "year": "2021",
        "percentage": 6.6202893
    },
    {
        "name": "Sweden",
        "isoCode": "SWE",
        "year": "2021",
        "percentage": 50.924007
    }
]
```

**Request:**(without country code):

`/energy/v1/renewables/current/`

**Response (**returns all countries**)**
```
[
    {
        "name": "Algeria",
        "isoCode": "DZA",
        "year": "2021",
        "percentage": 0.26136735
    },
    {
        "name": "Argentina",
        "isoCode": "ARG",
        "year": "2021",
        "percentage": 11.329249
    },
    {
        "name": "Australia",
        "isoCode": "AUS",
        "year": "2021",
        "percentage": 12.933532
    },
    ...
]
```


## 2. Endpoint: Historical percentages of renewables

The initial endpoint focuses on returning historical percentages of renewables in the energy mix, including individual levels, as well as mean values for individual or selections of countries.

    Method: GET
    Path: /energy/v1/renewables/history/{country?}{?begin=year&end=year?}

`{country?}` refers to an optional country 3-letter code.

`{?begin=year&end=year?}` data will be displayed only for the specified range of years between the 'begin' and 'end' values.

### Request and response examples

**Request:**

`/energy/v1/renewables/history/nor`

**Response**

```
[
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "1965",
        "percentage": 67.87996
    },
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "1966",
        "percentage": 65.3991
    },
    ...
]
```

**Request**(without country code)**:**

`/energy/v1/renewables/history/`

**Response**(returns mean percentages for all countries):
```
[
    {
        "name": "United Arab Emirates",
        "isoCode": "ARE",
        "percentage": 0.0444305504
    },
    {
        "name": "Argentina",
        "isoCode": "ARG",
        "percentage": 9.131337212280702
    },
    {
        "name": "Australia",
        "isoCode": "AUS",
        "percentage": 5.3000481596491245
    },
    ...
]
```

**Request:**

`/energy/v1/renewables/history/nor?begin=1960&end=1970`

**Response**

```
[
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "1965",
        "percentage": 67.87996
    },
    {
        "name": "Norway",
        "isoCode": "NOR",
        "year": "1966",
        "percentage": 65.3991
    },
    ...
]
```



**Request:**

`/energy/v1/renewables/history/nor?end=1967`

**Response**

```
[
  {
    "name": "Norway",
    "isoCode": "NOR",
    "year": "1965",
    "percentage": 67.87996
  },
  {
    "name": "Norway",
    "isoCode": "NOR",
    "year": "1966",
    "percentage": 65.3991
  },
  {
    "name": "Norway",
    "isoCode": "NOR",
    "year": "1967",
    "percentage": 66.591644
  }
]
```
**Request:**

`/energy/v1/renewables/history/?begin=1960&end=2000`

**Response**
```
[
  {
    "name": "Algeria",
    "isoCode": "DZA",
    "percentage": 1.3703186058888892
  },
  {
    "name": "Argentina",
    "isoCode": "ARG",
    "percentage": 7.4716496
  },
  {
    "name": "Australia",
    "isoCode": "AUS",
    "percentage": 4.910349591666669
  },
  {
    "name": "Austria",
    "isoCode": "AUT",
    "percentage": 27.5290
  }
	...
]
```


**Request:**

`/energy/v1/renewables/history/nor?begin=2020`

**Response**

```
[
  {
    "name": "Norway",
    "isoCode": "NOR",
    "year": "2020",
    "percentage": 70.96306
  },
  {
    "name": "Norway",
    "isoCode": "NOR",
    "year": "2021",
    "percentage": 71.558365
  }
]
```

**Request:**

`/energy/v1/renewables/history/nor?begin=2004&end=2006&sortByValue=true`

**Response**

```
[
  {
    "name": "Norway",
    "isoCode": "NOR",
    "year": "2005",
    "percentage": 69.73603
  },
  {
    "name": "Norway",
    "isoCode": "NOR",
    "year": "2006",
    "percentage": 66.73525
  },
  {
    "name": "Norway",
    "isoCode": "NOR",
    "year": "2004",
    "percentage": 64.23876
  }
]
```
## 3. Endpoint: Notification

The Notification Endpoint allows users to register webhooks that will be triggered when the country specified is requested every n (specified in `calls=n`) number of times. The minimum frequency that can be specified is 1. Users can register multiple webhooks, and webhook registrations are persistent, surviving service restarts through the use of a Firebase DB as backend.

### Registration of Webhook

    Method: POST

**Request:**

`/energy/v1/notifications/`

Content type: **`application/json`**

The request body should contain:

- The URL to be triggered upon the event
- The country for which the trigger applies (if empty, it applies to any invocation)
- The number of invocations after which a notification is triggered (it should re-occur every *number of invocations*)

**Example request body:**

```
{
    "url": "http://webhook.site/0aa53816-5e7b-4461-8c1e-d9732383bd0c",
    "country": "FIN",
    "calls": 5
}
```

**Response**

The response contains the unique ID for the registration, which can be used to view detail information or to delete the webhook registration.

**Example response body**:

```
{"webhook_id":"MqZstxmerxzmn"}
```

### Deletion of Webhook

    Method: DELETE

**Request**

`/energy/v1/notifications/MqZstxmerxzmn`

- `{MqZstxmerxzmn}` is the ID returned during the webhook registration

**Response**

Per RFC 7231 guidelines, the server provides a `202 Accepted` status code in response to a probable successful deletion request.

Upon removing a valid webhook ID, it is instantly eliminated from the server. Nonetheless, due to the periodic bulk updates by the firebase worker, confirming the deletion within the persistent system and preventing response delays becomes complex. As a result, the server replies with `202 Accepted` rather than `200 Ok`. It accurately communicates that the server has accepted the request but  hasn't completed the action yet, as there might be some delay or uncertainty in the process due to the periodic updates by the firebase worker.

### View registered webhook

    Method: GET

**Request:**

`/energy/v1/notifications/MqZstxmerxzmn`

- `{MqZstxmerxzmn}`is the ID for the webhook registration

**Response**

```
{
    "webhook_id": "MqZstxmerxzmn",
    "url": "http://webhook.site/0aa53816-5e7b-4461-8c1e-d9732383bd0c",
    "country": "FIN",
    "calls": 5
}
```

### View all registered webhooks

    Method: GET

**Request:**

`/energy/v1/notifications/`

**Response**

```
[
   {
      "webhook_id": "MqZstxmerxzmn",
      "url": "http://webhook.site/0aa53816-5e7b-4461-8c1e-d9732383bd0c",
      "country": "FIN",
      "calls": 5
   },
   {
      "webhook_id": "DiSoisivucios",
      "url": "http://webhook.site/0aa45256-5r6y-4461-8c1e-d73hf793yf73f",
      "country": "SWE",
      "calls": 2
   },
   ...
]
```

## 4. Endpoint: Status 

The Status Endpoint provides an overview of the health and status of various components within the service. It allows users to monitor the connectivity and functionality of external APIs, the Notification Database, and other aspects of the service.

### Request
    Method: GET
    Path: energy/v1/status/

### Response

The response contains the following information:

- **`countries_api`**: The HTTP status code for the *REST Countries API*, indicating the current state of the connection with the external API.
- **`notification_db`**: The HTTP status code for the *Notification DB* in Firebase, reflecting the status of the connection with the database used for storing webhook registrations.
- **`webhooks`**: The total number of registered webhooks in the service, giving users an idea of the current usage.
- **`version`**: The current version of the service (e.g., "v1"), useful for tracking updates and changes to the service.
- **`uptime`**: The time in seconds since the last service restart, providing insight into the stability and performance of the service.

**Example response:**

```
{
  "countries_api": 200,
  "notification_db": 200,
  "webhooks": 1,
  "version": "v1",
  "uptime": 11
}
```


