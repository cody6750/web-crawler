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

![Tracking Screen Shot][tracking-screenshot]

A fast high-level web crawling and web scraping application framework that includes multiple features using Go. Used to crawl websites concurrently and extract structured data from their pages. Hosted on a web server in a Docker container on AWS. This project includes a Dockerfile that is ready to go out of the box and can be hosted locally. Simply follow the instructions below to get started!
### Design
![Design][Design]


### Built With

* [Go](https://go.dev/)
* [gorilla](https://github.com/gorilla/mux)
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

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/cody6750/web-crawler
   ```

2. Get Go packages
   ```sh
	go get github.com/aws/aws-sdk-go
	go get github.com/gorilla/mux
	go get github.com/sirupsen/logrus 
   ```

<p align="right">(<a href="#top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Usage

The Web Crawler is designed to be deployed on AWS EC2 as a Docker container implementing a REST API, though it can be deployed locally by building and executing the exectuable or deploying the Docker container locally on your machine. This section will cover 2 of the ways to do so. Please note that these instructions are for Mac OS using a bash terminal.

### Build locally without Docker

1. Navigate to `webcrawler` repo location
2. Navigate to the `/web` directory, build the go exectuable . Set environment variables if you want to override anything.
```sh
go build -o app 
```
3. Run the go exectuable
```sh
./app
```

![Go build locally][go-build-locally-screenshot]


### Build locally with Docker

1. Navigate to the `discord-tracking-bot` repo location.

2. Navigate to the `/web` directory, build the go exectuable .
```sh
go build -o app 
```

5. Build the Docker image using the Dockerfile. Set environment variables if you want to override anything, in the Dockerfile.
```sh
docker build -t webcrawler .
```

6. Run the Docker image.
```sh
  docker run -d -p 9090:9090 --network=discord --name webcrawler webcrawler
```

7. Check for docker container
```sh
  docker ps -a
```

8. Check docker logs
```sh
  docker logs webcrawler
```

![Go build docker locally][go-build-docker-locally-screenshot]


<p align="right">(<a href="#top">back to top</a>)</p>


## Environment Variables
The Discord Tracking Bot uses environment variables to set configuration. Use Dockerfile or set through shell console.

Environment Variable | Default Value | Description
| :--- | ---: | :---:
`ALLOW_EMPTY_ITEM`  | false | Allows webcrawler to return scrape responses with empty items.
`AWS_WRITE_OUTPUT_TO_S3`  | false | Determines whether to write scrape responses to S3.
`AWS_MAX_RERIES`  | discord/token | If `AWS_WRITE_OUTPUT_TO_S3` is set to true, set maximum retry responses during creation of AWS session.
`AWS_REGION`  | us-east-1 | If `AWS_WRITE_OUTPUT_TO_S3` is set to true, region to configure AWS session.
`AWS_S3_BUCKET`  | webcrawler-results | If `AWS_WRITE_OUTPUT_TO_S3` is set to true, region to configure AWS session, S3 bucket to send scrape responses.
`CRAWL_DELAY`  | 5 | Delay between crawls per web scraper worker.
`HEADER_KEY`  | User-Agent | Header agent used during http request
`HEADER_VALUE`  |Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36 | Header agent value used during http request.
`LOG_LEVEL`  | INFO | Determines level of logs.
`IDLE_TIMEOUT`  |120 | Maximum amount of time to wait for the next request when keep-alives are enabled.
`MAX_DEPTH`  | 1 | Maximum crawl depth during an execution of a crawl.
`MAX_GO_ROUTINES`  | 10000 | Maximum go routines deployed during an execution of a crawl.
`MAX_VISITED_URLS`  | 20 | Maximum visited urls during an execution of a crawl.
`MAX_ITEMS_FOUND`  | 5000 | Maximum items extracted during an execution of a crawl.
`PORT`  | :9090 | Port used to expose web server.
`READ_TIMEOUT`  | 60 | Maximum duration for reading the entire request, including the body. 
`WEB_SCRAPER_WORKER_COUNT`  | 5| Path of tracking configs used to call `webcrawler`.
`WRITE_TIMEOUT`  | 60 | host name, used to send http request to `webcrawler`


## Start crawling
The web crawler is hosted on a web server that is exposed using a REST API. To call the web crawler, we must make an http request to the web server that host the web crawler.
1. Run the go exeutable locally or the web crawler docker container. Confirm that it is ready to recieve traffic.
2.  

## Features
The webcrawler includes various features:
* Ability to crawl multiple request concurrently
* REST API
* Json validation middleware
* Crawl depth restrictions
* Liveleness and readiness health checks
* Respects robots.txt
* Metrics
* Generates output files in JSON
* Sends output files to S3 bucket
* Item and URL validation
* Unit test
* Sends metrics to metrics channel



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

Project Link: [https://github.com/cody6750/web-crawler](https://github.com/cody6750/web-crawler)

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
[go-build-locally-screenshot]: media/build_locally.png
[go-build-docker-locally-screenshot]: media/build_docker.png
[admin-channels]: media/admin_channels.png
[tracking-channels]: media/tracking_channels.png
[tracking-configs]: media/tracking_configs.png
[metrics]: media/metrics.png
[bot-console]: media/bot_console.png
[logs]: media/logs.png
[Design]: media/design.png