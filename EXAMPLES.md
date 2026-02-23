# Aether Code Examples

This document showcases Aether's syntax through practical examples, from simple resources to complex multi-cloud deployments with AI agents.

---

## Example 1: Simple Web Server

Basic single-instance web server with security group.

```aether
// Simple web server deployment

resource network.vpc "main" {
  cidr = "10.0.0.0/16"
  region = "us-east"
  
  tags = {
    Name = "main-vpc"
    Environment = "production"
  }
}

resource network.subnet "public" {
  vpc = vpc.main
  cidr = "10.0.1.0/24"
  availability_zone = "us-east-1a"
  
  tags = {
    Name = "public-subnet"
  }
}

resource network.security_group "web" {
  vpc = vpc.main
  description = "Allow HTTP and HTTPS traffic"
  
  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource compute.instance "web_server" {
  machine_type = "small"
  os_image = "ubuntu-22.04"
  subnet = subnet.public
  security_groups = [sg.web]
  
  user_data = <<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    systemctl start nginx
  EOF
  
  tags = {
    Name = "web-server"
    Role = "webserver"
  }
}

output "server_ip" {
  value = instance.web_server.public_ip
  description = "Public IP address of the web server"
}
```

---

## Example 2: Variables and Script Logic

Using variables and scripting blocks for dynamic configuration.

```aether
// Configuration with variables

variable "environment" {
  type = string
  default = "development"
  description = "Deployment environment (development, staging, production)"
}

variable "instance_count" {
  type = int
  default = 2
  description = "Number of instances to create"
}

variable "enable_monitoring" {
  type = bool
  default = true
}

// Script block for dynamic configuration
script {
  // Determine machine type based on environment
  let machine_type = match var.environment {
    "production" => "large"
    "staging" => "medium"
    _ => "small"
  }
  
  // Calculate actual instance count with scaling rules
  let actual_count = var.environment == "production" 
    ? max(var.instance_count, 3)  // Minimum 3 in production
    : var.instance_count
  
  // Set region based on environment
  let region = var.environment == "production" ? "us-east" : "us-west"
  
  // Export for use in resource declarations
  export { machine_type, actual_count, region }
}

resource network.vpc "main" {
  cidr = "10.0.0.0/16"
  region = script.region
}

resource compute.instance "app" {
  count = script.actual_count
  
  machine_type = script.machine_type
  os_image = "ubuntu-22.04"
  region = script.region
  
  tags = {
    Name = "app-server-${count.index}"
    Environment = var.environment
  }
  
  // Conditional monitoring
  monitoring_enabled = var.enable_monitoring
}

output "instance_ips" {
  value = [for i in instance.app : i.public_ip]
}
```

---

## Example 3: Multi-Tier Web Application

Complete web application with load balancer, app servers, and database.

```aether
// Multi-tier web application

// Network infrastructure
resource network.vpc "app_vpc" {
  cidr = "10.0.0.0/16"
  region = "us-east"
  enable_dns = true
}

resource network.subnet "public" {
  count = 2
  vpc = vpc.app_vpc
  cidr = "10.0.${count.index}.0/24"
  availability_zone = "us-east-1${["a", "b"][count.index]}"
  public = true
}

resource network.subnet "private" {
  count = 2
  vpc = vpc.app_vpc
  cidr = "10.0.${count.index + 10}.0/24"
  availability_zone = "us-east-1${["a", "b"][count.index]}"
  public = false
}

// Security groups
resource network.security_group "alb" {
  vpc = vpc.app_vpc
  description = "Load balancer security group"
  
  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource network.security_group "app" {
  vpc = vpc.app_vpc
  description = "Application server security group"
  
  ingress {
    from_port = 8080
    to_port = 8080
    protocol = "tcp"
    security_groups = [sg.alb]
  }
}

resource network.security_group "database" {
  vpc = vpc.app_vpc
  description = "Database security group"
  
  ingress {
    from_port = 5432
    to_port = 5432
    protocol = "tcp"
    security_groups = [sg.app]
  }
}

// Load balancer
resource loadbalancer.application "web" {
  subnets = subnet.public[*]
  security_groups = [sg.alb]
  
  listener {
    port = 443
    protocol = "HTTPS"
    ssl_certificate = cert.web.arn
    
    default_action {
      type = "forward"
      target_group = target_group.app.id
    }
  }
}

resource loadbalancer.target_group "app" {
  port = 8080
  protocol = "HTTP"
  vpc = vpc.app_vpc
  
  health_check {
    path = "/health"
    interval = 30
    timeout = 5
    healthy_threshold = 2
    unhealthy_threshold = 3
  }
}

// Application servers
resource compute.instance "app" {
  count = 3
  
  machine_type = "medium"
  os_image = "ubuntu-22.04"
  subnet = subnet.private[count.index % 2]
  security_groups = [sg.app]
  
  user_data = <<-EOF
    #!/bin/bash
    # Install application
    apt-get update
    apt-get install -y docker.io
    docker run -d -p 8080:8080 \
      -e DATABASE_URL="${database.main.connection_string}" \
      myapp:latest
  EOF
  
  tags = {
    Name = "app-server-${count.index}"
    Role = "application"
  }
}

// Register instances with load balancer
resource loadbalancer.target_attachment "app" {
  count = length(instance.app)
  target_group = target_group.app.id
  target = instance.app[count.index].id
  port = 8080
}

// Database
resource database.instance "main" {
  engine = "postgres"
  engine_version = "15"
  instance_class = "medium"
  
  storage = 100
  storage_type = "ssd"
  
  database_name = "appdb"
  master_username = "admin"
  master_password = secret.db_password
  
  subnets = subnet.private[*]
  security_groups = [sg.database]
  
  backup_retention_days = 7
  backup_window = "03:00-04:00"
  maintenance_window = "sun:04:00-sun:05:00"
  
  multi_az = true
  
  tags = {
    Name = "app-database"
    Environment = "production"
  }
}

// Secrets
secret "db_password" {
  source = "vault"
  path = "secret/database/password"
}

secret "ssl_cert" {
  source = "aws_acm"
  domain = "example.com"
}

// Outputs
output "load_balancer_url" {
  value = "https://${lb.web.dns_name}"
  description = "Application URL"
}

output "database_endpoint" {
  value = database.main.endpoint
  sensitive = true
  description = "Database connection endpoint"
}
```

---

## Example 4: Using Modules

Reusable modules for common patterns.

```aether
// Using a web application module

module "frontend" {
  source = "stdlib/web-app@1.2.0"
  
  name = "frontend"
  environment = "production"
  
  instance_count = 5
  instance_type = "medium"
  
  vpc_cidr = "10.0.0.0/16"
  region = "us-east"
  
  database_enabled = true
  database_engine = "postgres"
  database_size = "large"
  
  monitoring_enabled = true
  auto_scaling = {
    min_size = 3
    max_size = 10
    target_cpu = 70
  }
}

module "backend" {
  source = "github.com/company/api-service@2.0.0"
  
  name = "api-backend"
  vpc = module.frontend.vpc
  
  instance_count = 3
  redis_enabled = true
  
  environment_variables = {
    DATABASE_URL = module.frontend.database_url
    CACHE_URL = "redis://localhost:6379"
    API_KEY = secret.api_key
  }
}

output "frontend_url" {
  value = module.frontend.load_balancer_url
}

output "api_endpoint" {
  value = module.backend.api_url
}
```

---

## Example 5: Multi-Cloud Deployment

Same application deployed to multiple cloud providers.

```aether
// Multi-cloud deployment

variable "primary_provider" {
  type = string
  default = "aws"
  description = "Primary cloud provider (aws, azure, gcp)"
}

variable "enable_failover" {
  type = bool
  default = true
  description = "Enable multi-cloud failover"
}

// AWS deployment
provider "aws" {
  region = "us-east-1"
}

module "app_aws" {
  source = "./modules/web-app"
  provider = aws
  
  name = "app-aws"
  instance_count = 5
  region = "us-east"
}

// Azure deployment (conditional failover)
provider "azure" {
  region = "eastus"
}

module "app_azure" {
  source = "./modules/web-app"
  provider = azure
  enabled = var.enable_failover
  
  name = "app-azure"
  instance_count = 3
  region = "us-east"
}

// GCP deployment (conditional failover)
provider "gcp" {
  project = "my-project"
  region = "us-east1"
}

module "app_gcp" {
  source = "./modules/web-app"
  provider = gcp
  enabled = var.enable_failover
  
  name = "app-gcp"
  instance_count = 3
  region = "us-east"
}

// Global load balancer for multi-cloud
resource loadbalancer.global "primary" {
  backends = compact([
    module.app_aws.load_balancer_url,
    var.enable_failover ? module.app_azure.load_balancer_url : null,
    var.enable_failover ? module.app_gcp.load_balancer_url : null
  ])
  
  routing_policy = "latency"
  health_check_interval = 30
}

output "global_url" {
  value = lb.primary.url
  description = "Global multi-cloud URL"
}
```

---

## Example 6: AI Agents in Action

Infrastructure with AI-driven optimization and security.

```aether
// AI-managed infrastructure

// Security analyzer agent
agent security_scanner {
  type = "analyzer"
  
  // Run before deployment
  checks = [
    "unencrypted_storage",
    "open_security_groups",
    "weak_passwords",
    "missing_backups",
    "compliance.cis_benchmark",
    "compliance.gdpr"
  ]
  
  severity_threshold = "medium"
  block_on_failure = true
  
  notification {
    slack = "#security-alerts"
    email = ["security@example.com"]
  }
}

// Cost optimization agent
agent cost_optimizer {
  type = "autonomous"
  scope = resource.compute.*
  
  goals = [
    "minimize_cost",
    "maintain_performance"
  ]
  
  constraints {
    max_cost_change = "$200/day"
    min_instances = 2
    max_instances = 20
  }
  
  approval {
    // Auto-approve small changes
    auto_approve = cost_change < "$50/day"
    
    // Require approval for larger changes
    required_for = cost_change >= "$50/day"
    notify = ["ops@example.com"]
    
    // Multi-approval for very large changes
    multi_approval {
      required = 2
      condition = cost_change >= "$500/day"
    }
  }
  
  learning_period = duration("14d")
  check_interval = duration("1h")
}

// Performance optimization agent
agent performance_tuner {
  type = "autonomous"
  scope = [resource.compute.*, resource.database.*]
  
  metrics = [
    "cpu_utilization",
    "memory_usage",
    "disk_io",
    "network_throughput",
    "request_latency"
  ]
  
  thresholds {
    high_cpu = 80
    low_cpu = 20
    target_latency_p99 = duration("200ms")
  }
  
  actions = [
    "scale_horizontally",
    "change_instance_type",
    "adjust_database_size"
  ]
  
  approval {
    auto_approve = true
    notify = ["ops@example.com"]
    rollback_on_error = true
  }
}

// Drift remediation agent
agent drift_watcher {
  type = "autonomous"
  
  check_interval = duration("15m")
  
  actions = [
    "revert_to_desired_state",
    "update_state_file",
    "alert_manual_review"
  ]
  
  approval {
    auto_approve = minor_drift == true
    required_for = major_drift == true
    notify = ["ops@example.com"]
  }
}

// Infrastructure with AI hints
resource compute.instance "web" {
  // @ai:optimize - Let AI tune this resource
  // @ai:secure - Apply security best practices
  
  count = 5
  machine_type = "medium"
  os_image = "ubuntu-22.04"
  
  auto_scaling {
    min = 3
    max = 15
    // AI agent will manage scaling within these bounds
    agent_managed = agent.cost_optimizer
  }
  
  tags = {
    Name = "web-server"
    AIManaged = "true"
  }
}

resource database.instance "main" {
  // @ai:optimize
  // @ai:backup
  
  engine = "postgres"
  instance_class = "large"
  
  // AI agent will adjust size based on usage
  agent_managed = agent.performance_tuner
  
  backup_retention_days = 30
  multi_az = true
}

// Agent monitoring output
output "agent_savings" {
  value = agent.cost_optimizer.total_savings
  description = "Total cost savings from AI optimization"
}

output "agent_actions" {
  value = agent.cost_optimizer.action_history
  description = "History of AI agent actions"
}
```

---

## Example 7: Testing Infrastructure

Comprehensive testing examples.

```aether
// Infrastructure tests

// Unit test - test resource configuration
test "web_server_has_correct_tags" {
  let server = resource.compute.instance.web_server
  
  assert server.tags["Environment"] == "production"
  assert server.tags["ManagedBy"] == "aether"
  assert length(server.tags) >= 2
}

test "security_group_blocks_ssh_from_internet" {
  let sg = resource.network.security_group.web
  
  for rule in sg.ingress {
    if rule.from_port == 22 {
      assert rule.cidr_blocks != ["0.0.0.0/0"]
    }
  }
}

test "database_has_backups_enabled" {
  let db = resource.database.instance.main
  
  assert db.backup_retention_days > 0
  assert db.backup_retention_days >= 7
}

// Property-based test - ensure invariants
property_test "all_storage_is_encrypted" {
  for resource in resources {
    if resource.type == "storage.bucket" {
      assert resource.encryption_enabled == true
    }
    if resource.type == "database.instance" {
      assert resource.storage_encrypted == true
    }
  }
}

property_test "no_open_security_groups" {
  for sg in resources where resource.type == "network.security_group" {
    for rule in sg.ingress {
      if rule.cidr_blocks == ["0.0.0.0/0"] {
        assert rule.from_port in [80, 443]  // Only HTTP/HTTPS from internet
      }
    }
  }
}

// Integration test - deploy and verify
integration_test "web_app_responds_correctly" {
  environment = "test"
  timeout = duration("10m")
  
  // Deploy infrastructure
  apply()
  
  // Wait for resources to be ready
  wait_for resource.compute.instance.web_server {
    condition = status == "running"
    timeout = duration("5m")
  }
  
  // Test HTTP endpoint
  let response = http.get("http://${output.server_ip}")
  assert response.status == 200
  assert response.body contains "Welcome"
  
  // Test database connectivity
  let db_check = database.query(
    output.database_endpoint,
    "SELECT 1"
  )
  assert db_check.success == true
  
  // Cleanup after test
  defer destroy()
}

integration_test "load_balancer_distributes_traffic" {
  environment = "test"
  
  apply()
  
  let responses = []
  for i in range(100) {
    let response = http.get(output.load_balancer_url)
    responses.append(response.headers["X-Server-ID"])
  }
  
  // Verify traffic is distributed across multiple servers
  let unique_servers = set(responses)
  assert length(unique_servers) >= 2
  
  defer destroy()
}

// Snapshot test - detect unintended changes
snapshot_test "infrastructure_plan_unchanged" {
  snapshot_file = "test/snapshots/production-plan.json"
  
  let current_plan = aether.plan()
  
  // Compare with saved snapshot
  assert_snapshot_match(current_plan, snapshot_file)
}
```

---

## Example 8: Complex Script Logic

Advanced scripting for dynamic infrastructure generation.

```aether
// Dynamic multi-region deployment with scripting

variable "regions" {
  type = list<string>
  default = ["us-east", "us-west", "eu-west"]
}

variable "services" {
  type = list<string>
  default = ["api", "web", "worker"]
}

script {
  // Generate VPC configuration for each region
  let vpcs = {}
  for region in var.regions {
    vpcs[region] = {
      cidr = "10.${index(var.regions, region)}.0.0/16"
      region = region
    }
  }
  
  // Calculate instance distribution
  function distribute_instances(total, region_count) {
    let base = int(total / region_count)
    let remainder = total % region_count
    
    let result = []
    for i in range(region_count) {
      result.append(base + (i < remainder ? 1 : 0))
    }
    return result
  }
  
  let instance_distribution = distribute_instances(15, length(var.regions))
  
  // Generate service configurations
  let services_config = {}
  for service in var.services {
    services_config[service] = {
      port = match service {
        "api" => 8080
        "web" => 80
        "worker" => null
      }
      public = service in ["api", "web"]
      instances_per_region = match service {
        "api" => 5
        "web" => 3
        "worker" => 2
      }
    }
  }
  
  export { vpcs, services_config, instance_distribution }
}

// Create VPCs dynamically
resource network.vpc "regional" {
  for_each = script.vpcs
  
  cidr = each.value.cidr
  region = each.value.region
  
  tags = {
    Name = "vpc-${each.key}"
    Region = each.key
  }
}

// Create services in each region
resource compute.instance "services" {
  for_each = {
    for combo in flatten([
      for region in var.regions : [
        for service in var.services : {
          key = "${region}-${service}"
          region = region
          service = service
        }
      ]
    ]) : combo.key => combo
  }
  
  count = script.services_config[each.value.service].instances_per_region
  
  machine_type = "medium"
  os_image = "ubuntu-22.04"
  region = each.value.region
  
  tags = {
    Name = "${each.value.service}-${each.value.region}-${count.index}"
    Service = each.value.service
    Region = each.value.region
  }
}

output "service_endpoints" {
  value = {
    for svc in var.services : svc => [
      for inst in instance.services 
        if inst.tags["Service"] == svc : inst.public_ip
    ]
  }
}
```

---

## Example 9: Kubernetes Cluster

Complete Kubernetes cluster deployment.

```aether
// Kubernetes cluster on multiple clouds

module "k8s_cluster" {
  source = "stdlib/kubernetes-cluster@2.0"
  
  name = "production"
  version = "1.28"
  
  // Control plane
  control_plane {
    instance_count = 3
    machine_type = "large"
    high_availability = true
  }
  
  // Worker nodes
  node_pools = [
    {
      name = "general"
      instance_count = 5
      machine_type = "large"
      auto_scaling = {
        min = 3
        max = 10
      }
      labels = {
        workload = "general"
      }
    },
    {
      name = "compute-intensive"
      instance_count = 2
      machine_type = "xlarge"
      taints = [
        {
          key = "workload"
          value = "compute"
          effect = "NoSchedule"
        }
      ]
    }
  ]
  
  // Networking
  network {
    vpc_cidr = "10.0.0.0/16"
    pod_cidr = "172.16.0.0/16"
    service_cidr = "172.17.0.0/16"
  }
  
  // Add-ons
  addons = {
    ingress_controller = "nginx"
    cert_manager = true
    metrics_server = true
    cluster_autoscaler = true
    prometheus = true
    grafana = true
  }
  
  region = "us-east"
}

output "kubeconfig" {
  value = module.k8s_cluster.kubeconfig
  sensitive = true
}

output "cluster_endpoint" {
  value = module.k8s_cluster.endpoint
}
```

---

## Example 10: Serverless Application

Complete serverless application with functions and event triggers.

```aether
// Serverless web application

resource storage.bucket "uploads" {
  region = "us-east"
  versioning = true
  
  lifecycle_rules = [
    {
      action = "delete"
      age_days = 90
    }
  ]
}

resource storage.bucket "processed" {
  region = "us-east"
  versioning = true
}

resource function.function "image_processor" {
  name = "image-processor"
  runtime = "python3.11"
  handler = "main.handler"
  code_path = "./functions/image-processor"
  
  memory = 1024
  timeout = 300
  
  environment_variables = {
    OUTPUT_BUCKET = bucket.processed.name
  }
  
  // Trigger on upload
  triggers = [
    {
      type = "storage"
      bucket = bucket.uploads.id
      events = ["object.created"]
      prefix = "images/"
    }
  ]
}

resource function.function "api" {
  name = "api-handler"
  runtime = "nodejs20"
  handler = "index.handler"
  code_path = "./functions/api"
  
  memory = 512
  timeout = 30
  
  environment_variables = {
    DATABASE_URL = database.main.connection_string
    BUCKET_NAME = bucket.processed.name
  }
}

resource api_gateway.rest_api "main" {
  name = "serverless-api"
  
  routes = [
    {
      path = "/images"
      method = "GET"
      function = function.api
    },
    {
      path = "/images/{id}"
      method = "GET"
      function = function.api
    },
    {
      path = "/upload"
      method = "POST"
      function = function.api
    }
  ]
}

resource database.serverless "main" {
  engine = "postgres"
  engine_version = "15"
  
  auto_scaling = {
    min_capacity = 2
    max_capacity = 16
  }
  
  auto_pause = true
  auto_pause_delay = duration("5m")
}

output "api_url" {
  value = api_gateway.main.url
}

output "upload_bucket" {
  value = bucket.uploads.name
}
```

---

These examples demonstrate Aether's power and flexibility—from simple single-resource deployments to complex multi-cloud architectures with AI-driven optimization. The hybrid declarative-scripting approach provides clarity for simple cases and power for complex scenarios.
