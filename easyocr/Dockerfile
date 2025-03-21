FROM python:3.9-slim
RUN apt-get update && apt-get install -y libgl1-mesa-glx libglib2.0-0 poppler-utils
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY app.py .
EXPOSE 8866
CMD ["python", "app.py"]