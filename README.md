# s3_api
Go API to check status of s3 buckets and their objects

Authentication required to reach all endpoints.
Hard coded credentials should be changed for any serious use.
  - admin/adminpass
  - user/userpass
An AWS account with some S3 buckets with some objects will be needed.

A credentials text file will need to be added to a .aws folder on the root of the project.
The credentials file should have the following contents:
  [default]
  aws_access_key_id = your_read_only_access_key_id
  aws_secret_access_key = your_read_only_access_key

Auth endpoint - /api/v1/auth

Other endpoints:
  - /api/v1/list-all-buckets
  - /api/v1/list-all-objects
  - /api/v1/get-all-buckets-encryption
  - /api/v1/get-specific-bucket-encryption
  - /api/v1/get-bucket-objects-encryption
  
Encryption endpoints require Admin token
