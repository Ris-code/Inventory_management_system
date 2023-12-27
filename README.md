# Inventory_management_system
> Empower your college clubs/societies with precision and ease – revolutionize your inventory management effortlessly!

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Screenshots](#screenshots)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)

## Introduction

A concise introduction to your project. Explain what problem it solves and why it's beneficial. Highlight any unique selling points.

## Features
- The application is developed for two types of users
  - Students
  - Club/Society Coordinators of College

## Features to the Students
### Student Authentication

![Dashboard](https://i.postimg.cc/jSNfPSS0/signin.jpg)
![Dashboard](https://i.postimg.cc/9FsTD7VY/signup.png)
 >- The application employs a secure user authentication mechanism using a combination of username and hashed passwords.
 >- During the sign-up process, the user's password undergoes one-way hashing using the bcrypt algorithm before being stored in the database.
 >- This ensures that sensitive user credentials are not stored in plaintext, adding a layer of security to protect user accounts.

### Club/Society Selection

![Feature 1](https://i.postimg.cc/NML1qRRh/choose-club.png)
>- Users can choose a club from the dynamically populated list on the webpage.
>- The displayed club options are updated in real-time based on the latest information from the SQL table.

### Inventory Interaction Feature

![Feature 2](https://i.postimg.cc/4y24Chv0/item-borrow.png)
>- The inventory dynamically displays items fetched from the backend, showcasing their available quantities.
>- Users have the flexibility to set the desired quantity for each item before adding them to their selection.
>- Added items can be easily removed with the "Remove Item" button for a seamless customization experience.
>- Once users finalize their selections, they can proceed to the borrow cart by clicking the "Borrow" button.

### Borrowing Workflow

![Feature 3](https://i.postimg.cc/3JDWMn6C/borrow-review-page.png)
![Feature 3](https://i.postimg.cc/yxDhC04g/borrow-email.png)
>- The selected items, along with their respective quantities, are conveniently displayed for review on this page.
>- It is mandatory for users to enter a return date before proceeding with the borrowing process.
>- Upon confirming the borrow action, the item data is updated in the user's document within a MongoDB collection.
>- Simultaneously, an email containing details such as the student's name, institute ID, and the borrowed items is sent to the corresponding club's email address.

### Item Return Workflow

![Feature 4](https://i.postimg.cc/bYnHgHTW/return-workflow.png)
![Feature 4](https://i.postimg.cc/CxRKCJCf/return-email.png)
>- Users can choose a club from a dynamic dropdown menu populated with available clubs.
>- Detailed information about previously borrowed items from the selected club is fetched and displayed.
>- Users can efficiently select individual items or use "Return All" for a bulk return.
>- Upon confirmation, the system updates the database and sends an informative return email to the club.

## Features to the Club/Society Coordinators

### Authentication


## Getting Started

Provide instructions on how to set up and run your project locally. Include prerequisites and installation steps.

### Prerequisites

List any software, libraries, or dependencies that need to be installed before running your project.

### Installation

1. Step-by-step installation instructions.

```bash
# Example installation command
npm install
