# Team Project Group 30

Our implementation of the team project uses HTML, JavaScript and CSS on the frontend. It has a Go backend using a PostgreSQL database instance, both running inside Docker containers.

The project can be started by entering the `/src` folder and using `make run`. It can then be stopped with `make stop`.

The frontend can be started by running `python -m http.server 8080` inside the `/Oaxaca` folder to run the website on `localhost:8080`.

Structure
=========
The backend code is in `/src`. This contains the Docker image information as well as the makefiles to run them. Frontend code is in `/Oaxaca`.
