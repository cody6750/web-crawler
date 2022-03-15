<div id="top"></div>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->



<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/cody6750/discord-tracking-bot">
    <img src="media/crawler.png" alt="Logo" width="240" height="240">
  </a>

<h1 align="center">Web Crawler</h1>
</div>



<!-- TABLE OF CONTENTS -->
## Table of Contents
<ol>
  <li>
    <a href="#about-the-project">About The Project</a>
    <ul>
      <li><a href="#built-with">Built With</a></li>
    </ul>
  </li>
  <li>
    <a href="#getting-started">Getting Started</a>
    <ul>
      <li><a href="#prerequisites">Prerequisites</a></li>
      <li><a href="#installation">Installation</a></li>
    </ul>
  </li>
  <li><a href="#usage">Usage</a></li>
  <li><a href="#environment variables">Environment Variables</a></li>
  <li><a href="#features">features</a></li>
  <li><a href="#contributing">Contributing</a></li>
  <li><a href="#license">License</a></li>
  <li><a href="#contact">Contact</a></li>
</ol>

<!-- ABOUT THE PROJECT -->
## About The Project

![Product Name Screen Shot][product-screenshot]
![Tracking Screen Shot][tracking-screenshot]

A fast high-level web crawling and web scraping application framework using Go. Used to crawl websites concurrently and extract structured data from their pages.

<p align="right">(<a href="#top">back to top</a>)</p>

### Built With

* [Go](https://go.dev/)
* [discordgo](https://github.com/bwmarrin/discordgo)
* [Logrus](https://github.com/sirupsen/logrus)
* [Docker](https://www.docker.com/)
* [AWS](https://aws.amazon.com/)
<p align="right">(<a href="#top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

1. This assumes you already have a working Go environment, if not   please see
  [this page](https://golang.org/doc/install) first.

2. This assumes you already have a working Docker environment, if not please see
  [this page](https://www.docker.com/get-started) first.

3. This assumes you already have a working Discord bot, if not please see
[this page](https://discord.com/developers/docs/intro) first.

4. If you are deploying the container to AWS, please configure your AWS credentials.
  Please see [this page](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/setup-credentials.html) for assistance.

5. The Discord Tracking Bot is designed to call upon the  [web crawler](https://github.com/cody6750/web-crawler) via REST API calls and return that       response in a formated structure to the designated discord channel. This assumes you already have a working webcrawler deployed. If not, please see [web  crawler](https://github.com/cody6750/web-crawler) for deployment instructions.

### Installation

1. In your Discord developer portal, create an API token for the Discord bot. If you are using AWS, you are able to store the API token into AWS secrets.

2. Clone the repo
   ```sh
   git clone https://github.com/cody6750/discord-tracking-bot
   ```

3. Get Go packages
   ```sh
	go get github.com/aws/aws-sdk-go
	go get github.com/bwmarrin/discordgo 
	go get github.com/sirupsen/logrus 
   ```

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Usage

The Discord Tracking bot was designed to be deployed on AWS EC2 as a Docker container, though it can be deployed locally by building and executing the go files or deploying the Docker container locally on your machine. This section will cover only 2 of the ways to do so. Please note that these instructions are for Mac OS using a bash terminal. 

### Set up Discord server

1. Set up `metrics` channel. Metrics will be sent to this channel.
2. Set up `logs` channel. Logs will be sent to this channel.
3. Set up `bot_console` channel. Bot will respond to commands in this channel.

![Admin channels][admin-channels]

4. Set up category. Naming must be in this convention `tracking_<CATEOGORY>`. For example `tracking_graphics_cards` or `tracking_consoles`. The bot uses the category to determine what to track in which channels. The category `graphics cards` or `consoles` is used in `/pkg/configs/<CATEGORY>` to determine which directory to use.
5. Set up channels. Naming must be in this convention `tracking_<ITEM_NAME>`. The item name is used to in the `/pkg/configs/<CATEGORY>/<ITEM_NAME>` to choose the config file.

![Tracking channels][tracking-channels]
![Tracking configs][tracking-configs]


### Set AWS Secret 
If you are planning on storing your Discord Api token in AWS secrets, please follow these instructions so that the bot can grab them. If not, please continue to the next section of the guide.
1. Create AWS Secret with your Discord Api token.

2. Set `LOCAL_RUN` as `false`
  ```sh
  set LOCAL_RUN=false
  ```

3. Set `DISCORD_TOKEN_AWS_SECRET_NAME` as the name of your AWS secret
  ```sh
  set DISCORD_TOKEN_AWS_SECRET_NAME=<SECRET_NAME>
  ```


### Build locally without Docker

1. Navigate to the `discord-tracking-bot` repo location.
2. Set `LOCAL_RUN` environment variable to `true`
  ```sh
  set LOCAL_RUN=true
  ```
3. Set `DISCORD_TOKEN` environment variable to your Discord api token.
  ```sh
  set DISCORD_TOKEN=<YOUR_DISCORD_TOKEN>
  ```
4. Navigate to the main directory, build the go exectuable 
```sh
go build -o app 
```
5. Run the go exectuable
```sh
./app
```

![Go build locally][go-build-locally-screenshot]


### Build locally with Docker

1. Navigate to the `discord-tracking-bot` repo location.

2. Set `LOCAL_RUN` environment variable to `true` in `discord-trackingbot/Dockerfile`
  ```sh
  ENV LOCAL_RUN="false"
  ```

3. Set `DISCORD_TOKEN` environment variable to your Discord api token in `discord-trackingbot/Dockerfile`
  ```sh
  ENV DISCORD_TOKEN=<YOUR_DISCORD_TOKEN>
  ```

4. Build the go exectuable 
```sh
go build -o app 
```

5. Build the Docker image using the Dockerfile. 
```sh
docker build -t discordbot .
```

6. Run the Docker bot.
Flags:
- `-d` OPTIONAL. Runs the container in detached mode.
- `-e` OPTIONAL. Sets environment variables.
- `-v` REQUIRED. Mounts volume to Docker container. Used to obtain configs and media from local machine.
- `--name` OPTIONAL. Sets docker container name.
```sh
  docker run -d -e LOCAL_RUN=${LOCAL_RUN} -e DISCORD_TOKEN=${DISCORD_TOKEN} -v /Users/cody.kieu/github/discord-tracking-bot/pkg:/pkg -v /Users/cody.kieu/github/discord-tracking-bot/media:/media --name discordbot discordbot
```

7. Check for docker container
```sh
  docker ps -a
```

8. Check docker logs
```sh
  docker logs discordbot
```

![Go build docker locally][go-build-docker-locally-screenshot]


<p align="right">(<a href="#top">back to top</a>)</p>



## Environment Variables
The Discord Tracking Bot uses environment variables to set configuration. Use Dockerfile or set through shell console.

Environment Variable | Default Value | Description
| :--- | ---: | :---:
`AWS_REGION`  | us-east-1 | If `LOCAL_RUN` is set to `false`, this region is used to grab the AWS Secret from.
`AWS_MAX_RETIRES`  | 5 | If `LOCAL_RUN` is set to `true . Maximum number of request to set up AWS session.
`DISCORD_TOKEN_AWS_SECRET_NAME`  | discord/token | If `LOCAL_RUN` is set to `true . AWS Secret name to grab that contains Discord Token.
`DISCORD_TOKEN`  | | Allows application to configure and connect to the Discord session.
`LOCAL_RUN`  | false | Determines whether or not to create AWS session and grab AWS secret.
`LOG_TO_DISCORD`  | true | Determines whether to send logs to Discord `log` channel.
`MEDIA_PATH`  | /media/ | Path of media folder within running platform.
`METRICS_TO_DISCORD`  | true | Determines whether to send metrics to Discord `metrics` channel.
`TRACKING_CONFIG_PATH`  | /pkg/configs/tracking/ | Path of tracking configs used to call `webcrawler`.
`TRACKING_CHANNELS_DELAY`  | 21600 | Used to determine tracking delay for all channels in seconds.
`WEBCRAWLER_HOST`  | 5 | host name, used to send http request to `webcrawler`
`WEBCRAWLER_PORT`  | 5 | host port, used to send http request to `webcrawler`

## Features

* Sends metrics to metrics channel
![metrics][metrics]

* Logs to both discord and stdout
![logs][logs]

* Ability to control bot through discord `bot_console` channel
![bot-console][bot-console]


<p align="right">(<a href="#top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Cody Kieu - cody6750@gmail.com

Project Link: [https://github.com/cody6750/discord-tracking-bot](https://github.com/cody6750/discord-tracking-bot)

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/cody6750/discord-tracking-bot.svg?style=for-the-badge
[contributors-url]: https://github.com/cody6750/discord-tracking-bot/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/cody6750/discord-tracking-bot.svg?style=for-the-badge
[forks-url]: https://github.com/cody6750/discord-tracking-bot/network/members
[stars-shield]: https://img.shields.io/github/stars/cody6750/discord-tracking-bot.svg?style=for-the-badge
[stars-url]:https://github.com/cody6750/discord-tracking-bot/stargazers
[issues-shield]: https://img.shields.io/github/issues/cody6750/discord-tracking-bot.svg?style=for-the-badge
[issues-url]: https://github.com/cody6750/discord-tracking-bot/pulls
[license-shield]: https://img.shields.io/github/license/cody6750/discord-tracking-bot.svg?style=for-the-badge
[license-url]: https://github.com/cody6750/discord-tracking-bot/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/cody-kieu-a6984a162/
[product-screenshot]: media/discord.png
[tracking-screenshot]: media/tracking.png
[go-build-locally-screenshot]: media/go_build_locally.png
[go-build-docker-locally-screenshot]: media/go_build_docker_locally.png
[admin-channels]: media/admin_channels.png
[tracking-channels]: media/tracking_channels.png
[tracking-configs]: media/tracking_configs.png
[metrics]: media/metrics.png
[bot-console]: media/bot_console.png
[logs]: media/logs.png
