import logging
import sys
import os
import requests
import tempfile

from contextlib import asynccontextmanager

from internal.colbert.embedder import Embedder
from internal.db.database import Database
from internal.parser.parse import Parse

from fastapi import FastAPI, HTTPException, BackgroundTasks, Request, status
from pydantic import BaseModel, HttpUrl


logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s"
        )

logger = logging.getLogger(__name__)

# if we develop this further move this to its own types/models folder
class UrlPayload(BaseModel):
    url: str
    category: str
    filename: str
    owner:int


#python time
# add trys to all this please

def parse_and_insert_pdf(owner:int, path: str, file_url:str, filename:str, category: str, db, parse):
    try:
        data, ok = parse.pdf(owner, path, file_url, filename, category)
        if not ok:
            logger.error("Background: Failed to parse PDF.")
            return

        resp, ok = db.insert_pdf(data)
        if not ok:
            logger.error("Background: Failed to insert PDF.")
        else:
            logger.info(f"Background: PDF inserted, ID: {resp}")
    except Exception as e:
        logger.exception(f"Background: Unexpected error during parse/insert {e}")
    finally:
        if os.path.exists(path):
            os.remove(path)



@asynccontextmanager
async def lifespan(app: FastAPI):
    try: 
        embedder = Embedder("model/model.onnx")

        db = Database("postgres://admin:password@localhost:9876/documents")
        if db.ping():
            logger.info("ping test successfull")
        else:
            logger.error("ping test failed, shutting down.")
            sys.exit(1)

        # init parser
        parse = Parse(embedder)

        # Store in app.state
        app.state.db = db
        app.state.parse = parse


        yield
    except Exception as e:
        logger.error(f"init failed error: {e}")
        sys.exit(1)

app = FastAPI(lifespan=lifespan)

    
@app.get("/")
async def root():
    return {"message": "hello world"}

@app.post("/upload-pdf", status_code=status.HTTP_202_ACCEPTED)
async def upload_pdf(payload: UrlPayload, background_tasks: BackgroundTasks, request: Request):
    try:
        print("payload url", payload.url)
            # Download the PDF
        response = requests.get(payload.url, timeout=10)
        response.raise_for_status()

        # Save to temp file
        with tempfile.NamedTemporaryFile(suffix=".pdf", delete=False) as tmp_file:
            tmp_file.write(response.content)
            tmp_path = tmp_file.name

        # Schedule background processing
        background_tasks.add_task(
                parse_and_insert_pdf,
                payload.owner,
                tmp_path,
                payload.url,
                payload.filename,
                payload.category,
                request.app.state.db,
                request.app.state.parse
                )

        return {
            "status": "accepted",
            "message": "PDF is being parsed and inserted in the background.",
            "url": str(payload.url)
        }

    except requests.RequestException as e:
        logger.exception("Failed to fetch PDF")
        raise HTTPException(status_code=502, detail=f"Error downloading PDF: {str(e)}")

    except Exception as e:
        logger.exception("Unexpected error")
        raise HTTPException(status_code=500, detail=f"Internal server error: {str(e)}")
