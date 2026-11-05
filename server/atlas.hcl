// Atlas (https://atlasgo.io) configuration for generating schema migrations from the GORM models.
// See the "Database Migrations" section in the README.

// The desired schema: GORM models + tools/atlasloader/legacy.sql
// The loader imports the server module via a replace directive; -mod=mod lets it follow
// server dependency updates without a separate `go mod tidy` in tools/atlasloader.
data "external_schema" "gorm" {
  program = ["go", "run", "-C", "tools/atlasloader", "-mod=mod", "."]
}

// Dev database Atlas uses to compute diffs. Override with --var dev_url=... to use an existing server.
variable "dev_url" {
  type    = string
  default = "docker://mysql/9.0.0/dev"
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = var.dev_url
  // tracking table of gormigrate, which applies the migrations at server startup
  exclude = ["migrations"]
  migration {
    dir = "file://backend/migration/sql"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
