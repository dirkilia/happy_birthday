CREATE TABLE IF NOT EXISTS employees
(
    id        INTEGER PRIMARY KEY,
    login     TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name     TEXT NOT NULL,
    surname TEXT NOT NULL, 
    patronymic TEXT NOT NULL,
    birthday TEXT NOT NULL,
    enable_notifications INTEGER NOT NULL,
    notify_of TEXT DEFAULT ""
);
