# Generic Forms API

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
