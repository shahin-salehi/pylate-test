FROM python:3.12-slim-bookworm

# Set working directory
WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy source code
COPY . .

# Use uvicorn to serve FastAPI app (assumes app is in server.py and named `app`)
CMD ["uvicorn", "server:app", "--host", "0.0.0.0", "--port", "8000"]
