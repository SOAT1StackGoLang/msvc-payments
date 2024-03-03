# msvc-payments

Microservice payments

## Payment Service

This project is a microservice written in Go that handles payment operations. It uses Redis as a datastore and a queue for processing payments. The service provides functionalities for creating, updating, and retrieving payments.

## Features

- **Create Payment**: This feature allows you to create a new payment. The payment is initially set to a 'pending' status and stored in the Redis datastore. It is also added to a 'pending' queue for further processing.

- **Update Payment**: This feature allows you to update the status of a payment. The payment status can be updated to 'paid' or 'failed'. Depending on the status, the payment is added to the respective queue ('paid' or 'failed') and a notification is published to the respective channel.

- **Get Payment**: This feature allows you to retrieve the details of a payment using its ID.

- **Process Payment**: This is an internal function that simulates the processing of a payment. It randomly sets the payment status to either 'paid' or 'failed', with a higher probability for 'paid'.

## Dependencies

- GoLang
- Redis

## How to Run

This project should not be run as a standalone service. It is part of a larger project that includes multiple microservices. Please refer to the main project for instructions on how to run all the microservices together.

[Tech-Challenge](https://github.com/SOAT1StackGoLang/tech-challenge/tree/main/microservices)

## Running the Service for Development

To run this project you need to have GoLang and Redis installed on your machine. You can install GoLang from [here](https://golang.org/doc/install) and Redis from [here](https://redis.io/download).

Example of launch.json from vscode for debugging:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/server",
            "env": {
                "KVSTORE_HOST": "localhost",
                "APP_LOG_LEVEL": "Debug"
            },
        }
    ]
}
```

Run without vscode:

```bash
go mod download
go run cmd/server/*.go
```

## API Endpoints (Simplified)

- **Create Payment**
  - Endpoint: `POST /payments`
  - Description: Creates a new payment.
  - Request body: A JSON object with the payment details (`CreatePaymentRequest`).

    ```json
    {
      "payment": {
        "ID": "<UUID>",
        "CreatedAt": "<time>",
        "UpdatedAt": "<time>",
        "Price": "<decimal>",
        "OrderID": "<UUID>",
        "Status": "<PaymentStatus>"
      }
    }
    ```

  - Response: A JSON object with the created payment's details (`CreatePaymentResponse`).

    ```json
    {
      "payment_id": "<UUID>",
      "status": "<PaymentStatus>"
    }
    ```

- **Get Payment**
  - Endpoint: `GET /payments/{payment_id}`
  - Description: Retrieves the details of a payment.
  - Request body: None.
  - Response: A JSON object with the payment's details (`GetPaymentResponse`).

    ```json
    {
      "payment": {
        "ID": "<UUID>",
        "CreatedAt": "<time>",
        "UpdatedAt": "<time>",
        "Price": "<decimal>",
        "OrderID": "<UUID>",
        "Status": "<PaymentStatus>"
      },
      "status": "<PaymentStatus>",
      "payment_error": "<string>"
    }
    ```

- **Update Payment**
  - Endpoint: `PUT /payments/{payment_id}`
  - Description: Updates the status of a payment.
  - Request body: A JSON object with the new payment status (`UpdatePaymentRequest`).

    ```json
    {
      "payment_id": "<UUID>",
      "payment_status": "<PaymentStatus>"
    }
    ```

  - Response: A JSON object with the updated payment's details (`UpdatePaymentResponse`).

    ```json
    {
      "payment_id": "<UUID>",
      "status": "<PaymentStatus>",
      "payment_error": "<string>"
    }
    ```

Please replace the request and response details with the correct ones for your service.

Please note that this is a simplified explanation of the project. For detailed information, please refer to the source code.
