

                                    __   ___  
      __ _ _ __ ___  _   _ _ __    / /_ / _ \
     / _` | '__/ _ \| | | | '_ \  | '_ \ (_) |
    | (_| | | | (_) | |_| | |_) | | (_) \__, |
     \__, |_|  \___/ \__,_| .__/   \___/  /_/
     |___/                |_|
    
     Transforming Code, One Upside-Down Innovation at a Time




## Project Structure

    .
    ├── deploy                          // related to service deployment
    │   ├── build.sh
    │   ├── Dockerfile
    │   └── README.md
    └── src                             // source code for the project
        ├── api
        ├── cmd                         // applications
        │   ├── app
        │   │   └── main.go
        │   └── stub
        └── internal                    // internal code for the project
            ├── db                      // interaction with data
            │   ├── db.go
            │   └── structs.go
            ├── utils                   // utility functions
            │   ├── time.go
            │   ├── url_parser.go
            │   └── url_parser_test.go
            └── web                     // http requests and responses for the web service
                ├── constants.go
                ├── handlers.go
                ├── middleware.go
                └── routes.go



## Project status

Currently in development

## Roadmap

See [milestones](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2023-workspace/idarmy/project-group69/-/milestones/1#tab-issues)

## Contributers

* Idar Løkset Mykløy
* Aksel Matashev
* Øystein Qvigstad

# Assignment 2

# Project overview

Our team has been assigned the task of creating a REST web application using Golang. The primary objective of this application is to offer clients access to information on renewable energy production developments within and across countries. We will achieve this goal by utilizing an existing web service in combination with our own data-centric web service, which will involve creating endpoints to expose our service. Additionally, our application will allow clients to register notifications using webhooks. The application will be dockerized, and we will deploy it using an IaaS system.

To develop this project, we will be using the REST Countries API and the Renewable Energy Dataset. The Renewable Energy Dataset will provide the percentage of renewable energy in a country's energy mix over time, serving as the foundation for our service.

Our team has effectively developed the necessary services by minimizing the number of requests made to these services and utilizing the most appropriate endpoints provided by the APIs. During the development process, we have also stubbed the services for testing purposes, ensuring that we did not use the API services in our tests.

We successfully deployed the final web service on our local OpenStack instance, SkyHigh. We initially developed the project on our local machines before moving to deployment. We are submitting both a URL to the deployed service and our code repository.

Throughout this project, we familiarized ourselves with various technologies introduced as part of the course, ensuring a smooth development process. We actively participated in or reviewed the lectures to understand these technologies.

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

### **Advanced Tasks:**

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

**TODO**

## 4. Endpoint: Status 

The Status Endpoint provides an overview of the health and status of various components within the service. It allows users to monitor the connectivity and functionality of external APIs, the Notification Database, and other aspects of the service.

### Request

`Method: GET
Path: energy/v1/status/`

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
   "countries_api": "<http status code for *REST Countries API*>",
   "notification_db": "<http status code for *Notification DB* in Firebase>",
   ...
   "webhooks": <number of registered webhooks>,
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}
```