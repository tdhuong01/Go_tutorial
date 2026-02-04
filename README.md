# My Go + React Project

## Cấu trúc thư mục
- `client/`: Frontend React app
- `server/`: Backend Go API

## Chạy dự án

### Backend (Server)
```bash
cd server
go run .
```
Server chạy trên http://localhost:8081

### Frontend (Client)
```bash
cd client
npm install
npm run dev
```
React app chạy trên http://localhost:5174

## API Endpoints
- POST /api/register: Đăng ký
- POST /api/login: Đăng nhập
- POST /api/logout: Đăng xuất
- GET /: Trang chủ (cần auth)

## Database
Sử dụng PostgreSQL, cấu hình trong server/.env