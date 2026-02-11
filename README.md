# 🔐 CIS Compliance Monitoring Platform

An automated security compliance monitoring platform that performs CIS benchmark checks on Linux systems, ingests compliance reports via a backend API, and visualizes results through a dashboard.

This project demonstrates a full-stack security monitoring pipeline including a system agent, backend API, and dashboard UI.

---

## 🚀 Project Overview

Modern organizations must continuously monitor system security configurations to maintain compliance with industry standards like CIS benchmarks.

This platform:

- Runs automated CIS security checks on Linux systems
- Collects compliance reports centrally
- Stores and processes security data
- Visualizes compliance status through a dashboard

---

## 🏗️ Architecture

CIS Agent (Go)
↓
FastAPI Backend (Python)
↓
Dashboard UI (Static/Next.js)


### Components:

### 1️⃣ CIS Agent
- Developed using Go
- Performs security checks aligned with CIS benchmarks
- Generates compliance reports in JSON format
- Sends reports to backend via REST API

Example checks:

- SSH configuration validation
- Audit daemon status
- IP forwarding settings
- System security configurations

---

### 2️⃣ Backend API (FastAPI)
Handles:

- Report ingestion
- Data storage
- API endpoints for dashboard
- Health monitoring

Key endpoints:
POST /api/v1/reports → Receive compliance reports
GET /hosts → List monitored systems
GET /reports/{host_id} → Retrieve reports
GET /health → Backend status

▶️ Running the Backend
cd backend
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn app:app --host 0.0.0.0 --port 8000


Access API docs:


▶️ Running CIS Agent
cd agent
go build -o cis-agent ./cmd/agent
./cis-agent


This generates and sends a compliance report to backend.


