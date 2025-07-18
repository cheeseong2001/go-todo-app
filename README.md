# Go-ToDo-App

## Introduction
This is my attempt at writing a To-do application. The goal of the project is to help me understand how microservices work, and how they communicate with one another. 

This is an exploratory project where I attempt to write code in Go, and containerise parts of the microservices with Docker Compose.

## Tech Stack
- Programming Language: **Go**
- Containerisation: **Docker with Dockerfiles and Docker Compose**

In the long run:
- [ ] Set up CI/CD for each service
- [ ] Add monitoring
- [ ] Deploy to Kubernetes / Docker Swarm

## Things to be done
- [ ] Auth-service
    - [x] API
        - [x] Register
        - [x] Login
        - [x] Generate JWT on successful login
    - [ ] Test cases

- [ ] Task-service
    - [ ] API
        - [x] Middleware for JWT authentication
        - [x] Create task
        - [x] List tasks or a specific task
        - [x] Delete task
        - [ ] Other features like update task, complete task, adding deadlines
    - [ ] Test cases

- [ ] Deployment-related tasks
    - [x] Containerise services
    - [x] Written `docker-compose.yml`
    - [x] Can deploy with `docker compose up -d`
    - [ ] Deploy to Kubernetes?
    - [ ] Add monitoring?
