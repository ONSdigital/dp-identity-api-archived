job "dp-identity-api" {
  datacenters = ["eu-west-1"]
  region      = "eu"
  type        = "service"

  group "web" {
    count = "{{WEB_TASK_COUNT}}"

    constraint {
      attribute = "${node.class}"
      value     = "web"
    }

    task "dp-identity-api-web" {
      driver = "exec"

      artifact {
        source = "s3::https://s3-eu-west-1.amazonaws.com/{{BUILD_BUCKET}}/dp-identity-api/{{REVISION}}.tar.gz"
      }

      artifact {
        source = "s3::https://s3-eu-west-1.amazonaws.com/{{DEPLOYMENT_BUCKET}}/dp-identity-api/{{REVISION}}.tar.gz"
      }

      config {
        command = "${NOMAD_TASK_DIR}/start-task"

        args    = [
          "${NOMAD_TASK_DIR}/dp-identity-api",
        ]
      }

      service {
        name = "dp-identity-api"
        port = "http"
        tags = ["web"]
        check {
          type     = "http"
          path     = "/healthcheck"
          interval = "10s"
          timeout  = "2s"
        }
      }

      resources {
        cpu    = "{{WEB_RESOURCE_CPU}}"
        memory = "{{WEB_RESOURCE_MEM}}"

        network {
          port "http" {}
        }
      }

      template {
        source      = "${NOMAD_TASK_DIR}/vars-template"
        destination = "${NOMAD_TASK_DIR}/vars"
      }

      vault {
        policies = ["dp-identity-api-web"]
      }
    }
  }

  group "publishing" {
    count = "{{PUBLISHING_TASK_COUNT}}"

    constraint {
      attribute = "${node.class}"
      value     = "publishing"
    }

    task "dp-identity-api-publishing" {
      driver = "exec"

      artifact {
        source = "s3::https://s3-eu-west-1.amazonaws.com/{{BUILD_BUCKET}}/dp-identity-api/{{REVISION}}.tar.gz"
      }

      artifact {
        source = "s3::https://s3-eu-west-1.amazonaws.com/{{DEPLOYMENT_BUCKET}}/dp-identity-api/{{REVISION}}.tar.gz"
      }

      config {
        command = "${NOMAD_TASK_DIR}/start-task"

        args    = [
          "${NOMAD_TASK_DIR}/dp-identity-api",
        ]
      }

      service {
        name = "dp-identity-api"
        port = "http"
        tags = ["publishing"]
        check {
          type     = "http"
          path     = "/healthcheck"
          interval = "10s"
          timeout  = "2s"
        }
      }

      resources {
        cpu    = "{{PUBLISHING_RESOURCE_CPU}}"
        memory = "{{PUBLISHING_RESOURCE_MEM}}"

        network {
          port "http" {}
        }
      }

      template {
        source      = "${NOMAD_TASK_DIR}/vars-template"
        destination = "${NOMAD_TASK_DIR}/vars"
      }

      vault {
        policies = ["dp-identity-api-publishing"]
      }
    }
  }

}