# Group 69 - Assignment 2

[TOC]

# Project overview

Our team has been assigned the task of creating a REST web application using Golang. The primary objective of this application is to offer clients access to information on renewable energy production developments within and across countries. We will achieve this goal by utilizing an existing web service in combination with our own data-centric web service, which will involve creating endpoints to expose our service. Additionally, our application will allow clients to register notifications using webhooks. The application will be dockerized, and we will deploy it using an IaaS system.

To develop this project, we will be using the REST Countries API and the Renewable Energy Dataset. The Renewable Energy Dataset will provide the percentage of renewable energy in a country's energy mix over time, serving as the foundation for our service.

Our team has effectively developed the necessary services by minimizing the number of requests made to these services and utilizing the most appropriate endpoints provided by the APIs. During the development process, we have also stubbed the services for testing purposes, ensuring that we did not use the API services in our tests.

We successfully deployed the final web service on our local OpenStack instance, SkyHigh. We initially developed the project on our local machines before moving to deployment. We are submitting both a URL to the deployed service and our code repository.

Throughout this project, we familiarized ourselves with various technologies introduced as part of the course, ensuring a smooth development process. We actively participated in or reviewed the lectures to understand these technologies.



## Project status

Currently in development

## Roadmap

See [milestones](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2023-workspace/idarmy/project-group69/-/milestones/1#tab-issues)

## Contributers

* Idar Løkset Mykløy
* Aksel Matashev
* Øystein Qvigstad




## Project Structure
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







# Endpoints

Our web service will have four resource root paths:

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


# Deployment

The service is deployed on an IaaS solution OpenStack using Docker.

**URL to the deployed service:** http://10.212.172.171:8080/energy/v1/

The service can also be deployed locally using the provided docker compose file.
For the local instance to work, a firebase secret key must be placed in the root folder of the project, and named `secret_key.json`.

Once the secret is provided, the service can be run with this command:
```sh
sudo docker compose up -d
```



## Automated testing

The project consists of unit and integration tests. To separate the tests from third-party services, ensure the stubbing service for the REST countries API is initiated on port 8081:

> go run ./cmd/stub/stub_countries_api.go

Once the stubbing service is active, initiate the automated tests using this command:

> go test ./...

By default, tests will be executed without the firebase update worker. To override this, modify the `testWithFirebase` boolean in `./internal/web/handlers_test.go`. However, this also requires the presence of an authentication key for the tests to pass.
