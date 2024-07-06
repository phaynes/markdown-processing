FROM texlive/texlive:latest

RUN apt-get update && \
    apt-get install -y gnupg && \
    apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 648ACFD622F3D138 0E98404D386FA1D9 DCC9EFBF77E11517 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN apt-get update && \
    apt-get install -y \
    fonts-liberation \
    pandoc \
    python3 \
    wget \
    texlive-fonts-recommended \
    texlive-latex-extra \
    texlive-xetex \
    texlive-publishers \
    texlive-bibtex-extra \
    biber && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN fc-cache -f -v

RUN mkdir -p /root/.pandoc/templates

COPY apa.csl /root/.pandoc/apa.csl
COPY apa7.latex /root/.pandoc/templates/

WORKDIR /app

COPY convert.py /app/convert.py

ENTRYPOINT ["python3", "/app/convert.py"]
