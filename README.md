
# Public GitHub Events

### This repo consists of 3 components:
* `events-collector` - fetches events from public GitHub REST & GraphQL apis.
* `events-api` - list and counts events, repos and users to the user.
* `mongodb` - persists the events, repos and users.

### Deployment instructions:
* Navigate to `deployment` dir.
* Edit `docker-compose.yaml` and  add your `GITHUB_TOKEN` env variable under the `collector` container (alongside `MONGO_DB_URL` and `MONGO_DB_PORT`).
* then run the docker compose: `docker-compose up --build -d`.
* Once all 3 containers are up & running, the `events-api` service is ready to get requests.

-----------------------------
## Events API Service
By default, events api service is deployed to http://localhost:8080/.

**All responses are in json format**

Events api service supports the following GET requests:

* `list` - lists entities collected by the `events-collector`
  accepts the following params:

| Parameter Key | Details                                            | Supported Values                             | Default Value |
|---------------|----------------------------------------------------|----------------------------------------------|---------------|
| dataType      | set data type to list                              | "events", "repos", "users"                   | "events"      |
| orderBy       | set the column to use in order to sort the results | columnNames, like "_id" or "last_updated_at" | "_id"         |
| orderType     | set the order type to apply                        | "ascending", "descending"                    | "ascending"   |
| limit         | set the num of returned entities in the result     | non-negative int. specify "0" for no limit   | 20            |

* `count` - count all entities collected by the `events-collector`
  accepts the following params:

| Parameter Key | Details                                            | Supported Values                             | Default Value |
|---------------|----------------------------------------------------|----------------------------------------------|---------------|
| dataType      | set data type to list                              | "events", "repos", "users"                   | "events"      |

## Examples

* "List all events" - http://localhost:8080/list?dataType=events&limit=0
* "Count all events" - http://localhost:8080/count?dataType=events
* "List the 20 most recent actors that were involved in the events that you collected" - http://localhost:8080/list?dataType=users&limit=20&orderBy=last_updated_at&orderType=descending
* "List the 20 most recent repositories that were involved in the events that you collected, including the amount of stars that each one of them has" - http://localhost:8080/list?dataType=repos&limit=20&orderBy=last_updated_at&orderType=descending

-----------------------------

## Environment Variables
| Key                             | Details                            | Supported Services            | Default Value |
|---------------------------------|------------------------------------|-------------------------------|---------------|
| MONGO_DB_URL                    | mongo db url to use                | events-collector, events-api  | localhost     |
| MONGO_DB_PORT                   | mongo db port to use               | events-collector, events-api  | 27017         |
| GITHUB_GRAPHQL_REQUEST_TIMEOUT  | timeout for github graphql queries | events-collector              | 30            |
| GITHUB_TOKEN                    | token to be used with graphql api  | events-collector              |               |
| FETCH_INTERVAL_MINUTES          | events fetch interval in minutes   | events-collector              | 1             |
| MAX_ITEMS                       | max items in batch                 | events-collector              | 100           |
| MAX_TIMEOUT_SECONDS             | max batch timeout in seconds       | events-collector              | 100           |
| EVENTS_DB                       | db name for storing events         | events-collector, events-api  | github        |
| EVENTS_COLLECTION               | collection name for storing events | events-collector, events-api  | events        |
| REPOS_DB                        | db name for storing repos          | events-collector, events-api  | github        |
| REPOS_COLLECTION                | collection name for storing repos  | events-collector, events-api  | repos         |
| USERS_DB                        | db name for storing users          | events-collector, events-api  | github        |
| USERS_COLLECTION                | collection name for storing users  | events-collector, events-api  | users         |