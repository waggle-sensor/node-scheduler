FROM waggle/plugin-base-light:0.1.0

COPY requirements.txt /app/
COPY kb.py server.py /app/
RUN pip3 install -r requirements.txt

WORKDIR /app/
ENTRYPOINT ["/usr/bin/python3 /app/server.py"]