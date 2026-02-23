# Aether AWS Provider Guide

Complete guide to deploying infrastructure on Amazon Web Services (AWS) using Aether.

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [AWS VPC - Networking](#aws-vpc---networking)
3. [AWS EC2 - Compute Instances](#aws-ec2---compute-instances)
4. [AWS S3 - Object Storage](#aws-s3---object-storage)
5. [AWS RDS - Databases](#aws-rds---databases)
6. [AWS Lambda - Serverless](#aws-lambda---serverless)
7. [AWS ELB - Load Balancing](#aws-elb---load-balancing)
8. [AWS IAM - Identity & Access Management](#aws-iam---identity--access-management)
9. [AWS ECS/EKS - Container Services](#aws-ecseks---container-services)
10. [AWS Route53 - DNS](#aws-route53---dns)
11. [AWS CloudWatch - Monitoring](#aws-cloudwatch---monitoring)
12. [Complete Examples](#complete-examples)

---

## Getting Started

### Provider Configuration

Configure the AWS provider in your Aether project:

```aether
provider aws {
  region = "us-east-1"
  profile = "default"  // AWS CLI profile
  
  // Optional: explicit credentials (use environment variables instead)
  // access_key = secret.aws_access_key
  // secret_key = secret.aws_secret_key
  
  default_tags = {
    ManagedBy = "Aether"
    Project = "my-project"
  }
}
```

### Multi-Region Setup

```aether
provider aws "primary" {
  region = "us-east-1"
  alias = "primary"
}

provider aws "secondary" {
  region = "us-west-2"
  alias = "secondary"
}

// Use specific provider
resource compute.instance "west_server" {
  provider = aws.secondary
  machine_type = "medium"
  os_image = "ubuntu-22.04"
}
```

---

## AWS VPC - Networking

### Basic VPC Setup

```aether
// Create VPC with public and private subnets

resource network.vpc "main" {
  cidr = "10.0.0.0/16"
  region = "us-east-1"
  
  enable_dns = true
  enable_dns_hostnames = true
  
  tags = {
    Name = "main-vpc"
    Environment = "production"
  }
}

resource network.subnet "public_1a" {
  vpc = vpc.main
  cidr = "10.0.1.0/24"
  availability_zone = "us-east-1a"
  public = true
  
  tags = {
    Name = "public-subnet-1a"
    Tier = "public"
  }
}

resource network.subnet "public_1b" {
  vpc = vpc.main
  cidr = "10.0.2.0/24"
  availability_zone = "us-east-1b"
  public = true
  
  tags = {
    Name = "public-subnet-1b"
    Tier = "public"
  }
}

resource network.subnet "private_1a" {
  vpc = vpc.main
  cidr = "10.0.10.0/24"
  availability_zone = "us-east-1a"
  public = false
  
  tags = {
    Name = "private-subnet-1a"
    Tier = "private"
  }
}

resource network.subnet "private_1b" {
  vpc = vpc.main
  cidr = "10.0.11.0/24"
  availability_zone = "us-east-1b"
  public = false
  
  tags = {
    Name = "private-subnet-1b"
    Tier = "private"
  }
}
```

### Internet Gateway and NAT Gateway

```aether
resource network.internet_gateway "main" {
  vpc = vpc.main
  
  tags = {
    Name = "main-igw"
  }
}

// NAT Gateway for private subnets
resource network.nat_gateway "az1" {
  subnet = subnet.public_1a
  
  tags = {
    Name = "nat-gateway-1a"
  }
}

resource network.nat_gateway "az2" {
  subnet = subnet.public_1b
  
  tags = {
    Name = "nat-gateway-1b"
  }
}
```

### Route Tables

```aether
// Public route table
resource network.route_table "public" {
  vpc = vpc.main
  
  route {
    cidr_block = "0.0.0.0/0"
    gateway = igw.main
  }
  
  tags = {
    Name = "public-rt"
  }
}

// Associate public subnets
resource network.route_table_association "public_1a" {
  subnet = subnet.public_1a
  route_table = rt.public
}

resource network.route_table_association "public_1b" {
  subnet = subnet.public_1b
  route_table = rt.public
}

// Private route tables (one per AZ for HA)
resource network.route_table "private_1a" {
  vpc = vpc.main
  
  route {
    cidr_block = "0.0.0.0/0"
    nat_gateway = nat.az1
  }
  
  tags = {
    Name = "private-rt-1a"
  }
}

resource network.route_table "private_1b" {
  vpc = vpc.main
  
  route {
    cidr_block = "0.0.0.0/0"
    nat_gateway = nat.az2
  }
  
  tags = {
    Name = "private-rt-1b"
  }
}
```

### Security Groups

```aether
// Web tier security group
resource network.security_group "web" {
  vpc = vpc.main
  description = "Security group for web servers"
  
  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTP from internet"
  }
  
  ingress {
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS from internet"
  }
  
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "All outbound traffic"
  }
  
  tags = {
    Name = "web-sg"
    Tier = "web"
  }
}

// Application tier security group
resource network.security_group "app" {
  vpc = vpc.main
  description = "Security group for application servers"
  
  ingress {
    from_port = 8080
    to_port = 8080
    protocol = "tcp"
    security_groups = [sg.web]
    description = "App traffic from web tier"
  }
  
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "All outbound traffic"
  }
  
  tags = {
    Name = "app-sg"
    Tier = "application"
  }
}

// Database tier security group
resource network.security_group "database" {
  vpc = vpc.main
  description = "Security group for database servers"
  
  ingress {
    from_port = 5432
    to_port = 5432
    protocol = "tcp"
    security_groups = [sg.app]
    description = "PostgreSQL from app tier"
  }
  
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "All outbound traffic"
  }
  
  tags = {
    Name = "database-sg"
    Tier = "database"
  }
}
```

### VPC Peering

```aether
resource network.vpc "app_vpc" {
  cidr = "10.1.0.0/16"
  region = "us-east-1"
  
  tags = {
    Name = "app-vpc"
  }
}

resource network.vpc_peering "main_to_app" {
  vpc = vpc.main
  peer_vpc = vpc.app_vpc
  
  auto_accept = true
  
  tags = {
    Name = "main-to-app-peering"
  }
}

// Add routes for peering
resource network.route "main_to_app" {
  route_table = rt.private_1a
  destination_cidr = vpc.app_vpc.cidr
  vpc_peering_connection = peering.main_to_app
}
```

### VPC Endpoints

```aether
// S3 VPC Endpoint (Gateway type)
resource network.vpc_endpoint "s3" {
  vpc = vpc.main
  service = "com.amazonaws.us-east-1.s3"
  type = "Gateway"
  
  route_tables = [rt.private_1a, rt.private_1b]
  
  tags = {
    Name = "s3-endpoint"
  }
}

// EC2 VPC Endpoint (Interface type)
resource network.vpc_endpoint "ec2" {
  vpc = vpc.main
  service = "com.amazonaws.us-east-1.ec2"
  type = "Interface"
  
  subnets = [subnet.private_1a, subnet.private_1b]
  security_groups = [sg.app]
  
  private_dns_enabled = true
  
  tags = {
    Name = "ec2-endpoint"
  }
}
```

---

## AWS EC2 - Compute Instances

### Basic EC2 Instance

```aether
resource compute.instance "web_server" {
  machine_type = "medium"  // t3.medium on AWS
  os_image = "ubuntu-22.04"
  region = "us-east-1"
  
  network {
    vpc = vpc.main
    subnet = subnet.public_1a
    security_groups = [sg.web]
    public_ip = true
  }
  
  storage {
    root_volume_size = 30  // GB
    root_volume_type = "ssd"  // gp3 on AWS
    
    additional_volumes = [
      {
        size = 100
        type = "ssd"
        device_name = "/dev/sdf"
        encrypted = true
      }
    ]
  }
  
  user_data = <<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y nginx
    echo "Hello from Aether!" > /var/www/html/index.html
    systemctl start nginx
    systemctl enable nginx
  EOF
  
  ssh_key = "my-keypair"
  
  monitoring_enabled = true
  
  tags = {
    Name = "web-server-01"
    Environment = "production"
    Role = "web"
  }
}

output "web_server_ip" {
  value = instance.web_server.public_ip
  description = "Public IP of web server"
}
```

### EC2 with Advanced Configuration

```aether
// Custom AMI
variable "custom_ami" {
  type = string
  description = "Custom AMI ID"
}

resource compute.instance "app_server" {
  ami = var.custom_ami  // AWS-specific AMI ID
  instance_type = "t3.large"  // Direct AWS instance type
  
  availability_zone = "us-east-1a"
  
  subnet = subnet.private_1a
  security_groups = [sg.app]
  
  // IAM instance profile
  iam_instance_profile = "ec2-app-role"
  
  // EBS optimization
  ebs_optimized = true
  
  // Detailed monitoring
  monitoring = true
  
  // Root volume configuration
  root_block_device {
    volume_type = "gp3"
    volume_size = 50
    iops = 3000
    throughput = 125
    encrypted = true
    delete_on_termination = true
  }
  
  // Additional EBS volumes
  ebs_block_device {
    device_name = "/dev/sdf"
    volume_type = "gp3"
    volume_size = 200
    iops = 5000
    throughput = 250
    encrypted = true
  }
  
  // Metadata options (IMDSv2)
  metadata_options {
    http_endpoint = "enabled"
    http_tokens = "required"  // Require IMDSv2
    http_put_response_hop_limit = 1
  }
  
  // Credit specification for T-series instances
  credit_specification {
    cpu_credits = "unlimited"
  }
  
  user_data = file("scripts/app_init.sh")
  
  tags = {
    Name = "app-server-01"
    Role = "application"
    Backup = "daily"
  }
  
  // Lifecycle configuration
  lifecycle {
    create_before_destroy = true
    ignore_changes = ["user_data"]
  }
}
```

### Auto Scaling Group

```aether
// Launch template
resource compute.launch_template "app" {
  name_prefix = "app-lt-"
  
  image = "ami-0c55b159cbfafe1f0"  // Ubuntu 22.04 in us-east-1
  instance_type = "t3.medium"
  
  vpc_security_groups = [sg.app]
  
  iam_instance_profile {
    name = "ec2-app-role"
  }
  
  block_device_mappings {
    device_name = "/dev/sda1"
    ebs {
      volume_size = 30
      volume_type = "gp3"
      encrypted = true
      delete_on_termination = true
    }
  }
  
  user_data = base64encode(<<-EOF
    #!/bin/bash
    apt-get update
    apt-get install -y docker.io
    docker run -d -p 8080:8080 myapp:latest
  EOF
  )
  
  tag_specifications {
    resource_type = "instance"
    tags = {
      Name = "app-asg-instance"
      Environment = "production"
    }
  }
}

// Auto Scaling Group
resource compute.auto_scaling_group "app" {
  name = "app-asg"
  
  launch_template = lt.app
  
  min_size = 2
  max_size = 10
  desired_capacity = 3
  
  vpc_zone_identifier = [subnet.private_1a, subnet.private_1b]
  
  target_group_arns = [tg.app.arn]
  
  health_check_type = "ELB"
  health_check_grace_period = 300
  
  // Scaling policies
  scaling_policy {
    name = "cpu-scale-up"
    policy_type = "TargetTrackingScaling"
    
    target_tracking_configuration {
      predefined_metric_type = "ASGAverageCPUUtilization"
      target_value = 70.0
    }
  }
  
  tags = {
    Name = "app-asg"
    Environment = "production"
    propagate_at_launch = true
  }
}
```

### Spot Instances

```aether
resource compute.instance "batch_processor" {
  instance_type = "c5.xlarge"
  
  // Spot instance configuration
  instance_market_options {
    market_type = "spot"
    
    spot_options {
      max_price = "0.15"  // Per hour
      spot_instance_type = "one-time"
      instance_interruption_behavior = "terminate"
    }
  }
  
  subnet = subnet.private_1a
  security_groups = [sg.app]
  
  user_data = <<-EOF
    #!/bin/bash
    # Run batch processing job
    aws s3 cp s3://my-bucket/job-data.zip /tmp/
    process-batch /tmp/job-data.zip
    aws s3 cp /tmp/results.csv s3://my-bucket/results/
  EOF
  
  tags = {
    Name = "batch-processor-spot"
    Type = "spot"
  }
}
```

---

## AWS S3 - Object Storage

### Basic S3 Bucket

```aether
resource storage.bucket "app_data" {
  region = "us-east-1"
  
  versioning = true
  encryption_enabled = true
  
  public_access = false
  
  tags = {
    Name = "app-data-bucket"
    Environment = "production"
  }
}

output "bucket_name" {
  value = bucket.app_data.name
}
```

### S3 Bucket with Lifecycle Rules

```aether
resource storage.bucket "logs" {
  region = "us-east-1"
  
  versioning = true
  encryption_enabled = true
  
  // Lifecycle management
  lifecycle_rules = [
    {
      action = "transition"
      age_days = 30
      storage_class = "STANDARD_IA"
      prefix = "logs/"
    },
    {
      action = "transition"
      age_days = 90
      storage_class = "GLACIER"
      prefix = "logs/"
    },
    {
      action = "delete"
      age_days = 365
      prefix = "logs/"
    },
    {
      action = "delete"
      age_days = 7
      prefix = "temp/"
    }
  ]
  
  tags = {
    Name = "logs-bucket"
    DataClassification = "internal"
  }
}
```

### S3 Bucket with Advanced Features

```aether
resource storage.bucket "static_website" {
  bucket_name = "my-website-${var.environment}"
  region = "us-east-1"
  
  // Static website hosting
  website {
    index_document = "index.html"
    error_document = "error.html"
    
    routing_rules = <<-JSON
    [
      {
        "Redirect": {
          "ReplaceKeyPrefixWith": "docs/"
        },
        "Condition": {
          "KeyPrefixEquals": "documents/"
        }
      }
    ]
    JSON
  }
  
  // Versioning
  versioning = true
  
  // Encryption
  encryption {
    sse_algorithm = "AES256"
    // Or use KMS
    // kms_master_key_id = kms_key.app.id
  }
  
  // Public access for website
  public_access = true
  
  // CORS configuration
  cors_rules = [
    {
      allowed_origins = ["https://example.com"]
      allowed_methods = ["GET", "HEAD"]
      allowed_headers = ["*"]
      max_age_seconds = 3600
    }
  ]
  
  // Logging
  logging {
    target_bucket = bucket.logs.name
    target_prefix = "s3-access-logs/"
  }
  
  // Access policy
  access_policy = <<-JSON
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "PublicReadGetObject",
        "Effect": "Allow",
        "Principal": "*",
        "Action": "s3:GetObject",
        "Resource": "arn:aws:s3:::my-website-${var.environment}/*"
      }
    ]
  }
  JSON
  
  tags = {
    Name = "static-website"
    Type = "website"
  }
}
```

### S3 Bucket Replication

```aether
// Source bucket
resource storage.bucket "primary" {
  bucket_name = "my-app-primary"
  region = "us-east-1"
  
  versioning = true
  
  replication_configuration {
    role = iam_role.replication.arn
    
    rules {
      id = "replicate-all"
      status = "Enabled"
      
      destination {
        bucket = bucket.replica.arn
        storage_class = "STANDARD_IA"
        
        // Optional: Replication Time Control
        replication_time {
          status = "Enabled"
          minutes = 15
        }
        
        // Optional: Metrics
        metrics {
          status = "Enabled"
          minutes = 15
        }
      }
    }
  }
  
  tags = {
    Name = "primary-bucket"
  }
}

// Replica bucket
resource storage.bucket "replica" {
  provider = aws.secondary  // Different region
  
  bucket_name = "my-app-replica"
  region = "us-west-2"
  
  versioning = true
  
  tags = {
    Name = "replica-bucket"
  }
}
```

### S3 Event Notifications

```aether
resource storage.bucket "uploads" {
  bucket_name = "app-uploads"
  region = "us-east-1"
  
  // Lambda notification
  notification {
    lambda_function = lambda.image_processor.arn
    events = ["s3:ObjectCreated:*"]
    filter_prefix = "uploads/"
    filter_suffix = ".jpg"
  }
  
  // SQS notification
  notification {
    queue = sqs.processing_queue.arn
    events = ["s3:ObjectCreated:*"]
    filter_prefix = "data/"
  }
  
  // SNS notification
  notification {
    topic = sns.bucket_notifications.arn
    events = ["s3:ObjectRemoved:*"]
  }
  
  tags = {
    Name = "uploads-bucket"
  }
}
```

---

## AWS RDS - Databases

### PostgreSQL Database

```aether
// DB Subnet Group
resource database.subnet_group "main" {
  name = "main-db-subnet-group"
  subnets = [subnet.private_1a, subnet.private_1b]
  
  tags = {
    Name = "main-db-subnet-group"
  }
}

// RDS Instance
resource database.instance "postgres" {
  engine = "postgres"
  engine_version = "15.3"
  
  instance_class = "large"  // db.t3.large
  
  storage = 100  // GB
  storage_type = "ssd"  // gp3
  storage_encrypted = true
  
  database_name = "myapp"
  master_username = "admin"
  master_password = secret.db_password
  
  // High availability
  multi_az = true
  
  // Network configuration
  vpc = vpc.main
  subnet_group = db_subnet_group.main
  security_groups = [sg.database]
  publicly_accessible = false
  
  // Backup configuration
  backup_retention_days = 7
  backup_window = "03:00-04:00"
  maintenance_window = "sun:04:00-sun:05:00"
  
  // Performance Insights
  performance_insights_enabled = true
  performance_insights_retention = 7
  
  // Enhanced Monitoring
  monitoring_interval = 60
  monitoring_role_arn = iam_role.rds_monitoring.arn
  
  // Auto minor version upgrades
  auto_minor_version_upgrade = true
  
  // Deletion protection
  deletion_protection = true
  
  tags = {
    Name = "postgres-main"
    Environment = "production"
  }
}

output "db_endpoint" {
  value = db.postgres.endpoint
  sensitive = true
}
```

### MySQL Database with Read Replicas

```aether
// Primary database
resource database.instance "mysql_primary" {
  identifier = "myapp-mysql-primary"
  
  engine = "mysql"
  engine_version = "8.0.33"
  
  instance_class = "xlarge"  // db.r6g.xlarge
  
  storage = 200
  storage_type = "provisioned-iops"
  iops = 10000
  
  database_name = "myapp"
  master_username = "admin"
  master_password = secret.mysql_password
  
  multi_az = true
  
  vpc = vpc.main
  subnet_group = db_subnet_group.main
  security_groups = [sg.database]
  
  backup_retention_days = 14
  backup_window = "03:00-04:00"
  
  // Parameter group for custom settings
  parameter_group = db_parameter_group.mysql_custom
  
  tags = {
    Name = "mysql-primary"
    Role = "primary"
  }
}

// Read replica in same region
resource database.read_replica "mysql_replica_1" {
  identifier = "myapp-mysql-replica-1"
  source_db = db.mysql_primary
  
  instance_class = "large"
  
  publicly_accessible = false
  
  tags = {
    Name = "mysql-replica-1"
    Role = "read-replica"
  }
}

// Read replica in different region
resource database.read_replica "mysql_replica_west" {
  provider = aws.secondary
  
  identifier = "myapp-mysql-replica-west"
  source_db = db.mysql_primary
  
  instance_class = "large"
  
  tags = {
    Name = "mysql-replica-west"
    Role = "read-replica"
    Region = "us-west-2"
  }
}
```

### Aurora Cluster

```aether
resource database.aurora_cluster "main" {
  cluster_identifier = "myapp-aurora-cluster"
  
  engine = "aurora-postgresql"
  engine_version = "15.3"
  
  database_name = "myapp"
  master_username = "admin"
  master_password = secret.aurora_password
  
  // Network
  vpc = vpc.main
  subnet_group = db_subnet_group.main
  security_groups = [sg.database]
  
  // Availability
  availability_zones = ["us-east-1a", "us-east-1b", "us-east-1c"]
  
  // Backup
  backup_retention_period = 14
  preferred_backup_window = "03:00-04:00"
  preferred_maintenance_window = "sun:04:00-sun:05:00"
  
  // Encryption
  storage_encrypted = true
  kms_key_id = kms_key.rds.arn
  
  // Backtrack (point-in-time recovery)
  backtrack_window = 72  // hours
  
  // Enhanced features
  enabled_cloudwatch_logs_exports = ["postgresql"]
  
  // Scaling
  scaling_configuration {
    auto_pause = true
    max_capacity = 16
    min_capacity = 2
    seconds_until_auto_pause = 300
  }
  
  tags = {
    Name = "aurora-main-cluster"
    Environment = "production"
  }
}

// Aurora cluster instances
resource database.aurora_instance "writer" {
  identifier = "myapp-aurora-writer"
  cluster = aurora.main
  
  instance_class = "r6g.xlarge"
  
  performance_insights_enabled = true
  
  tags = {
    Name = "aurora-writer"
    Role = "writer"
  }
}

resource database.aurora_instance "reader" {
  count = 2
  
  identifier = "myapp-aurora-reader-${count.index + 1}"
  cluster = aurora.main
  
  instance_class = "r6g.large"
  
  performance_insights_enabled = true
  
  tags = {
    Name = "aurora-reader-${count.index + 1}"
    Role = "reader"
  }
}
```

---

## AWS Lambda - Serverless

### Basic Lambda Function

```aether
resource compute.lambda_function "api_handler" {
  function_name = "api-handler"
  runtime = "python3.11"
  handler = "index.handler"
  
  // Code from local file
  source_code = file("lambda/api_handler.py")
  
  // Or from S3
  // source_bucket = bucket.lambda_code
  // source_key = "api_handler.zip"
  
  // IAM role
  role = iam_role.lambda_exec.arn
  
  // Memory and timeout
  memory = 512  // MB
  timeout = 30  // seconds
  
  // Environment variables
  environment {
    variables = {
      DB_HOST = db.postgres.endpoint
      TABLE_NAME = "users"
      LOG_LEVEL = "INFO"
    }
  }
  
  // VPC configuration (if accessing VPC resources)
  vpc_config {
    subnet_ids = [subnet.private_1a, subnet.private_1b]
    security_group_ids = [sg.lambda]
  }
  
  tags = {
    Name = "api-handler"
    Environment = "production"
  }
}

// Lambda permission for API Gateway
resource lambda.permission "api_gateway" {
  function = lambda.api_handler
  action = "lambda:InvokeFunction"
  principal = "apigateway.amazonaws.com"
  source_arn = "${api_gateway.main.execution_arn}/*/*"
}
```

### Lambda with Event Sources

```aether
// SQS trigger
resource queue "processing" {
  name = "processing-queue"
  
  tags = {
    Name = "processing-queue"
  }
}

resource lambda.event_source_mapping "sqs" {
  event_source_arn = queue.processing.arn
  function = lambda.processor
  batch_size = 10
  maximum_batching_window_in_seconds = 5
}

// S3 trigger (via S3 bucket notification)
resource storage.bucket "uploads" {
  notification {
    lambda_function = lambda.image_processor.arn
    events = ["s3:ObjectCreated:*"]
    filter_suffix = ".jpg"
  }
}

// DynamoDB Stream trigger
resource lambda.event_source_mapping "dynamodb" {
  event_source_arn = dynamodb_table.users.stream_arn
  function = lambda.stream_processor
  starting_position = "LATEST"
  batch_size = 100
}

// EventBridge (CloudWatch Events) trigger
resource eventbridge.rule "hourly" {
  name = "hourly-lambda-trigger"
  schedule_expression = "rate(1 hour)"
}

resource eventbridge.target "lambda" {
  rule = eventbridge_rule.hourly
  arn = lambda.cleanup.arn
}
```

### Lambda with Layers

```aether
// Lambda Layer
resource lambda.layer "dependencies" {
  layer_name = "python-dependencies"
  
  source_code = file("layers/dependencies.zip")
  
  compatible_runtimes = ["python3.11", "python3.10"]
  
  description = "Common Python dependencies"
}

// Lambda using layer
resource compute.lambda_function "app" {
  function_name = "app-function"
  runtime = "python3.11"
  handler = "app.handler"
  
  source_code = file("lambda/app.py")
  
  layers = [layer.dependencies.arn]
  
  role = iam_role.lambda_exec.arn
  
  tags = {
    Name = "app-function"
  }
}
```

### Lambda with Provisioned Concurrency

```aether
resource compute.lambda_function "high_traffic" {
  function_name = "high-traffic-api"
  runtime = "nodejs18.x"
  handler = "index.handler"
  
  source_code = file("lambda/api.zip")
  
  role = iam_role.lambda_exec.arn
  
  memory = 1024
  timeout = 30
  
  // Reserved concurrent executions
  reserved_concurrent_executions = 100
  
  tags = {
    Name = "high-traffic-api"
  }
}

// Provisioned concurrency
resource lambda.provisioned_concurrency "high_traffic" {
  function = lambda.high_traffic
  qualifier = lambda.high_traffic.version
  provisioned_concurrent_executions = 10
}

// Auto scaling for provisioned concurrency
resource lambda.provisioned_concurrency_scaling "high_traffic" {
  function = lambda.high_traffic
  
  min_capacity = 5
  max_capacity = 50
  
  target_utilization = 0.70
}
```

---

## AWS ELB - Load Balancing

### Application Load Balancer

```aether
// ALB
resource loadbalancer.application "web" {
  name = "web-alb"
  
  vpc = vpc.main
  subnets = [subnet.public_1a, subnet.public_1b]
  security_groups = [sg.alb]
  
  internal = false  // Internet-facing
  
  // Access logs
  access_logs {
    bucket = bucket.logs
    prefix = "alb-logs"
    enabled = true
  }
  
  tags = {
    Name = "web-alb"
    Environment = "production"
  }
}

// Target Group
resource loadbalancer.target_group "web" {
  name = "web-tg"
  
  port = 80
  protocol = "HTTP"
  vpc = vpc.main
  
  // Health check
  health_check {
    enabled = true
    path = "/health"
    port = "traffic-port"
    protocol = "HTTP"
    healthy_threshold = 2
    unhealthy_threshold = 3
    timeout = 5
    interval = 30
    matcher = "200"
  }
  
  // Stickiness
  stickiness {
    enabled = true
    type = "lb_cookie"
    cookie_duration = 86400  // 1 day
  }
  
  tags = {
    Name = "web-target-group"
  }
}

// HTTP Listener
resource loadbalancer.listener "http" {
  loadbalancer = lb.web
  port = 80
  protocol = "HTTP"
  
  default_action {
    type = "redirect"
    redirect {
      protocol = "HTTPS"
      port = "443"
      status_code = "HTTP_301"
    }
  }
}

// HTTPS Listener
resource loadbalancer.listener "https" {
  loadbalancer = lb.web
  port = 443
  protocol = "HTTPS"
  
  ssl_policy = "ELBSecurityPolicy-TLS-1-2-2017-01"
  certificate_arn = acm_certificate.main.arn
  
  default_action {
    type = "forward"
    target_group = tg.web
  }
}

// Listener Rules
resource loadbalancer.listener_rule "api" {
  listener = listener.https
  priority = 10
  
  condition {
    path_pattern = ["/api/*"]
  }
  
  action {
    type = "forward"
    target_group = tg.api
  }
}

resource loadbalancer.listener_rule "static" {
  listener = listener.https
  priority = 20
  
  condition {
    path_pattern = ["/static/*", "/images/*"]
  }
  
  action {
    type = "forward"
    target_group = tg.static
  }
}

// Register instances with target group
resource loadbalancer.target_group_attachment "web_1" {
  target_group = tg.web
  target_id = instance.web_1.id
  port = 80
}
```

### Network Load Balancer

```aether
resource loadbalancer.network "tcp" {
  name = "tcp-nlb"
  
  vpc = vpc.main
  subnets = [subnet.public_1a, subnet.public_1b]
  
  internal = false
  
  enable_cross_zone_load_balancing = true
  
  tags = {
    Name = "tcp-nlb"
  }
}

resource loadbalancer.target_group "tcp" {
  name = "tcp-tg"
  
  port = 3306
  protocol = "TCP"
  vpc = vpc.main
  
  health_check {
    enabled = true
    port = "traffic-port"
    protocol = "TCP"
    healthy_threshold = 3
    unhealthy_threshold = 3
    interval = 30
  }
  
  // Preserve client IP
  preserve_client_ip = true
  
  tags = {
    Name = "tcp-target-group"
  }
}

resource loadbalancer.listener "tcp" {
  loadbalancer = nlb.tcp
  port = 3306
  protocol = "TCP"
  
  default_action {
    type = "forward"
    target_group = tg.tcp
  }
}
```

---

## AWS IAM - Identity & Access Management

### IAM Roles and Policies

```aether
// EC2 instance role
resource iam.role "ec2_app" {
  name = "ec2-app-role"
  
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
  
  tags = {
    Name = "ec2-app-role"
  }
}

// IAM Policy
resource iam.policy "s3_access" {
  name = "s3-access-policy"
  description = "Allow access to specific S3 buckets"
  
  policy = <<-JSON
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject"
        ],
        "Resource": [
          "${bucket.app_data.arn}/*"
        ]
      },
      {
        "Effect": "Allow",
        "Action": [
          "s3:ListBucket"
        ],
        "Resource": [
          "${bucket.app_data.arn}"
        ]
      }
    ]
  }
  JSON
}

// Attach policy to role
resource iam.role_policy_attachment "ec2_s3" {
  role = role.ec2_app
  policy_arn = policy.s3_access.arn
}

// Instance profile for EC2
resource iam.instance_profile "ec2_app" {
  name = "ec2-app-profile"
  role = role.ec2_app
}
```

### Lambda Execution Role

```aether
resource iam.role "lambda_exec" {
  name = "lambda-execution-role"
  
  assume_role_policy = <<-JSON
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "lambda.amazonaws.com"
        },
        "Effect": "Allow"
      }
    ]
  }
  JSON
}

// Attach AWS managed policy for Lambda
resource iam.role_policy_attachment "lambda_basic" {
  role = role.lambda_exec
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// Custom policy for Lambda
resource iam.policy "lambda_custom" {
  name = "lambda-custom-policy"
  
  policy = <<-JSON
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:Query"
        ],
        "Resource": "${dynamodb_table.users.arn}"
      },
      {
        "Effect": "Allow",
        "Action": [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage"
        ],
        "Resource": "${queue.processing.arn}"
      }
    ]
  }
  JSON
}

resource iam.role_policy_attachment "lambda_custom" {
  role = role.lambda_exec
  policy_arn = policy.lambda_custom.arn
}
```

### IAM Users and Groups

```aether
// IAM Group
resource iam.group "developers" {
  name = "developers"
}

// IAM Group Policy
resource iam.group_policy "dev_access" {
  group = group.developers
  
  policy = <<-JSON
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "ec2:Describe*",
          "s3:ListBucket",
          "s3:GetObject"
        ],
        "Resource": "*"
      }
    ]
  }
  JSON
}

// IAM User
resource iam.user "developer" {
  name = "john.doe"
  
  tags = {
    Department = "Engineering"
    Team = "Backend"
  }
}

// Add user to group
resource iam.user_group_membership "developer" {
  user = user.developer
  groups = [group.developers]
}
```

---

## AWS ECS/EKS - Container Services

### ECS with Fargate

```aether
// ECS Cluster
resource compute.ecs_cluster "main" {
  name = "main-cluster"
  
  setting {
    name = "containerInsights"
    value = "enabled"
  }
  
  tags = {
    Name = "main-ecs-cluster"
    Environment = "production"
  }
}

// Task Definition
resource compute.ecs_task_definition "app" {
  family = "app-task"
  
  network_mode = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  
  cpu = "1024"  // 1 vCPU
  memory = "2048"  // 2 GB
  
  execution_role_arn = iam_role.ecs_execution.arn
  task_role_arn = iam_role.ecs_task.arn
  
  container_definitions = <<-JSON
  [
    {
      "name": "app",
      "image": "123456789.dkr.ecr.us-east-1.amazonaws.com/myapp:latest",
      "cpu": 1024,
      "memory": 2048,
      "essential": true,
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "DB_HOST",
          "value": "${db.postgres.endpoint}"
        },
        {
          "name": "ENVIRONMENT",
          "value": "production"
        }
      ],
      "secrets": [
        {
          "name": "DB_PASSWORD",
          "valueFrom": "${secret.db_password.arn}"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/app",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
  JSON
  
  tags = {
    Name = "app-task"
  }
}

// ECS Service
resource compute.ecs_service "app" {
  name = "app-service"
  cluster = ecs_cluster.main.id
  task_definition = ecs_task_definition.app.arn
  
  desired_count = 3
  launch_type = "FARGATE"
  
  network_configuration {
    subnets = [subnet.private_1a, subnet.private_1b]
    security_groups = [sg.ecs]
    assign_public_ip = false
  }
  
  load_balancer {
    target_group_arn = tg.app.arn
    container_name = "app"
    container_port = 8080
  }
  
  // Auto scaling
  enable_ecs_managed_tags = true
  propagate_tags = "SERVICE"
  
  // Deployment configuration
  deployment_configuration {
    maximum_percent = 200
    minimum_healthy_percent = 100
  }
  
  // Service discovery
  service_registries {
    registry_arn = service_discovery_service.app.arn
  }
  
  tags = {
    Name = "app-ecs-service"
  }
}

// ECS Auto Scaling
resource compute.ecs_service_scaling_target "app" {
  service = ecs_service.app
  min_capacity = 2
  max_capacity = 10
}

resource compute.ecs_service_scaling_policy "cpu" {
  name = "cpu-autoscaling"
  
  service_scaling_target = ecs_scaling_target.app
  
  policy_type = "TargetTrackingScaling"
  
  target_tracking_scaling_policy_configuration {
    predefined_metric_type = "ECSServiceAverageCPUUtilization"
    target_value = 70.0
  }
}
```

### EKS Cluster

```aether
// EKS Cluster
resource compute.eks_cluster "main" {
  name = "main-eks-cluster"
  role_arn = iam_role.eks_cluster.arn
  version = "1.28"
  
  vpc_config {
    subnet_ids = [
      subnet.private_1a,
      subnet.private_1b,
      subnet.public_1a,
      subnet.public_1b
    ]
    security_group_ids = [sg.eks_cluster]
    endpoint_private_access = true
    endpoint_public_access = true
    public_access_cidrs = ["0.0.0.0/0"]
  }
  
  // Enable control plane logging
  enabled_cluster_log_types = [
    "api",
    "audit",
    "authenticator",
    "controllerManager",
    "scheduler"
  ]
  
  // Encryption
  encryption_config {
    resources = ["secrets"]
    provider {
      key_arn = kms_key.eks.arn
    }
  }
  
  tags = {
    Name = "main-eks-cluster"
    Environment = "production"
  }
}

// EKS Node Group
resource compute.eks_node_group "main" {
  cluster_name = eks_cluster.main.name
  node_group_name = "main-node-group"
  node_role_arn = iam_role.eks_node.arn
  
  subnet_ids = [subnet.private_1a, subnet.private_1b]
  
  scaling_config {
    desired_size = 3
    max_size = 10
    min_size = 2
  }
  
  instance_types = ["t3.large"]
  
  // Launch template for node configuration
  launch_template {
    id = launch_template.eks_nodes.id
    version = launch_template.eks_nodes.latest_version
  }
  
  // Update config
  update_config {
    max_unavailable_percentage = 33
  }
  
  labels = {
    role = "general"
    environment = "production"
  }
  
  tags = {
    Name = "main-eks-nodes"
  }
}
```

---

## AWS Route53 - DNS

### Hosted Zone and Records

```aether
// Hosted Zone
resource dns.hosted_zone "main" {
  name = "example.com"
  
  tags = {
    Name = "example.com"
    Environment = "production"
  }
}

// A Record pointing to ALB
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

// CNAME Record
resource dns.record "api" {
  zone = hosted_zone.main
  name = "api.example.com"
  type = "CNAME"
  ttl = 300
  
  records = [lb.api.dns_name]
}

// MX Record
resource dns.record "mail" {
  zone = hosted_zone.main
  name = "example.com"
  type = "MX"
  ttl = 3600
  
  records = [
    "10 mail1.example.com",
    "20 mail2.example.com"
  ]
}

// TXT Record for SPF
resource dns.record "spf" {
  zone = hosted_zone.main
  name = "example.com"
  type = "TXT"
  ttl = 3600
  
  records = [
    "v=spf1 include:_spf.google.com ~all"
  ]
}
```

### Route53 Health Checks and Routing

```aether
// Health Check
resource route53.health_check "web" {
  type = "HTTPS"
  resource_path = "/health"
  fqdn = "www.example.com"
  port = 443
  
  request_interval = 30
  failure_threshold = 3
  
  measure_latency = true
  
  tags = {
    Name = "web-health-check"
  }
}

// Geolocation Routing
resource dns.record "www_us" {
  zone = hosted_zone.main
  name = "www.example.com"
  type = "A"
  
  set_identifier = "US-East"
  
  geolocation_routing_policy {
    continent = "NA"
    country = "US"
  }
  
  alias {
    name = lb.us_east.dns_name
    zone_id = lb.us_east.zone_id
    evaluate_target_health = true
  }
  
  health_check = health_check.web
}

resource dns.record "www_eu" {
  zone = hosted_zone.main
  name = "www.example.com"
  type = "A"
  
  set_identifier = "EU-West"
  
  geolocation_routing_policy {
    continent = "EU"
  }
  
  alias {
    name = lb.eu_west.dns_name
    zone_id = lb.eu_west.zone_id
    evaluate_target_health = true
  }
}

// Latency-based Routing
resource dns.record "api_us_east" {
  zone = hosted_zone.main
  name = "api.example.com"
  type = "A"
  
  set_identifier = "US-East-1"
  
  latency_routing_policy {
    region = "us-east-1"
  }
  
  alias {
    name = lb.api_us_east.dns_name
    zone_id = lb.api_us_east.zone_id
    evaluate_target_health = true
  }
}
```

---

## AWS CloudWatch - Monitoring

### CloudWatch Alarms

```aether
// CPU Alarm for EC2
resource monitoring.alarm "ec2_high_cpu" {
  name = "web-server-high-cpu"
  description = "Alert when CPU exceeds 80%"
  
  namespace = "AWS/EC2"
  metric_name = "CPUUtilization"
  
  dimensions = {
    InstanceId = instance.web_server.id
  }
  
  statistic = "Average"
  period = 300  // 5 minutes
  evaluation_periods = 2
  threshold = 80
  comparison_operator = "GreaterThanThreshold"
  
  alarm_actions = [sns_topic.alerts.arn]
  ok_actions = [sns_topic.alerts.arn]
  
  treat_missing_data = "notBreaching"
  
  tags = {
    Name = "ec2-high-cpu-alarm"
  }
}

// RDS Alarm
resource monitoring.alarm "rds_connections" {
  name = "db-high-connections"
  description = "Alert when database connections exceed 80"
  
  namespace = "AWS/RDS"
  metric_name = "DatabaseConnections"
  
  dimensions = {
    DBInstanceIdentifier = db.postgres.identifier
  }
  
  statistic = "Average"
  period = 300
  evaluation_periods = 2
  threshold = 80
  comparison_operator = "GreaterThanThreshold"
  
  alarm_actions = [sns_topic.db_alerts.arn]
}

// ALB Target Response Time
resource monitoring.alarm "alb_response_time" {
  name = "alb-slow-response"
  description = "Alert when response time exceeds 2 seconds"
  
  namespace = "AWS/ApplicationELB"
  metric_name = "TargetResponseTime"
  
  dimensions = {
    LoadBalancer = lb.web.arn_suffix
  }
  
  extended_statistic = "p99"
  period = 300
  evaluation_periods = 3
  threshold = 2.0
  comparison_operator = "GreaterThanThreshold"
  
  alarm_actions = [sns_topic.performance_alerts.arn]
}
```

### CloudWatch Logs

```aether
// Log Group
resource monitoring.log_group "app" {
  name = "/aws/app/logs"
  retention_in_days = 30
  
  kms_key_id = kms_key.logs.arn
  
  tags = {
    Name = "app-logs"
    Application = "myapp"
  }
}

// Metric Filter
resource monitoring.log_metric_filter "errors" {
  name = "error-count"
  log_group = log_group.app
  
  pattern = "[timestamp, requestid, level = ERROR*, msg]"
  
  metric_transformation {
    name = "ErrorCount"
    namespace = "MyApp"
    value = "1"
    default_value = 0
  }
}

// Alarm on custom metric
resource monitoring.alarm "app_errors" {
  name = "app-error-count"
  description = "Alert on application errors"
  
  namespace = "MyApp"
  metric_name = "ErrorCount"
  
  statistic = "Sum"
  period = 300
  evaluation_periods = 1
  threshold = 10
  comparison_operator = "GreaterThanThreshold"
  
  alarm_actions = [sns_topic.alerts.arn]
}
```

### CloudWatch Dashboard

```aether
resource monitoring.dashboard "main" {
  name = "production-dashboard"
  
  dashboard_body = <<-JSON
  {
    "widgets": [
      {
        "type": "metric",
        "properties": {
          "metrics": [
            ["AWS/EC2", "CPUUtilization", {"stat": "Average"}],
            ["AWS/RDS", "CPUUtilization", {"stat": "Average"}]
          ],
          "period": 300,
          "stat": "Average",
          "region": "us-east-1",
          "title": "CPU Utilization",
          "yAxis": {
            "left": {
              "min": 0,
              "max": 100
            }
          }
        }
      },
      {
        "type": "metric",
        "properties": {
          "metrics": [
            ["AWS/ApplicationELB", "RequestCount", {"stat": "Sum"}],
            [".", "TargetResponseTime", {"stat": "Average"}]
          ],
          "period": 300,
          "stat": "Average",
          "region": "us-east-1",
          "title": "ALB Metrics"
        }
      }
    ]
  }
  JSON
}
```

---

## Complete Examples

### Example 1: Three-Tier Web Application

Complete infrastructure for a scalable web application with web, app, and database tiers.

```aether
// Variables
variable "environment" {
  type = string
  default = "production"
}

variable "app_name" {
  type = string
  default = "myapp"
}

// Provider
provider aws {
  region = "us-east-1"
  
  default_tags = {
    Application = var.app_name
    Environment = var.environment
    ManagedBy = "Aether"
  }
}

// VPC and Networking
resource network.vpc "main" {
  cidr = "10.0.0.0/16"
  enable_dns = true
  enable_dns_hostnames = true
  
  tags = {
    Name = "${var.app_name}-vpc"
  }
}

// Public Subnets
resource network.subnet "public" {
  count = 2
  
  vpc = vpc.main
  cidr = "10.0.${count.index}.0/24"
  availability_zone = "us-east-1${["a", "b"][count.index]}"
  public = true
  
  tags = {
    Name = "${var.app_name}-public-${count.index + 1}"
    Tier = "public"
  }
}

// Private Subnets for App Tier
resource network.subnet "app" {
  count = 2
  
  vpc = vpc.main
  cidr = "10.0.${10 + count.index}.0/24"
  availability_zone = "us-east-1${["a", "b"][count.index]}"
  public = false
  
  tags = {
    Name = "${var.app_name}-app-${count.index + 1}"
    Tier = "app"
  }
}

// Private Subnets for Database Tier
resource network.subnet "db" {
  count = 2
  
  vpc = vpc.main
  cidr = "10.0.${20 + count.index}.0/24"
  availability_zone = "us-east-1${["a", "b"][count.index]}"
  public = false
  
  tags = {
    Name = "${var.app_name}-db-${count.index + 1}"
    Tier = "database"
  }
}

// Internet Gateway and NAT Gateways
resource network.internet_gateway "main" {
  vpc = vpc.main
  tags = { Name = "${var.app_name}-igw" }
}

resource network.nat_gateway "main" {
  count = 2
  subnet = subnet.public[count.index]
  tags = { Name = "${var.app_name}-nat-${count.index + 1}" }
}

// Security Groups
resource network.security_group "alb" {
  vpc = vpc.main
  description = "ALB security group"
  
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
  
  tags = { Name = "${var.app_name}-alb-sg" }
}

resource network.security_group "app" {
  vpc = vpc.main
  description = "App tier security group"
  
  ingress {
    from_port = 8080
    to_port = 8080
    protocol = "tcp"
    security_groups = [sg.alb]
  }
  
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  tags = { Name = "${var.app_name}-app-sg" }
}

resource network.security_group "db" {
  vpc = vpc.main
  description = "Database security group"
  
  ingress {
    from_port = 5432
    to_port = 5432
    protocol = "tcp"
    security_groups = [sg.app]
  }
  
  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  tags = { Name = "${var.app_name}-db-sg" }
}

// Application Load Balancer
resource loadbalancer.application "main" {
  name = "${var.app_name}-alb"
  vpc = vpc.main
  subnets = [subnet.public[0], subnet.public[1]]
  security_groups = [sg.alb]
  internal = false
  
  tags = { Name = "${var.app_name}-alb" }
}

resource loadbalancer.target_group "app" {
  name = "${var.app_name}-tg"
  port = 8080
  protocol = "HTTP"
  vpc = vpc.main
  
  health_check {
    path = "/health"
    healthy_threshold = 2
    unhealthy_threshold = 3
    timeout = 5
    interval = 30
  }
  
  tags = { Name = "${var.app_name}-tg" }
}

resource loadbalancer.listener "http" {
  loadbalancer = lb.main
  port = 80
  protocol = "HTTP"
  
  default_action {
    type = "forward"
    target_group = tg.app
  }
}

// Auto Scaling Group
resource compute.launch_template "app" {
  name_prefix = "${var.app_name}-lt-"
  image = "ami-0c55b159cbfafe1f0"
  instance_type = "t3.medium"
  
  vpc_security_groups = [sg.app]
  
  user_data = base64encode(<<-EOF
    #!/bin/bash
    # Install and start application
    apt-get update
    apt-get install -y docker.io
    docker run -d -p 8080:8080 \
      -e DB_HOST=${db.main.endpoint} \
      ${var.app_name}:latest
  EOF
  )
  
  tag_specifications {
    resource_type = "instance"
    tags = {
      Name = "${var.app_name}-app-instance"
    }
  }
}

resource compute.auto_scaling_group "app" {
  name = "${var.app_name}-asg"
  launch_template = lt.app
  
  min_size = 2
  max_size = 10
  desired_capacity = 3
  
  vpc_zone_identifier = [subnet.app[0], subnet.app[1]]
  target_group_arns = [tg.app.arn]
  
  health_check_type = "ELB"
  health_check_grace_period = 300
  
  scaling_policy {
    name = "cpu-scaling"
    policy_type = "TargetTrackingScaling"
    target_tracking_configuration {
      predefined_metric_type = "ASGAverageCPUUtilization"
      target_value = 70.0
    }
  }
  
  tags = { Name = "${var.app_name}-asg" }
}

// RDS Database
resource database.subnet_group "main" {
  name = "${var.app_name}-db-subnet-group"
  subnets = [subnet.db[0], subnet.db[1]]
  tags = { Name = "${var.app_name}-db-subnet-group" }
}

resource database.instance "main" {
  identifier = "${var.app_name}-db"
  engine = "postgres"
  engine_version = "15.3"
  instance_class = "large"
  
  storage = 100
  storage_type = "ssd"
  storage_encrypted = true
  
  database_name = var.app_name
  master_username = "admin"
  master_password = secret.db_password
  
  multi_az = true
  
  vpc = vpc.main
  subnet_group = db_subnet_group.main
  security_groups = [sg.db]
  publicly_accessible = false
  
  backup_retention_days = 7
  backup_window = "03:00-04:00"
  maintenance_window = "sun:04:00-sun:05:00"
  
  performance_insights_enabled = true
  
  deletion_protection = true
  
  tags = { Name = "${var.app_name}-db" }
}

// S3 Bucket for Assets
resource storage.bucket "assets" {
  bucket_name = "${var.app_name}-assets-${var.environment}"
  region = "us-east-1"
  
  versioning = true
  encryption_enabled = true
  
  lifecycle_rules = [
    {
      action = "transition"
      age_days = 30
      storage_class = "STANDARD_IA"
    }
  ]
  
  tags = { Name = "${var.app_name}-assets" }
}

// CloudWatch Alarms
resource monitoring.alarm "high_cpu" {
  name = "${var.app_name}-high-cpu"
  description = "Alert when ASG CPU exceeds 80%"
  
  namespace = "AWS/EC2"
  metric_name = "CPUUtilization"
  statistic = "Average"
  period = 300
  evaluation_periods = 2
  threshold = 80
  comparison_operator = "GreaterThanThreshold"
  
  alarm_actions = [sns_topic.alerts.arn]
}

resource monitoring.alarm "db_connections" {
  name = "${var.app_name}-db-connections"
  description = "Alert when DB connections exceed 80"
  
  namespace = "AWS/RDS"
  metric_name = "DatabaseConnections"
  
  dimensions = {
    DBInstanceIdentifier = db.main.identifier
  }
  
  statistic = "Average"
  period = 300
  evaluation_periods = 2
  threshold = 80
  comparison_operator = "GreaterThanThreshold"
  
  alarm_actions = [sns_topic.alerts.arn]
}

// Outputs
output "alb_dns" {
  value = lb.main.dns_name
  description = "Load balancer DNS name"
}

output "db_endpoint" {
  value = db.main.endpoint
  sensitive = true
  description = "Database endpoint"
}

output "s3_bucket" {
  value = bucket.assets.name
  description = "Assets S3 bucket"
}
```

### Example 2: Serverless API with Lambda and DynamoDB

```aether
// Provider
provider aws {
  region = "us-east-1"
}

// DynamoDB Table
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
  
  attribute {
    name = "email"
    type = "S"
  }
  
  global_secondary_index {
    name = "EmailIndex"
    hash_key = "email"
    projection_type = "ALL"
  }
  
  ttl {
    attribute_name = "expiresAt"
    enabled = true
  }
  
  point_in_time_recovery {
    enabled = true
  }
  
  server_side_encryption {
    enabled = true
  }
  
  tags = {
    Name = "users-table"
  }
}

// IAM Role for Lambda
resource iam.role "lambda" {
  name = "api-lambda-role"
  
  assume_role_policy = <<-JSON
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Principal": {
          "Service": "lambda.amazonaws.com"
        },
        "Effect": "Allow"
      }
    ]
  }
  JSON
}

resource iam.policy "lambda_dynamodb" {
  name = "lambda-dynamodb-policy"
  
  policy = <<-JSON
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ],
        "Resource": [
          "${dynamodb_table.users.arn}",
          "${dynamodb_table.users.arn}/index/*"
        ]
      },
      {
        "Effect": "Allow",
        "Action": [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        "Resource": "*"
      }
    ]
  }
  JSON
}

resource iam.role_policy_attachment "lambda_policy" {
  role = role.lambda
  policy_arn = policy.lambda_dynamodb.arn
}

// Lambda Functions
resource compute.lambda_function "get_user" {
  function_name = "get-user"
  runtime = "python3.11"
  handler = "index.handler"
  
  source_code = file("lambda/get_user.py")
  
  role = role.lambda.arn
  
  memory = 512
  timeout = 30
  
  environment {
    variables = {
      TABLE_NAME = dynamodb_table.users.name
    }
  }
  
  tags = {
    Name = "get-user-function"
  }
}

resource compute.lambda_function "create_user" {
  function_name = "create-user"
  runtime = "python3.11"
  handler = "index.handler"
  
  source_code = file("lambda/create_user.py")
  
  role = role.lambda.arn
  
  memory = 512
  timeout = 30
  
  environment {
    variables = {
      TABLE_NAME = dynamodb_table.users.name
    }
  }
  
  tags = {
    Name = "create-user-function"
  }
}

// API Gateway
resource api.gateway "main" {
  name = "users-api"
  protocol_type = "HTTP"
  
  cors_configuration {
    allow_origins = ["*"]
    allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_headers = ["*"]
    max_age = 300
  }
  
  tags = {
    Name = "users-api"
  }
}

// API Gateway Routes
resource api.gateway_integration "get_user" {
  api = api_gateway.main
  integration_type = "AWS_PROXY"
  integration_uri = lambda.get_user.invoke_arn
  integration_method = "POST"
  payload_format_version = "2.0"
}

resource api.gateway_route "get_user" {
  api = api_gateway.main
  route_key = "GET /users/{userId}"
  target = "integrations/${api_integration.get_user.id}"
}

resource api.gateway_integration "create_user" {
  api = api_gateway.main
  integration_type = "AWS_PROXY"
  integration_uri = lambda.create_user.invoke_arn
  integration_method = "POST"
  payload_format_version = "2.0"
}

resource api.gateway_route "create_user" {
  api = api_gateway.main
  route_key = "POST /users"
  target = "integrations/${api_integration.create_user.id}"
}

// API Gateway Stage
resource api.gateway_stage "production" {
  api = api_gateway.main
  name = "production"
  auto_deploy = true
  
  access_log_settings {
    destination_arn = cloudwatch_log_group.api.arn
    format = "$context.requestId $context.error.message $context.error.messageString"
  }
  
  tags = {
    Environment = "production"
  }
}

// Lambda Permissions for API Gateway
resource lambda.permission "get_user" {
  function = lambda.get_user
  action = "lambda:InvokeFunction"
  principal = "apigateway.amazonaws.com"
  source_arn = "${api_gateway.main.execution_arn}/*/*"
}

resource lambda.permission "create_user" {
  function = lambda.create_user
  action = "lambda:InvokeFunction"
  principal = "apigateway.amazonaws.com"
  source_arn = "${api_gateway.main.execution_arn}/*/*"
}

// CloudWatch Logs
resource monitoring.log_group "api" {
  name = "/aws/apigateway/users-api"
  retention_in_days = 30
}

// Outputs
output "api_endpoint" {
  value = "${api_gateway.main.api_endpoint}/${api_stage.production.name}"
  description = "API endpoint URL"
}

output "dynamodb_table" {
  value = dynamodb_table.users.name
  description = "DynamoDB table name"
}
```

---

## AWS Best Practices

### Tagging Strategy

```aether
// Standardized tagging across all resources
variable "common_tags" {
  type = map<string, string>
  default = {
    ManagedBy = "Aether"
    Project = "myapp"
    Environment = "production"
    CostCenter = "engineering"
    Owner = "platform-team"
  }
}

// Apply to resources
resource compute.instance "app" {
  // ... configuration ...
  
  tags = merge(var.common_tags, {
    Name = "app-server"
    Role = "application"
  })
}
```

### Security Best Practices

```aether
// Use secrets manager for sensitive data
resource secrets.secret "db_password" {
  name = "db-password"
  description = "Database master password"
  
  rotation_rules {
    automatically_after_days = 30
  }
}

// Enable encryption at rest
resource storage.bucket "data" {
  encryption {
    sse_algorithm = "aws:kms"
    kms_master_key_id = kms_key.s3.id
  }
}

// Require IMDSv2 on EC2 instances
resource compute.instance "app" {
  metadata_options {
    http_endpoint = "enabled"
    http_tokens = "required"  // IMDSv2
    http_put_response_hop_limit = 1
  }
}

// Enable VPC Flow Logs
resource network.flow_log "main" {
  vpc = vpc.main
  traffic_type = "ALL"
  
  log_destination = cloudwatch_log_group.vpc_flow_logs.arn
  log_destination_type = "cloud-watch-logs"
  
  tags = {
    Name = "vpc-flow-logs"
  }
}
```

### Cost Optimization

```aether
// Use Spot instances for non-critical workloads
resource compute.instance "batch" {
  instance_market_options {
    market_type = "spot"
    spot_options {
      max_price = "0.15"
    }
  }
}

// S3 Intelligent-Tiering
resource storage.bucket "archives" {
  lifecycle_rules = [
    {
      action = "transition"
      age_days = 0
      storage_class = "INTELLIGENT_TIERING"
    }
  ]
}

// RDS Reserved Instances (manual reservation)
// Set instance class for consistent sizing
resource database.instance "prod" {
  instance_class = "large"  // Consistent for RI
}
```

---

## Additional Resources

- [AWS Documentation](https://docs.aws.amazon.com/)
- [AWS Well-Architected Framework](https://aws.amazon.com/architecture/well-architected/)
- [Aether Syntax Guide](syntax.md)
- [Aether Provider System](providers.md)
- [Aether Examples](../EXAMPLES.md)
