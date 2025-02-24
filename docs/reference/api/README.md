# Nexlayer API Documentation

This document details the Nexlayer API, which enables programmatic interaction with the Nexlayer platform for managing deployments, feedback, and custom domains. The API is defined using OpenAPI 3.0.0 and is designed to support AI-driven tools, automation, and seamless integration into development workflows.

## Table of Contents

- [API Information](#api-information)
- [Servers](#servers)
- [Authentication](#authentication)
- [Endpoints](#endpoints)
  - [Start User Deployment](#1-start-user-deployment)
  - [Send Feedback](#2-send-feedback)
  - [Save Custom Domain](#3-save-custom-domain)
  - [Get Deployments](#4-get-deployments)
  - [Get Deployment Info](#5-get-deployment-info)
- [Components / Schemas](#components--schemas)
  - [getDeploymentsResponse](#getdeploymentsresponse)
  - [getDeploymentInfoResponse](#getdeploymentinforesponse)
  - [startUserDeploymentResponse](#startuserdeploymentresponse)
  - [saveCustomDomainRequestBody](#savecustomdomainrequestbody)
  - [saveCustomDomainResponse](#savecustomdomainresponse)
  - [feedback](#feedback)
- [Usage Examples](#usage-examples)

## API Information

- **OpenAPI Version:** `3.0.0`
- **Title:** Nexlayer API
- **Description:** API for managing deployments, feedback, and custom domains in the Nexlayer Application. It enables seamless integration with AI tools and automation for deployment workflows.
- **Version:** `1.0.0`

## Servers

- **Staging Server URL:** [https://app.staging.nexlayer.io](https://app.staging.nexlayer.io)

> Note: This is the staging environment. Production URLs may differ and should be configured accordingly.

## Authentication

All API endpoints require authentication using a Bearer token. Include the token in the `Authorization` header with each request:

```http
Authorization: Bearer <your-token>
```

Example using cURL:
```bash
curl -X GET "https://app.staging.nexlayer.io/getDeployments/<applicationID>" \
     -H "Authorization: Bearer your-token-here"
```

> **Important:** Requests without a valid Bearer token will be rejected with a 401 Unauthorized response.

## Endpoints

### 1. Start User Deployment

- **Path:** `/startUserDeployment/{applicationID}`
- **Method:** `POST`
- **Tags:** Deployment
- **Summary:** Start a user deployment by uploading a YAML configuration file
- **Description:** Initiates a deployment for a user's application using a YAML configuration file uploaded via `--data-binary`. The YAML must follow the Nexlayer schema (see Nexlayer YAML Schema Documentation for details). The applicationID parameter is optional; if omitted, the default application tied to the user's profile is used.
- **Operation ID:** `startUserDeployment`

#### Path Parameters

| Parameter | In | Type | Description |
|-----------|----|----|-------------|
| `applicationID` | path | string | The unique identifier of the application to deploy. Optional; defaults to the user's profile application if not specified. |

#### Request Body

- **Required:** Yes
- **Content Type:** `text/x-yaml`
- **Schema:** Binary string (YAML file)
- **Example:** `# See https://github.com/Nexlayer/templates/blob/main/new-readme.md for YAML schema`

#### Responses

##### 200 OK
- **Description:** Deployment started successfully
- **Content Type:** `application/json`
- **Schema Reference:** [startUserDeploymentResponse](#startuserdeploymentresponse)
- **Example:**
```json
{
  "message": "Deployment started successfully",
  "namespace": "fantastic-fox",
  "url": "https://fantastic-fox-my-mern-app.alpha.nexlayer.ai"
}
```

##### 400 Bad Request
- **Description:** Invalid YAML or missing required fields
- **Content Type:** `application/json`
- **Example:**
```json
{
  "error": "Invalid YAML format or missing required fields."
}
```

##### 500 Internal Server Error
- **Description:** An unexpected error occurred on the server
- **Content Type:** `application/json`
- **Example:**
```json
{
  "error": "An unexpected error occurred on the server."
}
```

### 2. Send Feedback

- **Path:** `/feedback`
- **Method:** `POST`
- **Tags:** Feedback
- **Summary:** Send feedback to Nexlayer
- **Description:** Submits user feedback about the Nexlayer application in JSON format. The request must include a 'text' field with the feedback message.
- **Operation ID:** `sendFeedback`

#### Request Body

- **Required:** Yes
- **Content Type:** `application/json`
- **Schema Reference:** [feedback](#feedback)
- **Example:**
```json
{
  "text": "Great tool, but needs more documentation!"
}
```

#### Responses

##### 200 OK
- **Description:** Feedback received successfully
- **Content Type:** `application/json`
- **Example:**
```json
{
  "message": "Feedback received successfully"
}
```

##### 400 Bad Request
- **Description:** Missing or invalid feedback
- **Content Type:** `application/json`
- **Example:**
```json
{
  "error": "Feedback text is required."
}
```

##### 500 Internal Server Error
- **Description:** An unexpected error occurred on the server
- **Content Type:** `application/json`
- **Example:**
```json
{
  "error": "An unexpected error occurred on the server."
}
```

### 3. Save Custom Domain

- **Path:** `/saveCustomDomain/{applicationID}`
- **Method:** `POST`
- **Tags:** Domain Management
- **Summary:** Save a custom domain for an application
- **Description:** Associates a custom domain with the specified application. The domain must be a valid string and properly configured in DNS settings.
- **Operation ID:** `saveCustomDomain`

#### Path Parameters

| Parameter | In | Required | Type | Description |
|-----------|----|----|----|----|
| `applicationID` | path | Yes | string | The unique identifier of the application. |

#### Request Body

- **Required:** Yes
- **Content Type:** `application/json`
- **Schema Reference:** [saveCustomDomainRequestBody](#savecustomdomainrequestbody)
- **Example:**
```json
{
  "domain": "mydomain.com"
}
```

#### Responses

##### 200 OK
- **Description:** Custom domain saved successfully
- **Content Type:** `application/json`
- **Schema Reference:** [saveCustomDomainResponse](#savecustomdomainresponse)
- **Example:**
```json
{
  "message": "Custom domain saved successfully"
}
```

##### 400 Bad Request
- **Description:** Invalid domain format
- **Content Type:** `application/json`
- **Example:**
```json
{
  "error": "Invalid domain format. Please provide a valid domain name."
}
```

##### 500 Internal Server Error
- **Description:** An unexpected error occurred on the server
- **Content Type:** `application/json`
- **Example:**
```json
{
  "error": "An unexpected error occurred on the server."
}
```

### 4. Get Deployments

- **Path:** `/getDeployments/{applicationID}`
- **Method:** `GET`
- **Tags:** Deployment
- **Summary:** Get all deployments for an application
- **Description:** Retrieves a list of all deployments for the specified application ID, including details like namespace, template ID, and status.
- **Operation ID:** `getDeployments`

#### Path Parameters

| Parameter | In | Required | Type | Description |
|-----------|----|----|----|----|
| `applicationID` | path | Yes | string | The unique identifier of the application. |

#### Responses

##### 200 OK
- **Description:** Deployments retrieved successfully
- **Content Type:** `application/json`
- **Schema Reference:** [getDeploymentsResponse](#getdeploymentsresponse)
- **Example:**
```json
{
  "deployments": [
    {
      "namespace": "ecstatic-frog",
      "templateID": "0001",
      "templateName": "K-d chat",
      "deploymentStatus": "running"
    }
  ]
}
```

##### 400 Bad Request
- **Description:** Invalid application ID
- **Content Type:** `application/json`
- **Example:**
```json
{
  "error": "Invalid application ID."
}
```

##### 500 Internal Server Error
- **Description:** An unexpected error occurred on the server
- **Content Type:** `application/json`
- **Example:**
```json
{
  "error": "An unexpected error occurred on the server."
}
```

### 5. Get Deployment Info

- **Path:** `/getDeploymentInfo/{namespace}/{applicationID}`
- **Method:** `GET`
- **Tags:** Deployment
- **Summary:** Get detailed info for a specific deployment
- **Description:** Retrieves detailed information about a specific deployment identified by its namespace and application ID, including status and template details.
- **Operation ID:** `getDeploymentInfo`

#### Path Parameters

| Parameter | In | Required | Type | Description |
|-----------|----|----|----|----|
| `namespace` | path | Yes | string | The namespace of the deployment. |
| `applicationID` | path | Yes | string | The unique identifier of the application. |

#### Responses

##### 200 OK
- **Description:** Deployment info retrieved successfully
- **Content Type:** `application/json`
- **Schema Reference:** [getDeploymentInfoResponse](#getdeploymentinforesponse)
- **Example:**
```json
{
  "deployment": {
    "namespace": "ecstatic-frog",
    "templateID": "0001",
    "templateName": "K-d chat",
    "deploymentStatus": "running"
  }
}
```

##### 400 Bad Request
- **Description:** Invalid namespace or application ID
- **Content Type:** `application/json`
- **Example:**
```json
{
  "error": "Invalid namespace or application ID."
}
```

##### 500 Internal Server Error
- **Description:** An unexpected error occurred on the server
- **Content Type:** `application/json`
- **Example:**
```json
{
  "error": "An unexpected error occurred on the server."
}
```

## Components / Schemas

### getDeploymentsResponse

- **Type:** `object`
- **Description:** Contains an array of deployments associated with the application.
- **Properties:**
  - **deployments** (array of objects, required)
    - **namespace** (string)  
      *Example:* `"ecstatic-frog"`
    - **templateID** (string)  
      *Example:* `"0001"`
    - **templateName** (string)  
      *Example:* `"K-d chat"`
    - **deploymentStatus** (string)  
      *Example:* `"running"`

### getDeploymentInfoResponse

- **Type:** `object`
- **Description:** Contains detailed information about a specific deployment.
- **Properties:**
  - **deployment** (object, required)
    - **namespace** (string)  
      *Example:* `"ecstatic-frog"`
    - **templateID** (string)  
      *Example:* `"0001"`
    - **templateName** (string)  
      *Example:* `"K-d chat"`
    - **deploymentStatus** (string)  
      *Example:* `"running"`

### startUserDeploymentResponse

- **Type:** `object`
- **Description:** Provides a response when a user deployment is initiated.
- **Properties:**
  - **message** (string, required)  
    *Example:* `"Deployment started successfully"`
  - **namespace** (string, required)  
    *Example:* `"fantastic-fox"`
  - **url** (string, required)  
    *Example:* `"https://fantastic-fox-my-mern-app.alpha.nexlayer.ai"`

### saveCustomDomainRequestBody

- **Type:** `object`
- **Description:** The request body for saving a custom domain.
- **Properties:**
  - **domain** (string, required)  
    *Example:* `"mydomain.com"`

### saveCustomDomainResponse

- **Type:** `object`
- **Description:** Provides a response confirming the custom domain has been saved.
- **Properties:**
  - **message** (string, required)  
    *Example:* `"Custom domain saved successfully"`

### feedback

- **Type:** `object`
- **Description:** Contains the feedback text sent by a user.
- **Properties:**
  - **text** (string, required)  
    *Example:* `"Sample text"`

## Usage Examples

### Starting a User Deployment

To start a deployment, send a `POST` request to:

```http
https://app.staging.nexlayer.io/startUserDeployment/{applicationID}
```

Include your YAML file as binary data in the request body with the content type `text/x-yaml`. A successful response (`200 OK`) will return JSON similar to:

```json
{
  "message": "Deployment started successfully",
  "namespace": "fantastic-fox",
  "url": "https://fantastic-fox-my-mern-app.alpha.nexlayer.ai"
}
```

### Sending Feedback

Send a POST request to:

```http
https://app.staging.nexlayer.io/feedback
```

Include your feedback as JSON in the request body with the content type `application/json`:

```json
{
  "text": "Your feedback here..."
}
```

### Saving a Custom Domain

Send a POST request to:

```http
https://app.staging.nexlayer.io/saveCustomDomain/{applicationID}
```

Include a JSON object in the request body:

```json
{
  "domain": "mydomain.com"
}
```

A successful response (200 OK) returns:

```json
{
  "message": "Custom domain saved successfully"
}
```

### Retrieving Deployments

Send a GET request to:

```http
https://app.staging.nexlayer.io/getDeployments/{applicationID}
```

A successful response returns a JSON object listing deployments:

```json
{
  "deployments": [
    {
      "namespace": "ecstatic-frog",
      "templateID": "0001",
      "templateName": "K-d chat",
      "deploymentStatus": "running"
    }
  ]
}
```

### Getting Deployment Information

Send a GET request to:

```http
https://app.staging.nexlayer.io/getDeploymentInfo/{namespace}/{applicationID}
```

A successful response returns detailed deployment information:

```json
{
  "deployment": {
    "namespace": "ecstatic-frog",
    "templateID": "0001",
    "templateName": "K-d chat",
    "deploymentStatus": "running"
  }
}
```

