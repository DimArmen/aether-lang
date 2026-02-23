# AWS Examples for Aether

This directory contains comprehensive examples of deploying infrastructure on AWS using Aether.

## Examples Overview

### 1. [AWS VPC Setup](./aws-vpc-setup/)
**Complexity:** Beginner  
**Resources:** VPC, Subnets, NAT Gateways, Security Groups, Flow Logs

A production-ready VPC configuration with:
- Multi-AZ deployment across 3 availability zones
- Public and private subnets for different tiers (web, app, database)
- NAT gateways for high availability
- Security groups with proper isolation
- VPC Flow Logs for monitoring

**Use Case:** Foundation for any AWS infrastructure deployment

```bash
cd aws-vpc-setup
aether plan
aether apply
```

### 2. [AWS EC2 Auto Scaling](./aws-ec2-autoscaling/)
**Complexity:** Intermediate  
**Resources:** EC2, Auto Scaling Group, Application Load Balancer, CloudWatch

A scalable web application deployment with:
- Launch template with best practices (IMDSv2, encryption)
- Auto Scaling Group with target tracking policies
- Application Load Balancer with HTTPS
- CloudWatch monitoring and alarms
- Automated deployment with user data scripts

**Use Case:** Traditional auto-scaling web applications

```bash
cd aws-ec2-autoscaling
aether plan
aether apply
```

### 3. [AWS Serverless API](./aws-serverless-api/)
**Complexity:** Intermediate  
**Resources:** Lambda, API Gateway, DynamoDB, CloudWatch

A complete serverless REST API with:
- Lambda functions for CRUD operations
- API Gateway HTTP API with CORS
- DynamoDB with GSIs and streams
- Lambda layers for shared dependencies
- X-Ray tracing enabled
- CloudWatch logging and alarms

**Use Case:** Serverless APIs, microservices, event-driven architectures

```bash
cd aws-serverless-api
aether plan
aether apply
```

### 4. [AWS Complete Web Application](./aws-complete-webapp/)
**Complexity:** Advanced  
**Resources:** VPC, ECS Fargate, Aurora, ElastiCache, ALB, CloudFront, Route53, S3

A production-grade multi-tier application with:
- Complete networking setup with VPC
- ECS Fargate for containerized applications
- Aurora PostgreSQL with read replicas
- ElastiCache Redis for caching
- Application Load Balancer with SSL
- CloudFront CDN for static assets
- Route53 for DNS management
- Secrets Manager for credentials
- Comprehensive monitoring and alarms

**Use Case:** Production web applications requiring high availability and scalability

```bash
cd aws-complete-webapp
aether plan
aether apply
```

## Prerequisites

### AWS Account Setup

1. **AWS Account**: You need an active AWS account
2. **AWS CLI**: Install and configure AWS CLI
   ```bash
   aws configure
   ```
3. **Credentials**: Set up AWS credentials using one of these methods:
   - AWS CLI profile (recommended)
   - Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
   - IAM role (for EC2 instances)

### Required Permissions

Your AWS user/role needs permissions for the services used in each example:
- **VPC Example**: VPC, Subnet, Internet Gateway, NAT Gateway, Security Group, Route Table
- **EC2 Example**: EC2, Auto Scaling, Elastic Load Balancing, CloudWatch
- **Serverless Example**: Lambda, API Gateway, DynamoDB, IAM, CloudWatch
- **Complete Example**: All of the above plus ECS, RDS, ElastiCache, S3, CloudFront, Route53

## Quick Start

1. **Initialize your project:**
   ```bash
   aether init my-aws-project
   cd my-aws-project
   ```

2. **Copy an example:**
   ```bash
   cp -r examples/aws-vpc-setup/* .
   ```

3. **Customize variables:**
   Edit the variables in `main.ae` to match your requirements.

4. **Preview changes:**
   ```bash
   aether plan
   ```

5. **Deploy:**
   ```bash
   aether apply
   ```

## Example Structure

Each example follows this structure:

```
example-name/
├── main.ae           # Main infrastructure definition
├── README.md         # Example-specific documentation
└── scripts/          # Helper scripts (if needed)
    └── init.sh
```

## Cost Considerations

### Estimated Monthly Costs (us-east-1)

| Example | Estimated Cost | Notes |
|---------|----------------|-------|
| VPC Setup | $90-150 | NAT Gateways are the main cost |
| EC2 Auto Scaling | $150-500 | Depends on instance count and size |
| Serverless API | $5-100 | Pay per request, scales to zero |
| Complete Web App | $500-2000+ | Full production stack |

**Note:** These are rough estimates. Actual costs depend on:
- Traffic volume
- Data transfer
- Resource utilization
- Region selection

### Cost Optimization Tips

1. **Use Spot Instances** for non-critical workloads
2. **Enable Auto Scaling** to match capacity with demand
3. **Use Reserved Instances** or Savings Plans for predictable workloads
4. **Implement proper monitoring** to identify unused resources
5. **Use S3 Intelligent-Tiering** for storage
6. **Consider Aurora Serverless** for variable database workloads

## Best Practices

### Security

1. **Never commit credentials** to version control
2. **Use Secrets Manager** for sensitive data
3. **Enable encryption at rest** for all data stores
4. **Use VPC endpoints** to avoid internet gateway costs
5. **Implement least privilege** IAM policies
6. **Enable VPC Flow Logs** for network monitoring
7. **Use Security Groups** instead of Network ACLs when possible

### High Availability

1. **Deploy across multiple AZs** (minimum 2, preferably 3)
2. **Use Auto Scaling** for EC2 instances
3. **Enable Multi-AZ** for RDS databases
4. **Use Application Load Balancer** for health checks
5. **Implement health checks** for all services
6. **Use Route53 health checks** for DNS failover

### Monitoring

1. **Enable CloudWatch** detailed monitoring
2. **Set up alarms** for critical metrics
3. **Use CloudWatch Dashboards** for visualization
4. **Enable X-Ray tracing** for distributed applications
5. **Centralize logs** in CloudWatch Logs
6. **Set up SNS notifications** for alerts

### Tagging Strategy

Always use consistent tags across resources:

```aether
tags = {
  Application = "myapp"
  Environment = "production"
  ManagedBy = "Aether"
  CostCenter = "engineering"
  Owner = "team-name"
}
```

## Troubleshooting

### Common Issues

**Issue: NAT Gateway connection timeout**
- Check route tables are properly configured
- Verify NAT Gateway is in a public subnet
- Ensure Internet Gateway is attached to VPC

**Issue: ECS tasks failing health checks**
- Verify security groups allow ALB to reach tasks
- Check application logs in CloudWatch
- Ensure health check endpoint returns 200

**Issue: Lambda function timeout**
- Increase timeout setting
- Check VPC configuration if accessing VPC resources
- Review CloudWatch Logs for errors

**Issue: RDS connection refused**
- Verify security groups allow traffic from application
- Check database is in available state
- Ensure subnet group has subnets in multiple AZs

### Getting Help

- Review AWS documentation: https://docs.aws.amazon.com/
- Check Aether docs: [../docs/aws-guide.md](../docs/aws-guide.md)
- AWS Support (if you have a support plan)

## Clean Up

To avoid ongoing charges, destroy resources when done:

```bash
aether destroy
```

**Warning:** This will delete all resources. Make sure you have backups of any important data.

For resources with `deletion_protection = true`, you'll need to:
1. Update the configuration to set `deletion_protection = false`
2. Apply the change
3. Then run destroy

## Additional Resources

- [AWS Well-Architected Framework](https://aws.amazon.com/architecture/well-architected/)
- [AWS Architecture Center](https://aws.amazon.com/architecture/)
- [Aether AWS Provider Guide](../docs/aws-guide.md)
- [Aether Syntax Reference](../docs/syntax.md)

## Contributing

Found an issue or want to add an example? Please contribute!

1. Fork the repository
2. Create your feature branch
3. Add your example with documentation
4. Submit a pull request

## License

These examples are provided under the MIT License. See LICENSE file for details.
