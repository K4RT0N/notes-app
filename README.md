# NotesApp

A simple web server app for creating and managing user notes

## Features

- user registration
- session-based authentication
- password hashing with Argon2id
- public/private notes
- authorization by note ownership
- PostgreSQL transactions
- layered architecture

## Tech Stack

- Go
- PostgreSQL
- Gorilla Mux
- Argon2id
- Docker

## Architecture

API -> Service -> DAO -> PostgreSQL

## Requirements

- Go 1.26+
- Docker

## Configuration

 - configured via .env file (see .env.example)

## Running

`./start.sh`
`go run .`
