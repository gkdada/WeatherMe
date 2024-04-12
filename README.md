# WeatherMe
Looks up weather at a given lat/long.

## Configuration

<p>Requires a configuration of:</p>
<ul>
<li>http server port (defaults to 8192),</li>
<li>Weather URL (defaults to "/weather" if not provided), and </li>
<li>a openweathermap.org API key. Application cannot fetch weather info without this key, so the server refuses to run if the key field is not provided.</li>
</ul>


## Usage

Once the server is running, it can be used to fetch weather info in the following way (assuming HttpPort: 4966 in config.json):

### Weather by lattitude/longitude.

    GET http://localhost:4966/weather?lat=37.3&long=-121.9

Where lattitude = 37.3 and logitude = -121.9

## Future development

One possibility is to add the api key in the incoming request so that it doesn't have to be added to configuration.

We can have the service lookup zip code since openweathermap.org API allows for zip code based weather lookup already.

Also, we can return the actual temperature at the given location - either in Fahrenheit or Celsius, depending on country of the location. 

The code can be encapsulated into a small little WeatherMe app, possibly with hard-configured lat/long or zip code to provide one-click weather summary on mobile devices.

