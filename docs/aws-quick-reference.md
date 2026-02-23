# AWS Quick Reference - Aether

Quick reference guide for common AWS resources in Aether.

## Provider Configuration

```aether
provider aws {
  region = "us-east-1"
  profile = "default"  // AWS CLI profile
  
  default_tags = {
    ManagedBy = "Aether"
    Environment = "production"
  }
}
```

## VPC & Networking

### VPC
```aether
resource network.vpc "main" {
  cidr = "10.0.0.0/16"
  enable_dns = true
  enable_dns_hostnames = true
}
```

### Subnet
```aether
resource network.subnet "public" {
  vpc = vpc.main
  cidr = "10.0.1.0/24"
  availability_zone = "us-east-1a"
  public = true
}
```

### Security Group
```aether
resource network.security_group "web" {
  vpc = vpc.main
  description = "Web security group"
  
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
```

### NAT Gateway
```aether
resource network.nat_gateway "main" {
  subnet = subnet.public
}
```

### Internet Gateway
```aether
resource network.internet_gateway "main" {
  vpc = vpc.main
}
```

## Compute

### EC2 Instance
```aether
resource compute.instance "web" {
  machine_type = "medium"  // t3.medium
  os_image = "ubuntu-22.04"
  
  network {
    subnet = subnet.public
    security_groups = [sg.web]
    public_ip = true
  }
  
  user_data = <<-EOF
    #!/bin/bash
    echo "Hello World"
  EOF
}
```

### Auto Scaling Group
```aether
resource compute.auto_scaling_group "app" {
  launch_template = lt.app
  
  min_size = 2
  max_size = 10
  desired_capacity = 3
  
  vpc_zone_identifier = [subnet.private_1a, subnet.private_1b]
  
  scaling_policy {
    name = "cpu-scaling"
    policy_type = "TargetTrackingScaling"
    target_tracking_configuration {
      predefined_metric_type = "ASGAverageCPUUtilization"
      target_value = 70.0
    }
  }
}
```

### Lambda Function
```aether
resource compute.lambda_function "api_handler" {
  function_name = "api-handler"
  runtime = "python3.11"
  handler = "index.handler"
  
  source_code = file("lambda/handler.py")
  
  role = iam_role.lambda_exec.arn
  
  memory = 512
  timeout = 30
  
  environment {
    variables = {
      TABLE_NAME = "users"
    }
  }
}
```

### ECS Fargate Service
```aether
resource compute.ecs_service "app" {
  name = "app-service"
  cluster = ecs_cluster.main.id
  task_definition = ecs_task_definition.app.arn
  
  desired_count = 3
  launch_type = "FARGATE"
  
  network_configuration {
    subnets = [subnet.private_1a, subnet.private_1b]
    security_groups = [sg.app]
  }
  
  load_balancer {
    target_group_arn = tg.app.arn
    container_name = "app"
    container_port = 8080
  }
}
```

## Storage

### S3 Bucket
```aether
resource storage.bucket "data" {
  bucket_name = "my-bucket"
  region = "us-east-1"
  
  versioning = true
  encryption_enabled = true
  
  lifecycle_rules = [
    {
      action = "transition"
      age_days = 30
      storage_class = "STANDARD_IA"
    },
    {
      action = "delete"
      age_days = 365
    }
  ]
}
```

### EBS Volume
```aether
resource storage.volume "data" {
  size = 100  // GB
  type = "ssd"  // gp3
  availability_zone = "us-east-1a"
  encrypted = true
}
```

## Database

### RDS PostgreSQL
```aether
resource database.instance "main" {
  engine = "postgres"
  engine_version = "15.3"
  
  instance_class = "large"  // db.t3.large
  
  storage = 100
  storage_type = "ssd"
  storage_encrypted = true
  
  database_name = "mydb"
  master_username = "admin"
  master_password = secret.db_password
  
  multi_az = true
  
  vpc = vpc.main
  subnet_group = db_subnet_group.main
  security_groups = [sg.database]
  
  backup_retention_days = 7
}
```

### Aurora Cluster
```aether
resource database.aurora_cluster "main" {
  cluster_identifier = "aurora-cluster"
  
  engine = "aurora-postgresql"
  engine_version = "15.3"
  
  database_name = "mydb"
  master_username = "admin"
  master_password = secret.db_password
  
  availability_zones = ["us-east-1a", "us-east-1b", "us-east-1c"]
  
  backup_retention_period = 14
  storage_encrypted = true
}
```

### DynamoDB Table
```aether
resource database.dynamodb_table "users" {
  name = "users"
  billing_mode = "PAY_PER_REQUEST"
  
  hash_key = "userId"
  range_key = "timestamp"
  
  attribute {
    name = "userId"
    type = "S"
  }
  
  attribute {
    name = "timestamp"
    type = "N"
  }
  
  global_secondary_index {
    name = "EmailIndex"
    hash_key = "email"
    projection_type = "ALL"
  }
}
```

### ElastiCache Redis
```aether
resource cache.replication_group "redis" {
  replication_group_id = "redis-cluster"
  
  engine = "redis"
  engine_version = "7.0"
  
  node_type = "cache.r6g.large"
  num_cache_clusters = 3
  
  subnet_group = cache_subnet_group.main
  security_groups = [sg.cache]
  
  at_rest_encryption_enabled = true
  transit_encryption_enabled = true
  
  automatic_failover_enabled = true
}
```

## Load Balancing

### Application Load Balancer
```aether
resource loadbalancer.application "web" {
  name = "web-alb"
  
  vpc = vpc.main
  subnets = [subnet.public_1a, subnet.public_1b]
  security_groups = [sg.alb]
  
  internal = false
}
```

### Target Group
```aether
resource loadbalancer.target_group "app" {
  name = "app-tg"
  
  port = 8080
  protocol = "HTTP"
  vpc = vpc.main
  
  health_check {
    path = "/health"
    healthy_threshold = 2
    unhealthy_threshold = 3
    timeout = 5
    interval = 30
    matcher = "200"
  }
}
```

### Listener
```aether
resource loadbalancer.listener "https" {
  loadbalancer = lb.web
  port = 443
  protocol = "HTTPS"
  
  ssl_policy = "ELBSecurityPolicy-TLS-1-2-2017-01"
  certificate_arn = acm_certificate.main.arn
  
  default_action {
    type = "forward"
    target_group = tg.app
  }
}
```

## IAM

### IAM Role
```aether
resource iam.role "ec2" {
  name = "ec2-role"
  
  assume_role_policy = <<-JSON
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "ec2.amazonaws.com"
        },
        "Effect": "Allow"
      }
    ]
  }
  JSON
}
```

### IAM Policy
```aether
resource iam.policy "s3_access" {
  name = "s3-access-policy"
  
  policy = <<-JSON
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "s3:GetObject",
          "s3:PutObject"
        ],
        "Resource": "${bucket.data.arn}/*"
      }
    ]
  }
  JSON
}
```

### Policy Attachment
```aether
resource iam.role_policy_attachment "ec2_s3" {
  role = role.ec2
  policy_arn = policy.s3_access.arn
}
```

## DNS

### Route53 Hosted Zone
```aether
resource dns.hosted_zone "main" {
  name = "example.com"
}
```

### DNS Record
```aether
resource dns.record "www" {
  zone = hosted_zone.main
  name = "www.example.com"
  type = "A"
  
  alias {
    name = lb.web.dns_name
    zone_id = lb.web.zone_id
    evaluate_target_health = true
  }
}
```

## Monitoring

### CloudWatch Alarm
```aether
resource monitoring.alarm "high_cpu" {
  alarm_name = "high-cpu-alarm"
  alarm_description = "Alert when CPU exceeds 80%"
  
  namespace = "AWS/EC2"
  metric_name = "CPUUtilization"
  
  dimensions = {
    InstanceId = instance.web.id
  }
  
  statistic = "Average"
  period = 300
  evaluation_periods = 2
  threshold = 80
  comparison_operator = "GreaterThanThreshold"
  
  alarm_actions = [sns_topic.alerts.arn]
}
```

### Log Group
```aether
resource monitoring.log_group "app" {
  name = "/aws/app/logs"
  retention_in_days = 30
  
  kms_key_id = kms_key.logs.arn
}
```

### Dashboard
```aether
resource monitoring.dashboard "main" {
  dashboard_name = "app-dashboard"
  
  dashboard_body = <<-JSON
  {
    "widgets": [
      {
        "type": "metric",
        "properties": {
          "metrics": [
            ["AWS/EC2", "CPUUtilization"]
          ],
          "period": 300,
          "stat": "Average",
          "region": "us-east-1",
          "title": "CPU Utilization"
        }
      }
    ]
  }
  JSON
}
```

## Secrets Manager

### Secret
```aether
resource secrets.secret "db_password" {
  name = "db-password"
  description = "Database password"
  
  generate_secret_string {
    password_length = 32
    exclude_characters = "/@\" "
  }
  
  rotation_rules {
    automatically_after_days = 30
  }
}
```

## API Gateway

### HTTP API
```aether
resource api.gateway "main" {
  name = "my-api"
  protocol_type = "HTTP"
  
  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "POST", "PUT", "DELETE"]
    allow_headers = ["*"]
    max_age = 300
  }
}
```

### API Route
```aether
resource api.gateway_route "get_user" {
  api_id = api_gateway.main.id
  route_key = "GET /users/{userId}"
  target = "integrations/${api_integration.get_user.id}"
}
```

## CloudFront

### Distribution
```aether
resource cdn.distribution "main" {
  enabled = true
  
  origin {
    domain_name = bucket.assets.bucket_regional_domain_name
    origin_id = "S3-Assets"
  }
  
  default_cache_behavior {
    target_origin_id = "S3-Assets"
    viewer_protocol_policy = "redirect-to-https"
    
    allowed_methods = ["GET", "HEAD", "OPTIONS"]
    cached_methods = ["GET", "HEAD"]
  }
  
  viewer_certificate {
    acm_certificate_arn = acm_certificate.main.arn
    ssl_support_method = "sni-only"
  }
}
```

## Instance Types

### Common Machine Type Mappings

| Aether Type | AWS Instance | vCPUs | RAM | Use Case |
|-------------|--------------|-------|-----|----------|
| small | t3.small | 2 | 2 GB | Dev/test |
| medium | t3.medium | 2 | 4 GB | Web servers |
| large | t3.large | 2 | 8 GB | App servers |
| xlarge | t3.xlarge | 4 | 16 GB | Databases |
| 2xlarge | t3.2xlarge | 8 | 32 GB | Heavy workloads |

### RDS Instance Classes

| Aether Class | AWS Class | vCPUs | RAM |
|--------------|-----------|-------|-----|
| small | db.t3.small | 2 | 2 GB |
| medium | db.t3.medium | 2 | 4 GB |
| large | db.t3.large | 2 | 8 GB |
| xlarge | db.r6g.xlarge | 4 | 32 GB |

## Common Patterns

### High Availability Pattern
```aether
// Deploy across 3 AZs
availability_zones = ["us-east-1a", "us-east-1b", "us-east-1c"]

// Multiple NAT gateways (one per AZ)
resource network.nat_gateway "nat" {
  count = 3
  subnet = subnet.public[count.index]
}

// Multi-AZ RDS
multi_az = true

// Auto Scaling with minimum 2 instances
min_size = 2
```

### Security Best Practices
```aether
// Enable encryption
storage_encrypted = true

// Use Secrets Manager
master_password = secret.db_password

// Require IMDSv2
metadata_options {
  http_tokens = "required"
}

// Enable deletion protection
deletion_protection = true

// Private subnets for data
publicly_accessible = false
```

### Cost Optimization
```aether
// Spot instances for batch workloads
instance_market_options {
  market_type = "spot"
}

// S3 lifecycle policies
lifecycle_rules = [
  {
    action = "transition"
    age_days = 30
    storage_class = "STANDARD_IA"
  }
]

// Auto Scaling based on demand
scaling_policy {
  policy_type = "TargetTrackingScaling"
  target_value = 70.0
}
```

## Environment Variables

Useful AWS environment variables:

```bash
# AWS Credentials
export AWS_ACCESS_KEY_ID="your-access-key"
export AWS_SECRET_ACCESS_KEY="your-secret-key"
export AWS_DEFAULT_REGION="us-east-1"

# Use specific profile
export AWS_PROFILE="production"

# Session token (for temporary credentials)
export AWS_SESSION_TOKEN="your-session-token"
```

## CLI Commands

Common AWS CLI commands for reference:

```bash
# List EC2 instances
aws ec2 describe-instances

# List S3 buckets
aws s3 ls

# Get RDS instances
aws rds describe-db-instances

# List Lambda functions
aws lambda list-functions

# Get ECS services
aws ecs list-services --cluster my-cluster

# CloudWatch logs
aws logs tail /aws/lambda/my-function --follow
```

## Resources

- [Full AWS Guide](../docs/aws-guide.md)
- [AWS Examples](README-AWS.md)
- [Aether Syntax](../docs/syntax.md)
- [AWS Documentation](https://docs.aws.amazon.com/)
