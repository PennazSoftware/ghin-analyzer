# AWS Infrastructure for Handicap applications

## AWS Infrastructure
Terraform is deployed via workspaces. There is a Dev and Prod workspace.
1. Navigate to ./build/terraform folder
2. Run `terraform workspace select dev` for development or `terraform workspace select prod` for production
3. Run `terraform plan` to see a list of pending changes
4. Run `terraform apply --auto-approve` to apply the changes to AWS
5. If changes were made to API Gateway it may be necessary to manually force a deploy.
6. Test the APIs to make sure they are working.

## Terraform Workspace Commands
- Create a new workspace: `terraform workspace new dev`
- List workspaces: `terraform workspace list`
- Select a workspace: `terraform workspace select dev`

## Terraform State information
- List all resources: `terraform state list`
- Show the details (including IDs) of a resource: `terraform state show aws_cognito_user_pool.crm_user_pool`
