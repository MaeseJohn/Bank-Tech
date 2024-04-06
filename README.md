# API Description

This API simplifies the process of purchasing and selling invoices between issuers and investors.

Issuers can create detailed invoices, including information such as price and due date. Investors, on the other hand, have the opportunity to acquire portions of these invoices. Once the purchase is completed (i.e., the entire invoice is acquired), it enters a state of manual approval. If the invoice is approved, the money is transferred to the corresponding issuer. However, if the invoice expires or an error occurs during the approval process, it is canceled, and the capital is returned to the investors.

## API Features:

### Public Access Endpoints:

* **GET /users :** Returns a complete list of all users.
* **GET /users/:email :** Retrieves user information associated with the provided email.
* **POST /users/login :** Allows users to log in by sending their credentials through the request body. This endpoint also returns a JSON Web Token (JWT).


### Issuer Exclusive Endpoints:

* **POST /invoice :** Enables an issuer to create a new invoice. It must provide details such as the invoice name, price, and due date. Other details will be automatically filled in.
* **POST /invoice/:id :** Allows manual approval of invoices in a pending state. This endpoint includes an additional middleware to verify that the issuer making the request is the legitimate owner of the invoice.

### Investor Exclusive Endpoints:

* **GET /invoices :** Returns a list of all available invoices.
* **POST /invoices :** Allows investors to acquire an invoice by providing the invoice ID and the desired amount of money through the request body. This action can only be performed if the invoice is open, the investor has sufficient funds, and the sum of the provided funds plus those accumulated from other investors does not exceed the total invoice price.

_**Note:** The investor and issuer exclusive endpoints are protected by two middlewares designed to verify user authentication through JWT and their status as an issuer or investor._
