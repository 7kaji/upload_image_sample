# upload_image_sample

## backend

AWS (S3) + Ruby

```bash
$ mv .env.sample .env
$ bundle install --path vendor/bundle
$ bundle exec rackup config.ru
```
GCP (GCS) + Go

```bash
$ mv .env.sample .env
$ go run app.go
```

## frontend

```bash
$ cd frontend
$ yarn install
$ yarn run start
```
[![Image from Gyazo](https://i.gyazo.com/4327d3d76b021f931f5159a1a8bc4639.png)](https://gyazo.com/4327d3d76b021f931f5159a1a8bc4639)
