## GOKOMODO HOME TEST

Project ini menggunakan bahasa program golang & database postgres

## Installation

1. Replace dan ganti nama config.yml.example menjadi config.yml
2. Setup konfigurasi database
3. Jalankan migrasi database menggunakan command `go run main.go migration up`
4. Jalankan service menggunakan command `go run main.go start`, akses menggunakan port `8081`
5. Unit test menggunakan command `go test main_test.go -v`

## Docker

1. Jalankan command `docker-compose up -d --build`
2. Jalankan command `docker exec gokomodo go run main.go migration up`
3. Akses menggunakan port `8081`

## API Docs

1. Import collection `Gokomodo.postman_collection.json` pada aplikasi Postman
2. Akses user tersedia menggunakan email dengan password `secret` :
   - seller@mail.com / Role seller
   - buyer@mail.com / Role buyer
