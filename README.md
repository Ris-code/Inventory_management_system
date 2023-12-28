# Inventory_management_system
> Empower your college clubs/societies with precision and ease – revolutionize your inventory management effortlessly!

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)

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
>- Added items can be easily removed with the `Remove Item` button for a seamless customization experience.
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
>- Users can efficiently select individual items or use `Return All` for a bulk return.
>- Upon confirmation, the system updates the database and sends an informative return email to the club.

## Features to the Club/Society Coordinators

### Coordinator Authentication
![Feature 5](https://i.postimg.cc/bvNNKKmF/club-signin.png)
>- Coordinators can sign in to their respective clubs using a unique code.
>- This unique code is predefined in the backend SQL table.

### Club/Society Dashboard
![Feature 5](https://i.postimg.cc/3RmHcTsG/club-home.png)
>- This serves as the club dashboard, offering options to view the borrower's list, add inventory, modify club profile, and edit the inventory list.

### Edit Inventory List
![Feature 6](https://i.postimg.cc/RV2tx1yz/edit-inventory-list.png)
>- This page displays the list of all current club inventories.
>- Users can edit the quantity of items, reflecting backend updates.
>- Additionally, inventories can be deleted here, removing them from both the frontend and backend.

### List of Borrowers
![Feature 7](https://i.postimg.cc/2jQhv9H9/list-of-borrower.png)
![Feature 7](https://i.postimg.cc/FzVS41dX/list-of-borrower-item-view.png)
>- This page presents a list of borrowers who have borrowed items from the specific club, with data fetched from the MongoDB backend.
>- By clicking on the `Item List`, users can view all the items borrowed by the respective student.

### Add inventory
![Feature 8](https://i.postimg.cc/ZKZbBVgc/add-inventory.png)
>- Clicking on `Add Inventory` opens a pop-up where users can input new items along with quantities, updating the backend data.
>- The system does not accept items that are already present in the inventory list.

### Change Club Profile
![Feature 9](https://i.postimg.cc/RZKLsCG9/change-profile.png)
>- Clicking on `Change Profile` opens a pop-up where users can input club description, club profile image, and club email ID for receiving mails.
>- Users only need to fill the fields they want to change.
>- Clicking `Update` will update the club profile in the backend.

## Database Design

### SQL Database name is Club

#### Database Table Description

- `Clubs` : This table is designed for the purpose to store clubs profile information.
  - `Description`:
    - `club_id`  varchar(10) PRIMARY KEY
    - `club` varchar(100)
    - `Info` varchar(2000)
    - `Img_link` varchar(255)
    - `email` varchar(500)
    - `unique_id` varchar(100)

- `Items` : This table is designed for the purpose to store all the items with a item id.
  - `Description`:
    - `item_id`  varchar(50) PRIMARY KEY
    - `item` varchar(500)
    - `club_id` varchar(10) FOREIGN KEY REFERENCE from Clubs(club_id)
    - `quantity` int

- `Student` : This table is designed for the purpose to store all the student information.
  - `Description`:
    - `username`  varchar(50) PRIMARY KEY
    - `name` varchar(100)
    - `password` varchar(100)
    - `Institute_id` varchar(50)

### MongoDB Document Design
```
{
  "_id": {
    "$oid": ""
  },
  "username": "",
  "name": "",
  "institute_id": "",
  "club_info": [
    {
      "club_id": "",
      "club": "",
      "borrow_status": "",
      "items": [
        {
          "name": "",
          "quantity": ,
          "return_date": ""
        }
      ]
    },
    {
      "club_id": "",
      "club": "",
      "borrow_status": "",
      "items": [
        {
          "name": "",
          "quantity": ,
          "return_date": ""
        }
      ]
    },  
  ]
}
```
- `_id`: The unique identifier for the document in the MongoDB collection, automatically generated.
- `username`: The username associated with the user (e.g., "as").
- `name`: The name of the user (e.g., "Ashu").
- `institute_id`: The institute ID of the user (e.g., "B21AI007").
- `club_info`: An array containing information about the user's involvement in different clubs.

  For each `club_info` object within the array:
  - `club_id`: The unique identifier for the club.
  - `club`: The name of the club (e.g., "Robotics Society").
  - `borrow_status`: Indicates whether the user has borrowed items from the club ("Yes" or empty for no borrowing).
  - `items`: An array containing details about the items borrowed by the user from the club.

    For each `items` object within the array:
    - `name`: The name of the borrowed item (e.g., "Temperature sensor").
    - `quantity`: The quantity of the borrowed item.
    - `return_date`: The date by which the borrowed item is expected to be returned.

## Getting Started

Instructions on how to set up and run your project locally. Include prerequisites and installation steps.

### Prerequisites

- Atleast have sql and mongoDB installed in local and if its on cloud then well and good.

### Installation

Clone the repository and get into the directory havong the project files.

#### 1. Download Dependencies

```bash
go mod download
```
#### 2. Build the Project

```bash
go build
```

#### 3. Run the Project

```bash
go run .
```


