# Sankaku Viewer

## Project Overview

The Sankaku Viewer is a simple HTTP server written in Golang that allows users to input a URL from the Sankaku Complex and receive a content from the provided URL. This also enables users to easily generate link previews when sharing Sankaku Viewer URL to social media platforms like Discord or Telegram.

## How It Works

1. Users visit the server's web page.
2. On the web page, there is an input text field where users can enter the URL of a Sankaku Complex post.
3. When the user submits the form, the server makes a request to the Sankaku API with the ID from the provided URL.
4. The Sankaku API returns the requested image URL or video URL.
5. The server then generates an HTML page with the image or video content. The HTML page contains Open Graph protocol and Twitter Card to generate a link preview using the received information and displays it to the user.

## Prerequisites

- Golang should be installed on your system. You can download it from the official Golang website: https://go.dev/dl/

## Installation

1. Clone this repository to your local machine:

```bash
git clone https://github.com/DevonTM/sankaku-viewer.git
cd sankaku-viewer
```

2. Build the binary using the build.sh script:

```bash
# Replace [os] and [arch] with your desired target OS and architecture.
./build.sh [os] [arch]
```

3. Start the server:

```bash
# Replace [flags] with your desired flags.
./sankaku-viewer [flags]
```

## Flags

- `-l`: Listen address. Defaults to `:8000`. Support unix sockets by using the `unix:` prefix.
- `-p`: Proxy address.
- `--user`: Sankaku username
- `--pass`: Sankaku password

## Contributing

Contributions to this project are welcome! If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
