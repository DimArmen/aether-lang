# Aether Provider Architecture

This document describes how Aether abstracts cloud providers and enables portable infrastructure code.

---

## Overview

Aether's multi-cloud abstraction allows you to write infrastructure code once and deploy to AWS, Azure, GCP, or any other cloud provider. The system consists of:

1. **Universal Resource Types**: Common abstractions across all providers
2. **Provider Plugins**: Cloud-specific implementations
3. **Translation Layer**: Maps abstract resources to provider APIs
4. **Provider-Specific Overrides**: Access cloud-unique features when needed

---

## Universal Resource Types

### Compute Resources

#### `compute.instance`

Virtual machine/compute instance across all clouds.

```aether
resource compute.instance "app" {
  machine_type = "small" | "medium" | "large" | "xlarge" | "2xlarge"
  os_image = string  // e.g., "ubuntu-22.04", "windows-server-2022"
  region = string
  
  cpu_credits = "standard" | "unlimited"  // Optional
  
  network {
    vpc = reference
    subnet = reference
    security_groups = [reference]
    public_ip = bool
  }
  
  storage {
    root_volume_size = int  // GB
    root_volume_type = "ssd" | "hdd" | "provisioned-iops"
    
    additional_volumes = [
      {
        size = int
        type = string
        device_name = string
        encrypted = bool
      }
    ]
  }
  
  user_data = string
  ssh_key = string
  
  tags = map<string, string>
  
  monitoring_enabled = bool
  auto_scaling_group = reference  // Optional
}
```

**Provider Mappings**:
- AWS: EC2 instance
- Azure: Virtual Machine
- GCP: Compute Engine instance

**Machine Type Mappings**:
| Aether | AWS | Azure | GCP |
|--------|-----|-------|-----|
| small | t3.small | Standard_B1s | e2-small |
| medium | t3.medium | Standard_B2s | e2-medium |
| large | t3.large | Standard_B4ms | e2-standard-2 |
| xlarge | t3.xlarge | Standard_D4s_v3 | e2-standard-4 |
| 2xlarge | t3.2xlarge | Standard_D8s_v3 | e2-standard-8 |

### Storage Resources

#### `storage.bucket`

Object storage bucket.

```aether
resource storage.bucket "data" {
  region = string
  
  versioning = bool
  encryption_enabled = bool
  
  public_access = bool
  
  lifecycle_rules = [
    {
      action = "delete" | "archive" | "transition"
      age_days = int
      prefix = string  // Optional
    }
  ]
  
  access_policy = string  // JSON policy
  
  cors_rules = [
    {
      allowed_origins = [string]
      allowed_methods = [string]
      allowed_headers = [string]
      max_age_seconds = int
    }
  ]
  
  tags = map<string, string>
}
```

**Provider Mappings**:
- AWS: S3 Bucket
- Azure: Blob Storage Container
- GCP: Cloud Storage Bucket

#### `storage.volume`

Block storage volume.

```aether
resource storage.volume "data" {
  size = int  // GB
  type = "ssd" | "hdd" | "provisioned-iops"
  region = string
  availability_zone = string
  
  encrypted = bool
  snapshot_id = string  // Optional, restore from snapshot
  
  iops = int  // For provisioned-iops type
  throughput = int  // MB/s
  
  tags = map<string, string>
}
```

### Network Resources

#### `network.vpc`

Virtual Private Cloud/Network.

```aether
resource network.vpc "main" {
  cidr = string  // e.g., "10.0.0.0/16"
  region = string
  
  enable_dns = bool
  enable_dns_hostnames = bool
  
  tags = map<string, string>
}
```

**Provider Mappings**:
- AWS: VPC
- Azure: Virtual Network (VNet)
- GCP: VPC Network

#### `network.subnet`

Subnet within a VPC.

```aether
resource network.subnet "app" {
  vpc = reference
  cidr = string  // e.g., "10.0.1.0/24"
  availability_zone = string
  
  public = bool  // Public vs private subnet
  
  route_table = reference  // Optional
  
  tags = map<string, string>
}
```

#### `network.security_group`

Firewall rules for resources.

```aether
resource network.security_group "web" {
  vpc = reference
  description = string
  
  ingress = [
    {
      from_port = int
      to_port = int
      protocol = "tcp" | "udp" | "icmp" | "-1"
      cidr_blocks = [string]
      security_groups = [reference]  // Source security groups
      description = string
    }
  ]
  
  egress = [
    {
      from_port = int
      to_port = int
      protocol = string
      cidr_blocks = [string]
      description = string
    }
  ]
  
  tags = map<string, string>
}
```

**Provider Mappings**:
- AWS: Security Group
- Azure: Network Security Group (NSG)
- GCP: Firewall Rules

### Database Resources

#### `database.instance`

Managed database instance.

```aether
resource database.instance "main" {
  engine = "postgres" | "mysql" | "mariadb" | "sqlserver" | "oracle"
  engine_version = string
  
  instance_class = "small" | "medium" | "large" | "xlarge"
  
  storage = int  // GB
  storage_type = "ssd" | "hdd" | "provisioned-iops"
  storage_encrypted = bool
  
  database_name = string
  master_username = string
  master_password = string | secret
  
  multi_az = bool
  
  vpc = reference
  subnets = [reference]
  security_groups = [reference]
  
  backup_retention_days = int
  backup_window = string  // e.g., "03:00-04:00"
  maintenance_window = string  // e.g., "sun:04:00-sun:05:00"
  
  publicly_accessible = bool
  
  parameter_group = reference  // Optional
  
  tags = map<string, string>
}
```

**Provider Mappings**:
- AWS: RDS Instance
- Azure: Azure Database
- GCP: Cloud SQL Instance

### Load Balancer Resources

#### `loadbalancer.application`

Application (Layer 7) load balancer.

```aether
resource loadbalancer.application "web" {
  vpc = reference
  subnets = [reference]
  security_groups = [reference]
  
  internal = bool  // Internal vs internet-facing
  
  listener = [
    {
      port = int
      protocol = "HTTP" | "HTTPS"
      ssl_certificate = reference  // For HTTPS
      
      default_action {
        type = "forward" | "redirect" | "fixed-response"
        target_group = reference  // For forward
        
        // For redirect
        redirect {
          protocol = string
          port = int
          status_code = int
        }
      }
    }
  ]
  
  tags = map<string, string>
}
```

#### `loadbalancer.target_group`

Target group for load balancer.

```aether
resource loadbalancer.target_group "app" {
  port = int
  protocol = "HTTP" | "HTTPS" | "TCP"
  vpc = reference
  
  health_check {
    path = string  // For HTTP/HTTPS
    port = int
    protocol = string
    interval = int  // seconds
    timeout = int
    healthy_threshold = int
    unhealthy_threshold = int
  }
  
  stickiness {
    enabled = bool
    type = "lb_cookie" | "app_cookie"
    duration = int  // seconds
  }
  
  tags = map<string, string>
}
```

### Function Resources

#### `function.function`

Serverless function.

```aether
resource function.function "handler" {
  name = string
  runtime = string  // e.g., "python3.11", "nodejs20", "go1.21"
  handler = string  // Entry point
  
  code_path = string  // Path to code
  code_s3_bucket = string  // Alternative: S3 location
  code_s3_key = string
  
  memory = int  // MB
  timeout = int  // seconds
  
  environment_variables = map<string, string>
  
  vpc = reference  // Optional
  subnets = [reference]
  security_groups = [reference]
  
  iam_role = reference
  
  triggers = [
    {
      type = "http" | "storage" | "queue" | "schedule"
      
      // For HTTP trigger
      path = string
      method = string
      
      // For storage trigger
      bucket = reference
      events = [string]
      prefix = string
      
      // For schedule trigger
      schedule = string  // Cron expression
    }
  ]
  
  tags = map<string, string>
}
```

---

## Provider Plugin Interface

### Provider Interface (Go)

```go
package provider

type Provider interface {
    // Initialize the provider with configuration
    Configure(config map[string]any) error
    
    // Get provider metadata
    GetSchema() (*ProviderSchema, error)
    GetName() string
    GetVersion() string
    
    // Resource operations
    ValidateResourceConfig(resourceType string, config map[string]any) error
    CreateResource(resourceType string, config map[string]any) (*Resource, error)
    ReadResource(resourceType string, id string) (*Resource, error)
    UpdateResource(resource *Resource, changes map[string]any) (*Resource, error)
    DeleteResource(resource *Resource) error
    
    // Import existing resources
    ImportResource(resourceType string, id string) (*Resource, error)
    
    // Provider-specific queries
    GetAvailabilityZones(region string) ([]string, error)
    GetResourcePrice(resourceType string, config map[string]any) (*PriceEstimate, error)
}

type Resource struct {
    ID         string
    Type       string
    Provider   string
    Properties map[string]any
    Metadata   ResourceMetadata
    State      ResourceState
}

type ResourceMetadata struct {
    CreatedAt time.Time
    UpdatedAt time.Time
    CreatedBy string
    Tags      map[string]string
}

type ResourceState string

const (
    StateCreating  ResourceState = "creating"
    StateAvailable ResourceState = "available"
    StateModifying ResourceState = "modifying"
    StateDeleting  ResourceState = "deleting"
    StateDeleted   ResourceState = "deleted"
    StateError     ResourceState = "error"
)

type ProviderSchema struct {
    Provider      string
    Version       string
    ResourceTypes map[string]ResourceTypeSchema
}

type ResourceTypeSchema struct {
    Type        string
    Description string
    Properties  map[string]PropertySchema
    Required    []string
}

type PropertySchema struct {
    Type        string  // string, int, bool, list, map, object
    Description string
    Required    bool
    Default     any
    Enum        []string
    Pattern     string  // Regex for validation
    MinLength   int
    MaxLength   int
    Minimum     int
    Maximum     int
}
```

### Provider Registration

```go
package main

import (
    "github.com/aether-lang/aether/provider"
    awsprovider "github.com/aether-lang/aether-provider-aws"
    azureprovider "github.com/aether-lang/aether-provider-azure"
    gcpprovider "github.com/aether-lang/aether-provider-gcp"
)

func init() {
    provider.Register("aws", awsprovider.New())
    provider.Register("azure", azureprovider.New())
    provider.Register("gcp", gcpprovider.New())
}
```

---

## Provider-Specific Overrides

When you need cloud-specific features:

```aether
resource compute.instance "app" {
  // Generic configuration
  machine_type = "medium"
  os_image = "ubuntu-22.04"
  region = "us-east"
  
  // AWS-specific overrides
  aws {
    instance_type = "t3.medium"  // Override generic type
    ebs_optimized = true
    iam_instance_profile = "app-role"
    placement_group = "my-group"
    
    // AWS-only feature
    hibernation_enabled = true
  }
  
  // Azure-specific overrides
  azure {
    vm_size = "Standard_B2s"
    priority = "Spot"  // Azure spot instances
    eviction_policy = "Deallocate"
    
    // Azure-only feature
    availability_set = "app-availability-set"
  }
  
  // GCP-specific overrides
  gcp {
    machine_type = "e2-medium"
    preemptible = true  // GCP preemptible instances
    
    // GCP-only feature
    scheduling {
      automatic_restart = false
      on_host_maintenance = "TERMINATE"
    }
  }
}
```

---

## Translation Layer

The translation layer converts universal resources to provider-specific API calls.

### Example: compute.instance → AWS EC2

```go
func (p *AWSProvider) CreateResource(resourceType string, config map[string]any) (*Resource, error) {
    switch resourceType {
    case "compute.instance":
        return p.createEC2Instance(config)
    // ... other resource types
    }
}

func (p *AWSProvider) createEC2Instance(config map[string]any) (*Resource, error) {
    // Translate generic machine_type to AWS instance type
    instanceType := p.translateMachineType(config["machine_type"].(string))
    
    // Translate generic os_image to AWS AMI
    ami, err := p.resolveAMI(config["os_image"].(string), config["region"].(string))
    if err != nil {
        return nil, err
    }
    
    // Build EC2 RunInstances request
    input := &ec2.RunInstancesInput{
        ImageId:      aws.String(ami),
        InstanceType: aws.String(instanceType),
        MinCount:     aws.Int32(1),
        MaxCount:     aws.Int32(1),
    }
    
    // Handle networking
    if network := config["network"]; network != nil {
        net := network.(map[string]any)
        input.SubnetId = aws.String(net["subnet"].(string))
        input.SecurityGroupIds = translateSecurityGroups(net["security_groups"])
    }
    
    // Handle AWS-specific overrides
    if awsConfig := config["aws"]; awsConfig != nil {
        aws := awsConfig.(map[string]any)
        if ebsOpt := aws["ebs_optimized"]; ebsOpt != nil {
            input.EbsOptimized = aws.Bool(ebsOpt.(bool))
        }
        if profile := aws["iam_instance_profile"]; profile != nil {
            input.IamInstanceProfile = &ec2.IamInstanceProfileSpecification{
                Name: aws.String(profile.(string)),
            }
        }
    }
    
    // Call AWS API
    result, err := p.ec2Client.RunInstances(context.TODO(), input)
    if err != nil {
        return nil, err
    }
    
    instance := result.Instances[0]
    
    // Convert to Aether resource
    return &Resource{
        ID:       *instance.InstanceId,
        Type:     "compute.instance",
        Provider: "aws",
        Properties: map[string]any{
            "public_ip":  instance.PublicIpAddress,
            "private_ip": instance.PrivateIpAddress,
            "state":      string(instance.State.Name),
        },
        State: StateCreating,
    }, nil
}

func (p *AWSProvider) translateMachineType(machineType string) string {
    mapping := map[string]string{
        "small":   "t3.small",
        "medium":  "t3.medium",
        "large":   "t3.large",
        "xlarge":  "t3.xlarge",
        "2xlarge": "t3.2xlarge",
    }
    return mapping[machineType]
}
```

---

## Provider SDK

For creating custom providers:

### Directory Structure

```
aether-provider-mycloud/
  go.mod
  provider.go
  resource_instance.go
  resource_storage.go
  resource_network.go
  schema.json
  README.md
```

### Example Provider Implementation

```go
package mycloud

import (
    "github.com/aether-lang/aether/provider"
)

type MyCloudProvider struct {
    client *MyCloudClient
    config map[string]any
}

func New() provider.Provider {
    return &MyCloudProvider{}
}

func (p *MyCloudProvider) Configure(config map[string]any) error {
    p.config = config
    
    apiKey := config["api_key"].(string)
    region := config["region"].(string)
    
    client, err := NewMyCloudClient(apiKey, region)
    if err != nil {
        return err
    }
    
    p.client = client
    return nil
}

func (p *MyCloudProvider) GetSchema() (*provider.ProviderSchema, error) {
    // Load from embedded schema.json
    return loadSchema()
}

func (p *MyCloudProvider) CreateResource(resourceType string, config map[string]any) (*provider.Resource, error) {
    switch resourceType {
    case "compute.instance":
        return p.createInstance(config)
    case "storage.bucket":
        return p.createBucket(config)
    default:
        return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
    }
}

// Implement other interface methods...
```

---

## Provider Discovery and Loading

### Dynamic Loading

```go
// Load providers from plugins directory
func LoadProviders() error {
    pluginsDir := filepath.Join(configDir, "plugins")
    
    files, err := os.ReadDir(pluginsDir)
    if err != nil {
        return err
    }
    
    for _, file := range files {
        if !strings.HasPrefix(file.Name(), "aether-provider-") {
            continue
        }
        
        pluginPath := filepath.Join(pluginsDir, file.Name())
        provider, err := plugin.Open(pluginPath)
        if err != nil {
            log.Warnf("Failed to load provider %s: %v", file.Name(), err)
            continue
        }
        
        // Load provider symbol
        symbolProvider, err := provider.Lookup("Provider")
        if err != nil {
            log.Warnf("Provider %s doesn't export Provider symbol", file.Name())
            continue
        }
        
        providerInstance := symbolProvider.(provider.Provider)
        name := providerInstance.GetName()
        
        provider.Register(name, providerInstance)
        log.Infof("Loaded provider: %s", name)
    }
    
    return nil
}
```

---

## Multi-Cloud Deployment Example

```aether
// Same infrastructure on multiple providers

variable "provider" {
  type = string
  default = "aws"
}

provider "aws" {
  region = "us-east-1"
}

provider "azure" {
  region = "eastus"
}

provider "gcp" {
  project = "my-project"
  region = "us-east1"
}

// Deploy to selected provider
resource compute.instance "app" {
  provider = var.provider
  
  machine_type = "medium"
  os_image = "ubuntu-22.04"
  
  tags = {
    Name = "app-server"
    Provider = var.provider
  }
}

// Deploy to all providers simultaneously
resource storage.bucket "backup_aws" {
  provider = aws
  region = "us-east"
}

resource storage.bucket "backup_azure" {
  provider = azure
  region = "us-east"
}

resource storage.bucket "backup_gcp" {
  provider = gcp
  region = "us-east"
}
```

---

This provider architecture enables true multi-cloud portability while preserving access to cloud-specific features when needed. The plugin system allows community-contributed providers for any cloud platform or service.
