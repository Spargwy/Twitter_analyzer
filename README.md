

Web-service for scrapping data from twitter

Instruction:
- 
- create .env file in root directory and add following env variables:
```
DB_PORT_IN_DOCKER=PortForYourDB
USERNAMES=UsernamesThatYouWantToExplore
BEARER_TOKEN=TokenFromYourTwitterApp
CONN="user=postgres password=postgres host=db sslmode=disable dbname=twitter"
```
You can leave the CONN variable unchanged.

- run the comand:

```sudo make docker-run```

Endpoints
-
- /handler - go here for get number of followers, following and tweets of users that you wrote in .env 