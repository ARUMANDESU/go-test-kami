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
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project
Test task for Kami

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

[![Go][go-shield]][go-url]    [![Docker][docker-shield]][docker-url]    [![Postgresql][postgresql-shield]][postgresql-url]    [![Taskfile][taskfile-shield]][taskfile-url]    


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



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[go-url]: https://golang.org/
[docker-url]: https://www.docker.com/
[taskfile-url]: https://taskfile.dev/
[postgresql-url]: https://www.postgresql.org/

[go-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[docker-shield]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[taskfile-shield]: https://img.shields.io/badge/Taskfile-00ADD8?style=for-the-badge&logo=go&logoColor=white
[postgresql-shield]: https://img.shields.io/badge/Postgresql-336791?style=for-the-badge&logo=postgresql&logoColor=white
