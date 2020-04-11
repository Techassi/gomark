# gomark

Web-based bookmark manager written in GO. CURRENTLY IN DEVELOPMENT (NOT FEATURE-COMPLETE / PRODUCTION READY)

## Features

-   Secure Login Flow with the optional use of 2FA via TOTP (via Google Authenticator e.g.)
-   Argon2 Password Hashing

## Current Development

-   API Routes
-   Scheduler, which does multi-threaded work in the background
-   Add image resizer

## Planned Features

-   Management of Bookmarks, Notes and Folders
-   Archive pages and save them on disk
    -   Style and JavaScipt donwloader
    -   Simplify HTML structure
-   User Management
-   Sharing of Bookmarks, Notes and Folders

## Todo

-   Remove Temp2FAToken from DB after login
