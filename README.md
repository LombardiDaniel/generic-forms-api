[![Licence](https://img.shields.io/github/license/LombardiDaniel/generic-forms-api?style=for-the-badge)](./LICENSE)
[![BuildStatus](https://img.shields.io/github/actions/workflow/status/LombardiDaniel/generic-forms-api/ci.yml?style=for-the-badge)](https://github.com/LombardiDaniel/generic-forms-api/actions)

# Generic Forms API

https://hub.docker.com/r/lombardi/generic-forms-api

Generic data collector for your startup idea.

```sh
curl -X 'PUT' \
  'http://forms.example.com/v1/entries' \
  -H 'accept: text/plain' \
  -H 'Content-Type: application/json' \
  -d '{
  "data": "example msg",
  "email": "email@example.com",
  "id": "project_name",
  "ts": "2006-01-02T15:04:05Z"
}'
```

This adds to the db (mongo). Note that the field `data` is any JS object, so customize it to your liking!

Then just:

```sh
curl -X 'GET' \
  'http://forms.example.com/v1/entries/ticktr' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer AUTH_TOKEN'
```

...and get the results.
