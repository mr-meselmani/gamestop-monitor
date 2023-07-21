# Gamestop Monitor

The **Gamestop Monitor** is a simple Go application that periodically monitors the availability of products on the Gamestop website. It collects data from the specified URL and evaluates the availability of products to notify users about any changes in product availability. This project utilizes the external dependency, [goquery](https://github.com/PuerkitoBio/goquery), to scrape and parse HTML data.

## Features

- Scrapes Gamestop website data from a specific URL.
- Evaluates the availability of products and generates messages if products become available.
- Periodically performs monitoring runs to keep track of changes.
- Uses Go's context package to manage concurrency and potential cancellations.

## Prerequisites

Before running the Gamestop Monitor, ensure you have the following installed:

- Go (Golang) - version 1.16 or later
- Internet connection to access the Gamestop website

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/gamestop-monitor.git
   cd gamestop-monitor
   ```

2. Fetch and download the dependencies:

   ```bash
   go mod download
   ```

## Usage

To start monitoring Gamestop products, run the main application:

```bash
go run cmd/main.go
```

The application will begin monitoring the Gamestop URL specified in the `main.go` file. It will periodically check for product availability and generate messages for newly available products. The monitoring runs for ten iterations and then stops, but you can customize the run duration or limit as needed.

## Customization

1. URL: The default URL to monitor is set in the `main.go` file. If you want to monitor a different Gamestop URL, modify the `url` parameter in the `main.go` file.

2. Run Duration: You can adjust the duration between monitoring runs by modifying the `time.Sleep()` value in the `main.go` file.


## License

This project is licensed under the [MIT License](LICENSE).

---

Thank you for using the Gamestop Monitor! We hope this tool helps you stay informed about the availability of your favorite Gamestop products. If you have any questions or feedback, please don't hesitate to reach out.

Happy monitoring! 🚀
