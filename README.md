# GoChat

A Fullstack real-time chat application built with Go, React and WebSockets, easily deployable with Docker.

## Technical Stack

* **Backend :** Go (Golang), Gorilla WebSocket
* **Frontend :** React 19, Vite, SCSS
* **Base de données :** PostgreSQL
* **DevOps :** Docker, Docker Compose, Nginx

## Installation

The project is fully containerized. You only need **Docker** and **Docker Compose** installed.

1.  **Clone the repository :**
    ```bash
    git clone https://github.com/felixp57/gochat.git
    cd gochat
    ```

2.  **Configure the environment :**
    Copy the example environment file to setup your configuration :
    ```bash
    cp .env.example .env
    ```

3.  **Launch the application :**
    ```bash
    docker-compose up --build
    ```

4.  **Access the application :**
    Open your browser and go to : **http://localhost**

## Project Structure

* `/backend` : Go API and WebSockets handling.
* `/frontend` : React (Vite) User Interface.
* `compose.yml` : Services orchestration (App, DB, Nginx).
* `init.sql` : Initialisation script for the database.
