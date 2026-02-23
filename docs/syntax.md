# Aether Language Syntax Specification

## Version 0.1.0

This document provides the formal specification for the Aether programming language syntax.

---

## Table of Contents

1. [Lexical Structure](#lexical-structure)
2. [Types](#types)
3. [Expressions](#expressions)
4. [Statements](#statements)
5. [Resource Declarations](#resource-declarations)
6. [Variables and Outputs](#variables-and-outputs)
7. [Modules](#modules)
8. [Agents](#agents)
9. [Scripts](#scripts)
10. [Comments and Annotations](#comments-and-annotations)

---

## Lexical Structure

### Keywords

Reserved keywords in Aether:

```
resource    variable    output      module      agent       provider
script      for         in          if          else        match
let         function    return      import      export      test
property_test   integration_test   assert      defer       wait_for
true        false       null        and         or          not
where       contains    duration    secret      backend     count
each        depends_on  on_error    on_failure  approval    type
```

### Identifiers

```ebnf
identifier = letter { letter | digit | "_" }
letter = "a" ... "z" | "A" ... "Z"
digit = "0" ... "9"
```

Examples: `my_resource`, `vpc123`, `CamelCase`, `snake_case`

### Literals

#### String Literals

```aether
"simple string"
"string with \"escapes\""
'single quoted string'

// Multi-line strings
<<-EOF
  Multi-line
  string content
EOF

// String interpolation
"Hello ${name}!"
"Total: ${count * price}"
```

#### Number Literals

```aether
42          // Integer
3.14        // Float
1e6         // Scientific notation
0x2A        // Hexadecimal
0o52        // Octal
0b101010    // Binary
```

#### Boolean Literals

```aether
true
false
```

#### Collection Literals

```aether
[1, 2, 3]                    // List
["a", "b", "c"]              // List of strings
{key = "value"}              // Map
{name = "John", age = 30}    // Map with multiple entries
```

### Operators

#### Arithmetic
```
+  -  *  /  %  **  (power)
```

#### Comparison
```
==  !=  <  >  <=  >=
```

#### Logical
```
and  or  not  &&  ||  !
```

#### Other
```
=   (assignment)
.   (member access)
?.  (optional chaining)
::  (namespace separator)
=>  (lambda arrow)
|>  (pipe operator)
```

### Punctuation

```
( )  [ ]  { }  ,  ;  :  ...
```

---

## Types

### Primitive Types

```aether
string      // UTF-8 text
number      // 64-bit float
int         // 64-bit integer
bool        // Boolean
duration    // Time duration
bytes       // Binary data
any         // Any type (use sparingly)
```

### Collection Types

```aether
list<T>              // Ordered collection
map<K, V>            // Key-value mapping
set<T>               // Unordered unique collection
```

### Type Declarations

```aether
type ServerConfig = {
  port: int
  host: string
  ssl_enabled: bool
}

type Environment = "development" | "staging" | "production"  // Union type

type Result<T, E> = {
  success: bool
  value?: T
  error?: E
}
```

### Type Annotations

```aether
let name: string = "Aether"
let count: int = 42
let items: list<string> = ["a", "b", "c"]

function greet(name: string): string {
  return "Hello ${name}"
}
```

### Type Inference

```aether
let name = "Aether"          // Inferred: string
let count = 42               // Inferred: int
let items = ["a", "b"]       // Inferred: list<string>
let config = {port = 80}     // Inferred: {port: int}
```

---

## Expressions

### Binary Expressions

```aether
a + b
a - b
a * b
a / b
a % b
a ** b        // Power

a == b
a != b
a < b
a > b
a <= b
a >= b

a and b
a or b
```

### Unary Expressions

```aether
-a            // Negation
not a         // Logical NOT
!a            // Logical NOT (alternative)
```

### Conditional Expression (Ternary)

```aether
condition ? true_value : false_value

// Example
let size = is_production ? "large" : "small"
```

### Member Access

```aether
object.property
object.method()
object?.optional_property    // Optional chaining
```

### Index Access

```aether
list[0]
map["key"]
list[1:3]        // Slice
```

### Function Call

```aether
function_name(arg1, arg2)
object.method(arg1, arg2)
```

### Lambda Expression

```aether
(x) => x * 2
(x, y) => x + y
() => 42
```

### Pipe Expression

```aether
value |> function1 |> function2 |> function3

// Example
"hello" 
  |> upper 
  |> trim 
  |> split(" ")
```

---

## Statements

### Variable Declaration

```aether
let name = "value"
let count = 42
let items = ["a", "b", "c"]

// With type annotation
let name: string = "value"
let count: int = 42
```

### Assignment

```aether
name = "new value"
count = count + 1
items[0] = "x"
```

### If Statement

```aether
if condition {
  // code
}

if condition {
  // code
} else {
  // code
}

if condition1 {
  // code
} else if condition2 {
  // code
} else {
  // code
}
```

### Match Expression

```aether
match value {
  pattern1 => result1
  pattern2 => result2
  _ => default_result
}

// Example
let size = match environment {
  "production" => "large"
  "staging" => "medium"
  _ => "small"
}
```

### For Loop

```aether
for item in items {
  // code
}

for index, item in enumerate(items) {
  // code
}

for key, value in map {
  // code
}

// With range
for i in range(10) {
  // code
}

for i in range(1, 10) {
  // code with i from 1 to 9
}
```

### While Loop

```aether
while condition {
  // code
}
```

### Function Definition

```aether
function name(param1, param2) {
  // code
  return result
}

// With types
function greet(name: string): string {
  return "Hello ${name}!"
}

// With default parameters
function greet(name: string = "World"): string {
  return "Hello ${name}!"
}
```

### Return Statement

```aether
return value
return  // return void
```

### Export Statement

```aether
export { name, count, items }
export name
```

---

## Resource Declarations

### Basic Resource

```aether
resource resource_type "name" {
  property1 = value1
  property2 = value2
}

// Example
resource compute.instance "web_server" {
  machine_type = "small"
  os_image = "ubuntu-22.04"
  region = "us-east"
}
```

### Resource with Count

```aether
resource compute.instance "web" {
  count = 3
  
  machine_type = "small"
  
  tags = {
    Name = "web-${count.index}"
  }
}

// Reference: instance.web[0], instance.web[1], instance.web[2]
```

### Resource with For Each

```aether
resource compute.instance "servers" {
  for_each = {
    web = "small"
    api = "medium"
    worker = "large"
  }
  
  machine_type = each.value
  
  tags = {
    Name = each.key
  }
}

// Reference: instance.servers["web"], instance.servers["api"]
```

### Resource with Dependencies

```aether
resource network.vpc "main" {
  cidr = "10.0.0.0/16"
}

resource network.subnet "app" {
  vpc = vpc.main        // Implicit dependency
  cidr = "10.0.1.0/24"
}

resource compute.instance "app" {
  subnet = subnet.app
  depends_on = [database.main]  // Explicit dependency
}
```

### Resource with Nested Blocks

```aether
resource network.security_group "web" {
  vpc = vpc.main
  
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
```

### Resource with Lifecycle Hooks

```aether
resource database.instance "main" {
  engine = "postgres"
  
  on_error {
    retry = 3
    backoff = "exponential"
  }
  
  on_failure {
    notify = ["ops@example.com"]
    rollback = true
  }
  
  on_success {
    run_script = "./scripts/init-db.sh"
  }
}
```

### Resource Reference

```aether
// Access resource properties
vpc.main.id
instance.web_server.public_ip
database.main.connection_string

// With count
instance.web[0].id
instance.web[*].id        // All IDs

// With for_each
instance.servers["web"].id
instance.servers[*].id    // All IDs
```

---

## Variables and Outputs

### Variable Declaration

```aether
variable "name" {
  type = string
  default = "default value"
  description = "Variable description"
}

variable "count" {
  type = int
  description = "Number of instances"
}

variable "enable_feature" {
  type = bool
  default = false
}

variable "regions" {
  type = list<string>
  default = ["us-east", "us-west"]
}

variable "tags" {
  type = map<string, string>
  default = {}
}
```

### Variable Reference

```aether
var.name
var.count
var.enable_feature
```

### Output Declaration

```aether
output "name" {
  value = expression
  description = "Output description"
  sensitive = false
}

// Example
output "web_server_ip" {
  value = instance.web_server.public_ip
  description = "Public IP of web server"
}

output "database_password" {
  value = secret.db_password
  sensitive = true
  description = "Database password"
}

output "all_instance_ips" {
  value = [for i in instance.web : i.public_ip]
  description = "All instance IPs"
}
```

### Output Reference

```aether
output.web_server_ip
output.database_password
```

---

## Modules

### Module Declaration

```aether
module "name" {
  source = "path/or/url"
  version = "1.0.0"
  
  // Module inputs
  input1 = value1
  input2 = value2
}

// Examples
module "web_app" {
  source = "./modules/web-app"
  
  instance_count = 3
  region = "us-east"
}

module "database" {
  source = "github.com/company/db-module@v2.1.0"
  
  engine = "postgres"
  size = "large"
}

module "stdlib" {
  source = "stdlib/networking@1.0"
  
  vpc_cidr = "10.0.0.0/16"
}
```

### Module Reference

```aether
module.web_app.load_balancer_url
module.database.connection_string
module.stdlib.vpc_id
```

### Module Definition (aether.mod)

```aether
module "web-app" {
  version = "1.2.0"
  author = "ops-team"
  license = "MIT"
  description = "Scalable web application module"
  
  dependencies = [
    "stdlib/networking@1.0",
    "stdlib/compute@2.1"
  ]
  
  inputs = {
    instance_count = {
      type = int
      description = "Number of instances"
      default = 3
    }
    region = {
      type = string
      description = "Deployment region"
    }
  }
  
  outputs = {
    load_balancer_url = {
      type = string
      description = "URL of the load balancer"
    }
  }
}
```

---

## Agents

### Agent Declaration

```aether
agent "name" {
  type = "assistant" | "analyzer" | "autonomous"
  
  // Agent-specific configuration
}
```

### Assistant Agent

```aether
agent ide_assistant {
  type = "assistant"
  
  features = [
    "code_completion",
    "documentation",
    "refactoring"
  ]
  
  model = "local"  // or "cloud"
}
```

### Analyzer Agent

```aether
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
  
  notification {
    slack = "#security"
    email = ["security@example.com"]
  }
}

agent cost_analyzer {
  type = "analyzer"
  
  checks = [
    "oversized_resources",
    "unused_resources",
    "unattached_volumes"
  ]
  
  generate_report = true
}
```

### Autonomous Agent

```aether
agent cost_optimizer {
  type = "autonomous"
  scope = resource.compute.*
  
  goals = [
    "minimize_cost",
    "maintain_performance"
  ]
  
  constraints {
    max_cost_change = "$100/day"
    min_instances = 2
    max_instances = 20
  }
  
  approval {
    auto_approve = cost_change < "$10/day"
    required_for = cost_change >= "$10/day"
    
    notify = ["ops@example.com"]
  }
  
  learning_period = duration("7d")
  check_interval = duration("1h")
}

agent drift_watcher {
  type = "autonomous"
  
  check_interval = duration("15m")
  
  on_drift {
    action = "revert" | "alert" | "update_state"
  }
  
  approval {
    auto_approve = minor_drift == true
    required_for = major_drift == true
  }
}
```

### Agent Events

```aether
agent custom_handler {
  type = "autonomous"
  
  on ResourceCreated {
    // Handle resource creation
  }
  
  on CostThresholdExceeded {
    // Handle cost threshold
  }
  
  on SecurityViolation {
    // Handle security issue
  }
}
```

---

## Scripts

### Script Block

```aether
script {
  // Full programming capabilities
  
  let value = 42
  let items = ["a", "b", "c"]
  
  for item in items {
    // process item
  }
  
  if condition {
    // code
  }
  
  function helper(x) {
    return x * 2
  }
  
  // Export values for use in resources
  export { value, items }
}
```

### Script with Functions

```aether
script {
  function calculate_subnet_cidrs(vpc_cidr, count) {
    // Complex calculation logic
    let base_octets = split(vpc_cidr, ".")
    let result = []
    
    for i in range(count) {
      result.append("${base_octets[0]}.${base_octets[1]}.${i}.0/24")
    }
    
    return result
  }
  
  let subnet_cidrs = calculate_subnet_cidrs("10.0.0.0/16", 4)
  
  export { subnet_cidrs }
}

// Use in resources
resource network.subnet "app" {
  count = length(script.subnet_cidrs)
  cidr = script.subnet_cidrs[count.index]
}
```

---

## Comments and Annotations

### Single-line Comments

```aether
// This is a single-line comment

let name = "value"  // Inline comment
```

### Multi-line Comments

```aether
/*
 * This is a multi-line comment
 * spanning multiple lines
 */

/**
 * Documentation comment
 * for resources or functions
 */
```

### AI Annotations

```aether
// @ai:optimize - AI should suggest optimizations
// @ai:secure - AI should apply security best practices
// @ai:explain - AI should add documentation
// @ai:test - AI should generate tests

// @ai:optimize
resource compute.instance "web" {
  machine_type = "small"
}

// @ai:secure
resource storage.bucket "data" {
  // AI will suggest encryption, access controls, etc.
}

// @ai:explain
function complex_calculation(x, y) {
  // Complex logic here
  // AI will add explanatory comments
}
```

---

## Built-in Functions

### String Functions

```aether
upper(s: string) -> string                    // Convert to uppercase
lower(s: string) -> string                    // Convert to lowercase
trim(s: string) -> string                     // Remove whitespace
split(s: string, sep: string) -> list<string> // Split string
join(list: list<string>, sep: string) -> string // Join strings
replace(s: string, old: string, new: string) -> string // Replace
contains(s: string, substr: string) -> bool   // Check substring
starts_with(s: string, prefix: string) -> bool
ends_with(s: string, suffix: string) -> bool
```

### Collection Functions

```aether
length(collection) -> int                     // Get length
contains(list, item) -> bool                  // Check membership
index(list, item) -> int                      // Find index
reverse(list) -> list                         // Reverse list
sort(list) -> list                           // Sort list
unique(list) -> list                         // Remove duplicates

filter(list, predicate) -> list              // Filter elements
map(list, transform) -> list                 // Transform elements
reduce(list, initial, reducer) -> any        // Reduce to single value

keys(map) -> list                            // Get map keys
values(map) -> list                          // Get map values
merge(map1, map2) -> map                     // Merge maps
```

### Numeric Functions

```aether
abs(x: number) -> number                     // Absolute value
min(a, b) -> number                          // Minimum
max(a, b) -> number                          // Maximum
floor(x: number) -> int                      // Floor
ceil(x: number) -> int                       // Ceiling
round(x: number) -> int                      // Round
pow(base: number, exp: number) -> number     // Power
sqrt(x: number) -> number                    // Square root
```

### Type Conversion

```aether
string(value) -> string                      // Convert to string
int(value) -> int                           // Convert to int
number(value) -> number                      // Convert to number
bool(value) -> bool                         // Convert to bool
```

### Cloud Functions

```aether
region_available(provider: string, region: string) -> bool
resource_exists(id: string) -> bool
get_resource_property(id: string, property: string) -> any
list_availability_zones(region: string) -> list<string>
```

### Utility Functions

```aether
range(n: int) -> list<int>                   // Generate range
range(start: int, end: int) -> list<int>
enumerate(list) -> list<{index, value}>

duration(s: string) -> duration              // Parse duration
  // Examples: "5m", "2h", "1d"

hash(data: string, algorithm: string) -> string
encrypt(data: string, key: string) -> bytes
decrypt(data: bytes, key: string) -> string
```

---

## Testing

### Unit Test

```aether
test "description" {
  // Test code
  assert condition
  assert expression == expected
}

// Example
test "web_server_has_correct_type" {
  let server = resource.compute.instance.web_server
  assert server.machine_type == "small"
}
```

### Property-Based Test

```aether
property_test "description" {
  for resource in resources where condition {
    assert property
  }
}

// Example
property_test "all_storage_encrypted" {
  for resource in resources where resource.type == "storage.*" {
    assert resource.encryption_enabled == true
  }
}
```

### Integration Test

```aether
integration_test "description" {
  environment = "test"
  timeout = duration("10m")
  
  // Deploy
  apply()
  
  // Test
  assert condition
  
  // Cleanup
  defer destroy()
}
```

---

## Provider Configuration

```aether
provider "name" {
  property1 = value1
  property2 = value2
}

// Examples
provider "aws" {
  region = "us-east-1"
  profile = "production"
}

provider "azure" {
  subscription_id = "..."
  tenant_id = "..."
  region = "eastus"
}

provider "gcp" {
  project = "my-project"
  region = "us-east1"
  credentials = file("credentials.json")
}
```

---

## Backend Configuration

```aether
backend "type" {
  property1 = value1
  property2 = value2
}

// Examples
backend "s3" {
  bucket = "my-state"
  key = "production/infrastructure"
  region = "us-east-1"
  encrypt = true
}

backend "local" {
  path = "./terraform.tfstate"
}
```

---

## Secrets

```aether
secret "name" {
  source = "vault" | "aws_secretsmanager" | "env"
  path = "secret/path"
}

// Examples
secret "db_password" {
  source = "vault"
  path = "secret/database/password"
}

secret "api_key" {
  source = "env"
  key = "API_KEY"
}

// Reference
resource.property = secret.db_password
```

---

## File Organization

### Single File

```aether
// main.ae
variable "environment" { ... }

resource compute.instance "web" { ... }

output "ip" { ... }
```

### Multiple Files

```
project/
  main.ae           # Main resources
  variables.ae      # Variable declarations
  outputs.ae        # Output declarations
  agents.ae         # Agent definitions
  providers.ae      # Provider configuration
  modules/
    web-app/
      main.ae
      variables.ae
      outputs.ae
```

---

This syntax specification serves as the formal definition for implementing the Aether parser and type checker. It balances declarative simplicity with scripting power while maintaining clear, readable infrastructure code.
