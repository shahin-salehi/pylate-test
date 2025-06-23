import logging
import sys

from contextlib import asynccontextmanager

from internal.colbert.embedder import Embedder
from internal.db.database import Database
from internal.parser.parse import Parse

from fastapi import FastAPI
from pydantic import BaseModel, HttpUrl


logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s [%(levelname)s] %(name)s: %(message)s"
        )

logger = logging.getLogger(__name__)

# if we develop this further move this to its own types/models folder
class UrlPayload(BaseModel):
    url: HttpUrl


"""
## parse test
    data, ok = parse.pdf("internal/parser/docs/Q125_Quarterly_report.pdf", "reports")
    if ok:
        logger.info("pdf parsed succesfully.")
    else:
        logger.error("Failed to parse PDF.")
        sys.exit(1)

## insert test

    resp, ok = db.insert_pdf(data)
    if not ok:
        logger.error("I just shitted my pants, ong")
        sys.exit(1)
    else:
        logger.info(f"pdf inserted id: {resp} ")


"""


#python time
# add trys to all this please
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

        yield
    except Exception as e:
        logger.error(f"init failed error: {e}")
        sys.exit(1)

app = FastAPI(lifespan=lifespan)

    
@app.get("/")
async def root():
    return {"message": "hello world"}

@app.post("/upload-pdf")
async def upload_pdf(payload: UrlPayload):
    print("doing work...")
    return {"received": str(payload.url)}

          

