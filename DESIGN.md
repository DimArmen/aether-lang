# Aether: Technical Design Document

## Language Design Philosophy

Aether is a **hybrid declarative-scripting language** with these core characteristics:

- **Declarative Core**: Resource definitions express desired state
- **Embedded Scripting**: Complex logic uses familiar programming constructs
- **Strong Type System**: Structural typing with inference
- **Interpreted with JIT**: Fast iteration with optimized execution
- **Multi-Cloud Native**: Abstract cloud providers from the start
- **AI-Integrated**: Agents are language-level constructs

## Syntax Overview

### Resource Declarations (Declarative)

```aether
resource compute.instance "web_server" {
  machine_type = "small"
  os_image = "ubuntu-22.04"
  region = "us-east"
  
  network {
    vpc = vpc.main
    security_groups = [sg.web]
  }
  
  tags = {
    Environment = "production"
    ManagedBy = "aether"
  }
}
```

### Scripting Blocks (Imperative)

```aether
script {
  // Full programming constructs available
  let instance_count = env.is_production ? 5 : 2
  let regions = ["us-east", "us-west", "eu-west"]
  
  for region in regions {
    if region.needs_deployment() {
      deploy_resources(region)
    }
  }
}
```

### Variables and Outputs

```aether
variable "environment" {
  type = string
  default = "development"
  description = "Deployment environment"
}

output "web_url" {
  value = load_balancer.web.public_ip
  description = "Public URL for web application"
}
```

### Modules

```aether
module "web_app" {
  source = "stdlib/web-app@v1.2"
  
  instance_count = 3
  database_tier = "high-performance"
  region = "us-east"
}
```

### AI Agent Definitions

```aether
agent cost_optimizer {
  type = "autonomous"
  scope = resource.compute.*
  
  goals = [
    "minimize_cost",
    "maintain_performance"
  ]
  
  constraints {
    max_cost_increase = "$100/day"
    min_instances = 2
  }
  
  approval {
    required_for = changes_exceeding("$50/day")
    notify = ["ops-team@example.com"]
  }
  
  learning_period = duration("7d")
}

agent security_scanner {
  type = "analyzer"
  
  checks = [
    "unencrypted_storage",
    "open_security_groups",
    "weak_passwords",
    "compliance.cis_benchmark"
  ]
  
  severity_threshold = "medium"
  block_on_failure = true
}
```

### Comments and AI Hints

```aether
// Single-line comment

/*
 * Multi-line comment
 */

// @ai:optimize - AI should suggest optimizations for this resource
// @ai:secure - AI should apply security best practices
// @ai:explain - AI should add documentation
resource storage.bucket "assets" {
  // @ai:optimize
  replication = "multi-region"
}
```

## Type System

### Primitive Types

- `string`: UTF-8 text
- `number`: 64-bit floating point
- `int`: 64-bit integer
- `bool`: true/false
- `duration`: Time duration (e.g., `duration("5m")`, `duration("2h")`)
- `bytes`: Binary data

### Collection Types

- `list<T>`: Ordered collection
- `map<K, V>`: Key-value mapping
- `set<T>`: Unordered unique collection

### Resource Types

Resources have typed properties based on provider schemas:

```aether
type ComputeInstance = {
  machine_type: string
  os_image: string
  region: string
  network: NetworkConfig
  tags: map<string, string>
}
```

### Type Inference

```aether
let count = 5                    // Inferred: int
let regions = ["us", "eu"]       // Inferred: list<string>
let config = {port = 80}         // Inferred: {port: int}
```

### Structural Typing

Types are compatible based on structure, not names:

```aether
type ServerConfig = {
  port: int
  host: string
}

type WebConfig = {
  port: int
  host: string
  ssl: bool
}

// WebConfig is compatible with ServerConfig (has all required fields)
function setup_server(config: ServerConfig) { ... }
setup_server({port = 443, host = "example.com", ssl = true})  // Valid
```

## Execution Model

### Phases

1. **Parse**: Source → AST
2. **Type Check**: Validate types, catch errors
3. **Plan**: Determine required changes (like Terraform plan)
4. **Agent Analysis**: Run Tier 2 agents (security, cost, etc.)
5. **Approval**: Present changes for confirmation
6. **Apply**: Execute changes with Tier 3 agents monitoring
7. **State Update**: Record new state

### Dependency Resolution

Resources automatically determine dependencies from references:

```aether
resource network.vpc "main" {
  cidr = "10.0.0.0/16"
}

resource network.subnet "app" {
  vpc = vpc.main           // Implicit dependency
  cidr = "10.0.1.0/24"
}

// Aether creates vpc.main before network.subnet.app
```

Explicit dependencies when needed:

```aether
resource compute.instance "app" {
  depends_on = [database.main]
}
```

### Parallel Execution

Independent resources are created in parallel for speed:

```aether
// These can be created simultaneously
resource storage.bucket "assets" { ... }
resource storage.bucket "logs" { ... }
resource storage.bucket "backups" { ... }
```

### Error Handling

```aether
resource database.instance "main" {
  on_error {
    retry = 3
    backoff = "exponential"
  }
  
  on_failure {
    notify = ["ops@example.com"]
    rollback = true
  }
}
```

## Multi-Cloud Provider Abstraction

### Universal Resource Types

```aether
// Generic compute instance
resource compute.instance "app" {
  machine_type = "small"      // Maps to t3.small (AWS), B1s (Azure), e2-small (GCP)
  os_image = "ubuntu-22.04"
  region = "us-east"
}

// Generic storage bucket
resource storage.bucket "data" {
  region = "us-east"
  versioning = true
  lifecycle_rules = [...]
}

// Generic database
resource database.instance "main" {
  engine = "postgres"
  version = "15"
  instance_class = "medium"
}
```

### Provider Configuration

```aether
provider "aws" {
  region = "us-east-1"
  profile = "production"
}

provider "azure" {
  region = "eastus"
  subscription = "..."
}

provider "gcp" {
  project = "my-project"
  region = "us-east1"
}

// Choose provider per resource
resource compute.instance "aws_app" {
  provider = aws
  machine_type = "small"
}

resource compute.instance "azure_app" {
  provider = azure
  machine_type = "small"
}
```

### Provider-Specific Features

```aether
resource compute.instance "app" {
  machine_type = "small"
  
  // AWS-specific features
  aws {
    instance_type = "t3.small"  // Override generic mapping
    iam_role = "..."
    ebs_optimized = true
  }
  
  // Azure-specific features
  azure {
    vm_size = "Standard_B1s"
    availability_set = "..."
  }
}
```

### Provider Plugin Architecture

Providers implement standardized interfaces:

```go
type Provider interface {
  Configure(config map[string]any) error
  Schema() ResourceSchemas
  Create(resourceType string, config map[string]any) (Resource, error)
  Read(resourceType string, id string) (Resource, error)
  Update(resource Resource, changes map[string]any) error
  Delete(resource Resource) error
}
```

## State Management

### State Structure

```json
{
  "version": 2,
  "resources": [
    {
      "id": "vpc.main",
      "type": "network.vpc",
      "provider": "aws",
      "properties": {
        "cidr": "10.0.0.0/16",
        "region": "us-east-1"
      },
      "metadata": {
        "created_at": "2026-02-23T10:30:00Z",
        "updated_at": "2026-02-23T10:30:00Z",
        "created_by": "user@example.com"
      },
      "dependencies": []
    }
  ],
  "outputs": {
    "vpc_id": "vpc-12345"
  }
}
```

### State Backends

```aether
backend "s3" {
  bucket = "my-terraform-state"
  key = "production/infrastructure"
  region = "us-east-1"
  
  lock_table = "aether-locks"  // DynamoDB table
  encrypt = true
}

backend "azure" {
  storage_account = "aetherstate"
  container = "tfstate"
  key = "production.tfstate"
}

backend "gcp" {
  bucket = "aether-state"
  prefix = "production"
}
```

### Automatic State Features

- **Locking**: Automatic acquisition and release
- **Encryption**: Sensitive values encrypted at rest
- **Versioning**: Every state change versioned for rollback
- **Drift Detection**: Periodic comparison with actual resources
- **Conflict Resolution**: Merge strategies for concurrent changes

## AI Agent Architecture

### Agent Types

#### Tier 1: Assistant (IDE Integration)

```typescript
interface Assistant {
  complete(context: CodeContext): Suggestion[]
  explain(selection: CodeRange): Explanation
  refactor(code: string, intent: string): CodeChange[]
  generateFromIntent(intent: string): string
}
```

#### Tier 2: Analyzer (Pre-Deployment)

```typescript
interface Analyzer {
  analyze(plan: DeploymentPlan): Analysis
  
  type Analysis = {
    issues: Issue[]
    suggestions: Suggestion[]
    compliance: ComplianceReport
  }
}
```

#### Tier 3: Autonomous (Runtime)

```typescript
interface AutonomousAgent {
  observe(metrics: Metrics): void
  decide(state: InfrastructureState): Action[]
  execute(action: Action): Result
  
  // Approval workflow
  requiresApproval(action: Action): bool
  requestApproval(action: Action): ApprovalRequest
}
```

### Agent Execution Flow

```
1. Code written → Assistant suggests improvements
2. Plan generated → Analyzers scan for issues
3. Issues found → Block or warn user
4. Approved → Apply changes
5. Runtime → Autonomous agents monitor
6. Agent detects optimization → Request approval
7. Approved → Agent executes change
8. State updated → Audit log recorded
```

### Agent Communication Protocol

Agents communicate via event bus:

```aether
event ResourceCreated {
  resource_id: string
  resource_type: string
  properties: map<string, any>
}

event CostThresholdExceeded {
  current_cost: number
  threshold: number
  resources: list<string>
}

agent cost_watcher {
  on ResourceCreated {
    // Track new resource costs
  }
  
  on CostThresholdExceeded {
    // Take action
  }
}
```

### Approval Policies

```aether
approval_policy "production" {
  // Auto-approve small changes
  auto_approve {
    cost_change < "$10/day"
    resource_count_change < 5
    no_deletions = true
  }
  
  // Require human approval for large changes
  require_approval {
    cost_change >= "$10/day"
    resource_deletions > 0
    security_group_changes = true
  }
  
  // Multi-person approval for critical changes
  require_multi_approval {
    approvers = 2
    resource_deletions > 10
    production_database_changes = true
  }
  
  notification {
    slack = "#ops-approvals"
    email = ["ops@example.com"]
  }
}
```

## Module System

### Module Structure

```
my-web-app/
  aether.mod       # Module metadata
  main.ae          # Main module file
  variables.ae     # Input variables
  outputs.ae       # Output values
  README.md        # Documentation
```

### Module Definition

```aether
// aether.mod
module "web-app" {
  version = "1.2.0"
  author = "ops-team"
  license = "MIT"
  
  dependencies = [
    "stdlib/networking@1.0",
    "stdlib/compute@2.1"
  ]
}
```

### Using Modules

```aether
module "app" {
  source = "github.com/company/web-app@1.2.0"
  
  instance_count = 5
  region = "us-east"
  database_size = "large"
}

// Access module outputs
output "app_url" {
  value = module.app.load_balancer_url
}
```

## Built-in Functions

### String Functions
- `upper(s: string) -> string`
- `lower(s: string) -> string`
- `trim(s: string) -> string`
- `split(s: string, sep: string) -> list<string>`
- `join(list: list<string>, sep: string) -> string`

### Collection Functions
- `length(collection: list<T> | map<K,V>) -> int`
- `contains(collection: list<T>, item: T) -> bool`
- `filter(list: list<T>, predicate: T -> bool) -> list<T>`
- `map(list: list<T>, transform: T -> U) -> list<U>`

### Cloud Functions
- `region_available(provider: string, region: string) -> bool`
- `resource_exists(id: string) -> bool`
- `get_resource_property(id: string, property: string) -> any`

### Crypto Functions
- `hash(data: string, algorithm: string) -> string`
- `encrypt(data: string, key: string) -> bytes`
- `decrypt(data: bytes, key: string) -> string`

## Testing Framework

### Unit Tests

```aether
test "web_server_config" {
  let server = resource.compute.instance.web_server
  
  assert server.machine_type == "small"
  assert server.region in ["us-east", "us-west"]
  assert length(server.tags) > 0
}
```

### Integration Tests

```aether
integration_test "deploy_web_app" {
  environment = "test"
  
  // Deploy infrastructure
  plan = aether.plan()
  result = aether.apply(plan)
  
  // Test deployed resources
  assert http.get(output.web_url).status == 200
  assert database.query("SELECT 1").success == true
  
  // Cleanup
  defer aether.destroy()
}
```

### Property-Based Tests

```aether
property_test "all_storage_encrypted" {
  for resource in resources where resource.type == "storage.*" {
    assert resource.encryption_enabled == true
  }
}
```

## CLI Command Reference

```bash
# Initialize project
aether init [--template=<name>]

# Validate syntax and types
aether validate

# Generate execution plan
aether plan [--out=plan.json]

# Apply changes
aether apply [plan.json]

# Destroy infrastructure
aether destroy [--target=<resource>]

# Import existing resources
aether import <resource_address> <resource_id>

# Manage state
aether state list
aether state show <resource>
aether state mv <source> <destination>
aether state rm <resource>

# Agent management
aether agent list
aether agent enable <agent_name>
aether agent disable <agent_name>
aether agent logs <agent_name>

# Testing
aether test [--filter=<pattern>]

# Format code
aether fmt

# Generate documentation
aether docs

# Migration tools
aether migrate terraform <dir>
aether migrate cloudformation <template>
```

## Implementation Language

**Recommendation: Go**

Rationale:
- Excellent standard library for parsing, networking, concurrency
- Strong type system aligns with Aether's design
- Fast compilation and execution
- Good cross-platform support
- Popular in DevOps tools (Terraform, Kubernetes built with Go)
- Easy to embed as CLI tool

Alternative: Rust (for maximum performance and safety, but steeper development curve)

## Performance Considerations

### JIT Compilation

Hot paths identified and compiled:
- Frequently used modules
- Complex script blocks
- Tight loops in resource generation

### Caching

- Provider schemas cached locally
- Module downloads cached
- Type checking results cached between runs
- State snapshots for fast rollback

### Parallelism

- Independent resources created in parallel
- Provider API calls concurrent
- State backend operations batched

## Security

### Secrets Management

```aether
secret "database_password" {
  source = "vault"
  path = "secret/database/prod"
}

resource database.instance "main" {
  password = secret.database_password
}
```

### Sandboxing

- Modules run in restricted environment
- File system access controlled
- Network access restricted
- Agent permissions scoped

### Audit Logging

Every operation logged:
- User/service account
- Timestamp
- Operation type
- Resources affected
- Approval chain

---

This design provides a comprehensive foundation for implementing Aether. The hybrid syntax balances declarative simplicity with scripting power, while multi-tier AI agents provide unprecedented intelligence in infrastructure management.
