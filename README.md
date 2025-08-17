# Zota
Key-Value Store which randomly breaks

====================================================================================

To Start the project
- clone the repository
- fill the environment variables in the .env file
- Run with command < go run main.go > or with the run of the IDE

=====================================================================================

PostgreSQL was chosen as the database for data storage.
A struct model (Store) is used, which provides methods for Put, Get, Delete, and Dump.
An isUnlucky function has been created, which returns an error in ~30% of cases to simulate unstable system.
The project uses pgxmock for unit testing, tests cover both successful operations and error cases.