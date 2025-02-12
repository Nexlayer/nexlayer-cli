# Nexlayer API Reference

This document describes the Nexlayer API endpoints used by the CLI.

## Endpoints

### Deployment

#### POST /startUserDeployment
Start a new deployment from a YAML template.

**Request:**
```json
{
  "template": {
    "application": {
      "name": "string",
      "url": "string",
      "pods": [...]
    }
  }
}
```

**Response:**
```json
{
  "deploymentId": "string",
  "status": "string"
}
```

#### GET /getDeployments
List all deployments.

**Response:**
```json
{
  "deployments": [
    {
      "id": "string",
      "name": "string",
      "status": "string",
      "url": "string"
    }
  ]
}
```

### Domain Management

#### POST /saveCustomDomain
Configure a custom domain for an application.

**Request:**
```json
{
  "applicationName": "string",
  "domain": "string"
}
```

**Response:**
```json
{
  "status": "string",
  "domain": "string"
}
```

## Error Responses

All endpoints may return these standard error responses:

```json
{
  "error": {
    "code": "string",
    "message": "string",
    "details": {}
  }
}
```

## Authentication

All requests require a Bearer token in the Authorization header:
```
Authorization: Bearer <token>
```
