FROM texlive/texlive:latest

# Install necessary packages and import GPG keys
RUN apt-get update && \
    apt-get install -y gnupg wget && \
    wget -qO- https://ftp-master.debian.org/keys/archive-key-11.asc | gpg --dearmor > /usr/share/keyrings/debian-archive-keyring.gpg && \
    wget -qO- https://ftp-master.debian.org/keys/archive-key-12.asc | gpg --dearmor > /usr/share/keyrings/debian-archive-keyring-12.gpg && \
    wget -O debian-archive-keyring.deb http://ftp.debian.org/debian/pool/main/d/debian-archive-keyring/debian-archive-keyring_2023.3+deb12u1_all.deb && \
    dpkg -i debian-archive-keyring.deb && \
    rm debian-archive-keyring.deb && \
    echo "deb [signed-by=/usr/share/keyrings/debian-archive-keyring.gpg] http://deb.debian.org/debian testing main" > /etc/apt/sources.list && \
    echo "deb [signed-by=/usr/share/keyrings/debian-archive-keyring.gpg] http://deb.debian.org/debian testing-updates main" >> /etc/apt/sources.list && \
    echo "deb [signed-by=/usr/share/keyrings/debian-archive-keyring.gpg] http://deb.debian.org/debian-security testing-security main" >> /etc/apt/sources.list && \
    rm -f /etc/apt/sources.list.d/debian.sources && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install required packages
RUN apt-get update && \
    apt-get install -y \
    fonts-liberation \
    pandoc \
    python3 \
    texlive-fonts-recommended \
    texlive-latex-extra \
    texlive-xetex \
    texlive-publishers \
    texlive-bibtex-extra \
    texlive-latex-recommended \
    biber && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN fc-cache -f -v

RUN mkdir -p /usr/local/share/pandoc/data
RUN mkdir -p /scripts

# Copy Pandoc resource files to the new directory
COPY apa.csl /usr/local/share/pandoc/data/
COPY apa7.latex /usr/local/share/pandoc/data/

# Copy the conversion script
COPY convert.py /scripts/convert.py

# Set the working directory to /data
WORKDIR /data

# Update the ENTRYPOINT to use the new script location
ENTRYPOINT ["python3", "/scripts/convert.py"]
