# Smart Spore Hub

Smart Spore Hub is a web application for managing spore data, providing an API
for sending emails.

## Usage

### Running the application

To run the application, execute the following command:

```bash
go run main.go
```

This will start the server on the host and port specified in the configuration.
By default, it runs on `0.0.0.0:8080`.

### Configuration

The application can be configured using command-line flags, environment variables,
and a `.env` file.

#### Command-Line Flags

The following command-line flags are available:

- `--port`: Port to listen on (default: `8080`)
- `--host`: Host to listen on (default: `0.0.0.0`)
- `--SMTP_HOST`: SMTP Host (default: `smtp.gmail.com`)
- `--SMTP_PORT`: SMTP Port (default: `587`)
- `--SMTP_USERNAME`: SMTP Username
- `--SMTP_PASSWORD`: SMTP Password
- `--SMTP_EMAIL`: SMTP Email

Example:

```bash
go run main.go --port 9000 --host 127.0.0.1
```

#### Environment Variables

The application also supports configuration through environment variables.
The corresponding environment variables for the flags are:

- `PORT`
- `HOST`
- `SMTP_HOST`
- `SMTP_PORT`
- `SMTP_USERNAME`
- `SMTP_PASSWORD`
- `SMTP_EMAIL`

#### .env File

You can also use a `.env` file to configure the application.
The application will automatically load the `.env` file
if it is present in the same directory as theexecutable.

Example `.env` file:

```bash
PORT=9000
HOST=127.0.0.1
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your_username
SMTP_PASSWORD=your_password
SMTP_EMAIL=your_email@example.com
```

## API Endpoints

The following API endpoints are available:

- `/health`: Health check endpoint (unprotected)
- `POST /api/v1/mailer/send`: Send email (protected, requires authentication)

### Mailer Service

The `/api/v1/mailer/send` endpoint requires authentication.
You need to configure the authentication middleware to use this endpoint.
Here's how the `/api/v1/mailer/send` endpoint works:

1. **Request:** It expects a `POST` request with a JSON body.
The JSON body should have the following format:

    ```json
    {
      "to": ["recipient1@example.com", "recipient2@example.com"],
      "subject": "Your Email Subject",
      "body": "Your email content.",
      "is_html": false
    }
    ```

    - `to`: An array of email addresses to send the email to.
    - `subject`: The subject of the email.
    - `body`: The content of the email.
    - `is_html`: A boolean value indicating whether the email body is HTML
    or plain text.

2. **Response:**  Upon successful email delivery, the endpoint returns a JSON
response with the following format:

    ```json
    {
      "message": "Email sent successfully",
      "success": true
    }
    ```

    If there is an error, the endpoint returns a JSON response with an
    error message and a `success` value of `false`.
    The HTTP status code will also indicate the type of error
    (e.g., 400 for a bad request, 500 for a server error).

    Example error response:

    ```json
    {
      "message": "Invalid request body",
      "success": false
    }
    ```

#### Whatsapp Service

The `/api/v1/whatsapp/send` endpoint requires authentication.
You need to configure the authentication middleware to use this endpoint.
Here's how the `/api/v1/whatsapp/send` endpoint works:

1. **Request:** It expects a `POST` request with a JSON body.
The JSON body should have the following format:

    ```json
    {
     "phone_number": "254712345678",
     "message": "Hello, dev!"
    }
    ```
2. **Response:** Upon successful WhatsApp message delivery, the endpoint returns a JSON
response with the following format:

    ```json
    {
      "message": "WhatsApp message sent successfully",
      "success": true
    }
    ```

    If there is an error, the endpoint returns a JSON response with an
    error message and a `success` value of `false`.
    The HTTP status code will also indicate the type of error
    (e.g., 400 for a bad request, 500 for a server error).

    Example error response:

    ```json
    {
      "message": "Failed to send WhatsApp message",
      "success": false
    }
    ```
