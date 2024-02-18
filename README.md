# WeatherMe
Looks up weather at a given lat/long or a given zip code (with optional country code)

## Configuration

<p>Requires a configuration of:</p>
<ul>
<li>http server port (defaults to 8192),</li>
<li>Weather URL (defaults to "/weather" if not provided), and </li>
<li>a openweathermap.org API key. Application cannot fetch weather info without this key, so the server refuses to run if the key field is not provided.</li>
</ul>


## Usage

Once the server is running, it can be used to fetch weather info in one of the three following ways. All examples assume using localhost, 4966 for http port, and "/weather" for weather URL.

### Weather by lattitude/longitude.

    GET http://localhost:4966/weather?lat=37.3&long=-121.9

Where lattitude = 37.3 and logitude = -121.9

### Weather by zip code (US zip codes only)

    GET http://localhost:4966/weather?zip=90210

### Weather by zip code and country code 

    GET http://localhost:4966/weather?zip=2020&country=AU

Where zip = 2020 and country = AU (Austrilia. The zip code is for Sydney Internatinal Airport)

## Future development

One possibility is to add the api key in the incoming request so that it doesn't have to be added to configuration.

The code can be encapsulated into a small little WeatherMe app, possibly with hard-configured zip code to provide one-click weather summary on mobile devices.

