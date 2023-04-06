

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
