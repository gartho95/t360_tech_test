**HOW TO RUN**

Rename ```.env.example``` to ```.env```

IF you have a google project ID with a pubsub set up:

fill in the ```GOOGLE_CLOUD_PROJECT_ID``` in the .env

elsewise keep it empty and the code will return a mocked example


in the project directory run:

```go mod tidy```
```go run main.go```



With more time I would've liked to add unit testing, and proper error handling everywhere as I am aware the way it currently handles it in the concurrency is not good.

Below are the notes I made prior to writing any code, to hopefully give a small idea of my thought process:

//essentials
//get list of test VRMS (vehicle registrations) and store locally in project, seems that url containing VRMS calls data server side, quick webcrawler
//create env to set base_url and google cloud creds
//create structs for post request and responses
//function to build map of VRMS and their associated company
//function to make requests that handles timeout
//concurrent handling of requests, one per datasource then one per VRM for each datasource (one concurrent loop for company with another for vehicle within company)
//mock function to handle publishing positive_searches, once each concurrent process if it has succeeded

//nice to have
//handle non 200 responses from requests
//full google cloud intergration
//clean concurrent error handling and logging
//dockerisation
//github actions to pull env from staging/production config
//unit tests
//mock batch publishing positive_searches
