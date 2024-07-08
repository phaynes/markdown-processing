FROM python:3.9-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY apa7_to_biblatex.py .

CMD ["python", "apa7_to_biblatex.py"]
