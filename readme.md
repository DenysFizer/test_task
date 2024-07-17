# Spy Cat Agency Management Application

## Overview
This application is designed to manage the operations of the Spy Cat Agency (SCA), simplifying their spying work processes. It allows the management of spy cats, missions they undertake, and targets they are assigned to.


## Technologies Used
- **Backend Framework**: Echo
- **Database**: PostgreSQL
- **Database Driver**: pgx
- **Migrations**: migrations (as cli tool)

## Getting Started

### Prerequisites
- Docker
- Docker Compose
- Makefile

### Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/DenysFizer/test_task
   cd spy-cat-agency    
2. Run docker compose command:
    ```bash
   docker compose up
3. Run migrations using makefile
   ```bash
   make migrations-up
   ```
   or 
    ```bash
    make migrations-down
   ```
   to down all created migrations

### Config
Create a config.json file in the root directory with the following structure:
 ```
  {
  "db_host": "your_host",
  "db_port": your_port,
  "db_user": "user",
  "db_password": "password",
  "db_name": "db_name",
  "APIKey": "your-cat-api-key"
 }
   ```
### Postman collection file 
#located in root directiory as
```
test.postman_collection.json
```