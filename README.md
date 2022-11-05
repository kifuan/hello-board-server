# Hello Board

Backend server for [Hello Board](https://github.com/kifuan/hello-board).

# Requirements

+ Go 19+

+ MySQL 5+

  You'd better use MySQL 8.

# Installation

1. Clone this repository

   ```bash
git clone https://github.com/kifuan/hello-board-server.git
   ```
   
2. Prepare `.env`

   A template `.env` is already provided. Just copy it to `.env`:

   ```bash
   cp .env.template .env
   ```

   Then you should edit it. Here are explanations for some variables.

   + `DSN`: data source name. See the [go-sql-driver repo](https://github.com/go-sql-driver/mysql#dsn-data-source-name) for more information.
   + `MAIL_ACCOUNT`: it will be used for sending notifications.
   + `MAIL_PASSWORD`: or authorization code. You can get the code on any email platform.
   + `OWNER_NOTICE`: `true` or `false`. It means whether or not you will receive a notice when someone posts a new message on your site.
   + `OWNER_EMAIL`: it cannot be used by users. You should enter `OWNER_SECRET` when posting a message.
   + `OWNER_SECRET`: it should be an email-like complex string. The **regex** we use to test in the front end is `/.+@.+\..+/`. Any string which follows this can be `OWNER_SECRET`.
   + `UNSUBSCRIBE_SALT`: it should be a complex string. Different from `OWNER_SECRET`, it can be any string you want. We use it to calculate `md5` for the unsubscribe key. 

3. Prepare your database

   Create your database and source `init.sql`:

   ```sql
   CREATE DATABASE YOUR_NAME;
   USE YOUR_NAME;
   SOURCE init.sql;
   ```

4. Simply build and deploy

   ```bash
   go get
   go build
   ./hello-board-server
   ```

   You can use any deployment tool you like.

# License

This project is licensed under the MIT License.

