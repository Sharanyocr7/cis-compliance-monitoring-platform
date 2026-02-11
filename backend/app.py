from datetime import datetime
from typing import Any, Dict, List, Optional
import json

from fastapi import FastAPI, Header, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from sqlalchemy import Column, DateTime, Integer, String, Text, create_engine
from sqlalchemy.orm import declarative_base, sessionmaker


# ===============================
# CONFIGURATION
# ===============================
DB_URL = "sqlite:///./cis.db"
API_KEY = "change-me-strong-key"  # must match agent config.yaml

engine = create_engine(DB_URL, connect_args={"check_same_thread": False})
SessionLocal = sessionmaker(bind=engine, autocommit=False, autoflush=False)
Base = declarative_base()


# ===============================
# DATABASE MODEL
# ===============================
class ReportRow(Base):
    __tablename__ = "reports"

    id = Column(Integer, primary_key=True, autoincrement=True)
    host_id = Column(String, index=True)
    hostname = Column(String, index=True)
    ts_utc = Column(String, index=True)
    payload_json = Column(Text)
    created_at = Column(DateTime, default=datetime.utcnow)


Base.metadata.create_all(bind=engine)


# ===============================
# FASTAPI APP
# ===============================
app = FastAPI(title="CIS Compliance Backend", version="1.1")

# CORS for browser dashboard (demo-friendly; restrict in production)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# ===============================
# Pydantic Payload Schema
# ===============================
class IngestPayload(BaseModel):
    agent: Dict[str, Any]
    host: Dict[str, Any]
    packages: List[Dict[str, Any]] = []
    cis_results: List[Dict[str, Any]] = []


# ===============================
# HEALTH CHECK
# ===============================
@app.get("/health")
def health():
    return {"status": "ok", "time": datetime.utcnow().isoformat()}


# ===============================
# INTERNAL SAVE FUNCTION
# ===============================
def save_report(payload: IngestPayload):
    host = payload.host
    agent = payload.agent

    db = SessionLocal()
    try:
        row = ReportRow(
            host_id=str(host.get("host_id", "")),
            hostname=str(host.get("hostname", "")),
            ts_utc=str(agent.get("ts_utc", "")),
            payload_json=payload.model_dump_json(),
        )
        db.add(row)
        db.commit()
        db.refresh(row)
        return {"ok": True, "report_id": row.id}
    finally:
        db.close()


# ===============================
# INGEST ENDPOINTS
# ===============================
@app.post("/api/v1/reports")
def ingest_v1(payload: IngestPayload, x_api_key: Optional[str] = Header(default=None)):
    if x_api_key != API_KEY:
        raise HTTPException(status_code=401, detail="Invalid API key")
    return save_report(payload)


@app.post("/ingest")
def ingest_alias(payload: IngestPayload, x_api_key: Optional[str] = Header(default=None)):
    if x_api_key != API_KEY:
        raise HTTPException(status_code=401, detail="Invalid API key")
    return save_report(payload)


# ===============================
# LIST HOSTS
# ===============================
@app.get("/hosts")
def list_hosts():
    db = SessionLocal()
    try:
        rows = db.query(ReportRow.host_id, ReportRow.hostname).distinct().all()
        return [{"host_id": r[0], "hostname": r[1]} for r in rows]
    finally:
        db.close()


# ===============================
# LIST REPORTS (META) FOR A HOST
# ===============================
@app.get("/reports/{host_id}")
def list_reports(host_id: str, limit: int = 10):
    db = SessionLocal()
    try:
        rows = (
            db.query(ReportRow)
            .filter(ReportRow.host_id == host_id)
            .order_by(ReportRow.id.desc())
            .limit(limit)
            .all()
        )
        return [
            {
                "id": r.id,
                "ts_utc": r.ts_utc,
                "created_at": r.created_at.isoformat(),
            }
            for r in rows
        ]
    finally:
        db.close()


# ===============================
# GET FULL REPORT (FOR DASHBOARD)
# ===============================
@app.get("/report/{report_id}")
def get_report(report_id: int):
    db = SessionLocal()
    try:
        row = db.query(ReportRow).filter(ReportRow.id == report_id).first()
        if not row:
            raise HTTPException(status_code=404, detail="Not found")

        return {
            "id": row.id,
            "host_id": row.host_id,
            "hostname": row.hostname,
            "ts_utc": row.ts_utc,
            "created_at": row.created_at.isoformat(),
            "payload": json.loads(row.payload_json),
        }
    finally:
        db.close()
