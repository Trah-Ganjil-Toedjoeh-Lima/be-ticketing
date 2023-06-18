Staging Branch: development
[![Docker Image CI](https://github.com/Trah-Ganjil-Toedjoeh-Lima/be-ticketing/actions/workflows/docker-image.yml/badge.svg?branch=development)](https://github.com/Trah-Ganjil-Toedjoeh-Lima/be-ticketing/actions/workflows/docker-image.yml)

## About This Project
This repository is used to developing the ticketing website for Gadjah Mada Chamber Orchestra. This application is based on my previous project ([ticketing-gmco](https://github.com/frchandra/ticketing-gmco)) rewriten with Go language and Gin framework. There are some improvements from the previous project that planned to added in this application.



Tech Stack:
 - Docker version 20.10.22 and docker compose version 2.14.1 
 - Golang version 1.19 (with Gin framework and Gorm ORM)
 - PostgreSQL version 15
 - Redis version 7
 - MinIO
 - ELK stack: https://github.com/deviantony/docker-elk
 - Traefik version 2.9.9




Chandra Herdiputra, January 2023



# Software Testing
A [Postman collection](https://documenter.getpostman.com/view/16816087/2s93si1qAW) has been set up for testing this back-end API.

## How to test?
1. Open the Postman collection
2. Open the 'Environments' tab on the left sidebar
3. Click on 'online'
4. Change the 'BaseUrl' variable to your local URL (explained below)
   - Your local URL should be: http://localhost:8080, therefore set 'BaseUrl' to http://localhost:8080.
5. Have fun testing!

## How to set up? (Docker, WSL2/Ubuntu)
0. Have Docker installed on your system.
1. Clone this repository
2. Navigate to the repository folder
3. Rename 'env.example' to '.env' or create an .env file of your own, following the example.
4. Run these commands in your terminal of choice:

   Building the app:
   ```
   docker-compose -f ./docker-compose.dev.yml build
   ```
   Running the app:
   ```
   docker-compose -f ./docker-compose.dev.yml up db app cache minio
   ```
6. Wait for everything to complete.
7. Done!
8. Navigate to http://localhost:8080. This port number is seen in 'docker-compose.dev.yml', in line 40 and 41:
   ```
   ports:
     - "8080:8080"
   ```

   Updated by: Daffa Romero, 18 June 2023
