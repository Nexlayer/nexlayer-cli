# Nexlayer API Documentation

This document details the Nexlayer API as defined by the OpenAPI 3.0.0 specification. It includes endpoints, parameters, request bodies, responses, and data schemas.

---

## Table of Contents

- [API Information](#api-information)
- [Servers](#servers)
- [Endpoints](#endpoints)
  - [Start User Deployment](#1-start-user-deployment)
  - [Send Feedback](#2-send-feedback)
  - [Save Custom Domain](#3-save-custom-domain)
  - [Get Deployments](#4-get-deployments)
  - [Get Deployment Info](#5-get-deployment-info)
- [Components / Schemas](#components--schemas)
  - [getDeploymentsResponse](#getdeploymentsresponse)
  - [getDeploymentInfoResponse](#getdeploymentinforesponse)
  - [startTemplateDeploymentResponse](#starttemplatedeploymentresponse)
  - [checkSiteStatusResponse](#checksitestatusresponse)
  - [startUserDeploymentResponse](#startuserdeploymentresponse)
  - [saveCustomDomainRequestBody](#savecustomdomainrequestbody)
  - [saveCustomDomainResponse](#savecustomdomainresponse)
  - [feedback](#feedback)
  - [startUserDeploymentRequestBody](#startuserdeploymentrequestbody)
- [Authentication](#authentication)

---

## API Information

- **OpenAPI Version:** `3.0.0`
- **Title:** Nexlayer API
- **Description:** API for the Nexlayer Application
- **Version:** `1.0.0`

---

## Servers

- **Server URL:** [https://app.staging.nexlayer.io](https://app.staging.nexlayer.io)

---

## Endpoints

### 1. Start User Deployment

- **Path:** `/startUserDeployment/{applicationID?}`
- **Method:** `POST`
- **Summary:** Start User Deployment
- **Description:** Start a deployment given an application ID. Accepts a YAML file uploaded using `--data-binary`.

#### Path Parameters

| Parameter         | In   | Type   | Description              |
| ----------------- | ---- | ------ | ------------------------ |
| `applicationID?`  | path | string | The application ID. This is a dynamic path parameter as indicated by the trailing `?`. |

#### Request Body

- **Required:** Yes
- **Content Type:** `text/x-yaml`
- **Schema Reference:** [startUserDeploymentRequestBody](#startuserdeploymentrequestbody)

#### Responses

- **200 OK**
  - **Description:** OK
  - **Content Type:** `application/json`
  - **Schema Reference:** [startUserDeploymentResponse](#startuserdeploymentresponse)
- **500 Internal Server Error**
  - **Description:** Internal Server Error

---

### 2. Send Feedback

- **Path:** `/feedback`
- **Method:** `POST`
- **Summary:** Send Feedback
- **Description:** Send feedback to Nexlayer.

#### Request Body

- **Required:** Yes
- **Content Type:** `application/json`
- **Schema Reference:** [feedback](#feedback)

#### Responses

- **200 OK**
  - **Description:** OK
- **500 Internal Server Error**
  - **Description:** Internal Server Error

---

### 3. Save Custom Domain

- **Path:** `/saveCustomDomain/{applicationID}`
- **Method:** `POST`
- **Summary:** Save Custom Domain
- **Description:** Save a custom domain to user profile.

#### Path Parameters

| Parameter       | In   | Required | Type   | Description           |
| --------------- | ---- | -------- | ------ | --------------------- |
| `applicationID` | path | Yes      | string | The application ID.   |

#### Request Body

- **Required:** Yes
- **Content Type:** `application/json`
- **Schema Reference:** [saveCustomDomainRequestBody](#savecustomdomainrequestbody)

#### Responses

- **200 OK**
  - **Description:** OK
  - **Content Type:** `application/json`
  - **Schema Reference:** [saveCustomDomainResponse](#savecustomdomainresponse)
- **400 Bad Request**
  - **Description:** Bad Request
- **500 Internal Server Error**
  - **Description:** Internal Server Error

---

### 4. Get Deployments

- **Path:** `/getDeployments/{applicationID}`
- **Method:** `GET`
- **Summary:** Get Deployments
- **Description:** Get all user deployments.

#### Path Parameters

| Parameter       | In   | Required | Type   | Description           |
| --------------- | ---- | -------- | ------ | --------------------- |
| `applicationID` | path | Yes      | string | The application ID.   |

#### Responses

- **200 OK**
  - **Description:** OK
  - **Content Type:** `application/json`
  - **Schema Reference:** [getDeploymentsResponse](#getdeploymentsresponse)
- **400 Bad Request**
  - **Description:** Bad Request
- **500 Internal Server Error**
  - **Description:** Internal Server Error

---

### 5. Get Deployment Info

- **Path:** `/getDeploymentInfo/{namespace}/{applicationID}`
- **Method:** `GET`
- **Summary:** Get Deployment Info
- **Description:** Get information around a specific deployment.

#### Path Parameters

| Parameter       | In   | Required | Type   | Description                       |
| --------------- | ---- | -------- | ------ | --------------------------------- |
| `namespace`     | path | Yes      | string | The deployment namespace.         |
| `applicationID` | path | Yes      | string | The application ID.               |

#### Responses

- **200 OK**
  - **Description:** OK
  - **Content Type:** `application/json`
  - **Schema Reference:** [getDeploymentInfoResponse](#getdeploymentinforesponse)
- **400 Bad Request**
  - **Description:** Bad Request
- **500 Internal Server Error**
  - **Description:** Internal Server Error

---

## Components / Schemas

### getDeploymentsResponse

- **Type:** `object`
- **Description:** Contains an array of deployments.
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

---

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

---

### startTemplateDeploymentResponse

- **Type:** `object`
- **Description:** Provides a response when a template deployment is initiated.
- **Properties:**
  - **message** (string, required)  
    *Example:* `"Deployment started successfully"`
  - **namespace** (string, required)  
    *Example:* `"ecstatic-frog"`
  - **url** (string, required)  
    *Example:* `"https://ecstatic-frog-kd-chat.alpha.nexlayer.ai"`

---

### checkSiteStatusResponse

- **Type:** `object`
- **Description:** Provides the status of the site.
- **Properties:**
  - **message** (string, required)  
    *Example:* `"UP"`

---

### startUserDeploymentResponse

- **Type:** `object`
- **Description:** Provides a response when a user deployment is initiated.
- **Properties:**
  - **message** (string, required)  
    *Example:* `"Deployment started successfully"`
  - **namespace** (string, required)  
    *Example:* `"fantastic-fox"`
  - **url** (string, required)  
    *Example:* `"https://fantastic-fox-my-mern-app.alpha.nexlayer.io"`

---

### saveCustomDomainRequestBody

- **Type:** `object`
- **Description:** The request body for saving a custom domain.
- **Properties:**
  - **domain** (string, required)  
    *Example:* `"mydomain.com"`

---

### saveCustomDomainResponse

- **Type:** `object`
- **Description:** Provides a response confirming the custom domain has been saved.
- **Properties:**
  - **message** (string, required)  
    *Example:* `"Custom domain saved successfully"`

---

### feedback

- **Type:** `object`
- **Description:** Contains the feedback text sent by a user.
- **Properties:**
  - **text** (string, required)  
    *Example:* `"Sample text"`

---

### startUserDeploymentRequestBody

- **Type:** `string`
- **Format:** `binary`
- **Description:** A binary string representing the YAML file used to initiate a user deployment.

---

## Authentication

All API endpoints require authentication using a Bearer token. Include the token in the `Authorization` header with each request:

```http
Authorization: Bearer <your-token>
```

Example using cURL:
```bash
curl -X GET "https://app.staging.nexlayer.io/getDeployments/{applicationID}" \
     -H "Authorization: Bearer your-token-here"
```

> **Important:** Requests without a valid Bearer token will be rejected with a 401 Unauthorized response.

---

## Usage Examples

### Starting a User Deployment

To start a deployment, send a `POST` request to:

```http
https://app.staging.nexlayer.io/startUserDeployment/{applicationID?}
```

Include your YAML file as binary data in the request body with the content type `text/x-yaml`.

A successful response (`200 OK`) will return JSON similar to:

```json
{
  "message": "Deployment started successfully",
  "namespace": "fantastic-fox",
  "url": "https://fantastic-fox-my-mern-app.alpha.nexlayer.io"
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

