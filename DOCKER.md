# Docker Deployment Guide

## คำสั่ง Docker พื้นฐาน

### 1. Build และ Run ทั้งหมด (Postgres + App)

```bash
docker-compose up -d --build
```

### 2. ดู logs

```bash
# ดู logs ทั้งหมด
docker-compose logs -f

# ดู logs เฉพาะ app
docker-compose logs -f app

# ดู logs เฉพาะ postgres
docker-compose logs -f postgres
```

### 3. หยุดการทำงาน

```bash
docker-compose down
```

### 4. หยุดและลบ volumes (ข้อมูลทั้งหมด)

```bash
docker-compose down -v
```

## Build เฉพาะ Docker Image

หากต้องการ build เฉพาะ image โดยไม่ใช้ docker-compose:

```bash
# Build image
docker build -t scadul-app:latest .

# Run container
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=123456 \
  -e DB_NAME=app_db \
  --name scadul_app \
  scadul-app:latest
```

## การตั้งค่า Environment Variables

แก้ไขค่าต่างๆ ใน `docker-compose.yml` หรือสร้างไฟล์ `.env`:

```env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123456
DB_NAME=app_db

# Application (ปรับตามที่ต้องการ)
APP_PORT=8080
```

## คำสั่งที่มีประโยชน์

```bash
# Rebuild เฉพาะ app service
docker-compose up -d --build app

# เข้าไปใน container
docker exec -it scadul_app sh

# ดูสถานะ containers
docker-compose ps

# Restart service
docker-compose restart app
```

## โครงสร้าง Dockerfile

Dockerfile นี้ใช้ multi-stage build:
- **Stage 1 (builder)**: Build Go application
- **Stage 2 (final)**: รัน application บน Alpine Linux (ขนาดเล็ก)

ขนาดของ image จะประมาณ 20-30 MB แทน 800+ MB
