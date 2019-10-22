### DEPLOYMENT for TEST
To run the server for test, you can follow these steps:
1. Make sure you have docker installed, or a mongo db server accessible through
port localhost:27017. If you have docker installed and running on your computer
just use docker-compose file in ./extra/dev-env folder to spin up a mongodb





### TODOs
1. Each user must have only one session, creating a new session must close all other sessions.
2. Create music db and search for musics on that
3. Provide a web-hook for crawlers to pushed what they found.