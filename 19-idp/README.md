# Identity Provider (IdP)

## What is Identity Provider (IdP)?

- Simply put, an Identity Provider (IdP) **manages and maintains identity data** for users
- It's often used in conjunction with **Sigle Sign On (SSO)**
  - It gives a user a **single login & password (and optional MFA capability)**
- **OpenID Connect (OIDC)** and **Security Assertion Markup Language (SAML)** are authentication mechanisms, they don't store login/password information themselves
- You'd still need to validate the login, password, and potentially MFA token with a separate mechanism
  - Users can be in a database, in LDAP, Microsoft Active Directory, or others

## Why Implement OpenID Connect (OIDC)?

- It's great **learning experience**
  - Exposure to a lot of technologies: REST API, OAuth, JWT, JWK
- You're **often exposed to an IdP**, and it's worth understanding the inner workings
- You can **build you own** IdP authorization server, client, or application
  - Understanding how the how flow works **will help you** when you need to build one of these components

# OpenID Connect (OIDC)

## What is OpenID Connect (OIDC)?

- **OIDC** stands for **OpenID Connect**
- It's simple **identity layer on top of the OAuth 2.0 protocol**
- OIDC can **verify the identity** of a user using an **Authorization Server (AS)**
- OIDC uses **REST endpoints**, so it's **easily implemented**

## OIDC Flow

- Authorization Code Flow
  - For web applications that can store a _client_secret_
  - This is the flow we're going to implement
- Implicit Flow
  - For frontends/mobile apps that can't store a _client_secret_
- Hybrid Flow
  - Combines above two flows
  - Immediate access to an ID token

```mermaid
sequenceDiagram
  participant web as Web Application
  participant user as User
  participant auth as Authorization Server

  user ->> web: 1. Access Protected Resource or Login Request
  web -->> user: 2. Redirect to Authorization Server
  user ->> auth: 3. Authorization Code Request (/authorize endpoint)
  user ->> auth: 4. Login Prompt
  auth -->> user: 5. Redirect to the Application with "?code=" Query Parameter
  user ->> web: 6. URL with Code Parameter
  web ->> auth: 7. Exchange Code for Token (/token endpoint)
  auth -->> web: 8. Access & ID Tokens
```

```mermaid
sequenceDiagram
  participant user as User
  participant auth as Authorization Server

  user ->> auth: 1. /authorize?client_id=123&response_type=code&redirect_uri=https://myapp.com/callback&scope=openid&state=randomString (/authorize endpoint)
  auth -->> user: 2. Redirect to login
  user ->> auth: 3. /login (submit credentials)
```

```mermaid
sequenceDiagram
  participant web as Web Application
  participant user as User
  participant auth as Authorization Server

  user ->> web: 1. Access Protected Resource or Login Request
  web ->> auth: 2. POST /token grant_type=authorization_code&client_id=123&client_secret=456&redirect_uri=https://myapp.com/callback&code=789 (/token endpoint)
  auth ->> user: 3. Redirect: https://myapp.com/callback?code=789&state=randomString (/authorize endpoint)
```

```mermaid
sequenceDiagram
  participant web as Web Application
  participant user as User
  participant auth as Authorization Server

  user ->> web: 1. Access Protected Resource or Login Request
  web ->> auth: 2. GET /jwks.json (/token endpoint)
```

---

# Challenge

- Write **OpenID Connect (OIDC) Implementation**
- Start project (see below)

## Project Setup

- Clone the project from [GitHub](https://github.com/wardviaene/golang-for-devops-course.git)

```bash
git clone #PROJECT_URL#
cd #PROJECT_NAME/oidc-start#
```
