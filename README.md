# Web Analyzer

This application analyzes web pages, extracting key information like titles, headings, links, and more. It's built with Go and offers both local development and a hosted instance for easy testing.

## Table of Contents

- [Features](##features)
- [Building and Running](#building-and-running)
- [Project Structure](#project-structure)
- [Design Processes and Decisions](#design-processes-and-decisions)
- [Suggested Future Improvements](#suggested-future-improvements)

## Features

The Web Analyzer provides the following insights about a given web page:

- **Page Title:** Extracts the main title of the page.
- **HTML Version:**  Determines the version of HTML used (e.g., HTML5).
- **Headlines:**  Lists all headings (H1 to H6) found on the page.
- **Links:**  Extracts both internal and external links present in the HTML.
- **Inaccessible Links:**  Identifies and counts links that are currently unreachable.
- **Login Form Detection:**  Detects the presence of a login form on the page.

## Building and running

### Quick Start (Hosted Instance)

Try the Web Analyzer without installing anything: [http://154.12.255.134:8080/](http://154.12.255.134:8080/)

### Prerequisites

- Go (version 1.20 or higher)
- Docker (optional, for containerized development)

### Local Development

### Clone the Repository

```
git clone https://github.com/your-username/web-analyzer.git
cd web-analyzer
```

### Run with Go (Development mode)

```bash
go run ./cmd/
```

### Build and Run Executable

```bash
go build ./cmd/
./web-analyzer
```
### Docker (Optional)

```bash
docker-compose build
docker-compose up
```
Now you can access the application in your web browser at http://localhost:8080/.

### Running the tests

```
go test -v ./...
```

## Project structure

The project structure is organized as follows:


```
├── cmd
│   ├── config.go
│   └── main.go
├── config
│   └── config.yml
├── internal
│   ├── pageanalyzer
│   │   ├── htmlextract
│   │   ├── pageanalayzer.go
│   │   └── pageanalyzer_test.go
│   ├── pagedownloader
│   │   └── pagedownloader.go
│   └── server
│       ├── analzerhandler.go
│       ├── middelware.go
│       └── server.go
├── pkg
│   └── slicetools
│       ├── filter.go
│       └── filter_test.go
└── templates
    ├── form.html
    └── results.html
├── go.mod
├── go.sum
```


- `cmd`: Contains the main package, which initializes the logger and configuration, starts the server, and handles graceful shutdown.
- `config`: Holds the application configuration file.
- `internal`: Includes the core packages of the web application.
  - `pageanalyzer`: The heart of the application, responsible for analyzing the HTML content and processing the results from the extractors.
    - `htmlextract`: Provides simple extractors that use goquery to extract information from the HTML based on specific tags and attributes.
  - `pagedownloader`: Handles fetching the HTML content of a web page given a URL.
  - `server`: Contains the HTTP router, handlers, and middleware.
- `pkg`: Holds utility package used across the application.
  - `slicetools`: Offers helpful functions for working with slices.
- templates: Stores the HTML templates for the web application.


## Design Processes and Decisions

In the development of the Web Analyzer application, the focus was on creating a modular and maintainable codebase. The following key design processes and decisions were made during the implementation:

### Decoupling Page Analyzer Logic

The page analyzer logic is decoupled from the rest of the application and treated as an internal service. It is placed in a separate package to ensure loose coupling and facilitate future modifications without impacting other parts of the system. This separation of concerns allows for easier maintenance and extensibility.

### Separating Page Downloader Logic

The page downloader logic, responsible for fetching the HTML content of a web page, is separated from the page analyzer. This decision was made considering that downloading is an I/O operation, and in the future, concurrent downloading can be implemented within the package for performance improvements. By keeping the downloader logic separate, changes can be made without affecting the analyzer.

### Handler Responsibility

The handler is responsible for handling the incoming requests and generating the appropriate responses. It relies on the page downloader and page analyzer to perform their respective tasks. Any changes needed in the logic of the downloader or analyzer can be made within their own packages, following the principle of encapsulation and modularity.

### Internal Link Definition

In the current implementation, an internal link is defined as a link with exactly the same host as the original URL. This decision was made to simplify the logic and provide a clear distinction between internal and external links. For example, if the entered URL is `https://example.com`, then `https://example.com/temp` and `https://www.example.com/temp` are considered internal links, while `https://example.org` and `https://sub.example.com` are considered external links.

### Concurrent Inaccessible Link Counting

Initially, the logic for counting inaccessible links was implemented sequentially. However, it was noticed that it was slow for certain URLs. To improve performance, the logic was made concurrent using Goroutines and simple locks. While this adds a bit of complexity to the codebase, it leverages one of Go's strengths, concurrency, to enhance the performance of this I/O-bound task.

### Static HTML Rendering

During development, it was observed that for some URLs, the HTML rendering is dynamic, and the simple downloader only retrieves the initial static HTML. To handle such URLs properly, a page downloader that uses a headless browser (like [chromedp](https://github.com/chromedp/chromedp)) would be needed to fetch the complete HTML. However, for the sake of simplicity, the decision was made to stick with the simple webpage downloader for now. In the future, adding a more advanced downloader alongside the simple one and providing a configuration option to choose which downloader to use based on the specific requirements can be considered.

### Error Handling

When the page downloader encounters an `ErrNotFound` error, it is passed to the handler layer. The handler checks if it should return a 404 Not Found status code to the user with a "Page Not Found" error or handle it differently. If there are any other errors from the analyzer or downloader, a 500 Internal Server Error is returned. For future improvements, checking for more specific errors in the analyzer and downloader and passing them to the handler to provide more informative error messages to the user can be considered.

### Configuration and Template Paths

Currently, the paths to the template files and configuration file are hardcoded for simplicity. However, in the future, command-line flags can be used to allow passing the paths to the binary file, making it more flexible and configurable.

These design processes and decisions were made to strike a balance between simplicity, performance, and maintainability. As the application evolves, there is room for further improvements and enhancements based on the specific requirements and feedback received.

## Suggested Future Improvements

Here are some suggested future improvements for the Web Analyzer application:

- Add a cache layer to cache the results for a URL, improving performance for repeated requests.
- Introduce an abstraction layer to the extractor logic, allowing for easier addition of new extract queries.
- Implement concurrent page downloading in the `pagedownloader` package to enhance performance.
- Provide separate static and dynamic page downloaders and choose the appropriate one based on configuration.
- Handle more specific status codes when passing errors from the downloader to the handler for better error reporting.
- Expand test coverage by adding tests for the `pageanalyzer` and `pagedownloader` packages, utilizing test HTTP clients.
- Accept template and configuration file paths as command-line arguments to the binary, instead of hardcoding them in the code.
- Enhance the user interface and create more sophisticated templates for an improved user experience.

These improvements focus on performance optimization, extensibility, error handling, testing, configurability, and user experience. By implementing these enhancements, the Web Analyzer application can become more robust, efficient, and user-friendly.

