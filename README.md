<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>


<!-- PROJECT LOGO -->
<br />
<div align="center">

<h3 align="center">go-test-kami</h3>

</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
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
    <li><a href="#testing">Testing</a></li>
    <li><a href="#curl-snippets">Curl snippets</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project
Test task for Kami

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

[![Go][go-shield]][go-url]    [![Docker][docker-shield]][docker-url]    [![Postgresql][postgresql-shield]][postgresql-url]    [![Taskfile][taskfile-shield]][taskfile-url]    [![TestContainers][testcontainers-shield]][testcontainers-url]


<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- GETTING STARTED -->
## Getting Started
### Prerequisites

* Go version 1.23.0
* [Docker](https://www.docker.com/get-started)
* [Taskfile 3](https://taskfile.dev/installation/) 
```sh
go version
docker --version
task --version
```

### Installation

1. Clone the repo
    ```sh
    git clone https://github.com/ARUMANDESU/go-test-kami.git
    ```
2. Change directory
    ```sh
    cd go-test-kami
    ```
3. Write the environment variables in the `.env` file
    ```sh
    touch .env
    ```
    ```sh
    nano .env #or use your favorite text editor
    ```
    ```dotenv
    ENV=local #local, test, dev, prod
    HTTP_PORT=8080
    DATABASE_URL=postgresql://user:password@localhost:5432/dbname
    ```
4. Run the service
    ```sh
    task r
    ```
    or run on docker(with docker compose)
    ```sh
    task rc
    ```
    to compose build
    ```sh
    task rcb
    ```
  
  

## Testing

1. Coverage 
    ```sh
    task tc
    ```
2. Unit tests
    ```sh
    task tu
    ```
3. Integration tests
    ```sh
    task ti
    ```


<p align="right">(<a href="#readme-top">back to top</a>)</p>


## Curl snippets

1. Create reservation
  ```sh
  curl --location -i --request POST 'http://localhost:8080/v1/reservations' \
  --data '{
      "room_id": "018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b6d",
      "start_time": "2021-09-01T10:00:00Z",
      "end_time": "2021-09-01T11:00:00Z"
  }'
  ```
2. Get reservation
  ```sh
  curl --location -i --request GET 'http://localhost:8080/v1/reservations/018c0f7e-7b9b-7f4b-8e4b-4b5e4b5e4b6d'
  ```

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[go-url]: https://golang.org/
[docker-url]: https://www.docker.com/
[taskfile-url]: https://taskfile.dev/
[postgresql-url]: https://www.postgresql.org/
[testcontainers-url]: https://www.testcontainers.org/

[go-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[docker-shield]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[taskfile-shield]: https://img.shields.io/badge/Taskfile-00ADD8?style=for-the-badge&logo=go&logoColor=white
[postgresql-shield]: https://img.shields.io/badge/Postgresql-336791?style=for-the-badge&logo=postgresql&logoColor=white
[testcontainers-shield]: https://img.shields.io/badge/TestContainers-00ADD8?style=for-the-badge&logo=data:image/png;base64,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

