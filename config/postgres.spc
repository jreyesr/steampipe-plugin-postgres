connection "postgres" {
  plugin = "jreyesr/postgres"

  # Write a connection string, in the form that is expected by the pgx package:
  # https://pkg.go.dev/github.com/jackc/pgx/v5#hdr-Establishing_a_Connection
  # Required
  # Can also be set with the `DATABASE_URL` environment variable
  # connection_string = "postgres://username:password@localhost:5432/database_name"

  # The remote DB's schema that this plugin will expose
  # If you leave this unset, it'll default to `public`
  # schema = "public"

  # List of tables that will be exposed from the remote DB.
  # No dynamic tables will be created if this arg is empty or not set.
  # Wildcard based searches are supported.
  # For example:
  #  - "*" will expose every table in the remote DB
  #  - "auth-*" will expose tables whose names start with "auth-"
  #  - "users" will only expose the specific table "users"
  # You can have several items (for example, ["auth-*", "users"] will expose 
  # all the tables that start with "auth-", PLUS the table "users")
  # Defaults to all custom resources
  tables_to_expose = ["*"]
}