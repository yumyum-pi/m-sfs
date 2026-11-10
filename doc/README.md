# Static Asset Management Service — Project Scope

## 1. Project Overview

The Static Asset Management Service is a **multi-tenant API application written in Golang** for centrally managing public-facing static assets.

The application will address the current manual process where static files are directly placed on servers by the DevOps team without centralized ownership, metadata management, tenant association, or controlled access.

The service will provide authorized APIs for managing static assets while **Nginx remains responsible for serving the actual files publicly**.

The frontend or administrative interface is outside the scope of this project and will be developed independently.

---

# 2. Problem Statement

Currently, static assets such as:

* Icons
* Banners
* Email images
* Application images
* PDF documents
* Other supported static files

are manually deployed to servers.

This creates several challenges:

* No centralized asset management.
* No clear ownership of assets.
* Difficulty identifying which asset belongs to which tenant.
* Manual dependency on DevOps for asset deployment.
* No centralized metadata repository.
* No controlled API-based upload or deletion process.
* Limited auditability of asset operations.

The proposed application will provide a centralized, multi-tenant solution for managing these assets.

---

# 3. Project Objectives

The application will:

1. Centralize management of static assets.
2. Associate every asset with an account/tenant.
3. Provide JWT-authorized REST APIs for asset management.
4. Validate uploaded files before storage.
5. Store asset metadata in a database.
6. Store physical files in a configurable local directory.
7. Allow Nginx to serve stored files directly.
8. Maintain audit records for administrative operations.
9. Enforce tenant-level isolation.
10. Eliminate manual file deployment as part of normal asset management.

---

# 4. High-Level Architecture

```text id="dfx8pa"
                         ┌─────────────────────┐
                         │  Admin Frontend     │
                         │ (Separate Project)  │
                         └──────────┬──────────┘
                                    │
                              JWT Protected
                                    │
                                    ▼
                    ┌───────────────────────────┐
                    │ Static Asset Management   │
                    │        Service            │
                    │         Golang            │
                    │                           │
                    │ • Upload                  │
                    │ • Search                  │
                    │ • Soft Delete             │
                    │ • Validation              │
                    │ • Authorization           │
                    │ • Audit Logging           │
                    └─────────────┬─────────────┘
                                  │
                     Write / Delete / Manage
                                  │
                                  ▼
                    ┌───────────────────────────┐
                    │      Local Storage        │
                    │                           │
                    │       target_dir          │
                    └─────────────┬─────────────┘
                                  │
                                  │ Read Only
                                  ▼
                          ┌──────────────┐
                          │    Nginx     │
                          │ Static Server│
                          └──────┬───────┘
                                 │
                              Public
                                 │
                                 ▼
                            Consumers
```

## Architectural Principle

The system will separate responsibilities into two planes:

### Management Plane — Golang Application

Responsible for:

* Authentication
* Authorization
* File upload
* File validation
* Metadata management
* Tenant isolation
* Search
* Soft deletion
* Audit logging

### Serving Plane — Nginx

Responsible for:

* Serving static files
* Handling public requests
* HTTP connection management
* Static file performance
* Filesystem reads

The Golang application will **not serve public static files**.

---

# 5. Multi-Tenancy

## Tenant Identifier

The system will use:

```text id="5p1uyq"
account_id
```

as the tenant identifier.

The `account_id` will be obtained from the JWT token.

### Tenant Isolation

Tenant isolation will be enforced at the API layer.

A tenant user must only be able to access assets belonging to their own account.

Example:

```text id="miycn4"
JWT
 └── account_id: acme

Request
 └── GET /api/v1/assets?account_id=other

System
 └── Automatically enforce account_id from JWT
```

The user-provided tenant/account filter must not allow access to another tenant's assets.

### Authorization Principle

```text id="g1qkn8"
Authenticated User
       │
       ▼
Extract account_id from JWT
       │
       ▼
Apply account filter
       │
       ▼
Execute operation
```

---

# 6. Asset Storage

## Storage Location

Assets will be stored on the same server as:

* Golang application
* Nginx

The storage directory will be configurable.

Example:

```json id="h1d6fd"
{
  "target_dir": "/var/static/assets"
}
```

## Directory Structure

A flat directory structure will be used.

Example:

```text id="8h4e6w"
/var/static/assets
│
├── acme_550e8400.png
├── acme_1a2b3c4d.svg
├── account123_98a7b6.webp
├── account456_5d4c3b.pdf
```

No tenant-specific subdirectories will be created.

---

# 7. File Naming Convention

Each uploaded file will receive a UUID.

The stored filename format will be:

```text id="j9h4u4"
{account_id}_{asset_uuid}.{extension}
```

Example:

```text id="9kdp3j"
acme_550e8400-e29b-41d4-a716-446655440000.png
```

The original uploaded filename will remain unchanged in the database metadata.

Example:

| Property          | Value                                           |
| ----------------- | ----------------------------------------------- |
| Original Filename | `company-logo-final.png`                        |
| Stored Filename   | `acme_550e8400-e29b-41d4-a716-446655440000.png` |
| Account ID        | `acme`                                          |

---

# 8. Supported File Types

The application will have default supported file types.

### Default Supported Types

* SVG
* PNG
* JPEG
* WebP
* GIF
* PDF

The supported file types will be configurable.

## Configuration Override

The configuration will **replace the default list**, not extend it.

Example default:

```text id="qmp6an"
SVG
PNG
JPEG
WebP
GIF
PDF
```

Configuration:

```json id="1x5rmx"
{
  "supported_files": [
    "png",
    "webp"
  ]
}
```

Final supported types:

```text id="q5ppq8"
PNG
WebP
```

SVG, JPEG, GIF, and PDF will no longer be accepted.

---

# 9. File Size Limits

The maximum allowed file size will be configurable.

Default:

```text id="x49fw6"
5 MB
```

Example:

```json id="h05csf"
{
  "max_file_size": "5MB"
}
```

The application must reject files exceeding the configured size limit.

---

# 10. File Validation

Every uploaded file must undergo validation before permanent storage.

Validation will include:

### 10.1 File Size

Validate against configured maximum file size.

### 10.2 File Extension

Validate against configured supported extensions.

Example:

```text id="hue0vl"
.png
.svg
.webp
```

### 10.3 MIME Type

Validate the detected MIME type.

### 10.4 File Signature / Magic Bytes

The application must inspect the actual file content to validate its type.

Validation must not rely solely on:

* File extension
* Client-provided MIME type

Example:

```text id="7adqbf"
malicious.exe
      ↓ renamed
image.png
```

The application must reject files whose actual content does not match the expected file type.

---

# 11. SVG Security

SVG files will be sanitized before permanent storage.

The system must remove or reject potentially unsafe SVG content before storing the asset.

The sanitization process must occur before the file is moved to the final asset directory.

---

# 12. File Upload

Files will be uploaded through a REST API.

The upload process will follow this workflow:

```text id="2ub60b"
Client Upload
      │
      ▼
Temporary Directory
      │
      ▼
Validate File Size
      │
      ▼
Validate Extension
      │
      ▼
Validate MIME Type
      │
      ▼
Validate File Signature
      │
      ▼
Sanitize SVG (if applicable)
      │
      ▼
Generate UUID
      │
      ▼
Generate Stored Filename
      │
      ▼
Create Database Record
      │
      ├──── Failure ────► Rollback / Cleanup
      │
      ▼
Atomically Move File
      │
      ├──── Failure ────► Rollback DB Record
      │
      ▼
Success
```

The application must ensure consistency between:

* Filesystem
* Database

No orphaned database records or unmanaged files should remain after failed operations.

---

# 13. Temporary Storage

A temporary directory will be configurable.

Example:

```json id="hh4xus"
{
  "temp_dir": "/var/static/tmp",
  "target_dir": "/var/static/assets"
}
```

Uploaded files must first be placed in the temporary directory.

Temporary files must not be directly exposed through Nginx.

The application should include a cleanup mechanism for abandoned or failed temporary files.

---

# 14. Asset Management APIs

All management APIs will be versioned.

Base API path:

```text id="v9asxd"
/api/v1
```

All management APIs require JWT authorization.

---

## 14.1 Upload Asset

```http id="aaxmnp"
POST /api/v1/assets
```

Responsibilities:

* Authenticate user.
* Authorize user.
* Extract `account_id`.
* Upload file.
* Validate file.
* Generate asset UUID.
* Generate stored filename.
* Store metadata.
* Move file to target directory.
* Create audit log.

---

## 14.2 Search Assets

```http id="q9sjk4"
GET /api/v1/assets
```

Search will support:

### Search

* Original filename
* Stored filename

### Filters

* File type
* Created by
* Created date
* Updated date

Tenant/account isolation will automatically be enforced using the JWT.

---

## 14.3 Get Asset Metadata

This needs to be explicitly included for the frontend/API consumer to retrieve details of a single asset.

Proposed endpoint:

```http id="1r9bqk"
GET /api/v1/assets/{id}
```

This endpoint will return metadata only.

It will **not serve the actual file**.

Please note: I consider this a necessary companion to Search. If you do not want it, it should be explicitly excluded.

---

## 14.4 Delete Asset

```http id="2foyg2"
DELETE /api/v1/assets/{id}
```

Deletion will use a soft-delete model.

The frontend will handle deletion confirmation.

After deletion:

* Asset will no longer appear in normal search results.
* Asset metadata will be retained.
* Audit history will remain available.

### Important technical point

You selected **soft delete**, but we have not yet explicitly defined whether the physical file should remain on disk.

My recommended interpretation is:

```text id="4v5ptg"
DELETE
   │
   ▼
Soft delete DB record
   │
   ▼
Remove physical file from target directory
   │
   ▼
Public URL returns 404
```

This retains metadata/auditability without leaving publicly accessible deleted assets.

I recommend we use this behavior in the final technical specification.

---

# 15. Asset Updates

Asset updates and replacement are explicitly **out of scope**.

The following operations will not be supported:

* Replacing an existing file.
* Changing the file associated with an asset.
* Changing the stored filename.
* Changing the public asset URL.

To change an asset, users must:

1. Upload a new asset.
2. Update the external application/reference to use the new asset URL.
3. Delete the old asset when no longer required.

This ensures asset URLs remain immutable.

---

# 16. Public Static File Access

Static files will be publicly accessible through Nginx.

Example:

```text id="e13lzs"
https://static.example.com/assets/acme_550e8400.png
```

Public static URLs will:

* Not require JWT.
* Not pass through the Golang application.
* Be served directly by Nginx.

The API server and static file server will have separate URLs/domains.

Example:

```text id="rdc74d"
Management API
https://asset-manager.example.com/api/v1/assets

Static Assets
https://static.example.com/assets/acme_550e8400.png
```

---

# 17. Authorization

All management APIs will require JWT authentication.

JWT claims will include the information required to identify:

* User
* Account/Tenant
* Role/authorization context

Authorization will be role-based.

The authorization layer must ensure:

```text id="khslmc"
User
 │
 ├── Authenticated?
 │
 ▼
Extract JWT Claims
 │
 ▼
Validate Role
 │
 ▼
Extract account_id
 │
 ▼
Enforce Tenant Isolation
 │
 ▼
Allow / Deny Operation
```

The exact role definitions should integrate with the organization's existing authorization model where applicable.

---

# 18. Database

## Production

```text id="auq3az"
MySQL
```

## Local Development

```text id="xj40v1"
SQLite
```

The application should abstract database operations sufficiently to support both environments.

---

# 19. Asset Metadata

The following metadata will be stored.

| Field                | Description                     |
| -------------------- | ------------------------------- |
| `id`                 | UUID of the asset               |
| `account_id`         | Asset owner / tenant identifier |
| `original_file_name` | Original uploaded filename      |
| `stored_file_name`   | Generated filename              |
| `file_path`          | Physical storage path           |
| `file_size`          | File size in bytes              |
| `extension`          | File extension                  |
| `mime_type`          | Validated MIME type             |
| `created_at`         | Asset creation timestamp        |
| `created_by`         | User who uploaded the asset     |
| `updated_at`         | Last metadata update timestamp  |
| `deleted_at`         | Soft deletion timestamp         |

The checksum field has been explicitly removed from scope.

---

# 20. Duplicate Files

Duplicate files will be allowed.

The system will not reject uploads based on:

* Same original filename.
* Same binary content.
* Same file size.

Because each asset receives a UUID, stored filenames will remain unique.

Example:

```text id="t1r2r3"
acme_uuid-1.png
acme_uuid-2.png
```

Both may contain identical file content.

---

# 21. Search and Pagination

Search APIs will support pagination.

Default page size:

```text id="g93ws4"
20
```

The API will also support sorting.

Example:

```http id="dcv8w4"
GET /api/v1/assets?page=1&page_size=20
```

Example with sorting:

```http id="fwc8g1"
GET /api/v1/assets?page=1&page_size=20&sort_by=created_at&sort_order=desc
```

Sorting fields should include relevant metadata fields, subject to an explicit allow-list.

Recommended initial fields:

* `created_at`
* `updated_at`
* `original_file_name`
* `file_size`

The API should not directly pass arbitrary user-provided fields into SQL sorting.

---

# 22. API Response Standard

All APIs will follow a common response format.

## Successful Response

```json id="9a4sh1"
{
  "success": true,
  "data": {},
  "message": "Asset uploaded successfully"
}
```

## Error Response

```json id="b6n2kl"
{
  "success": false,
  "error": {
    "code": "INVALID_FILE_TYPE",
    "message": "The uploaded file type is not supported"
  }
}
```

The final API specification should define:

* Error codes
* HTTP status codes
* Validation errors
* Authorization errors
* Pagination response structure

---

# 23. Application Configuration

Configuration will be loaded from:

```text id="zb9zmy"
config.json
```

Configuration is loaded only during application startup.

Hot configuration reload is out of scope.

The application must be restarted after configuration changes.

## Configuration Areas

### File Configuration

* `target_dir`
* `temp_dir`
* `max_file_size`
* `supported_files`

### JWT Configuration

JWT-related settings required for token validation.

### Database Configuration

Connection configuration for:

* MySQL
* SQLite

### Logging Configuration

Supported logging level.

Example:

```text id="5d5c2o"
DEBUG
INFO
WARN
ERROR
```

### Conceptual Configuration Example

```json id="y9krsn"
{
  "files": {
    "target_dir": "/var/static/assets",
    "temp_dir": "/var/static/tmp",
    "max_file_size": "5MB",
    "supported_files": [
      "svg",
      "png",
      "jpeg",
      "webp",
      "gif",
      "pdf"
    ]
  },
  "jwt": {},
  "database": {},
  "logging": {
    "level": "INFO"
  }
}
```

The actual configuration schema can be finalized during technical design.

---

# 24. Application Logging

The application will support configurable operational logging.

Logging levels:

* DEBUG
* INFO
* WARN
* ERROR

Application logs may include:

* Application startup/shutdown.
* Upload success/failure.
* Validation failures.
* Database errors.
* File operation failures.
* Authorization failures.
* Unexpected system errors.

Sensitive information such as JWT tokens must not be logged.

---

# 25. Audit Logging

The application will maintain audit logs separately from operational logs.

Audit events should include:

* Asset uploaded.
* Asset deleted.
* Failed management operations where relevant.

Suggested audit metadata:

| Field          | Description                       |
| -------------- | --------------------------------- |
| `id`           | Audit event identifier            |
| `asset_id`     | Related asset                     |
| `account_id`   | Related tenant                    |
| `action`       | Upload/Delete                     |
| `performed_by` | User identifier                   |
| `timestamp`    | Action timestamp                  |
| `metadata`     | Additional contextual information |

Audit logs must remain available even when an asset is soft deleted.

---

# 26. Health Checks

The application will expose health endpoints.

## Liveness

```http id="jqjog1"
GET /health
```

Indicates whether the application process is running.

## Readiness

```http id="dvk3ir"
GET /ready
```

Indicates whether the application is ready to process requests.

Readiness may validate dependencies such as:

* Database connectivity.
* Required filesystem availability.

---

# 27. File Permissions

The Golang application and Nginx will share access to the asset storage directory.

The deployment must enforce:

```text id="6zj3pb"
Golang Application
   │
   ├── Read
   ├── Write
   └── Delete

Asset Directory
   │
   ▼

Nginx
   │
   └── Read Only
```

Nginx must not have permission to:

* Modify files.
* Delete files.
* Upload files.

Temporary upload directories must not be publicly accessible.

---

# 28. Soft Delete Behaviour

The database will implement soft deletion.

A deleted asset should:

```text id="2msi0d"
DELETE Request
      │
      ▼
Authorization Check
      │
      ▼
Mark Asset as Deleted
      │
      ▼
Remove from Search Results
      │
      ▼
Remove Public File
      │
      ▼
Create Audit Log
```

This approach preserves:

* Asset metadata.
* Deletion history.
* Auditability.

While ensuring deleted public assets are no longer accessible.

---

# 29. In Scope

The following are included in the project scope.

### Application

* Golang REST API.
* Multi-tenant architecture.
* JWT authentication.
* Role-based authorization.
* Tenant isolation.

### Asset Management

* File upload.
* File validation.
* SVG sanitization.
* Asset metadata management.
* Asset search.
* Asset metadata retrieval.
* Soft delete.

### Storage

* Local filesystem storage.
* Configurable target directory.
* Configurable temporary directory.
* Atomic file movement.

### Database

* MySQL production support.
* SQLite development support.
* Asset metadata persistence.
* Audit log persistence.

### API

* API versioning.
* Standard response format.
* Pagination.
* Sorting.
* Filtering.

### Operations

* Health checks.
* Application logging.
* Audit logging.

### Security

* JWT authorization.
* Role-based access.
* Tenant isolation.
* Magic byte validation.
* SVG sanitization.
* Filesystem permission separation.

---

# 30. Out of Scope

The following are explicitly outside this project's scope.

### Frontend

The administration UI will be developed separately.

### Static File Serving

The Golang application will not serve static assets.

Nginx will be responsible for public file delivery.

### Nginx Management

DevOps will manage:

* Nginx installation.
* Nginx configuration.
* Nginx deployment.
* Nginx optimization.

### Asset Updates

The following are out of scope:

* Replacing existing assets.
* Changing an asset's file.
* Changing an existing asset URL.

### Application-Level File Cache

No in-memory file caching will be implemented.

### CDN

CDN integration is out of scope.

### Image Processing

The following are out of scope:

* Image resizing.
* Image cropping.
* Thumbnail generation.
* Format conversion.
* Image compression.

### Object Storage

The initial version will not support:

* Amazon S3
* Azure Blob Storage
* Google Cloud Storage

### Hot Configuration Reload

Configuration changes require application restart.

### Duplicate Detection

Duplicate files will be allowed.

---

# 31. Key Non-Functional Requirements

| Area                  | Requirement                    |
| --------------------- | ------------------------------ |
| Language              | Golang                         |
| Database Production   | MySQL                          |
| Database Development  | SQLite                         |
| Authentication        | JWT                            |
| Authorization         | Role-based                     |
| Tenant Isolation      | Account ID derived from JWT    |
| Storage               | Local filesystem               |
| File Serving          | Nginx                          |
| Default Max File Size | 5 MB                           |
| API Versioning        | Required                       |
| Pagination            | Required                       |
| Default Page Size     | 20                             |
| Audit Logging         | Required                       |
| Application Logging   | Required                       |
| Health Checks         | Required                       |
| Configuration         | `config.json`                  |
| Config Reload         | Startup only                   |
| File Validation       | Extension + MIME + magic bytes |
| SVG                   | Sanitized                      |
| Delete                | Soft delete                    |
| Asset Update          | Not supported                  |

---

# 32. Key Risks and Dependencies

## Dependencies

| Dependency                     | Owner                                   |
| ------------------------------ | --------------------------------------- |
| Nginx configuration            | DevOps                                  |
| Shared filesystem permissions  | DevOps                                  |
| JWT claim structure            | Existing authentication system          |
| Role definitions               | Existing authorization system / Product |
| MySQL availability             | Infrastructure                          |
| Account ID availability in JWT | Authentication system                   |

## Risks

### Public Asset Deletion

Soft deletion will remove public access to the file.

External systems using the URL may receive:

```text id="aw2l1h"
404 Not Found
```

The frontend confirmation is important but cannot detect external references.

### Shared Filesystem

The Go application and Nginx must maintain correct permissions and ownership.

### SVG Security

SVG sanitization must be carefully implemented because improperly sanitized SVGs can introduce security risks.

### Filesystem and Database Consistency

The upload and deletion processes must handle partial failures and rollback correctly.

