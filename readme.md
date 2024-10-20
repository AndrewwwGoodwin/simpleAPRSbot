# Simple APRS bot

### Video

[![APRS Bot Demonstration](https://img.youtube.com/vi/2dZiYyuAWDY/0.jpg)](https://www.youtube.com/watch?v=2dZiYyuAWDY)


### Description
This is my first attempt at a simple APRS bot written in GoLang

Expandable by simply adding a command file, and pointing to it via commandRegistry in main.go

http://www.aprs.org/ <br>
https://aprs.fi/ <br>
http://www.aprs.org/doc/APRS101.PDF <br>

Weather Info from: https://openweathermap.org/

https://osu.ppy.sh/ <br>
https://osu.ppy.sh/docs/index.html

Special shout out to the people @ https://github.com/ebarkie/aprs.

## Running with docker-compose
Copy docker-compose-template.yml to docker-compose.yml

Edit the docker-compose file with your api keys

docker-compose up

### Required Environment Variables
APRS_CALL= <br>
APRS_PASS= <br>

### Optional Environment Vars
APRS_FI_API_KEY= <br>
OWM_API_KEY= <br>
OSU_CLIENT_ID= <br>
OSU_CLIENT_SECRET= <br>