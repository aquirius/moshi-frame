# Sprout - Backend

### Sprout is a private hydroponic monitoring and nutient correction software.

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Installation](#installation)

---

## Introduction

Sprout aims at monitoring a hydroponic system, tracks the nutrient concentration and corrects it according to the crop which is planted. 

**Example**:  
This project is a web application designed to help users manage greenhouse environments by tracking crop growth, adjusting environmental parameters, and notifying users of important events.

---

## Features

- Real-time monitoring of environmental conditions (e.g., temperature, humidity)
- Manage and track multiple greenhouses and crops.
- Notification system for important updates.
- Easy-to-use interface for environmental adjustments.
- Support for indoor and outdoor greenhouse environments.
- Enrich nutrient solution according to crop

---

## Installation

**Prerequiries**
    golang 1.23
    docker
    make

1. **Clone the repository**:
   ```bash
   git clone https://github.com/aquirius/sprout-frame.git
   cd sprout-frame

2. **Add .env file**:
    ```txt
    # .env
    BACKEND_HOST=
    BACKEND_PORT=

    MYSQL_HOST=
    MYSQL_PORT=
    MYSQL_USER=
    MYSQL_PASSWORD=
    MYSQL_NAME=
    MYSQL_NETWORK=

    REDIS_HOST=
    REDIS_PORT=
    REDIS_USER=
    REDIS_PASSWORD=
    REDIS_DB=

3. **Docker compose**:
   ```bash
   docker-compose up -d

4. **create tables**:
   ```bash
   make load-schemas

5. **run it**:
   ```bash
   make run
