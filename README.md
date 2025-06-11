# plaYoo - Find Your Gaming Teammates 🎮

## Overview
plaYoo is a pet project designed to help gamers find teammates for online multiplayer games. The platform allows users to create and join gaming events, making it easier to connect with like-minded players.

## Features
- 🎯 **Event Creation System**: Create events with custom parameters (game, time, max members)
- 🔎 **Search For Events**: Find events based on your gaming preferences and schedule
- 🔔 **Notifications**: Get notifications about the start of events

## Tech Stack
- **Backend**: Golang (Fiber)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Infra**: Docker + Docker Compose

### Prerequisites
- Docker Desktop for Windows
- Git for Windows
- Windows Terminal (recommended)

### Installation
1. Clone the repository:
```powershell
git clone https://github.com/kiryuhakipyatok/playooo
cd plaYoo
```
2. Configure environment or config:
```powershell
Copy-Item .env.example .env
```
OR 
```powershell
Copy-Item config-example.yaml config.yaml
```
Edit the .env file or config.yaml with your configuration.
```
3. Start services:
```powershell
docker-compose up --build
```
4. The application should now be running at http://localhost:1111.

After starting the application, you can access the API documentation at:
http://localhost:1111/swagger
