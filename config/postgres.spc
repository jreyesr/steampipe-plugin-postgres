connection "postgres" {
  plugin = "jreyesr/postgres"

  # Write a connection string, in the form that is expected by the pgx package:
  # https://pkg.go.dev/github.com/jackc/pgx/v5#hdr-Establishing_a_Connection
  # Required
  # connection_string = "postgres://username:password@localhost:5432/database_name"

  # The remote DB's schema that this plugin will expose
  # If you leave this unset, it'll default to `public`
  # schema = "public"
}