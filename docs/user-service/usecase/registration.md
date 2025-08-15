# User Registration Use Case Sequence Diagram

This document describes the user registration flow using a sequence diagram that can be rendered at [sequencediagram.org](https://sequencediagram.org).

## Sequence Diagram

```sequence
title User Registration Use Case Flow

Client->Controller: POST /users/ (CreateAccount)
note over Controller: Validate request body\n- FullName (required)\n- Email (required, valid format)\n- Password (required, min 8 chars)

Controller->Controller: Validate password strength
Controller->Controller: Validate email format
Controller->Controller: Hash password with bcrypt

Controller->UseCase: CreateAccount(ctx, account)
UseCase->Repository: Create(ctx, account)
Repository->Database: INSERT INTO accounts
Database-->Repository: Account created with ID
Repository-->UseCase: Created account

UseCase->Repository: FindByAccountID(ctx, accountID)
Repository->Database: SELECT * FROM verifications WHERE account_id = ?
Database-->Repository: No existing verification
Repository-->UseCase: No verification found

UseCase->UseCase: NewEmailVerification(accountID)
UseCase->Repository: CreateVerification(ctx, verification)
Repository->Database: INSERT INTO verifications
Database-->Repository: Verification created
Repository-->UseCase: Verification saved

UseCase->EmailService: SendVerificationEmail(ctx, email, name, token)
EmailService->EmailService: IsEnabled()
EmailService->EmailService: Send email with verification link
note over EmailService: Email contains:\n- Verification token\n- Account details\n- Verification instructions

EmailService-->UseCase: Email sent (or error logged)
UseCase-->Controller: Created account with verification pending
Controller-->Client: 201 Created + Success message

note over Client,Database: Email Verification Flow (Separate Request)

Client->Controller: GET /users/verify?token=<token>
Controller->UseCase: VerifyEmail(ctx, token)
UseCase->Repository: FindByToken(ctx, token)
Repository->Database: SELECT * FROM verifications WHERE token = ?
Database-->Repository: Verification record
Repository-->UseCase: Verification found

UseCase->UseCase: IsTokenValid(token)
UseCase->UseCase: MarkAsCompleted()
UseCase->Repository: UpdateVerification(ctx, verification)
Repository->Database: UPDATE verifications SET completed = true
Database-->Repository: Verification updated
Repository-->UseCase: Verification updated

UseCase->Repository: FindByID(ctx, accountID)
Repository->Database: SELECT * FROM accounts WHERE id = ?
Database-->Repository: Account record
Repository-->UseCase: Account found

UseCase->UseCase: Set EmailVerified = true
UseCase->Repository: Update(ctx, account)
Repository->Database: UPDATE accounts SET email_verified = true
Database-->Repository: Account updated
Repository-->UseCase: Account updated

UseCase-->Controller: Email verification successful
Controller-->Client: 200 OK + Success message

note over Client,Database: Resend Verification Flow (Optional)

Client->Controller: POST /users/resend-verification
Controller->UseCase: ResendVerificationEmail(ctx, email)
UseCase->Repository: FindByEmail(ctx, email)
Repository->Database: SELECT * FROM accounts WHERE email = ?
Database-->Repository: Account record
Repository-->UseCase: Account found

UseCase->UseCase: Check if already verified
UseCase->Repository: FindByAccountID(ctx, accountID)
Repository->Database: SELECT * FROM verifications WHERE account_id = ?
Database-->Repository: Existing verification
Repository-->UseCase: Verification found

UseCase->UseCase: GenerateToken()
UseCase->Repository: UpdateVerification(ctx, verification)
Repository->Database: UPDATE verifications SET token = ?
Database-->Repository: Verification updated
Repository-->UseCase: Verification updated

UseCase->EmailService: SendVerificationEmail(ctx, email, name, newToken)
EmailService-->UseCase: Email sent
UseCase-->Controller: Verification email resent
Controller-->Client: 200 OK + Success message
```

## Key Components

### Actors
- **Client**: External system or user making HTTP requests
- **Controller**: HTTP request handler (Gin controller)
- **UseCase**: Business logic layer
- **Repository**: Data access layer
- **Database**: PostgreSQL database
- **EmailService**: Email infrastructure service

### Main Flows

1. **Account Creation**: Client submits registration data, system validates, creates account, and sends verification email
2. **Email Verification**: Client clicks verification link, system validates token and marks email as verified
3. **Resend Verification**: Client requests new verification email, system generates new token and resends

### Error Handling

- Input validation errors return 400 Bad Request
- Database errors return 500 Internal Server Error
- Email sending failures are logged but don't fail account creation
- Invalid verification tokens return 400 Bad Request

### Security Features

- Password hashing using bcrypt
- Email format validation
- Password strength requirements (minimum 8 characters)
- Verification token generation and validation
- Email verification requirement before account activation

## Database Tables

- **accounts**: Stores user account information
- **verifications**: Stores email verification tokens and status

## Email Templates

Email templates are located in `src/email/` directory and include:
- `verification_email.html` - HTML version of verification email
- `verification_email.txt` - Plain text version of verification email
