<h1 align="center">
  Kick-off league
  <br>
</h1>

<p align="center">
  <a href="#overview">Overview</a> •
  <a href="#features">Feature</a> •
  <a href="#built-with">Built With</a> •
  <a href="#architecture">Architecture</a> •
  <a href="#endpoints">Endpoints</a> •
  <a href="#installation">Installation</a> •
  <a href="#license">License</a>
</p>

## Overview

The Kickoff League platform is designed to facilitate and enhance the experience for both ***competition organizers*** and ***general users***(who love football/futsal) interested in participating in or following sports competitions, This platform aims to make it easy for users to discover, join, and keep track of competitions. It also provides organizers with tools to manage and update competition details effectively, enabling participants and spectators to stay up-to-date.


## Features

### General Users
1. **Discover Competitions**: General users, organizers, and even unregistered visitors can browse available soccer competitions.
2. **Register for Competitions**: Users can apply to join competitions, with eligibility and conditions checked during the registration process.
3. **Create a Team**: Users can establish their own teams and manage team members.
4. **Manage Team Members**: Users can add, remove, or update members within their teams.
5. **Profile Management**: Users can create and update their profiles, making it easy to keep their information current.
6. **View Competition Statistics**: Users can view their own competition statistics and track their progress.

### Competition Organizers
7. **Create Competitions**: Organizers can create new competitions and make them available for others to join.
8. **Update Competition Details**: Organizers can update information about each competition, including the competition's status, to ensure accurate and timely information is available.

### Search and Discovery
9. **User and Organizer Search**: All users, including unregistered visitors, can search for other users within the system, whether they are general users or competition organizers, allowing for easy networking and discovery.

## Built With

- [![Go][GO.dev]][GO-url]
- [![Gorm][GORM.io]][GORM-url]
- [![Gin-gonic][Gin-badge]][Gin-url]
- [![Postgresql][Postgresql-badge]][Postgresql-url]
- [![Docker][Docker-badge]][Docker-url]

## Endpoints

### Auth Endpoints

- **POST** `/auth/register/organizer`
- **POST** `/auth/login`
- **POST** `/auth/register/normal`
- **POST** `/auth/logout`

### View Endpoints

- **GET** `/view/competition`
- **GET** `/view/competition/:id`
- **GET** `/view/match/:matchID`
- **GET** `/view/normalUsers`
- **GET** `/view/normalUsers/:id`
- **GET** `/view/organizer`
- **GET** `/view/organizer/:organizerID`
- **GET** `/view/teams`
- **GET** `/view/teams/:id`
- **GET** `/view/users`
- **GET** `/view/users/:id`

### Image Management Endpoints

- **DELETE** `/image/banner/:competitionID`
- **PATCH** `/image/banner/:competitionID`
- **DELETE** `/image/team/profile/:teamID`
- **PATCH** `/image/team/profile/:teamID`
- **PATCH** `/image/profile`
- **DELETE** `/image/profile`
- **PATCH** `/image/cover`
- **DELETE** `/image/cover`

### Organizer Endpoints

- **POST** `/organizer/competition`
- **PATCH** `/organizer/competition/finish/:id`
- **PATCH** `/organizer/competition/cancel/:id`
- **PATCH** `/organizer/competition/open/:id`
- **PATCH** `/organizer/competition/start/:id`
- **PUT** `/organizer/competition/joinCode/add/:competitionID`
- **PUT** `/organizer/match/:id`
- **PUT** `/organizer`
- **DELETE** `/organizer/competition/:competitionID`

### User Endpoints

- **POST** `/user/competition/join`
- **GET** `/user/nextMatch`
- **PATCH** `/user/normalUser`
- **GET** `/user/teams`
- **POST** `/user/team`
- **DELETE** `/user/team/:teamID`
- **POST** `/user/sendAddMemberRequest`
- **PATCH** `/user/acceptAddMemberRequest`
- **PATCH** `/user/ignoreAddMemberRequest`
- **GET** `/user/requests`
- **POST** `/user/competition/join` 
- **GET** `/user/nextMatch` 
- **PATCH** `/user/normalUser` 
- **GET** `/user/teams` 
- **POST** `/user/team` 
- **DELETE** `/user/team/:teamID` 
- **POST** `/user/sendAddMemberRequest` 
- **PATCH** `/user/acceptAddMemberRequest` 
- **PATCH** `/user/ignoreAddMemberRequest` 
- **GET** `/user/requests` 

## Installation

1. Clone the repository:

   ```bash
   https://github.com/KritAsawaniramol/kick-off-league_server.git
   cd kick-off-league_server
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Set up environment variables

   ```bash
   cp config_example.yaml config.yaml
   ```

4. Run Docker Compose for PostgreSQL, pgAdmin

    ```bash
    docker-compose up -d
    ```

    Optional:
    Accessing a Container's Shell:

    ```bash
    docker exec -it <container name> bash
    ```

    Stopping Services:

    ```bash
    docker-compose down
    ```

5. Run the project:

   ```bash
   go run main.go
   ```

Here’s how you could add instructions to run the application with Docker in your README:

### Running the Application with Docker

To run the application using Docker, follow these steps:

1. **Build the Docker Image**

   ```bash
   docker build -t kickoffbackend:latest .
   ```

2. **Update docker-compose.yml**
    uncomment this section
    ```bash
    # kickoffbackend:
    #   image: kickoffbackend:latest
    #   container_name: kickoffbackend
    #   ports:
    #     - "8080:8080"
    #   restart: unless-stopped
    ```

3. **Run the Docker Container**

   ```bash
   docker-compose up -d
   ```

## License

Distributed under the MIT License. See LICENSE for more information.

---

> GitHub [@kritAsawaniramol](https://github.com/kritAsawaniramol) &nbsp;&middot;&nbsp;



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[Docker-url]: https://www.docker.com/
[Docker-badge]: https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white
[GO-url]: https://go.dev/
[GO.dev]: https://img.shields.io/badge/golang-00ADD8?&style=for-the-badge&logo=go&logoColor=white
[GORM-url]: https://gorm.io/
[GORM.io]: https://img.shields.io/badge/gorm-ORM-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Gin-url]: https://gin-gonic.com/
[Gin-badge]: https://img.shields.io/badge/gin-008ECF?style=for-the-badge&logo=gin&logoColor=white
[Postgresql-badge]: https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white
[Postgresql-url]: https://www.postgresql.org/
