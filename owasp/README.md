# Real Security Issues

This document explains the security issues found by the security scanner in our application.

[Initial Report Link](./2024-03-09-ZAP-Report-Before.html)

[After Fixing Report Link](./2024-03-09-ZAP-Report-After.html)

## X-Content-Type-Options Header Missing

The X-Content-Type-Options Header Missing issue is raised when the application doesn't include the `X-Content-Type-Options` header in the response. This header prevents the browser from MIME-sniffing a response away from the declared content-type. This reduces exposure to drive-by download attacks and user-uploaded content that, by clever naming, could be treated as executable or dynamic HTML files.

In our application, we didn't include the `X-Content-Type-Options` header in the response. This could expose the application to MIME-sniffing attacks, which could lead to drive-by download attacks and user-uploaded content being treated as executable or dynamic HTML files.

To fix this issue, we added the `X-Content-Type-Options` header to the response with the value `nosniff`. This prevents the browser from MIME-sniffing the response away from the declared content-type, which reduces exposure to drive-by download attacks and user-uploaded content that could be treated as executable or dynamic HTML files.

## False Positive Alerts Explanation

This document explains why each type of alert raised by the security scanner is a false positive in our application.

### Path Traversal

The Path Traversal alert is raised when the application allows the user to access files outside the intended directory. This can lead to various security vulnerabilities if the user can access sensitive files on the server.

In our application, we don't allow the user to access files directly. We do not access files anywhere in the application, and we don't use user-supplied input to access files. This prevents path traversal vulnerabilities.

### SQL Injection

The SQL Injection alert is raised when user input is used in a SQL query without proper validation and sanitization. This can lead to various security vulnerabilities if the user input includes SQL injection payloads.

In our application, we use grom to interact with the database, which automatically sanitizes user input to prevent SQL injection. We also use prepared statements to prevent SQL injection vulnerabilities. This ensures that user input is properly validated and sanitized before it's used in a SQL query.

More info at the [GORM documentation](https://gorm.io/docs/security.html)

### Application Error Disclosure

The Application Error Disclosure alert is raised when the application reveals sensitive information in error messages. This can lead to various security vulnerabilities if the error messages include sensitive information like stack traces or file paths.

In our application, we've taken measures to prevent sensitive information from being revealed in error messages. We use a custom error handler to handle all errors and return generic error messages to the user. This ensures that sensitive information is not revealed in error messages, or we documented the error messages to be shown to the user.

### Format String Error

The Format String Error alert is raised when user input is used in a format string function like `printf`. This can lead to various security vulnerabilities if the user input includes format specifiers.

In our application, we don't use any user-supplied input as a format string. The user input is properly validated and sanitized before it's used, which prevents format string vulnerabilities.

### User-Agent Fuzzing

The User-Agent Fuzzing alert is raised when the application behaves differently based on the User-Agent string in the request. This can lead to various security vulnerabilities if the application reveals sensitive information or behaves insecurely for certain User-Agents.

In our application, we don't use the User-Agent string to change the behavior of the application. The User-Agent string is only used for logging and analytics, and it doesn't affect the security of the application.

## Security Issues Summary

The table below summarizes the security issues found by the security scanner in our application.

| Issue | Severity | Number of Occurrences |
| --- | --- | --- |
| Path Traversal | High | 3 |
| SQL Injection | High | 4 |
| Format String Error | Medium | 1 |
| Application Error Disclosure | Low | 9 |
| X-Content-Type-Options Header Missing | Low | 5 |

The security scanner found a total of 22 security issues in our application. We've explained why each type of alert is a false positive in our application, and we've taken measures to prevent these security issues from being exploited.

## Conclusion

In this document, we've explained the security issues found by the security scanner in our application and why each type of alert is a false positive. We've also taken measures to prevent these security issues from being exploited.

The security scanner has helped us identify potential security vulnerabilities in our application, and we've taken steps to fix these issues and improve the security of our application.
