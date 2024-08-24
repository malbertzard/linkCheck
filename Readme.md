# LinkCheck

This Go project is a command-line tool that crawls a given website, checks for broken links (HTTP status 4xx or 5xx), and reports where those broken links were found. The tool export the results to a CSV file.

## Todo
-[ ] Add concurrency to make it go faster

## Features

- **Crawls a website** starting from a given URL.
- **Optional CSV export** of the results.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/malbertzard/linkCheck.git
   cd linkCheck
   ```

2. **Build the application:**

   You can use the provided Makefile for building the project:

   ```bash
   make build
   ```

3. **Run the application:**

   ```bash
   ./linkCheck -url "https://example.com"
   ```

## Usage

The application can be run with the following command-line arguments:

```bash
./linkCheck -url "https://example.com" -output "broken_links.csv" -concurrency 20
```

### Command-Line Options

- `-url` **(required)**: The starting URL for crawling the website.
- `-output` **(optional)**: The file path to export broken links as a CSV. If not provided, results will be printed to the console.

### Examples

1. **Basic Usage**:
   ```bash
   ./linkCheck -url "https://example.com"
	   ```

2. **Exporting Results to CSV**:
   ```bash
   ./linkCheck -url "https://example.com" -output "broken_links.csv"
   ```

## Output

### Console Output

If you do not specify the `-output` flag, the results will be printed to the console:

```plaintext
Broken links found:
https://example.com/broken-link: 404 (found at: https://example.com/page1)
https://example.com/another-broken-link: 500 (found at: https://example.com/page2)
```

### CSV Output

If you specify the `-output` flag, the results will be exported as a CSV file with the following structure:

```csv
Link,Status Code,Found At
https://example.com/broken-link,404,https://example.com/page1
https://example.com/another-broken-link,500,https://example.com/page2
```

## Development

### Dependencies

This project requires Go to be installed on your machine. All dependencies are part of the Go standard library except for the `golang.org/x/net/html` package, which is used for parsing HTML.

### Building from Source

You can build the project from source using the provided `Makefile`:

Should be pretty clear
