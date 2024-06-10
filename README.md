# happy_birthday API

Create .env file in the root of the project directory like this
## .env

    PORT=8080
    APP_ENV=local

    DB_URL=path/to/database.db
    JWT_SECRET=your-super-secret-jwt-signature-key

First run migrations

## Run migrations

    task mig

Then run app

## Run the app

    task run

# REST API

## Register

### Request body

`POST /register`
    
    {
        "login": "test",
        "password": "test",
        "first_name": "test",
        "surname": "test",
        "patronymic": "test",
        "birthday": "2000-1-1",
        "enable_notifications": 0
    }

### Response

    {
        "status": "OK"
    }

## Login

### Request body

`POST /login`

    {
        "login": "test",
        "password": "test"
    }

### Response if user didn't set birthdays he wants to be notified of

    {
        "message": "Logged in successfully"
    }

### Response if user set birthdays he wants to be notified of

    [
        {
            "id": 1,
            "login": "test",
            "password": "$2a$10$Pf3KhmumqZ9qXkPPyWS9DOHn9JPUjYdWjXR6g649nR9SIOLEhyU.O",
            "first_name": "test",
            "surname": "test",
            "patronymic": "test",
            "birthday": "2000-1-1",
            "enable_notifications": 0,
            "notify_of": ""
        }
    ]

## User info

`GET /info`

### Response

    {
        "id": 1,
        "login": "test",
        "password": "$2a$10$1FU3S/7jgbcmEMXuldE0vOVSeE57mHNxPzi/DR1X/eVvYlXuL.gi6",
        "first_name": "test",
        "surname": "test",
        "patronymic": "test",
        "birthday": "2000-1-1",
        "enable_notifications": 1,
        "notify_of": "test;test1;test2"
    }

## Get list of all employees 

`GET /employees`

### Response

    [
        {
            "id": 1,
            "login": "test",
            "password": "$2a$10$Pf3KhmumqZ9qXkPPyWS9DOHn9JPUjYdWjXR6g649nR9SIOLEhyU.O",
            "first_name": "test",
            "surname": "test",
            "patronymic": "test",
            "birthday": "2000-1-1",
            "enable_notifications": 0,
            "notify_of": ""
        },
        {...}
    ]

## Get list of all employees today birthdays

`GET /emptdaybdays`

### Response

    [
        {
            "id": 1,
            "login": "test",
            "password": "$2a$10$Pf3KhmumqZ9qXkPPyWS9DOHn9JPUjYdWjXR6g649nR9SIOLEhyU.O",
            "first_name": "test",
            "surname": "test",
            "patronymic": "test",
            "birthday": "2000-1-1",
            "enable_notifications": 0,
            "notify_of": ""
        },
        {...}
    ]

## Set employees user want to be notified about

### Request body

`POST /setnotify`

    {
        "login": "test",
        "notify_of": "test;test1;test2"
    }

### Response

    {
        "status": "OK"
    }

## Logout

`POST /logout`